package combinationorder

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

const combinationActivityType = 4

var (
	errBadInput         = errors.New("拼团下单参数错误")
	errAddress          = errors.New("收货地址不存在或无权访问")
	errUnavailable      = errors.New("拼团活动不可用")
	errBuyingClosed     = errors.New("拼团已结束")
	errBuyingFull       = errors.New("拼团人数已满")
	errAlreadyJoined    = errors.New("不可重复参团")
	errIdempotencyClash = errors.New("幂等键已被其他订单使用")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/order/group/check", h.check)
	r.POST("/order/group/create", h.create)
}

type input struct {
	ProductGroupID uint64 `json:"product_group_id"`
	GroupBuyingID  uint64 `json:"group_buying_id"`
	CartNum        int    `json:"cart_num"`
	AddressID      uint64 `json:"address_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

type quote struct {
	ProductGroupID uint64  `gorm:"column:product_group_id"`
	ProductID      uint64  `gorm:"column:product_id"`
	MerchantID     uint64  `gorm:"column:merchant_id"`
	StoreID        uint64  `gorm:"column:store_id"`
	StoreName      string  `gorm:"column:store_name"`
	MerchantName   string  `gorm:"column:merchant_name"`
	Title          string  `gorm:"column:title"`
	CoverURL       string  `gorm:"column:cover_url"`
	Price          float64 `gorm:"column:price"`
	BuyingCountNum int     `gorm:"column:buying_count_num"`
	BuyingNum      int     `gorm:"column:buying_num"`
	DurationHours  int     `gorm:"column:duration_hours"`
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

type buyingCreate struct {
	ID             uint64 `gorm:"column:group_buying_id"`
	ProductGroupID uint64 `gorm:"column:product_group_id"`
	Status         int    `gorm:"column:status"`
	BuyingCountNum int    `gorm:"column:buying_count_num"`
	BuyingNum      int    `gorm:"column:buying_num"`
	YetBuyingNum   int    `gorm:"column:yet_buying_num"`
	MerchantID     uint64 `gorm:"column:mer_id"`
	EndTime        int64  `gorm:"column:end_time"`
	IsDel          int    `gorm:"column:is_del"`
}

type buying struct {
	ID             uint64 `gorm:"column:group_buying_id"`
	ProductGroupID uint64 `gorm:"column:product_group_id"`
	Status         int    `gorm:"column:status"`
	BuyingCountNum int    `gorm:"column:buying_count_num"`
	IsDel          int    `gorm:"column:is_del"`
	EndTime        int64  `gorm:"column:end_time"`
}

func (h *Handler) check(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil || normalize(&in) != nil {
		writeErr(c, errBadInput)
		return
	}
	q, err := h.loadQuote(c.Request.Context(), h.db, in.ProductGroupID, false)
	if err != nil {
		writeErr(c, err)
		return
	}
	if in.CartNum > q.BuyingNum {
		writeErr(c, errBadInput)
		return
	}
	if in.GroupBuyingID != 0 {
		if _, err := h.loadBuying(c.Request.Context(), h.db, in.GroupBuyingID, q.ProductGroupID, false); err != nil {
			writeErr(c, err)
			return
		}
	}
	response.OK(c, quoteResponse(q, in.GroupBuyingID, in.CartNum))
}

func (h *Handler) create(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil || normalize(&in) != nil || in.AddressID == 0 {
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
			if created.ActivityType != combinationActivityType {
				return errIdempotencyClash
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		q, err := h.loadQuote(c.Request.Context(), tx, in.ProductGroupID, true)
		if err != nil {
			return err
		}
		if in.CartNum > q.BuyingNum {
			return errBadInput
		}
		buyingID, leader, err := h.reserveBuying(c.Request.Context(), tx, q, uid, in.GroupBuyingID)
		if err != nil {
			return err
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
		total := q.Price * float64(in.CartNum)
		if err := tx.Table("qixi_crm_b_group_order").Create(map[string]any{
			"order_no": newOrderNo("CG"), "user_id": uid, "total_amount": total, "discount_amount": 0, "freight_amount": 0, "pay_amount": total, "total_quantity": in.CartNum,
			"recipient_snapshot": string(snapshot), "pay_status": "pending", "idempotency_key": in.IdempotencyKey, "activity_type": combinationActivityType, "points_amount": 0,
		}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Create(map[string]any{
			"group_order_id": created.ID, "order_no": newOrderNo("CO"), "merchant_id": q.MerchantID, "merchant_name_snapshot": q.MerchantName, "store_id": q.StoreID, "store_name_snapshot": q.StoreName,
			"user_id": uid, "total_amount": total, "discount_amount": 0, "freight_amount": 0, "pay_amount": total, "total_quantity": in.CartNum, "recipient_snapshot": string(snapshot), "status": "pending_pay", "activity_type": combinationActivityType, "points_amount": 0,
		}).Error; err != nil {
			return err
		}
		var child struct {
			ID uint64 `gorm:"column:id"`
		}
		if err := tx.Table("qixi_crm_b_order").Select("id").Where("group_order_id = ? AND user_id = ?", created.ID, uid).First(&child).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order_item").Create(map[string]any{"order_id": child.ID, "product_id": q.ProductID, "sku_key": "", "title_snapshot": q.Title, "cover_url_snapshot": q.CoverURL, "spec_snapshot": "{}", "unit_price": q.Price, "quantity": in.CartNum, "refund_quantity": 0}).Error; err != nil {
			return err
		}
		nickname := strings.TrimSpace(addr.Recipient)
		if nickname == "" {
			nickname = fmt.Sprintf("用户%d", uid)
		}
		if err := tx.Table("qixi_crm_b_combination_member").Create(map[string]any{"group_buying_id": buyingID, "product_group_id": q.ProductGroupID, "status": 0, "is_initiator": boolInt(leader), "is_leader": boolInt(leader), "order_id": child.ID, "uid": uid, "nickname": nickname, "is_del": 0}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_order_activity").Create(map[string]any{"group_order_id": created.ID, "activity_type": combinationActivityType, "activity_id": buyingID, "related_activity_id": q.ProductGroupID, "quantity": in.CartNum, "status": "reserved"}).Error
	})
	if err != nil {
		var existing group
		lookup := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&existing).Error
		if lookup != nil || existing.ActivityType != combinationActivityType {
			writeErr(c, err)
			return
		}
		created = existing
	}
	response.OK(c, gin.H{"group_order_id": created.ID, "group_order_sn": created.OrderNo, "pay_price": created.PayAmount, "total_num": created.TotalQuantity, "pay_status": created.PayStatus, "paid": created.PayStatus == "paid"})
}

func normalize(in *input) error {
	if in.ProductGroupID == 0 {
		return errBadInput
	}
	if in.CartNum == 0 {
		in.CartNum = 1
	}
	if in.CartNum < 1 {
		return errBadInput
	}
	return nil
}

func (h *Handler) loadQuote(ctx context.Context, db *gorm.DB, id uint64, lock bool) (quote, error) {
	q := quote{}
	query := db.WithContext(ctx).Table("qixi_crm_b_combination_group AS g").Select("g.product_group_id,g.product_id,p.merchant_id,p.store_id,p.store_name,p.merchant_name,p.title,p.cover_url,g.price,g.buying_count_num,g.buying_num,g.time AS duration_hours").Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = g.product_id").Where("g.product_group_id = ? AND g.status = 1 AND g.is_show = 1 AND g.is_del = 0 AND g.product_status = 1 AND g.action_status = 1 AND g.start_time <= ? AND g.end_time >= ? AND p.sale_status = 1", id, time.Now(), time.Now())
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.First(&q).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return q, errUnavailable
		}
		return q, err
	}
	if q.BuyingCountNum < 2 || q.BuyingNum < 1 || q.Price < 0 {
		return q, errUnavailable
	}
	return q, nil
}

func (h *Handler) loadBuying(ctx context.Context, db *gorm.DB, id, productGroupID uint64, lock bool) (buying, error) {
	b := buying{}
	query := db.WithContext(ctx).Table("qixi_crm_b_combination_buying").Where("group_buying_id = ? AND product_group_id = ?", id, productGroupID)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b, errBuyingClosed
		}
		return b, err
	}
	if b.Status != 0 || b.IsDel != 0 || (b.EndTime > 0 && b.EndTime <= time.Now().Unix()) {
		return b, errBuyingClosed
	}
	return b, nil
}

func (h *Handler) reserveBuying(ctx context.Context, tx *gorm.DB, q quote, uid, existingID uint64) (uint64, bool, error) {
	if existingID == 0 {
		created := buyingCreate{ProductGroupID: q.ProductGroupID, Status: 0, BuyingCountNum: q.BuyingCountNum, BuyingNum: q.BuyingNum, YetBuyingNum: 0, MerchantID: q.MerchantID, EndTime: time.Now().Add(time.Duration(q.DurationHours) * time.Hour).Unix(), IsDel: 0}
		if err := tx.Table("qixi_crm_b_combination_buying").Create(&created).Error; err != nil {
			return 0, false, err
		}
		return created.ID, true, nil
	}
	b, err := h.loadBuying(ctx, tx, existingID, q.ProductGroupID, true)
	if err != nil {
		return 0, false, err
	}
	var members int64
	if err := tx.Table("qixi_crm_b_combination_member").Where("group_buying_id = ? AND is_del = 0", b.ID).Count(&members).Error; err != nil {
		return 0, false, err
	}
	if members >= int64(b.BuyingCountNum) {
		return 0, false, errBuyingFull
	}
	var joined int64
	if err := tx.Table("qixi_crm_b_combination_member").Where("group_buying_id = ? AND uid = ? AND is_del = 0", b.ID, uid).Count(&joined).Error; err != nil {
		return 0, false, err
	}
	if joined > 0 {
		return 0, false, errAlreadyJoined
	}
	return b.ID, false, nil
}

func quoteResponse(q quote, buyingID uint64, quantity int) gin.H {
	return gin.H{"product_group_id": q.ProductGroupID, "group_buying_id": buyingID, "product_id": q.ProductID, "store_name": q.StoreName, "image": q.CoverURL, "mer_id": q.MerchantID, "mer_name": q.MerchantName, "cart_num": quantity, "price": q.Price, "pay_price": q.Price * float64(quantity), "activity_type": combinationActivityType}
}
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
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
	case errors.Is(err, errBadInput), errors.Is(err, errAddress), errors.Is(err, errUnavailable), errors.Is(err, errBuyingClosed), errors.Is(err, errBuyingFull), errors.Is(err, errAlreadyJoined), errors.Is(err, errIdempotencyClash):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "拼团下单失败")
	}
}
