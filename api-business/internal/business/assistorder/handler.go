package assistorder

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const assistActivityType = 3

var (
	errBadInput         = errors.New("助力下单参数错误")
	errAddress          = errors.New("收货地址不存在或无权访问")
	errUnavailable      = errors.New("助力单暂不可下单")
	errSoldOut          = errors.New("助力活动库存不足")
	errIdempotencyClash = errors.New("幂等键已被其他订单使用")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/order/assist/check", h.check)
	r.POST("/order/assist/create", h.create)
}

type input struct {
	ProductAssistSetID uint64 `json:"product_assist_set_id"`
	CartNum            int    `json:"cart_num"`
	AddressID          uint64 `json:"address_id"`
	IdempotencyKey     string `json:"idempotency_key"`
}

type quote struct {
	SetID        uint64  `gorm:"column:set_id"`
	AssistID     uint64  `gorm:"column:assist_id"`
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	StoreID      uint64  `gorm:"column:store_id"`
	StoreName    string  `gorm:"column:store_name"`
	MerchantName string  `gorm:"column:merchant_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:assist_price"`
	Stock        int     `gorm:"column:stock"`
}

type address struct {
	ID        uint64 `gorm:"column:id"`
	Recipient string `gorm:"column:recipient"`
	Mobile    string `gorm:"column:mobile"`
	Province  string `gorm:"column:province"`
	City      string `gorm:"column:city"`
	District  string `gorm:"column:district"`
	Detail    string `gorm:"column:detail"`
	PostCode  int    `gorm:"column:post_code"`
}

type group struct {
	ID            uint64  `gorm:"column:id"`
	OrderNo       string  `gorm:"column:order_no"`
	PayStatus     string  `gorm:"column:pay_status"`
	ActivityType  int     `gorm:"column:activity_type"`
	TotalQuantity int     `gorm:"column:total_quantity"`
	PayAmount     float64 `gorm:"column:pay_amount"`
}

func (h *Handler) check(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, errBadInput)
		return
	}
	if normalizeQuantity(&in) != nil {
		writeErr(c, errBadInput)
		return
	}
	q, err := h.loadQuote(c.Request.Context(), h.db, uint64(middleware.UID(c)), in.ProductAssistSetID, false)
	if err != nil {
		writeErr(c, err)
		return
	}
	if q.Stock < 1 {
		writeErr(c, errSoldOut)
		return
	}
	response.OK(c, quoteResponse(q))
}

func (h *Handler) create(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, errBadInput)
		return
	}
	if normalizeQuantity(&in) != nil || in.AddressID == 0 {
		writeErr(c, errBadInput)
		return
	}
	in.IdempotencyKey = strings.TrimSpace(in.IdempotencyKey)
	if len(in.IdempotencyKey) < 12 || len(in.IdempotencyKey) > 128 {
		writeErr(c, errBadInput)
		return
	}
	uid := uint64(middleware.UID(c))
	var created group
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err == nil {
			if created.ActivityType != assistActivityType {
				return errIdempotencyClash
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		q, err := h.loadQuote(c.Request.Context(), tx, uid, in.ProductAssistSetID, true)
		if err != nil {
			return err
		}
		reserved := tx.Table("qixi_crm_b_assist_set").Where("product_assist_set_id = ? AND status = ? AND is_del = 0", q.SetID, 10).Update("status", 11)
		if reserved.Error != nil {
			return reserved.Error
		}
		if reserved.RowsAffected != 1 {
			return errUnavailable
		}
		updated := tx.Table("qixi_crm_b_assist").
			Where("product_assist_id = ? AND stock >= 1 AND is_del = 0", q.AssistID).
			UpdateColumn("stock", gorm.Expr("stock - 1"))
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return errSoldOut
		}
		var addr address
		if err := tx.Table("qixi_crm_b_user_address").Where("id = ? AND user_id = ?", in.AddressID, uid).First(&addr).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errAddress
			}
			return err
		}
		snapshot, err := json.Marshal(addr)
		if err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_group_order").Create(map[string]any{
			"order_no": newOrderNo("AG"), "user_id": uid, "total_amount": q.Price, "discount_amount": 0,
			"freight_amount": 0, "pay_amount": q.Price, "total_quantity": 1, "recipient_snapshot": string(snapshot),
			"pay_status": "pending", "idempotency_key": in.IdempotencyKey, "activity_type": assistActivityType, "points_amount": 0,
		}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Create(map[string]any{
			"group_order_id": created.ID, "order_no": newOrderNo("AO"), "merchant_id": q.MerchantID,
			"merchant_name_snapshot": q.MerchantName, "store_id": q.StoreID, "store_name_snapshot": q.StoreName,
			"user_id": uid, "total_amount": q.Price, "discount_amount": 0, "freight_amount": 0, "pay_amount": q.Price,
			"total_quantity": 1, "recipient_snapshot": string(snapshot), "status": "pending_pay", "activity_type": assistActivityType, "points_amount": 0,
		}).Error; err != nil {
			return err
		}
		var child struct {
			ID uint64 `gorm:"column:id"`
		}
		if err := tx.Table("qixi_crm_b_order").Select("id").Where("group_order_id = ? AND user_id = ?", created.ID, uid).First(&child).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order_item").Create(map[string]any{
			"order_id": child.ID, "product_id": q.ProductID, "sku_key": "", "title_snapshot": q.Title,
			"cover_url_snapshot": q.CoverURL, "spec_snapshot": "{}", "unit_price": q.Price, "quantity": 1, "refund_quantity": 0,
		}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_order_activity").Create(map[string]any{
			"group_order_id": created.ID, "activity_type": assistActivityType, "activity_id": q.SetID,
			"related_activity_id": q.AssistID, "quantity": 1, "status": "reserved",
		}).Error
	})
	if err != nil {
		var existing group
		lookupErr := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").
			Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&existing).Error
		if lookupErr != nil || existing.ActivityType != assistActivityType {
			writeErr(c, err)
			return
		}
		created = existing
	}
	response.OK(c, gin.H{"group_order_id": created.ID, "group_order_sn": created.OrderNo, "pay_price": created.PayAmount, "total_num": created.TotalQuantity, "pay_status": created.PayStatus, "paid": created.PayStatus == "paid"})
}

func normalizeQuantity(in *input) error {
	if in.ProductAssistSetID == 0 {
		return errBadInput
	}
	if in.CartNum == 0 {
		in.CartNum = 1
	}
	if in.CartNum != 1 {
		return errBadInput
	}
	return nil
}

func (h *Handler) loadQuote(ctx context.Context, db *gorm.DB, uid, setID uint64, lock bool) (quote, error) {
	q := quote{}
	query := db.WithContext(ctx).Table("qixi_crm_b_assist_set AS s").
		Select("s.product_assist_set_id AS set_id,a.product_assist_id AS assist_id,a.product_id,p.merchant_id,p.store_id,p.store_name,p.merchant_name,p.title,p.cover_url,a.assist_price,a.stock").
		Joins("JOIN qixi_crm_b_assist AS a ON a.product_assist_id = s.product_assist_id").
		Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = a.product_id").
		Where("s.product_assist_set_id = ? AND s.uid = ? AND s.status = ? AND s.is_del = 0", setID, uid, 10).
		Where("a.status = 1 AND a.is_show = 1 AND a.product_status = 1 AND a.action_status = 1 AND a.is_del = 0 AND a.start_time <= ? AND a.end_time >= ?", time.Now(), time.Now()).
		Where("p.sale_status = 1")
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.First(&q).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return q, errUnavailable
		}
		return q, err
	}
	return q, nil
}

func quoteResponse(q quote) gin.H {
	return gin.H{
		"product_assist_set_id": q.SetID, "product_assist_id": q.AssistID, "product_id": q.ProductID,
		"store_name": q.StoreName, "image": q.CoverURL, "mer_id": q.MerchantID, "mer_name": q.MerchantName,
		"cart_num": 1, "price": q.Price, "pay_price": q.Price, "stock": q.Stock, "activity_type": assistActivityType,
	}
}

func newOrderNo(prefix string) string {
	var nonce [6]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s%d%x", prefix, time.Now().UnixMilli(), nonce)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errBadInput), errors.Is(err, errAddress), errors.Is(err, errUnavailable), errors.Is(err, errSoldOut), errors.Is(err, errIdempotencyClash):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "助力下单失败")
	}
}
