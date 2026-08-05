package order

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/business/funding"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	merchantstock "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CallbackHandler struct {
	db              *gorm.DB
	configs         *paymentconfig.Store
	platformConfigs *cloudconfig.Service
	allowMock       bool
}

func NewCallbackHandler(db *gorm.DB, configs *paymentconfig.Store, allowMock bool, platformConfigs ...*cloudconfig.Service) *CallbackHandler {
	var platformConfig *cloudconfig.Service
	if len(platformConfigs) > 0 {
		platformConfig = platformConfigs[0]
	}
	return &CallbackHandler{db: db, configs: configs, platformConfigs: platformConfig, allowMock: allowMock}
}
func (h *CallbackHandler) Register(r gin.IRoutes) {
	r.POST("/pay/wechat", h.Wechat)
	r.POST("/refund/wechat", h.WechatRefund)
	if h.allowMock {
		r.POST("/pay/mock", h.Mock)
		r.POST("/refund/mock", h.MockRefund)
	}
}

// MockRefund is a sandbox-only provider simulator. It deliberately accepts no
// amount, channel or order number from the caller; every value is reloaded
// from locked internal rows. The route is not registered in production.
func (h *CallbackHandler) MockRefund(c *gin.Context) {
	if !h.allowMock {
		response.Fail(c, http.StatusNotFound, "接口不存在")
		return
	}
	var in struct {
		RefundID uint64 `json:"refund_id"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.RefundID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.markMockRefunded(c.Request.Context(), in.RefundID); err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"refund_id": in.RefundID, "status": "refunded"})
}

// WechatRefund accepts only a signed REFUND.SUCCESS notification. It never
// trusts an operator action or a browser payload as proof that funds arrived.
func (h *CallbackHandler) WechatRefund(c *gin.Context) {
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "invalid request"})
		return
	}
	refund, config, err := h.decryptWechatRefundCallback(c.Request.Context(), c.Request.Header, body)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, wechatpayv3.ErrInvalidSignature) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"code": "FAIL", "message": "invalid callback"})
		return
	}
	if refund.MchID != config.MchID || refund.Status != "SUCCESS" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "invalid refund"})
		return
	}
	if err := h.markWechatRefunded(c.Request.Context(), refund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "refund validation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "成功"})
}
func (h *CallbackHandler) Mock(c *gin.Context) {
	if !h.allowMock {
		response.Fail(c, http.StatusNotFound, "接口不存在")
		return
	}
	var in struct {
		GroupOrderID uint64 `json:"group_order_id"`
		UID          uint64 `json:"uid"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.GroupOrderID == 0 || in.UID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	created, err := PayMock(c.Request.Context(), h.db, in.UID, in.GroupOrderID)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, createdResponse(created, true))
}

// Wechat accepts only an authenticated WeChat v3 callback. Browser JSON,
// mock HMAC payloads and client-provided amounts are intentionally rejected.
func (h *CallbackHandler) Wechat(c *gin.Context) {
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "invalid request"})
		return
	}
	transaction, config, err := h.decryptWechatCallback(c.Request.Context(), c.Request.Header, body)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, wechatpayv3.ErrInvalidSignature) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"code": "FAIL", "message": "invalid callback"})
		return
	}
	if transaction.TradeState != "SUCCESS" || transaction.MchID != config.MchID || transaction.AppID != config.AppID {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "invalid transaction"})
		return
	}
	if err := h.markWechatPaid(c.Request.Context(), transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "payment validation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "成功"})
}

func (h *CallbackHandler) decryptWechatCallback(ctx context.Context, headers http.Header, body []byte) (wechatpayv3.CallbackTransaction, wechatpayv3.Config, error) {
	platformValues, err := h.loadPlatformValues(ctx)
	if err != nil {
		return wechatpayv3.CallbackTransaction{}, wechatpayv3.Config{}, err
	}
	platformConfig := wechatConfig(platformValues, true)
	if err := wechatpayv3.VerifyCallbackSignature(platformConfig, headers, body, time.Now().UTC()); err != nil {
		return wechatpayv3.CallbackTransaction{}, wechatpayv3.Config{}, err
	}
	candidates := []wechatpayv3.Config{platformConfig}
	if h.configs != nil {
		stores, err := h.configs.LoadStores(ctx)
		if err != nil {
			return wechatpayv3.CallbackTransaction{}, wechatpayv3.Config{}, err
		}
		for _, values := range stores {
			candidate := wechatConfig(values, false)
			// WeChat signs all notifications with its platform key. Store
			// accounts own their APIv3 key but do not inherit platform payment
			// credentials or callback URLs.
			candidate.PublicKeyID, candidate.PublicKeyPEM = platformConfig.PublicKeyID, platformConfig.PublicKeyPEM
			candidate.PlatformCertSerial, candidate.PlatformCertPEM = platformConfig.PlatformCertSerial, platformConfig.PlatformCertPEM
			candidates = append(candidates, candidate)
		}
	}
	for _, candidate := range candidates {
		transaction, err := wechatpayv3.VerifyAndDecryptCallback(candidate, headers, body, time.Now().UTC())
		if err == nil {
			return transaction, candidate, nil
		}
		if errors.Is(err, wechatpayv3.ErrInvalidSignature) || errors.Is(err, wechatpayv3.ErrInvalidConfig) {
			return wechatpayv3.CallbackTransaction{}, wechatpayv3.Config{}, err
		}
	}
	return wechatpayv3.CallbackTransaction{}, wechatpayv3.Config{}, wechatpayv3.ErrInvalidCallback
}

func (h *CallbackHandler) decryptWechatRefundCallback(ctx context.Context, headers http.Header, body []byte) (wechatpayv3.CallbackRefund, wechatpayv3.Config, error) {
	platformValues, err := h.loadPlatformValues(ctx)
	if err != nil {
		return wechatpayv3.CallbackRefund{}, wechatpayv3.Config{}, err
	}
	platformConfig := wechatConfig(platformValues, true)
	if err := wechatpayv3.VerifyCallbackSignature(platformConfig, headers, body, time.Now().UTC()); err != nil {
		return wechatpayv3.CallbackRefund{}, wechatpayv3.Config{}, err
	}
	candidates := []wechatpayv3.Config{platformConfig}
	if h.configs != nil {
		stores, err := h.configs.LoadStores(ctx)
		if err != nil {
			return wechatpayv3.CallbackRefund{}, wechatpayv3.Config{}, err
		}
		for _, values := range stores {
			candidate := wechatConfig(values, false)
			candidate.PublicKeyID, candidate.PublicKeyPEM = platformConfig.PublicKeyID, platformConfig.PublicKeyPEM
			candidate.PlatformCertSerial, candidate.PlatformCertPEM = platformConfig.PlatformCertSerial, platformConfig.PlatformCertPEM
			candidates = append(candidates, candidate)
		}
	}
	for _, candidate := range candidates {
		refund, err := wechatpayv3.VerifyAndDecryptRefundCallback(candidate, headers, body, time.Now().UTC())
		if err == nil {
			return refund, candidate, nil
		}
		if errors.Is(err, wechatpayv3.ErrInvalidSignature) || errors.Is(err, wechatpayv3.ErrInvalidConfig) {
			return wechatpayv3.CallbackRefund{}, wechatpayv3.Config{}, err
		}
	}
	return wechatpayv3.CallbackRefund{}, wechatpayv3.Config{}, wechatpayv3.ErrInvalidCallback
}

func (h *CallbackHandler) loadPlatformValues(ctx context.Context) (paymentconfig.Values, error) {
	if h.platformConfigs != nil {
		values, err := h.platformConfigs.Values(ctx, "payment")
		return paymentconfig.Values(values), err
	}
	if h.configs == nil {
		return nil, paymentconfig.ErrNotConfigured
	}
	return h.configs.Load(ctx)
}

func (h *CallbackHandler) markWechatPaid(ctx context.Context, notify wechatpayv3.CallbackTransaction) error {
	if handled, err := funding.MarkWechatPaid(ctx, h.db, notify); handled {
		return err
	}
	return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment paymentRow
		// Read the payment first only to locate its group. The group lock is the
		// transaction serialization point shared with cancellation/expiry, avoiding
		// payment-first vs group-first deadlocks.
		if err := tx.Where("transaction_no = ? AND channel = ?", notify.OutTradeNo, "wechat").First(&payment).Error; err != nil {
			return err
		}
		if decimalCents(payment.Amount) != notify.AmountCents {
			return ErrOrderNotPayable
		}
		var group groupRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", payment.GroupOrderID).First(&group).Error; err != nil {
			return err
		}
		if decimalCents(group.PayAmount) != notify.AmountCents {
			return ErrOrderNotPayable
		}
		if group.ActivityType == pointsOrderActivityType {
			return ErrOrderNotPayable
		}

		var existing paymentCallbackRow
		err := tx.Where("channel = ? AND provider_event_id = ?", "wechat", notify.EventID).First(&existing).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		payload, err := json.Marshal(json.RawMessage(notify.RawPayload))
		if err != nil {
			return err
		}
		callback := paymentCallbackRow{Channel: "wechat", ProviderEventID: notify.EventID, TransactionNo: notify.OutTradeNo, Payload: payload, Verified: true, ProcessedAt: &notify.SuccessTime}
		if err := tx.Create(&callback).Error; err != nil {
			return err
		}

		if payment.Status == "succeeded" && group.PayStatus == "paid" {
			return nil
		}
		if group.PayStatus != "pending" {
			return ErrOrderNotPayable
		}
		if err := merchantstock.ReservationsReady(tx, group.ID); err != nil {
			return err
		}
		now := notify.SuccessTime
		updated := tx.Model(&groupRow{}).Where("id = ? AND pay_status = ?", group.ID, "pending").Updates(map[string]any{"pay_status": "paid", "pay_channel": "wechat", "paid_at": now})
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrOrderNotPayable
		}
		if err := tx.Model(&orderRow{}).Where("group_order_id = ? AND status = ?", group.ID, "pending_pay").Updates(map[string]any{"status": "paid", "paid_at": now}).Error; err != nil {
			return err
		}
		if err := enqueueStockActionForGroup(tx, "confirm", group.ID); err != nil {
			return err
		}
		if err := settleOrderActivity(tx, group); err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_coupon_user").Where("user_id = ? AND used_order_id = ? AND status = ?", group.UserID, group.ID, "locked").Update("status", "used").Error; err != nil {
			return err
		}
		if err := issueVerificationsForPaidGroup(tx, group.ID); err != nil {
			return err
		}
		return tx.Model(&paymentRow{}).Where("id = ?", payment.ID).Updates(map[string]any{"status": "succeeded", "provider_transaction_no": notify.TransactionID, "callback_idempotency_key": notify.EventID, "paid_at": now}).Error
	})
}

// markWechatRefunded is the only transition into the money-success terminal
// state for WeChat refunds. It locks the refund, its original payment and the
// group order in one transaction, records the provider event exactly once and
// refuses callbacks that do not match all server-owned identifiers and sums.
func (h *CallbackHandler) markWechatRefunded(ctx context.Context, notify wechatpayv3.CallbackRefund) error {
	return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID       uint64  `gorm:"column:id"`
			OrderID  uint64  `gorm:"column:order_id"`
			RefundNo string  `gorm:"column:refund_no"`
			Amount   float64 `gorm:"column:amount"`
			Status   string  `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_refund").Where("refund_no = ?", notify.OutRefundNo).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&refund).Error; err != nil {
			return err
		}
		if refund.ID == 0 || refund.RefundNo != notify.OutRefundNo || decimalCents(refund.Amount) != notify.RefundCents {
			return ErrOrderNotPayable
		}
		var transaction struct {
			ID       uint64  `gorm:"column:id"`
			RefundID uint64  `gorm:"column:refund_id"`
			Channel  string  `gorm:"column:channel"`
			Amount   float64 `gorm:"column:amount"`
			Status   string  `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND channel = ?", refund.ID, "wechat").Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&transaction).Error; err != nil {
			return err
		}
		if transaction.ID == 0 || decimalCents(transaction.Amount) != notify.RefundCents {
			return ErrOrderNotPayable
		}
		var payment struct {
			ID                uint64  `gorm:"column:id"`
			GroupOrderID      uint64  `gorm:"column:group_order_id"`
			TransactionNo     string  `gorm:"column:transaction_no"`
			ProviderReference string  `gorm:"column:provider_transaction_no"`
			Amount            float64 `gorm:"column:amount"`
			Status            string  `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_order AS o").Select("p.id,p.group_order_id,p.transaction_no,p.provider_transaction_no,p.amount,p.status").Joins("JOIN qixi_crm_b_payment_transaction AS p ON p.group_order_id = o.group_order_id AND p.channel = 'wechat'").Where("o.id = ?", refund.OrderID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&payment).Error; err != nil {
			return err
		}
		if payment.ID == 0 || payment.Status != "succeeded" || payment.TransactionNo != notify.OutTradeNo || payment.ProviderReference != notify.TransactionID || decimalCents(payment.Amount) != notify.TotalCents {
			return ErrOrderNotPayable
		}
		var existing refundCallbackRow
		err := tx.Where("channel = ? AND provider_event_id = ?", "wechat", notify.EventID).First(&existing).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		payload, err := json.Marshal(json.RawMessage(notify.RawPayload))
		if err != nil {
			return err
		}
		if err := tx.Create(&refundCallbackRow{Channel: "wechat", ProviderEventID: notify.EventID, ProviderRefundNo: notify.RefundID, OutTradeNo: notify.OutTradeNo, Payload: payload, Verified: true, ProcessedAt: &notify.SuccessTime}).Error; err != nil {
			return err
		}
		if refund.Status == "refunded" && transaction.Status == "succeeded" {
			return nil
		}
		if refund.Status != "refunding" || (transaction.Status != "created" && transaction.Status != "processing") {
			return ErrOrderNotPayable
		}
		if err := tx.Table("qixi_crm_b_refund_transaction").Where("id = ?", transaction.ID).Updates(map[string]any{"status": "succeeded", "provider_refund_no": notify.RefundID, "completed_at": notify.SuccessTime}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ? AND status = ?", refund.ID, "refunding").Update("status", "refunded").Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Where("id = ? AND status = ?", refund.OrderID, "aftersale").Update("status", "cancelled").Error; err != nil {
			return err
		}
		if err := enqueueStockRestockForRefund(tx, refund.ID); err != nil {
			return err
		}
		if err := enqueueSettlementReversalForRefund(tx, refund.ID); err != nil {
			return err
		}
		var activeOrders int64
		if err := tx.Table("qixi_crm_b_order").Where("group_order_id = ? AND status <> ?", payment.GroupOrderID, "cancelled").Count(&activeOrders).Error; err != nil {
			return err
		}
		if activeOrders == 0 {
			if err := tx.Table("qixi_crm_b_group_order").Where("id = ? AND pay_status = ?", payment.GroupOrderID, "paid").Update("pay_status", "refunded").Error; err != nil {
				return err
			}
			if err := tx.Table("qixi_crm_b_payment_transaction").Where("id = ? AND status = ?", payment.ID, "succeeded").Update("status", "refunded").Error; err != nil {
				return err
			}
		}
		return tx.Table("qixi_crm_b_refund_event").Create(map[string]any{"refund_id": refund.ID, "from_status": "refunding", "to_status": "refunded", "actor_type": "system", "actor_id": 0, "reason": "微信退款回调验签成功", "idempotency_key": "wechat-refund:" + notify.EventID}).Error
	})
}

func (h *CallbackHandler) markMockRefunded(ctx context.Context, refundID uint64) error {
	return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID      uint64 `gorm:"column:id"`
			OrderID uint64 `gorm:"column:order_id"`
			Status  string `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ?", refundID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&refund).Error; err != nil {
			return err
		}
		if refund.ID == 0 {
			return ErrOrderNotPayable
		}
		providerEventID := "mock-refund:" + strconv.FormatUint(refund.ID, 10)
		var existing refundCallbackRow
		err := tx.Where("channel = ? AND provider_event_id = ?", "mock", providerEventID).First(&existing).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		var transaction struct {
			ID     uint64 `gorm:"column:id"`
			Status string `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_refund_transaction").Where("refund_id = ? AND channel = ?", refund.ID, "mock").Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&transaction).Error; err != nil {
			return err
		}
		if transaction.ID == 0 || (transaction.Status != "created" && transaction.Status != "processing") || refund.Status != "refunding" {
			return ErrOrderNotPayable
		}
		var payment struct {
			ID           uint64 `gorm:"column:id"`
			GroupOrderID uint64 `gorm:"column:group_order_id"`
			Status       string `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_b_order AS o").Select("p.id,p.group_order_id,p.status").Joins("JOIN qixi_crm_b_payment_transaction AS p ON p.group_order_id = o.group_order_id AND p.channel = 'mock'").Where("o.id = ?", refund.OrderID).Clauses(clause.Locking{Strength: "UPDATE"}).Scan(&payment).Error; err != nil {
			return err
		}
		if payment.ID == 0 || payment.Status != "succeeded" {
			return ErrOrderNotPayable
		}
		now := time.Now().UTC()
		payload, err := json.Marshal(map[string]any{"mode": "sandbox", "refund_id": refund.ID})
		if err != nil {
			return err
		}
		if err := tx.Create(&refundCallbackRow{Channel: "mock", ProviderEventID: providerEventID, ProviderRefundNo: "mock-refund-" + strconv.FormatUint(refund.ID, 10), OutTradeNo: "mock-group-" + strconv.FormatUint(payment.GroupOrderID, 10), Payload: payload, Verified: true, ProcessedAt: &now}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_refund_transaction").Where("id = ?", transaction.ID).Updates(map[string]any{"status": "succeeded", "provider_refund_no": "mock-refund-" + strconv.FormatUint(refund.ID, 10), "completed_at": now}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_refund").Where("id = ? AND status = ?", refund.ID, "refunding").Update("status", "refunded").Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Where("id = ? AND status = ?", refund.OrderID, "aftersale").Update("status", "cancelled").Error; err != nil {
			return err
		}
		if err := enqueueStockRestockForRefund(tx, refund.ID); err != nil {
			return err
		}
		if err := enqueueSettlementReversalForRefund(tx, refund.ID); err != nil {
			return err
		}
		var activeOrders int64
		if err := tx.Table("qixi_crm_b_order").Where("group_order_id = ? AND status <> ?", payment.GroupOrderID, "cancelled").Count(&activeOrders).Error; err != nil {
			return err
		}
		if activeOrders == 0 {
			if err := tx.Table("qixi_crm_b_group_order").Where("id = ? AND pay_status = ?", payment.GroupOrderID, "paid").Update("pay_status", "refunded").Error; err != nil {
				return err
			}
			if err := tx.Table("qixi_crm_b_payment_transaction").Where("id = ? AND status = ?", payment.ID, "succeeded").Update("status", "refunded").Error; err != nil {
				return err
			}
		}
		return tx.Table("qixi_crm_b_refund_event").Create(map[string]any{"refund_id": refund.ID, "from_status": "refunding", "to_status": "refunded", "actor_type": "system", "actor_id": 0, "reason": "sandbox 模拟退款回调成功", "idempotency_key": providerEventID}).Error
	})
}

type paymentCallbackRow struct {
	ID              uint64     `gorm:"column:id"`
	Channel         string     `gorm:"column:channel"`
	ProviderEventID string     `gorm:"column:provider_event_id"`
	TransactionNo   string     `gorm:"column:transaction_no"`
	Payload         []byte     `gorm:"column:payload"`
	Verified        bool       `gorm:"column:verified"`
	ProcessedAt     *time.Time `gorm:"column:processed_at"`
}

func (paymentCallbackRow) TableName() string { return "qixi_crm_b_payment_callback" }

type refundCallbackRow struct {
	ID               uint64     `gorm:"column:id"`
	Channel          string     `gorm:"column:channel"`
	ProviderEventID  string     `gorm:"column:provider_event_id"`
	ProviderRefundNo string     `gorm:"column:provider_refund_no"`
	OutTradeNo       string     `gorm:"column:out_trade_no"`
	Payload          []byte     `gorm:"column:payload"`
	Verified         bool       `gorm:"column:verified"`
	ProcessedAt      *time.Time `gorm:"column:processed_at"`
}

func (refundCallbackRow) TableName() string { return "qixi_crm_b_refund_callback" }
