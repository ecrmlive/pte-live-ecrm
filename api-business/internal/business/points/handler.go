package points

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

const pointsActivityType = 20

var (
	errBadInput = errors.New("积分商品参数错误")
	errNotFound = errors.New("积分商品不存在或已下架")
	errStock    = errors.New("积分商品库存不足")
	errPoints   = errors.New("当前积分不足")
	errAddress  = errors.New("收货地址不存在或无权访问")
	errOrder    = errors.New("积分订单不存在或无权访问")
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/points/products", h.list)
	r.GET("/points/products/:id", h.detail)
}
func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.GET("/integral", h.integral)
	r.POST("/order/v3/check", h.check)
	r.POST("/order/v3/create", h.create)
	r.POST("/order/points/pay/:id", h.pay)
}

type product struct {
	ProductID      uint64  `gorm:"column:product_id"`
	MerchantID     uint64  `gorm:"column:merchant_id"`
	StoreID        uint64  `gorm:"column:store_id"`
	MerchantName   string  `gorm:"column:merchant_name"`
	StoreName      string  `gorm:"column:store_name"`
	Title          string  `gorm:"column:title"`
	CoverURL       string  `gorm:"column:cover_url"`
	OriginalPrice  float64 `gorm:"column:original_price"`
	PointsRequired int64   `gorm:"column:points_required"`
	Stock          int     `gorm:"column:stock"`
	SaleStatus     int     `gorm:"column:sale_status"`
}

func (product) TableName() string { return "qixi_crm_b_points_product_view" }

type input struct {
	ProductID      uint64 `json:"product_id"`
	CartNum        int    `json:"cart_num"`
	AddressID      uint64 `json:"address_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
type account struct {
	Points int64 `gorm:"column:points"`
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

func (address) TableName() string { return "qixi_crm_b_user_address" }

type group struct {
	ID             uint64  `gorm:"column:id"`
	OrderNo        string  `gorm:"column:order_no"`
	UserID         uint64  `gorm:"column:user_id"`
	PayStatus      string  `gorm:"column:pay_status"`
	PointsAmount   int64   `gorm:"column:points_amount"`
	ActivityType   int     `gorm:"column:activity_type"`
	IdempotencyKey string  `gorm:"column:idempotency_key"`
	TotalQuantity  int     `gorm:"column:total_quantity"`
	PayAmount      float64 `gorm:"column:pay_amount"`
}

func (group) TableName() string { return "qixi_crm_b_group_order" }

type order struct {
	ID           uint64 `gorm:"column:id"`
	GroupOrderID uint64 `gorm:"column:group_order_id"`
	ProductID    uint64 `gorm:"column:product_id"`
	Quantity     int    `gorm:"column:quantity"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	q := h.db.WithContext(c.Request.Context()).Model(&product{}).Where("sale_status = ?", 1)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c)
		return
	}
	var rows []product
	if err := q.Order("product_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		fail(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, p := range rows {
		list = append(list, productResponse(p))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}
func (h *Handler) detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.loadProduct(c.Request.Context(), id, false)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, productResponse(p))
}
func (h *Handler) integral(c *gin.Context) {
	var a account
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account").Select("COALESCE(points,0) AS points").Where("user_id = ?", middleware.UID(c)).Scan(&a).Error; err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"integral": a.Points})
}
func (h *Handler) check(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, errBadInput)
		return
	}
	p, n, points, err := h.quote(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	balance, err := h.userPoints(c.Request.Context(), uint64(middleware.UID(c)))
	if err != nil {
		fail(c)
		return
	}
	response.OK(c, gin.H{"product_id": p.ProductID, "store_name": p.StoreName, "mer_name": p.MerchantName, "image": p.CoverURL, "cart_num": n, "integral": points, "user_integral": balance, "stock": p.Stock, "pay_price": 0, "activity_type": pointsActivityType})
}
func (h *Handler) create(c *gin.Context) {
	var in input
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, errBadInput)
		return
	}
	uid := uint64(middleware.UID(c))
	in.IdempotencyKey = strings.TrimSpace(in.IdempotencyKey)
	if in.AddressID == 0 || len(in.IdempotencyKey) < 12 || len(in.IdempotencyKey) > 128 {
		writeErr(c, errBadInput)
		return
	}
	p, n, need, err := h.quote(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	balance, err := h.userPoints(c.Request.Context(), uid)
	if err != nil {
		fail(c)
		return
	}
	if balance < need {
		writeErr(c, errPoints)
		return
	}
	var addr address
	if err := h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", in.AddressID, uid).First(&addr).Error; err != nil {
		writeErr(c, errAddress)
		return
	}
	snapshot, _ := json.Marshal(addr)
	var created group
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error
		if err == nil {
			if created.ActivityType != pointsActivityType {
				return errBadInput
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		created = group{OrderNo: orderNo("PG"), UserID: uid, PayStatus: "pending", PointsAmount: need, ActivityType: pointsActivityType, IdempotencyKey: in.IdempotencyKey, TotalQuantity: n, PayAmount: 0}
		if err := tx.Table("qixi_crm_b_group_order").Create(map[string]any{"order_no": created.OrderNo, "user_id": uid, "total_amount": 0, "discount_amount": 0, "freight_amount": 0, "pay_amount": 0, "total_quantity": n, "recipient_snapshot": string(snapshot), "pay_status": "pending", "idempotency_key": in.IdempotencyKey, "activity_type": pointsActivityType, "points_amount": need}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ? AND idempotency_key = ?", uid, in.IdempotencyKey).First(&created).Error; err != nil {
			return err
		}
		var child struct {
			ID uint64 `gorm:"column:id"`
		}
		if err := tx.Table("qixi_crm_b_order").Create(map[string]any{"group_order_id": created.ID, "order_no": orderNo("PS"), "merchant_id": p.MerchantID, "merchant_name_snapshot": p.MerchantName, "store_id": p.StoreID, "store_name_snapshot": p.StoreName, "user_id": uid, "total_amount": 0, "discount_amount": 0, "freight_amount": 0, "pay_amount": 0, "total_quantity": n, "recipient_snapshot": string(snapshot), "status": "pending_pay", "activity_type": pointsActivityType, "points_amount": need}).Error; err != nil {
			return err
		}
		if err := tx.Table("qixi_crm_b_order").Select("id").Where("group_order_id = ? AND user_id = ?", created.ID, uid).Take(&child).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_order_item").Create(map[string]any{"order_id": child.ID, "product_id": p.ProductID, "sku_key": "", "title_snapshot": p.Title, "cover_url_snapshot": p.CoverURL, "spec_snapshot": "{}", "unit_price": 0, "quantity": n, "refund_quantity": 0}).Error
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"group_order_id": created.ID, "group_order_sn": created.OrderNo, "pay_price": 0, "integral": created.PointsAmount, "paid": created.PayStatus == "paid"})
}
func (h *Handler) pay(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uid := uint64(middleware.UID(c))
	if id == 0 {
		writeErr(c, errBadInput)
		return
	}
	var out group
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ? AND activity_type = ?", id, uid, pointsActivityType).First(&out).Error; err != nil {
			return errOrder
		}
		if out.PayStatus == "paid" {
			return nil
		}
		if out.PayStatus != "pending" {
			return errBadInput
		}
		if out.PointsAmount <= 0 {
			return errBadInput
		}
		var lines []struct {
			ProductID uint64 `gorm:"column:product_id"`
			Quantity  int    `gorm:"column:quantity"`
		}
		if err := tx.Table("qixi_crm_b_order_item AS i").Select("i.product_id,i.quantity").Joins("JOIN qixi_crm_b_order AS o ON o.id=i.order_id").Where("o.group_order_id = ?", out.ID).Find(&lines).Error; err != nil {
			return err
		}
		for _, line := range lines {
			updated := tx.Table("qixi_crm_b_points_product_view").Where("product_id = ? AND sale_status = 1 AND stock >= ?", line.ProductID, line.Quantity).Update("stock", gorm.Expr("stock - ?", line.Quantity))
			if updated.Error != nil {
				return updated.Error
			}
			if updated.RowsAffected != 1 {
				return errStock
			}
		}
		updated := tx.Table("qixi_crm_b_member_account").Where("user_id = ? AND points >= ?", uid, out.PointsAmount).Update("points", gorm.Expr("points - ?", out.PointsAmount))
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return errPoints
		}
		key := fmt.Sprintf("points-order:%d", out.ID)
		if err := tx.Table("qixi_crm_b_asset_ledger").Create(map[string]any{"user_id": uid, "asset_type": "points", "amount": -out.PointsAmount, "reference_type": "points_order", "reference_id": strconv.FormatUint(out.ID, 10), "idempotency_key": key}).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		if err := tx.Table("qixi_crm_b_group_order").Where("id = ? AND pay_status = 'pending'", out.ID).Updates(map[string]any{"pay_status": "paid", "paid_at": now}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_order").Where("group_order_id = ? AND status = 'pending_pay'", out.ID).Updates(map[string]any{"status": "paid", "paid_at": now}).Error
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"group_order_id": out.ID, "group_order_sn": out.OrderNo, "paid": true, "pay_price": 0, "integral": out.PointsAmount})
}
func (h *Handler) quote(ctx context.Context, in input) (product, int, int64, error) {
	if in.ProductID == 0 {
		return product{}, 0, 0, errBadInput
	}
	n := in.CartNum
	if n <= 0 {
		n = 1
	}
	if n > 100 {
		return product{}, 0, 0, errBadInput
	}
	p, err := h.loadProduct(ctx, in.ProductID, false)
	if err != nil {
		return product{}, 0, 0, err
	}
	if p.Stock < n {
		return product{}, 0, 0, errStock
	}
	if p.PointsRequired <= 0 {
		return product{}, 0, 0, errBadInput
	}
	need := p.PointsRequired * int64(n)
	return p, n, need, nil
}
func (h *Handler) loadProduct(ctx context.Context, id uint64, lock bool) (product, error) {
	var p product
	q := h.db.WithContext(ctx)
	if lock {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := q.Where("product_id=? AND sale_status=1", id).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return p, errNotFound
	}
	return p, err
}
func (h *Handler) userPoints(ctx context.Context, uid uint64) (int64, error) {
	var a account
	err := h.db.WithContext(ctx).Table("qixi_crm_b_member_account").Select("COALESCE(points,0) AS points").Where("user_id=?", uid).Scan(&a).Error
	return a.Points, err
}
func productResponse(p product) gin.H {
	return gin.H{"id": p.ProductID, "product_id": p.ProductID, "mer_id": p.MerchantID, "mer_name": p.MerchantName, "store_name": p.StoreName, "title": p.Title, "image": p.CoverURL, "price": "0.00", "ot_price": fmt.Sprintf("%.2f", p.OriginalPrice), "integral": p.PointsRequired, "sales": 0, "stock": p.Stock, "product_type": 1}
}
func pagination(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}
	return p, l
}
func orderNo(prefix string) string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s%d%x", prefix, time.Now().UnixMilli(), b)
}
func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "积分商城服务异常")
}
func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errNotFound), errors.Is(err, errOrder):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, errBadInput), errors.Is(err, errStock), errors.Is(err, errPoints), errors.Is(err, errAddress):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		fail(c)
	}
}
