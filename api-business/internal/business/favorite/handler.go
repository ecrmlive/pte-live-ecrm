// Package favorite implements the C-end product/store collection workflow.
// Collections are private user assets and are intentionally stored only in the
// business database; the handler only reads the published consumer projections.
package favorite

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/favorites", h.list)
	r.GET("/favorites/products/:id", h.productStatus)
	r.PUT("/favorites/products/:id", h.addProduct)
	r.DELETE("/favorites/products/:id", h.deleteProduct)
	r.GET("/favorites/stores/:id", h.storeStatus)
	r.PUT("/favorites/stores/:id", h.addStore)
	r.DELETE("/favorites/stores/:id", h.deleteStore)
}

func (h *Handler) list(c *gin.Context) {
	uid := uint64(middleware.UID(c))
	switch c.DefaultQuery("type", "product") {
	case "product":
		var rows []productFavorite
		err := h.db.WithContext(c.Request.Context()).
			Table("qixi_crm_b_product_favorite AS f").
			Select("p.product_id,p.merchant_id,p.store_id,p.merchant_name,p.store_name,p.category_id,p.title,p.cover_url,p.price,p.original_price,p.sales,p.stock,f.created_at").
			Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = f.product_id AND p.sale_status = 1").
			Where("f.user_id = ?", uid).Order("f.created_at DESC").Scan(&rows).Error
		if err != nil {
			internal(c)
			return
		}
		list := make([]gin.H, 0, len(rows))
		for _, row := range rows {
			list = append(list, row.view())
		}
		response.OK(c, gin.H{"type": "product", "list": list})
	case "store":
		var rows []storeFavorite
		err := h.db.WithContext(c.Request.Context()).
			Table("qixi_crm_b_user_follow_store AS f").
			Select("s.store_id,s.merchant_id,s.store_name,s.store_app_id,COUNT(all_follow.user_id) AS follower_count,f.created_at").
			Joins("JOIN qixi_crm_b_store_view AS s ON s.store_id = f.store_id AND s.status = 1").
			Joins("LEFT JOIN qixi_crm_b_user_follow_store AS all_follow ON all_follow.store_id = s.store_id").
			Where("f.user_id = ?", uid).Group("s.store_id,s.merchant_id,s.store_name,s.store_app_id,f.created_at").Order("f.created_at DESC").Scan(&rows).Error
		if err != nil {
			internal(c)
			return
		}
		list := make([]gin.H, 0, len(rows))
		for _, row := range rows {
			list = append(list, row.view())
		}
		response.OK(c, gin.H{"type": "store", "list": list})
	default:
		response.Fail(c, http.StatusBadRequest, "收藏类型错误")
	}
}

func (h *Handler) productStatus(c *gin.Context) {
	h.status(c, "qixi_crm_b_product_favorite", "product_id", true)
}
func (h *Handler) storeStatus(c *gin.Context) {
	h.status(c, "qixi_crm_b_user_follow_store", "store_id", false)
}

func (h *Handler) status(c *gin.Context, table, column string, product bool) {
	id, ok := pathID(c)
	if !ok {
		badID(c)
		return
	}
	if product {
		if !h.productExists(c, id) {
			response.Fail(c, http.StatusNotFound, "商品不存在或已下架")
			return
		}
	} else if !h.storeExists(c, id) {
		response.Fail(c, http.StatusNotFound, "店铺不存在或已停用")
		return
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Table(table).Where("user_id = ? AND "+column+" = ?", middleware.UID(c), id).Count(&count).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"id": id, "followed": count > 0})
}

func (h *Handler) addProduct(c *gin.Context) {
	h.add(c, "qixi_crm_b_product_favorite", "product_id", true)
}
func (h *Handler) addStore(c *gin.Context) {
	h.add(c, "qixi_crm_b_user_follow_store", "store_id", false)
}

func (h *Handler) add(c *gin.Context, table, column string, product bool) {
	id, ok := pathID(c)
	if !ok {
		badID(c)
		return
	}
	if product {
		if !h.productExists(c, id) {
			response.Fail(c, http.StatusNotFound, "商品不存在或已下架")
			return
		}
	} else if !h.storeExists(c, id) {
		response.Fail(c, http.StatusNotFound, "店铺不存在或已停用")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table(table).Clauses(clause.OnConflict{DoNothing: true}).Create(map[string]any{"user_id": middleware.UID(c), column: id}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"id": id, "followed": true})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	h.remove(c, "qixi_crm_b_product_favorite", "product_id")
}
func (h *Handler) deleteStore(c *gin.Context) {
	h.remove(c, "qixi_crm_b_user_follow_store", "store_id")
}

func (h *Handler) remove(c *gin.Context, table, column string) {
	id, ok := pathID(c)
	if !ok {
		badID(c)
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table(table).Where("user_id = ? AND "+column+" = ?", middleware.UID(c), id).Delete(&collectionRow{}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"id": id, "followed": false})
}

func (h *Handler) productExists(c *gin.Context, id uint64) bool {
	var count int64
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").Where("product_id = ? AND sale_status = 1", id).Count(&count).Error == nil && count == 1
}
func (h *Handler) storeExists(c *gin.Context, id uint64) bool {
	var count int64
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_store_view").Where("store_id = ? AND status = 1", id).Count(&count).Error == nil && count == 1
}

type productFavorite struct {
	ProductID     uint64   `gorm:"column:product_id"`
	MerchantID    uint64   `gorm:"column:merchant_id"`
	StoreID       uint64   `gorm:"column:store_id"`
	MerchantName  string   `gorm:"column:merchant_name"`
	StoreName     string   `gorm:"column:store_name"`
	CategoryID    uint64   `gorm:"column:category_id"`
	Title         string   `gorm:"column:title"`
	CoverURL      string   `gorm:"column:cover_url"`
	Price         float64  `gorm:"column:price"`
	OriginalPrice *float64 `gorm:"column:original_price"`
	Sales         int      `gorm:"column:sales"`
	Stock         int      `gorm:"column:stock"`
}

func (p productFavorite) view() gin.H {
	return gin.H{"id": p.ProductID, "mer_id": p.MerchantID, "store_id": p.StoreID, "mer_name": p.MerchantName, "store_name": p.StoreName, "shop_name": p.StoreName, "category_id": p.CategoryID, "title": p.Title, "image": p.CoverURL, "price": p.Price, "ot_price": p.OriginalPrice, "sales": p.Sales, "stock": p.Stock}
}

type storeFavorite struct {
	StoreID       uint64 `gorm:"column:store_id"`
	MerchantID    uint64 `gorm:"column:merchant_id"`
	StoreName     string `gorm:"column:store_name"`
	StoreAppID    string `gorm:"column:store_app_id"`
	FollowerCount int64  `gorm:"column:follower_count"`
}

func (s storeFavorite) view() gin.H {
	return gin.H{"store_id": s.StoreID, "mer_id": s.MerchantID, "store_name": s.StoreName, "merchant_app_id": s.StoreAppID, "follower_count": s.FollowerCount}
}

type collectionRow struct{}

func pathID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
func badID(c *gin.Context)    { response.Fail(c, http.StatusBadRequest, "ID 错误") }
func internal(c *gin.Context) { response.Fail(c, http.StatusInternalServerError, "收藏服务异常") }
