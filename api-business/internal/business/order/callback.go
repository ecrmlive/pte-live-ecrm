package order

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cloudconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/wechatpayv3"
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
	if h.allowMock {
		r.POST("/pay/mock", h.Mock)
	}
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
	return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment paymentRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("transaction_no = ? AND channel = ?", notify.OutTradeNo, "wechat").First(&payment).Error; err != nil {
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
		if err := tx.Table("qixi_crm_b_coupon_user").Where("user_id = ? AND used_order_id = ? AND status = ?", group.UserID, group.ID, "locked").Update("status", "used").Error; err != nil {
			return err
		}
		return tx.Model(&paymentRow{}).Where("id = ?", payment.ID).Updates(map[string]any{"status": "succeeded", "provider_transaction_no": notify.TransactionID, "callback_idempotency_key": notify.EventID, "paid_at": now}).Error
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
