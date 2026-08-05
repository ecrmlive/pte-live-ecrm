package reservation

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const reservationActivityType = 30

var (
	errBadInput         = errors.New("预约参数错误")
	errAddress          = errors.New("收货地址不存在或无权访问")
	errUnavailable      = errors.New("预约商品或时段不可用")
	errSlotFull         = errors.New("预约时段已满")
	errActivityMismatch = errors.New("幂等键已被其他订单使用")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/order/reservation/check", h.check)
	r.POST("/order/reservation/create", h.create)
}

type input struct {
	ProductID      uint64 `json:"product_id"`
	SlotID         uint64 `json:"slot_id"`
	Date           string `json:"date"`
	AddressID      uint64 `json:"address_id"`
	Mark           string `json:"mark"`
	IdempotencyKey string `json:"idempotency_key"`
}

type quote struct {
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	StoreID      uint64  `gorm:"column:store_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	ShowDays     int     `gorm:"column:show_reservation_days"`
	SlotID       uint64  `gorm:"column:slot_id"`
	StartTime    string  `gorm:"column:start_time"`
	EndTime      string  `gorm:"column:end_time"`
	SlotStock    int     `gorm:"column:slot_stock"`
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
	q, err := h.quote(c.Request.Context(), h.db, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, quoteResponse(q, in.Date))
}

func (h *Handler) create(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, errBadInput)
		return
	}
	uid := uint64(middleware.UID(c))
	in.IdempotencyKey = strings.TrimSpace(in.IdempotencyKey)
	if in.ProductID == 0 || in.SlotID == 0 || in.AddressID == 0 || len(in.IdempotencyKey) < 12 || len(in.IdempotencyKey) > 128 {
		writeErr(c, errBadInput)
		return
	}
	if _, err := h.quote(c.Request.Context(), h.db, in); err != nil {
		writeErr(c, err)
		return
	}
	var created group
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err == nil {
			if created.ActivityType != reservationActivityType {
				return errActivityMismatch
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		q, err := h.lockQuote(c.Request.Context(), tx, in)
		if err != nil {
			return err
		}
		var booked int64
		if err := tx.Table("qixi_crm_b_reservation_booking").Where("product_id = ? AND slot_id = ? AND booking_date = ? AND status = 1", q.ProductID, q.SlotID, in.Date).Count(&booked).Error; err != nil {
			return err
		}
		if booked >= int64(q.SlotStock) {
			return errSlotFull
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
			"order_no": orderNo("RG"), "user_id": uid, "total_amount": q.Price, "discount_amount": 0,
			"freight_amount": 0, "pay_amount": q.Price, "total_quantity": 1, "recipient_snapshot": string(snapshot),
			"pay_status": "pending", "idempotency_key": in.IdempotencyKey, "activity_type": reservationActivityType,
			"points_amount": 0, "remark": strings.TrimSpace(in.Mark),
		}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Create(map[string]any{
			"group_order_id": created.ID, "order_no": orderNo("RO"), "merchant_id": q.MerchantID,
			"merchant_name_snapshot": q.MerchantName, "store_id": q.StoreID, "store_name_snapshot": q.StoreName,
			"user_id": uid, "total_amount": q.Price, "discount_amount": 0, "freight_amount": 0, "pay_amount": q.Price,
			"total_quantity": 1, "recipient_snapshot": string(snapshot), "remark": strings.TrimSpace(in.Mark), "status": "pending_pay",
			"activity_type": reservationActivityType, "points_amount": 0,
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
		if err := tx.Table("qixi_crm_b_reservation_booking").Create(map[string]any{
			"product_id": q.ProductID, "slot_id": q.SlotID, "booking_date": in.Date, "order_id": child.ID,
			"user_id": uid, "status": 1, "verify_code": verifyCode(),
		}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_reservation_slot").Where("attr_reservation_id = ?", q.SlotID).UpdateColumn("use_num", gorm.Expr("use_num + 1")).Error
	})
	if err != nil {
		var existing group
		lookupErr := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&existing).Error
		if lookupErr != nil || existing.ActivityType != reservationActivityType {
			writeErr(c, err)
			return
		}
		created = existing
	}
	response.OK(c, gin.H{"group_order_id": created.ID, "group_order_sn": created.OrderNo, "pay_price": created.PayAmount, "total_num": created.TotalQuantity, "pay_status": created.PayStatus, "paid": created.PayStatus == "paid"})
}

func (h *Handler) quote(ctx context.Context, db *gorm.DB, in input) (quote, error) {
	if in.ProductID == 0 || in.SlotID == 0 || !validDate(in.Date) {
		return quote{}, errBadInput
	}
	var q quote
	err := db.WithContext(ctx).Table("qixi_crm_b_reservation_activity AS a").
		Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = a.product_id").
		Joins("JOIN qixi_crm_b_reservation_slot AS s ON s.product_id = a.product_id AND s.attr_reservation_id = ?", in.SlotID).
		Select("p.product_id,p.merchant_id,p.store_id,p.merchant_name,p.store_name,p.title,p.cover_url,p.price,a.show_reservation_days,s.attr_reservation_id AS slot_id,s.start_time,s.end_time,s.stock AS slot_stock").
		Where("a.product_id = ? AND a.status = 1 AND p.sale_status = 1", in.ProductID).First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return quote{}, errUnavailable
	}
	if err != nil {
		return quote{}, err
	}
	if q.SlotStock <= 0 || !dateInRange(in.Date, q.ShowDays) {
		return quote{}, errUnavailable
	}
	var booked int64
	if err := db.WithContext(ctx).Table("qixi_crm_b_reservation_booking").Where("product_id = ? AND slot_id = ? AND booking_date = ? AND status = 1", q.ProductID, q.SlotID, in.Date).Count(&booked).Error; err != nil {
		return quote{}, err
	}
	if booked >= int64(q.SlotStock) {
		return quote{}, errSlotFull
	}
	return q, nil
}

func (h *Handler) lockQuote(ctx context.Context, tx *gorm.DB, in input) (quote, error) {
	var slot struct {
		ProductID uint64 `gorm:"column:product_id"`
		SlotID    uint64 `gorm:"column:slot_id"`
	}
	if err := tx.WithContext(ctx).Table("qixi_crm_b_reservation_slot").Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("product_id,attr_reservation_id AS slot_id").Where("attr_reservation_id = ? AND product_id = ?", in.SlotID, in.ProductID).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return quote{}, errUnavailable
		}
		return quote{}, err
	}
	return h.quote(ctx, tx, in)
}

func quoteResponse(q quote, date string) gin.H {
	return gin.H{"product_id": q.ProductID, "slot_id": q.SlotID, "date": date, "time_part": q.StartTime + "-" + q.EndTime,
		"store_name": q.StoreName, "image": q.CoverURL, "mer_id": q.MerchantID, "mer_name": q.MerchantName,
		"pay_price": q.Price, "verify_hint": "支付后凭核销码到店核销", "activity_type": reservationActivityType}
}

func validDate(raw string) bool {
	_, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(raw), time.Local)
	return err == nil
}
func dateInRange(raw string, days int) bool {
	if days <= 0 {
		days = 7
	}
	d, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(raw), time.Local)
	if err != nil {
		return false
	}
	today := time.Now().In(time.Local).Truncate(24 * time.Hour)
	return !d.Before(today) && !d.After(today.AddDate(0, 0, days-1))
}
func orderNo(prefix string) string {
	var entropy [4]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s%d%x", prefix, time.Now().UnixMilli(), entropy)
}
func verifyCode() string {
	var entropy [4]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return fmt.Sprintf("RV%d", time.Now().UnixNano())
	}
	return "RV" + strconv.FormatUint(uint64(entropy[0])<<24|uint64(entropy[1])<<16|uint64(entropy[2])<<8|uint64(entropy[3]), 10)
}
func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errUnavailable):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, errBadInput), errors.Is(err, errAddress), errors.Is(err, errSlotFull), errors.Is(err, errActivityMismatch):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "预约服务异常")
	}
}
