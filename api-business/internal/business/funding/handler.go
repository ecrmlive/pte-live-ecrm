// Package funding owns C-end recharge and SVIP purchase orders.  It keeps
// product orders, wallet credits and provider callbacks in separate ledgers.
package funding

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/wechatpayv3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidInput       = errors.New("请求参数不合法")
	ErrOrderNotFound      = errors.New("资金订单不存在或无权访问")
	ErrOrderNotPayable    = errors.New("当前资金订单不可支付")
	ErrPlanUnavailable    = errors.New("该权益方案暂不可购买")
	ErrLifetimeActive     = errors.New("永久 SVIP 已生效，无需重复购买")
	ErrTrialAlreadyUsed   = errors.New("体验 SVIP 仅可领取一次")
	ErrPaymentUnavailable = errors.New("微信 H5 支付尚未配置")
)

var idempotencyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$`)

type Handler struct {
	db              *gorm.DB
	platformConfigs *cloudconfig.Service
	wechatClient    *wechatpayv3.Client
}

func NewHandler(db *gorm.DB, platformConfigs *cloudconfig.Service) *Handler {
	return &Handler{db: db, platformConfigs: platformConfigs, wechatClient: &wechatpayv3.Client{}}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/funding/recharge-plans", h.RechargePlans)
	r.POST("/funding/recharge-orders", h.CreateRecharge)
	r.GET("/funding/recharge-orders/:id", h.GetRecharge)
	r.POST("/funding/recharge-orders/:id/pay", h.PayRecharge)
	r.GET("/funding/svip-plans", h.SVIPPlans)
	r.GET("/funding/svip-status", h.SVIPStatus)
	r.POST("/funding/svip-orders", h.CreateSVIP)
	r.GET("/funding/svip-orders/:id", h.GetSVIP)
	r.POST("/funding/svip-orders/:id/pay", h.PaySVIP)
}

type rechargePlan struct {
	ID          uint64  `gorm:"column:id" json:"id"`
	Name        string  `gorm:"column:name" json:"name"`
	Amount      float64 `gorm:"column:amount" json:"amount"`
	BonusAmount float64 `gorm:"column:bonus_amount" json:"bonus_amount"`
}

func (rechargePlan) TableName() string { return "qixi_crm_b_recharge_plan" }

type svipPlan struct {
	ID           uint64  `gorm:"column:id" json:"id"`
	Name         string  `gorm:"column:name" json:"name"`
	Price        float64 `gorm:"column:price" json:"price"`
	PlanType     string  `gorm:"column:plan_type" json:"plan_type"`
	DurationDays *uint   `gorm:"column:duration_days" json:"duration_days"`
	Benefits     string  `gorm:"column:benefits" json:"benefits"`
}

func (svipPlan) TableName() string { return "qixi_crm_b_svip_plan" }

type rechargeOrder struct {
	ID             uint64     `gorm:"column:id"`
	RechargeNo     string     `gorm:"column:recharge_no"`
	UserID         uint64     `gorm:"column:user_id"`
	Amount         float64    `gorm:"column:amount"`
	BonusAmount    float64    `gorm:"column:bonus_amount"`
	Status         string     `gorm:"column:status"`
	IdempotencyKey string     `gorm:"column:idempotency_key"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	PaidAt         *time.Time `gorm:"column:paid_at"`
}

func (rechargeOrder) TableName() string { return "qixi_crm_b_recharge_order" }

type svipOrder struct {
	ID             uint64     `gorm:"column:id"`
	OrderNo        string     `gorm:"column:order_no"`
	UserID         uint64     `gorm:"column:user_id"`
	PlanID         uint64     `gorm:"column:plan_id"`
	PlanName       string     `gorm:"column:plan_name"`
	PlanType       string     `gorm:"column:plan_type"`
	DurationDays   *uint      `gorm:"column:duration_days"`
	Amount         float64    `gorm:"column:amount"`
	Status         string     `gorm:"column:status"`
	IdempotencyKey string     `gorm:"column:idempotency_key"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	PaidAt         *time.Time `gorm:"column:paid_at"`
}

func (svipOrder) TableName() string { return "qixi_crm_b_svip_order" }

type paymentRow struct {
	ID              uint64     `gorm:"column:id"`
	OrderType       string     `gorm:"column:order_type"`
	FundingOrderID  uint64     `gorm:"column:funding_order_id"`
	UserID          uint64     `gorm:"column:user_id"`
	Channel         string     `gorm:"column:channel"`
	OutTradeNo      string     `gorm:"column:out_trade_no"`
	Amount          float64    `gorm:"column:amount"`
	Status          string     `gorm:"column:status"`
	ProviderPayload []byte     `gorm:"column:provider_payload"`
	PaidAt          *time.Time `gorm:"column:paid_at"`
}

func (paymentRow) TableName() string { return "qixi_crm_b_funding_payment" }

type paymentIntent struct {
	Status       string  `json:"status"`
	Channel      string  `json:"channel"`
	PaymentMode  string  `json:"payment_mode"`
	OutTradeNo   string  `json:"out_trade_no"`
	PayPrice     float64 `json:"pay_price"`
	MWEBURL      string  `json:"mweb_url,omitempty"`
	ExpiresAt    string  `json:"expires_at,omitempty"`
	FundingType  string  `json:"funding_type"`
	FundingOrder uint64  `json:"funding_order_id"`
}

type prepayPayload struct {
	PaymentMode string `json:"payment_mode"`
	MWEBURL     string `json:"mweb_url"`
	ExpiresAt   string `json:"expires_at"`
}

func (h *Handler) RechargePlans(c *gin.Context) {
	var plans []rechargePlan
	if err := h.db.WithContext(c.Request.Context()).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&plans).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询充值方案失败")
		return
	}
	response.OK(c, gin.H{"list": plans, "custom_amount_enabled": true, "min_amount": 1, "max_amount": 1000000})
}

func (h *Handler) SVIPPlans(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	var plans []svipPlan
	if err := h.db.WithContext(c.Request.Context()).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&plans).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询 SVIP 方案失败")
		return
	}
	usedTrial, err := h.hasPaidTrial(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询 SVIP 购买资格失败")
		return
	}
	items := make([]gin.H, 0, len(plans))
	for _, plan := range plans {
		items = append(items, gin.H{"id": plan.ID, "name": plan.Name, "price": plan.Price, "plan_type": plan.PlanType, "duration_days": plan.DurationDays, "benefits": plan.Benefits, "purchasable": !(plan.PlanType == "trial" && usedTrial)})
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) SVIPStatus(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	var row struct {
		Status    string     `gorm:"column:status"`
		ExpiresAt *time.Time `gorm:"column:expires_at"`
	}
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_svip").Where("user_id = ?", uid).Take(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusInternalServerError, "查询 SVIP 状态失败")
		return
	}
	active := err == nil && (row.Status == "lifetime" || row.ExpiresAt != nil && row.ExpiresAt.After(time.Now().UTC()))
	response.OK(c, gin.H{"active": active, "status": row.Status, "expires_at": row.ExpiresAt})
}

type rechargeRequest struct {
	PlanID         uint64  `json:"plan_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func (h *Handler) CreateRecharge(c *gin.Context) {
	var req rechargeRequest
	if c.ShouldBindJSON(&req) != nil || !validKey(req.IdempotencyKey) {
		fundingFail(c, ErrInvalidInput)
		return
	}
	uid := uint64(middleware.UID(c))
	order, err := h.createRecharge(c.Request.Context(), uid, req)
	if err != nil {
		fundingFail(c, err)
		return
	}
	response.OK(c, rechargeResponse(order))
}

func (h *Handler) createRecharge(ctx context.Context, uid uint64, req rechargeRequest) (rechargeOrder, error) {
	req.IdempotencyKey = strings.TrimSpace(req.IdempotencyKey)
	var existing rechargeOrder
	if err := h.db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", uid, req.IdempotencyKey).Take(&existing).Error; err == nil {
		return existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return rechargeOrder{}, err
	}
	amount, bonus := roundedAmount(req.Amount), 0.0
	if req.PlanID != 0 {
		var plan rechargePlan
		if err := h.db.WithContext(ctx).Where("id = ? AND status = ?", req.PlanID, 1).Take(&plan).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return rechargeOrder{}, ErrPlanUnavailable
			}
			return rechargeOrder{}, err
		}
		amount, bonus = plan.Amount, plan.BonusAmount
	}
	if amount < 1 || amount > 1000000 || roundedAmount(req.Amount) != req.Amount && req.PlanID == 0 {
		return rechargeOrder{}, ErrInvalidInput
	}
	created := rechargeOrder{RechargeNo: fundingNo("R"), UserID: uid, Amount: amount, BonusAmount: bonus, Status: "pending", IdempotencyKey: req.IdempotencyKey}
	if err := h.db.WithContext(ctx).Create(&created).Error; err != nil {
		if findErr := h.db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", uid, req.IdempotencyKey).Take(&existing).Error; findErr == nil {
			return existing, nil
		}
		return rechargeOrder{}, err
	}
	return created, nil
}

func (h *Handler) GetRecharge(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var order rechargeOrder
	if err := h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", id, middleware.UID(c)).Take(&order).Error; err != nil {
		fundingFail(c, ErrOrderNotFound)
		return
	}
	response.OK(c, rechargeResponse(order))
}

func (h *Handler) PayRecharge(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	intent, err := h.pay(c.Request.Context(), uint64(middleware.UID(c)), "recharge", id, c.ClientIP())
	if err != nil {
		fundingFail(c, err)
		return
	}
	response.OK(c, intent)
}

type svipRequest struct {
	PlanID         uint64 `json:"plan_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *Handler) CreateSVIP(c *gin.Context) {
	var req svipRequest
	if c.ShouldBindJSON(&req) != nil || req.PlanID == 0 || !validKey(req.IdempotencyKey) {
		fundingFail(c, ErrInvalidInput)
		return
	}
	order, err := h.createSVIP(c.Request.Context(), uint64(middleware.UID(c)), req)
	if err != nil {
		fundingFail(c, err)
		return
	}
	response.OK(c, svipResponse(order))
}

func (h *Handler) createSVIP(ctx context.Context, uid uint64, req svipRequest) (svipOrder, error) {
	var existing svipOrder
	if err := h.db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", uid, strings.TrimSpace(req.IdempotencyKey)).Take(&existing).Error; err == nil {
		return existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return svipOrder{}, err
	}
	var plan svipPlan
	if err := h.db.WithContext(ctx).Where("id = ? AND status = ?", req.PlanID, 1).Take(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return svipOrder{}, ErrPlanUnavailable
		}
		return svipOrder{}, err
	}
	if plan.Price < 0 || !validSVIPPlan(plan) {
		return svipOrder{}, ErrPlanUnavailable
	}
	if plan.PlanType == "trial" {
		used, err := h.hasPaidTrial(ctx, uid)
		if err != nil {
			return svipOrder{}, err
		}
		if used {
			return svipOrder{}, ErrTrialAlreadyUsed
		}
	}
	var state struct {
		Status string `gorm:"column:status"`
	}
	if err := h.db.WithContext(ctx).Table("qixi_crm_b_user_svip").Where("user_id = ?", uid).Take(&state).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return svipOrder{}, err
	}
	if state.Status == "lifetime" {
		return svipOrder{}, ErrLifetimeActive
	}
	created := svipOrder{OrderNo: fundingNo("S"), UserID: uid, PlanID: plan.ID, PlanName: plan.Name, PlanType: plan.PlanType, DurationDays: plan.DurationDays, Amount: plan.Price, Status: "pending", IdempotencyKey: strings.TrimSpace(req.IdempotencyKey)}
	if err := h.db.WithContext(ctx).Create(&created).Error; err != nil {
		if findErr := h.db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", uid, created.IdempotencyKey).Take(&existing).Error; findErr == nil {
			return existing, nil
		}
		return svipOrder{}, err
	}
	return created, nil
}

func (h *Handler) GetSVIP(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var order svipOrder
	if err := h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", id, middleware.UID(c)).Take(&order).Error; err != nil {
		fundingFail(c, ErrOrderNotFound)
		return
	}
	response.OK(c, svipResponse(order))
}

func (h *Handler) PaySVIP(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	intent, err := h.pay(c.Request.Context(), uint64(middleware.UID(c)), "svip", id, c.ClientIP())
	if err != nil {
		fundingFail(c, err)
		return
	}
	response.OK(c, intent)
}

func (h *Handler) pay(ctx context.Context, uid uint64, orderType string, orderID uint64, clientIP string) (paymentIntent, error) {
	if strings.TrimSpace(clientIP) == "" {
		return paymentIntent{}, ErrInvalidInput
	}
	config, err := h.wechatConfig(ctx)
	if err != nil {
		return paymentIntent{}, err
	}
	var orderNo string
	var amount float64
	var status string
	if orderType == "recharge" {
		var order rechargeOrder
		if err := h.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, uid).Take(&order).Error; err != nil {
			return paymentIntent{}, ErrOrderNotFound
		}
		orderNo, amount, status = order.RechargeNo, order.Amount, order.Status
	} else {
		var order svipOrder
		if err := h.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, uid).Take(&order).Error; err != nil {
			return paymentIntent{}, ErrOrderNotFound
		}
		orderNo, amount, status = order.OrderNo, order.Amount, order.Status
	}
	if status == "paid" {
		return paymentIntent{Status: "paid", Channel: "wechat", PaymentMode: "mweb", OutTradeNo: orderNo, PayPrice: amount, FundingType: orderType, FundingOrder: orderID}, nil
	}
	if status != "pending" || amount <= 0 {
		return paymentIntent{}, ErrOrderNotPayable
	}
	payment, create, err := h.beginPayment(ctx, uid, orderType, orderID, orderNo, amount)
	if err != nil {
		return paymentIntent{}, err
	}
	if !create {
		intent := paymentIntent{Status: "pending", Channel: "wechat", PaymentMode: "mweb", OutTradeNo: payment.OutTradeNo, PayPrice: payment.Amount, FundingType: orderType, FundingOrder: orderID}
		var payload prepayPayload
		if json.Unmarshal(payment.ProviderPayload, &payload) == nil {
			intent.MWEBURL, intent.ExpiresAt = payload.MWEBURL, payload.ExpiresAt
		}
		return intent, nil
	}
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	prepay, err := h.wechatClient.H5Prepay(ctx, config, wechatpayv3.H5Request{Description: fundingDescription(orderType), OutTradeNo: orderNo, Attach: fmt.Sprintf("funding:%s:%d", orderType, orderID), AmountCents: cents(amount), PayerClientIP: clientIP, ExpireAt: expiresAt})
	if err != nil {
		_ = h.db.WithContext(ctx).Model(&paymentRow{}).Where("id = ? AND status = ?", payment.ID, "processing").Update("status", "failed").Error
		return paymentIntent{}, err
	}
	payload, err := json.Marshal(prepayPayload{PaymentMode: "mweb", MWEBURL: prepay.MWEBURL, ExpiresAt: expiresAt.Format(time.RFC3339)})
	if err != nil {
		return paymentIntent{}, err
	}
	if err := h.db.WithContext(ctx).Model(&paymentRow{}).Where("id = ? AND status = ?", payment.ID, "processing").Updates(map[string]any{"status": "created", "provider_payload": payload}).Error; err != nil {
		return paymentIntent{}, err
	}
	return paymentIntent{Status: "pending", Channel: "wechat", PaymentMode: "mweb", OutTradeNo: orderNo, PayPrice: amount, MWEBURL: prepay.MWEBURL, ExpiresAt: expiresAt.Format(time.RFC3339), FundingType: orderType, FundingOrder: orderID}, nil
}

func (h *Handler) beginPayment(ctx context.Context, uid uint64, kind string, orderID uint64, orderNo string, amount float64) (paymentRow, bool, error) {
	var payment paymentRow
	create := false
	err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_type = ? AND funding_order_id = ? AND channel = ?", kind, orderID, "wechat").Take(&payment).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			payment = paymentRow{OrderType: kind, FundingOrderID: orderID, UserID: uid, Channel: "wechat", OutTradeNo: orderNo, Amount: amount, Status: "processing"}
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
			create = true
			return nil
		}
		if err != nil {
			return err
		}
		if payment.UserID != uid || cents(payment.Amount) != cents(amount) {
			return ErrOrderNotPayable
		}
		switch payment.Status {
		case "succeeded":
			return nil
		case "created":
			return nil
		case "failed":
			if err := tx.Model(&payment).Updates(map[string]any{"status": "processing", "provider_payload": nil}).Error; err != nil {
				return err
			}
			payment.Status = "processing"
			payment.ProviderPayload = nil
			create = true
			return nil
		default:
			return ErrOrderNotPayable
		}
	})
	return payment, create, err
}

func (h *Handler) wechatConfig(ctx context.Context) (wechatpayv3.Config, error) {
	if h.platformConfigs == nil {
		return wechatpayv3.Config{}, ErrPaymentUnavailable
	}
	values, err := h.platformConfigs.Values(ctx, "payment")
	if err != nil || !paymentconfig.ChannelReady(paymentconfig.Values(values), "wechat") {
		return wechatpayv3.Config{}, ErrPaymentUnavailable
	}
	config := wechatpayv3.Config{AppID: values["wechat_app_id"], H5AppID: values["wechat_h5_app_id"], H5SiteURL: values["wechat_h5_site_url"], MchID: values["wechat_mch_id"], MerchantSerialNo: values["wechat_serial_no"], MerchantPrivateKey: values["wechat_private_key"], MerchantCertPEM: values["wechat_merchant_cert"], APIv3Key: values["wechat_api_v3_key"], NotifyURL: values["wechat_notify_url"]}
	if !config.ValidForH5() {
		return wechatpayv3.Config{}, ErrPaymentUnavailable
	}
	return config, nil
}

func (h *Handler) hasPaidTrial(ctx context.Context, uid uint64) (bool, error) {
	var count int64
	err := h.db.WithContext(ctx).Model(&svipOrder{}).Where("user_id = ? AND plan_type = ? AND status = ?", uid, "trial", "paid").Count(&count).Error
	return count > 0, err
}
func validSVIPPlan(plan svipPlan) bool {
	if plan.PlanType == "lifetime" {
		return plan.Price > 0
	}
	return plan.DurationDays != nil && *plan.DurationDays > 0 && plan.Price > 0 && (plan.PlanType == "trial" || plan.PlanType == "period")
}
func validKey(value string) bool          { return idempotencyPattern.MatchString(strings.TrimSpace(value)) }
func roundedAmount(value float64) float64 { return math.Round(value*100) / 100 }
func cents(value float64) int64           { return int64(math.Round(value * 100)) }
func fundingNo(prefix string) string {
	raw := make([]byte, 8)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Sprintf("%s%d", prefix, time.Now().UTC().UnixNano())
	}
	return fmt.Sprintf("%s%s%s", prefix, time.Now().UTC().Format("060102150405"), strings.ToUpper(hex.EncodeToString(raw)))
}
func fundingDescription(kind string) string {
	if kind == "recharge" {
		return "CRM Live 账户充值"
	}
	return "CRM Live SVIP 开通"
}
func rechargeResponse(order rechargeOrder) gin.H {
	return gin.H{"id": order.ID, "recharge_no": order.RechargeNo, "amount": order.Amount, "bonus_amount": order.BonusAmount, "credited_amount": roundedAmount(order.Amount + order.BonusAmount), "status": order.Status, "created_at": order.CreatedAt, "paid_at": order.PaidAt}
}
func svipResponse(order svipOrder) gin.H {
	return gin.H{"id": order.ID, "order_no": order.OrderNo, "plan_id": order.PlanID, "plan_name": order.PlanName, "plan_type": order.PlanType, "duration_days": order.DurationDays, "amount": order.Amount, "status": order.Status, "created_at": order.CreatedAt, "paid_at": order.PaidAt}
}
func parseID(c *gin.Context) (uint64, bool) {
	var id uint64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil || id == 0 {
		fundingFail(c, ErrInvalidInput)
		return 0, false
	}
	return id, true
}
func fundingFail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrOrderNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrPaymentUnavailable):
		response.Fail(c, http.StatusConflict, err.Error())
	case errors.Is(err, ErrOrderNotPayable), errors.Is(err, ErrPlanUnavailable), errors.Is(err, ErrLifetimeActive), errors.Is(err, ErrTrialAlreadyUsed):
		response.Fail(c, http.StatusConflict, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "资金服务暂时不可用")
	}
}
