package presellorder

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

const presellActivityType = 2

var (
	errBadInput      = errors.New("预售下单参数错误")
	errAddress       = errors.New("收货地址不存在或无权访问")
	errUnavailable   = errors.New("预售活动不可用")
	errSoldOut       = errors.New("预售库存不足")
	errFinalNotFound = errors.New("尾款单不存在或无权访问")
	errFinalNotOpen  = errors.New("尚未到尾款支付时间")
	errFinalClosed   = errors.New("尾款单不可支付")
	errBalance       = errors.New("余额不足")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/order/presell/check", h.check)
	r.POST("/order/presell/create", h.create)
	r.GET("/presell/finals", h.listFinals)
	r.GET("/presell/finals/:id", h.getFinal)
	r.POST("/presell/pay/:id", h.payFinal)
	r.POST("/order/presell/final/pay/:id", h.payFinal)
}

type input struct {
	ProductPresellID uint64 `json:"product_presell_id"`
	CartNum          int    `json:"cart_num"`
	AddressID        uint64 `json:"address_id"`
	IdempotencyKey   string `json:"idempotency_key"`
}
type quote struct {
	ID           uint64  `gorm:"column:product_presell_id"`
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	StoreID      uint64  `gorm:"column:store_id"`
	StoreName    string  `gorm:"column:store_name"`
	MerchantName string  `gorm:"column:merchant_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	DownPrice    float64 `gorm:"column:down_price"`
	FinalPrice   float64 `gorm:"column:final_price"`
	Stock        int     `gorm:"column:stock"`
	PresellType  int     `gorm:"column:presell_type"`
	FinalStart   string  `gorm:"column:final_start_time"`
	FinalEnd     string  `gorm:"column:final_end_time"`
}
type addr struct {
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
type finalRow struct {
	ID         uint64    `gorm:"column:presell_order_id" json:"presell_order_id"`
	OrderNo    string    `gorm:"column:presell_order_sn" json:"presell_order_sn"`
	UserID     uint64    `gorm:"column:uid"`
	OrderID    uint64    `gorm:"column:order_id"`
	ActivityID uint64    `gorm:"column:product_presell_id"`
	FinalStart time.Time `gorm:"column:final_start_time" json:"final_start_time"`
	FinalEnd   time.Time `gorm:"column:final_end_time" json:"final_end_time"`
	Paid       int       `gorm:"column:paid" json:"paid"`
	Status     int       `gorm:"column:status"`
	PayPrice   float64   `gorm:"column:pay_price" json:"pay_price"`
	StoreName  string    `gorm:"column:store_name" json:"store_name"`
}

func (h *Handler) check(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || normalize(&in) != nil {
		writeErr(c, errBadInput)
		return
	}
	q, e := h.loadQuote(c.Request.Context(), h.db, in.ProductPresellID, false)
	if e != nil {
		writeErr(c, e)
		return
	}
	if q.Stock < in.CartNum {
		writeErr(c, errSoldOut)
		return
	}
	response.OK(c, quoteOut(q, in.CartNum))
}
func (h *Handler) create(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || normalize(&in) != nil || in.AddressID == 0 {
		writeErr(c, errBadInput)
		return
	}
	in.IdempotencyKey = strings.TrimSpace(in.IdempotencyKey)
	if len(in.IdempotencyKey) < 12 || len(in.IdempotencyKey) > 128 {
		writeErr(c, errBadInput)
		return
	}
	uid := uint64(middleware.UID(c))
	var made group
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if e := tx.Table("qixi_crm_b_group_order").Where("user_id=? AND idempotency_key=?", uid, in.IdempotencyKey).First(&made).Error; e == nil {
			if made.ActivityType != presellActivityType {
				return errBadInput
			}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		q, e := h.loadQuote(c.Request.Context(), tx, in.ProductPresellID, true)
		if e != nil {
			return e
		}
		updated := tx.Table("qixi_crm_b_presell").Where("product_presell_id=? AND stock>=? AND is_del=0", q.ID, in.CartNum).UpdateColumn("stock", gorm.Expr("stock - ?", in.CartNum))
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return errSoldOut
		}
		var a addr
		if e := tx.Table("qixi_crm_b_user_address").Where("id=? AND user_id=?", in.AddressID, uid).First(&a).Error; e != nil {
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return errAddress
			}
			return e
		}
		snap, e := json.Marshal(a)
		if e != nil {
			return e
		}
		total := q.Price * float64(in.CartNum)
		first := total
		if q.PresellType == 2 {
			first = q.DownPrice * float64(in.CartNum)
		}
		if e = tx.Table("qixi_crm_b_group_order").Create(map[string]any{"order_no": sn("PG"), "user_id": uid, "total_amount": total, "discount_amount": 0, "freight_amount": 0, "pay_amount": first, "total_quantity": in.CartNum, "recipient_snapshot": string(snap), "pay_status": "pending", "idempotency_key": in.IdempotencyKey, "activity_type": presellActivityType, "points_amount": 0}).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_group_order").Where("user_id=? AND idempotency_key=?", uid, in.IdempotencyKey).First(&made).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_order").Create(map[string]any{"group_order_id": made.ID, "order_no": sn("PO"), "merchant_id": q.MerchantID, "merchant_name_snapshot": q.MerchantName, "store_id": q.StoreID, "store_name_snapshot": q.StoreName, "user_id": uid, "total_amount": total, "discount_amount": 0, "freight_amount": 0, "pay_amount": first, "total_quantity": in.CartNum, "recipient_snapshot": string(snap), "status": "pending_pay", "activity_type": presellActivityType, "points_amount": 0}).Error; e != nil {
			return e
		}
		var child struct {
			ID uint64 `gorm:"column:id"`
		}
		if e = tx.Table("qixi_crm_b_order").Select("id").Where("group_order_id=?", made.ID).First(&child).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_order_item").Create(map[string]any{"order_id": child.ID, "product_id": q.ProductID, "sku_key": "", "title_snapshot": q.Title, "cover_url_snapshot": q.CoverURL, "spec_snapshot": "{}", "unit_price": q.Price, "quantity": in.CartNum, "refund_quantity": 0}).Error; e != nil {
			return e
		}
		return tx.Table("qixi_crm_b_order_activity").Create(map[string]any{"group_order_id": made.ID, "activity_type": presellActivityType, "activity_id": q.ID, "related_activity_id": q.PresellType, "quantity": in.CartNum, "status": "reserved"}).Error
	})
	if err != nil {
		var old group
		if e := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_group_order").Where("user_id=? AND idempotency_key=?", uid, in.IdempotencyKey).First(&old).Error; e != nil || old.ActivityType != presellActivityType {
			writeErr(c, err)
			return
		}
		made = old
	}
	response.OK(c, gin.H{"group_order_id": made.ID, "group_order_sn": made.OrderNo, "pay_price": made.PayAmount, "total_num": made.TotalQuantity, "pay_status": made.PayStatus, "paid": made.PayStatus == "paid"})
}
func (h *Handler) listFinals(c *gin.Context) {
	page, limit := finalPageParams(c)
	uid := uint64(middleware.UID(c))
	now := time.Now()

	// 过期单必须先走现有状态迁移和库存归还，不能仅从列表隐藏。
	var expired []finalRow
	if e := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_presell_order AS p").Select("p.presell_order_id,p.presell_order_sn,p.uid,p.order_id,p.product_presell_id,p.final_start_time,p.final_end_time,p.paid,p.status,p.pay_price,o.store_name_snapshot AS store_name").Joins("JOIN qixi_crm_b_order AS o ON o.id=p.order_id").Where("p.uid=? AND p.paid=0 AND p.status=1 AND p.final_end_time<?", uid, now).Find(&expired).Error; e != nil {
		writeErr(c, e)
		return
	}
	for _, row := range expired {
		if e := h.expire(c.Request.Context(), uid, row.ID); e != nil {
			writeErr(c, e)
			return
		}
	}

	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_presell_order AS p").Joins("JOIN qixi_crm_b_order AS o ON o.id=p.order_id").Where("p.uid=? AND p.paid=0 AND p.status=1 AND p.final_start_time<=? AND p.final_end_time>=?", uid, now, now)
	var total int64
	if e := q.Count(&total).Error; e != nil {
		writeErr(c, e)
		return
	}
	rows := []finalRow{}
	if e := q.Select("p.presell_order_id,p.presell_order_sn,p.uid,p.order_id,p.product_presell_id,p.final_start_time,p.final_end_time,p.paid,p.status,p.pay_price,o.store_name_snapshot AS store_name").Order("p.presell_order_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; e != nil {
		writeErr(c, e)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func finalPageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

func (h *Handler) getFinal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, e := h.findFinal(c.Request.Context(), uint64(middleware.UID(c)), id)
	if e == nil && time.Now().After(row.FinalEnd) {
		if closeErr := h.expire(c.Request.Context(), uint64(middleware.UID(c)), id); closeErr != nil {
			e = closeErr
		} else {
			e = errFinalClosed
		}
	}
	if e != nil {
		writeErr(c, e)
		return
	}
	response.OK(c, row)
}
func (h *Handler) payFinal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in struct {
		PayType string `json:"pay_type"`
	}
	_ = c.ShouldBindJSON(&in)
	if in.PayType != "balance" {
		writeErr(c, errBadInput)
		return
	}
	uid := uint64(middleware.UID(c))
	var out finalRow
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		row, e := h.findFinalTx(tx, uid, id)
		if e != nil {
			return e
		}
		if row.Paid == 1 {
			return nil
		}
		now := time.Now()
		if row.Status != 1 {
			return errFinalClosed
		}
		if now.Before(row.FinalStart) {
			return errFinalNotOpen
		}
		if now.After(row.FinalEnd) {
			return errFinalClosed
		}
		up := tx.Table("qixi_crm_b_member_account").Where("user_id=? AND balance>=?", uid, row.PayPrice).UpdateColumn("balance", gorm.Expr("balance - ?", row.PayPrice))
		if up.Error != nil {
			return up.Error
		}
		if up.RowsAffected != 1 {
			return errBalance
		}
		key := fmt.Sprintf("presell-final:%d", row.ID)
		if e = tx.Exec("INSERT INTO qixi_crm_b_asset_ledger (user_id,asset_type,amount,reference_type,reference_id,idempotency_key) VALUES (?,'balance',?,'presell_final',?,?)", uid, -row.PayPrice, strconv.FormatUint(row.ID, 10), key).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_presell_order").Where("presell_order_id=? AND paid=0 AND status=1", row.ID).Updates(map[string]any{"paid": 1, "pay_type": 1, "pay_time": now}).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_order").Where("id=? AND user_id=? AND status=?", row.OrderID, uid, "awaiting_final").Update("status", "fulfilling").Error; e != nil {
			return e
		}
		out = row
		out.Paid = 1
		return nil
	})
	if err != nil {
		if errors.Is(err, errFinalClosed) {
			_ = h.expire(c.Request.Context(), uid, id)
		}
		writeErr(c, err)
		return
	}
	response.OK(c, out)
}
func (h *Handler) expire(ctx context.Context, uid, id uint64) error {
	return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row, e := h.findFinalTx(tx, uid, id)
		if e != nil {
			return e
		}
		if row.Paid == 1 || row.Status != 1 {
			return nil
		}
		if time.Now().Before(row.FinalEnd) {
			return errFinalNotOpen
		}
		if e = tx.Table("qixi_crm_b_presell_order").Where("presell_order_id=? AND paid=0 AND status=1", id).Update("status", -1).Error; e != nil {
			return e
		}
		var qty int
		if e = tx.Table("qixi_crm_b_order").Select("total_quantity").Where("id=?", row.OrderID).Scan(&qty).Error; e != nil {
			return e
		}
		if e = tx.Table("qixi_crm_b_order").Where("id=?", row.OrderID).Update("status", "final_timeout").Error; e != nil {
			return e
		}
		return tx.Table("qixi_crm_b_presell").Where("product_presell_id=?", row.ActivityID).UpdateColumn("stock", gorm.Expr("stock + ?", qty)).Error
	})
}
func (h *Handler) findFinal(ctx context.Context, uid, id uint64) (finalRow, error) {
	return h.findFinalTx(h.db.WithContext(ctx), uid, id)
}
func (h *Handler) findFinalTx(db *gorm.DB, uid, id uint64) (finalRow, error) {
	var row finalRow
	e := db.Table("qixi_crm_b_presell_order AS p").Select("p.presell_order_id,p.presell_order_sn,p.uid,p.order_id,p.product_presell_id,p.final_start_time,p.final_end_time,p.paid,p.status,p.pay_price,o.store_name_snapshot AS store_name").Joins("JOIN qixi_crm_b_order AS o ON o.id=p.order_id").Where("p.presell_order_id=? AND p.uid=?", id, uid).First(&row).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return row, errFinalNotFound
	}
	return row, e
}
func (h *Handler) loadQuote(ctx context.Context, db *gorm.DB, id uint64, lock bool) (quote, error) {
	q := quote{}
	x := db.WithContext(ctx).Table("qixi_crm_b_presell AS p").Select("p.product_presell_id,p.product_id,v.merchant_id,v.store_id,v.store_name,v.merchant_name,v.title,v.cover_url,p.price,p.down_price,p.final_price,p.stock,p.presell_type,p.final_start_time,p.final_end_time").Joins("JOIN qixi_crm_b_product_view AS v ON v.product_id=p.product_id").Where("p.product_presell_id=? AND p.status=1 AND p.is_show=1 AND p.is_del=0 AND p.product_status=1 AND p.action_status=1 AND p.start_time<=? AND p.end_time>=? AND v.sale_status=1", id, time.Now(), time.Now())
	if lock {
		x = x.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if e := x.First(&q).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return q, errUnavailable
		}
		return q, e
	}
	if q.PresellType != 1 && q.PresellType != 2 {
		return q, errUnavailable
	}
	if q.Price <= 0 || (q.PresellType == 2 && (q.DownPrice <= 0 || q.FinalPrice <= 0)) {
		return q, errUnavailable
	}
	return q, nil
}
func normalize(in *input) error {
	if in.ProductPresellID == 0 {
		return errBadInput
	}
	if in.CartNum == 0 {
		in.CartNum = 1
	}
	if in.CartNum < 1 || in.CartNum > 99 {
		return errBadInput
	}
	return nil
}
func quoteOut(q quote, n int) gin.H {
	pay := q.Price * float64(n)
	if q.PresellType == 2 {
		pay = q.DownPrice * float64(n)
	}
	return gin.H{"product_presell_id": q.ID, "product_id": q.ProductID, "store_name": q.StoreName, "image": q.CoverURL, "mer_id": q.MerchantID, "mer_name": q.MerchantName, "cart_num": n, "price": q.Price, "down_price": q.DownPrice, "final_price": q.FinalPrice, "pay_price": pay, "total_price": q.Price * float64(n), "stock": q.Stock, "activity_type": presellActivityType, "presell_type": q.PresellType, "final_start_time": q.FinalStart, "final_end_time": q.FinalEnd}
}
func sn(prefix string) string {
	var b [6]byte
	if _, e := rand.Read(b[:]); e != nil {
		return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s%d%x", prefix, time.Now().UnixMilli(), b)
}
func writeErr(c *gin.Context, e error) {
	switch {
	case errors.Is(e, errBadInput), errors.Is(e, errAddress), errors.Is(e, errUnavailable), errors.Is(e, errSoldOut), errors.Is(e, errFinalNotFound), errors.Is(e, errFinalNotOpen), errors.Is(e, errFinalClosed), errors.Is(e, errBalance):
		response.Fail(c, http.StatusBadRequest, e.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "预售操作失败")
	}
}
