// Package cart owns C-end cart reads and writes in qixi_crm_business.
//
// Product data is deliberately read from qixi_crm_b_product_view: the cart is
// not permitted to reach into a merchant-owned table. Stock is checked here
// for feedback only; the order transaction remains the final stock authority.
package cart

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errStock = errors.New("库存不足或商品已下架")

// Handler uses only the business database tables.
type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/cart", h.List)
	r.POST("/cart", h.Add)
	r.PUT("/cart/:id", h.Update)
	r.DELETE("/cart/:id", h.Delete)
}

type cartRow struct {
	ID        uint64 `gorm:"column:id"`
	UserID    uint64 `gorm:"column:user_id"`
	StoreID   uint64 `gorm:"column:store_id"`
	ProductID uint64 `gorm:"column:product_id"`
	SKUKey    string `gorm:"column:sku_key"`
	Quantity  int    `gorm:"column:quantity"`
}

func (cartRow) TableName() string { return "qixi_crm_b_cart" }

type productRow struct {
	ProductID    uint64  `gorm:"column:product_id"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	StoreID      uint64  `gorm:"column:store_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	Stock        int     `gorm:"column:stock"`
	SaleStatus   int8    `gorm:"column:sale_status"`
}

func (productRow) TableName() string { return "qixi_crm_b_product_view" }

type cartView struct {
	CartID       uint64  `gorm:"column:cart_id"`
	ProductID    uint64  `gorm:"column:product_id"`
	SKUKey       string  `gorm:"column:sku_key"`
	Quantity     int     `gorm:"column:quantity"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	Stock        int     `gorm:"column:stock"`
	SaleStatus   int8    `gorm:"column:sale_status"`
}

type addRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	CartNum   int    `json:"cart_num"`
	// ProductAttrUnique is kept for app compatibility. A SKU consumption
	// projection will replace this numeric fallback when SKU migration lands.
	ProductAttrUnique string `json:"product_attr_unique"`
}

type updateRequest struct {
	CartNum *int `json:"cart_num" binding:"required"`
}

type bucket struct {
	MerchantID   uint64  `json:"mer_id"`
	MerchantName string  `json:"mer_name"`
	Subtotal     float64 `json:"subtotal"`
	Items        []gin.H `json:"items"`
}

func (h *Handler) List(c *gin.Context) {
	rows, err := h.list(c, uint64(middleware.UID(c)))
	if err != nil {
		fail(c, err)
		return
	}

	buckets := map[uint64]*bucket{}
	totalNum := 0
	totalPrice := 0.0
	for _, row := range rows {
		item := itemResponse(row)
		current := buckets[row.MerchantID]
		if current == nil {
			current = &bucket{
				MerchantID:   row.MerchantID,
				MerchantName: row.MerchantName,
				Items:        []gin.H{},
			}
			buckets[row.MerchantID] = current
		}
		current.Items = append(current.Items, item)
		if isAvailable(row) {
			subtotal := row.Price * float64(row.Quantity)
			current.Subtotal += subtotal
			totalNum += row.Quantity
			totalPrice += subtotal
		}
	}

	list := make([]bucket, 0, len(buckets))
	for _, value := range buckets {
		list = append(list, *value)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].MerchantID < list[j].MerchantID })
	response.OK(c, gin.H{"list": list, "total_num": totalNum, "total_price": totalPrice})
}

func (h *Handler) Add(c *gin.Context) {
	var req addRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.CartNum == 0 {
		req.CartNum = 1
	}
	if req.CartNum < 1 {
		response.Fail(c, http.StatusBadRequest, "购买数量无效")
		return
	}

	uid := uint64(middleware.UID(c))
	skuKey := resolveSKUKey(req.ProductID, req.ProductAttrUnique)

	var cartID uint64
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var product productRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND sale_status = 1", req.ProductID).
			First(&product).Error; err != nil {
			return err
		}
		if err := verifyMerchantContext(c, product); err != nil {
			return err
		}
		if product.Stock < req.CartNum {
			return errStock
		}

		var row cartRow
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND product_id = ? AND sku_key = ?", uid, req.ProductID, skuKey).
			First(&row).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			row = cartRow{
				UserID: uid, StoreID: product.StoreID, ProductID: req.ProductID,
				SKUKey: skuKey, Quantity: req.CartNum,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			if row.Quantity+req.CartNum > product.Stock {
				return errStock
			}
			if err := tx.Model(&cartRow{}).
				Where("id = ? AND user_id = ?", row.ID, uid).
				Update("quantity", row.Quantity+req.CartNum).Error; err != nil {
				return err
			}
		}
		cartID = row.ID
		return nil
	})
	if err != nil {
		writeError(c, err)
		return
	}

	view, err := h.findView(c, uid, cartID)
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, itemResponse(view))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "购物车 ID 错误")
		return
	}
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.CartNum == nil || *req.CartNum < 1 {
		response.Fail(c, http.StatusBadRequest, "购买数量无效")
		return
	}

	uid := uint64(middleware.UID(c))
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row cartRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", id, uid).First(&row).Error; err != nil {
			return err
		}
		var product productRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND sale_status = 1", row.ProductID).First(&product).Error; err != nil {
			return err
		}
		if err := verifyMerchantContext(c, product); err != nil {
			return err
		}
		if product.Stock < *req.CartNum {
			return errStock
		}
		return tx.Model(&cartRow{}).Where("id = ? AND user_id = ?", id, uid).
			Update("quantity", *req.CartNum).Error
	})
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "购物车 ID 错误")
		return
	}
	uid := uint64(middleware.UID(c))
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var row cartRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", id, uid).First(&row).Error; err != nil {
			return err
		}
		var product productRow
		if err := tx.Where("product_id = ?", row.ProductID).First(&product).Error; err != nil {
			return err
		}
		if err := verifyMerchantContext(c, product); err != nil {
			return err
		}
		result := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&cartRow{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) list(c *gin.Context, uid uint64) ([]cartView, error) {
	var rows []cartView
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_cart AS c").
		Select(cartViewColumns).
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id = c.product_id").
		Where("c.user_id = ?", uid).Order("c.id DESC").Scan(&rows).Error
	return rows, err
}

func (h *Handler) findView(c *gin.Context, uid, id uint64) (cartView, error) {
	var row cartView
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_cart AS c").
		Select(cartViewColumns).
		Joins("INNER JOIN qixi_crm_b_product_view AS p ON p.product_id = c.product_id").
		// cart_id is a response alias, not a column on qixi_crm_b_cart. Using
		// First(&cartView) makes GORM append ORDER BY c.cart_id and breaks the
		// post-create readback on MySQL.
		Where("c.id = ? AND c.user_id = ?", id, uid).Limit(1).Scan(&row).Error
	return row, err
}

const cartViewColumns = "c.id AS cart_id, c.product_id, c.sku_key, c.quantity, " +
	"p.merchant_id, p.merchant_name, p.store_name, p.title, p.cover_url, " +
	"p.price, p.stock, p.sale_status"

func resolveSKUKey(productID uint64, unique string) string {
	unique = strings.TrimSpace(unique)
	if unique == "" {
		return strconv.FormatUint(productID, 10)
	}
	return unique
}

func isAvailable(row cartView) bool {
	return row.SaleStatus == 1 && row.Stock >= row.Quantity
}

var errStoreContext = errors.New("请先进入对应店铺后再加入购物车")

// verifyMerchantContext binds a merchant product to the server-issued C-end
// store-context token. X-AppId alone is client input and must never authorize
// a cart write.
func verifyMerchantContext(c *gin.Context, product productRow) error {
	if product.MerchantID == 0 && product.StoreID == 0 {
		return nil
	}
	claims := middleware.ClaimsFrom(c)
	appID := strings.TrimSpace(c.GetHeader("X-AppId"))
	if claims == nil || claims.AuthContext != authjwt.ContextStore ||
		claims.MerchantID != uint(product.MerchantID) || claims.StoreID != uint(product.StoreID) ||
		claims.MerchantAppID == "" || appID != claims.MerchantAppID {
		return errStoreContext
	}
	return nil
}

func itemResponse(row cartView) gin.H {
	failed := 0
	if !isAvailable(row) {
		failed = 1
	}
	return gin.H{
		"cart_id":             row.CartID,
		"product_id":          row.ProductID,
		"product_attr_unique": row.SKUKey,
		"cart_num":            row.Quantity,
		"mer_id":              row.MerchantID,
		"mer_name":            row.MerchantName,
		"store_name":          row.StoreName,
		"title":               row.Title,
		"image":               row.CoverURL,
		"price":               row.Price,
		"stock":               row.Stock,
		"is_fail":             failed,
	}
}

func writeError(c *gin.Context, err error) {
	if errors.Is(err, errStock) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, errStoreContext) {
		response.Fail(c, http.StatusForbidden, err.Error())
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "商品或购物车项不存在")
		return
	}
	fail(c, err)
}

func fail(c *gin.Context, _ error) {
	response.Fail(c, http.StatusInternalServerError, "购物车服务异常")
}
