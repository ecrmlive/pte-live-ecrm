// Package history implements private C-end product browsing history.
package history

import (
	"net/http"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/history/products/:id", h.recordProduct)
	r.GET("/history", h.list)
	r.POST("/history/delete/:id", h.deleteOne)
	r.POST("/history/batch/delete", h.deleteBatch)
}

type productView struct {
	ProductID uint64 `gorm:"column:product_id"`
	StoreID   uint64 `gorm:"column:store_id"`
}

type record struct {
	ID        uint64    `gorm:"column:id"`
	UserID    uint64    `gorm:"column:user_id"`
	ProductID uint64    `gorm:"column:product_id"`
	StoreID   uint64    `gorm:"column:store_id"`
	ViewedAt  time.Time `gorm:"column:viewed_at"`
}

type listRow struct {
	HistoryID uint64    `gorm:"column:history_id"`
	ProductID uint64    `gorm:"column:product_id"`
	StoreID   uint64    `gorm:"column:store_id"`
	StoreName string    `gorm:"column:store_name"`
	Title     string    `gorm:"column:title"`
	CoverURL  string    `gorm:"column:cover_url"`
	Price     float64   `gorm:"column:price"`
	Sales     int       `gorm:"column:sales"`
	ViewedAt  time.Time `gorm:"column:viewed_at"`
}

func (h *Handler) recordProduct(c *gin.Context) {
	productID, ok := positiveID(c.Param("id"))
	if !ok {
		bad(c, "商品 ID 错误")
		return
	}
	var product productView
	err := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_product_view").
		Select("product_id,store_id").
		Where("product_id = ? AND sale_status = ?", productID, 1).
		Take(&product).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, http.StatusNotFound, "商品不存在或已下架")
		return
	}
	if err != nil {
		internal(c)
		return
	}
	uid := uint64(middleware.UID(c))
	now := time.Now()
	var saved record
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("qixi_crm_b_user_browse_history").
			Where("user_id = ? AND product_id = ?", uid, productID).
			Order("id DESC").Clauses(clause.Locking{Strength: "UPDATE"}).Take(&saved).Error
		if err == gorm.ErrRecordNotFound {
			return tx.Table("qixi_crm_b_user_browse_history").Create(map[string]any{
				"user_id": uid, "product_id": productID, "store_id": product.StoreID, "viewed_at": now,
			}).Error
		}
		if err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_user_browse_history").Where("id = ? AND user_id = ?", saved.ID, uid).
			Updates(map[string]any{"store_id": product.StoreID, "viewed_at": now}).Error
	})
	if err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"product_id": productID, "recorded": true})
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pagination(c)
	uid := uint64(middleware.UID(c))
	query := h.db.WithContext(c.Request.Context()).
		Table("qixi_crm_b_user_browse_history AS h").
		Select("h.id AS history_id,h.product_id,h.store_id,p.store_name,p.title,p.cover_url,p.price,p.sales,h.viewed_at").
		Joins("JOIN qixi_crm_b_product_view AS p ON p.product_id = h.product_id AND p.sale_status = 1").
		Where("h.user_id = ?", uid)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		internal(c)
		return
	}
	rows := make([]listRow, 0)
	if err := query.Order("h.viewed_at DESC,h.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		internal(c)
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.view())
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) deleteOne(c *gin.Context) {
	id, ok := positiveID(c.Param("id"))
	if !ok {
		bad(c, "浏览记录 ID 错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").
		Where("id = ? AND user_id = ?", id, middleware.UID(c)).Delete(&record{}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"history_id": id, "deleted": true})
}

type batchRequest struct {
	HistoryIDs []uint64 `json:"history_ids"`
}

func (h *Handler) deleteBatch(c *gin.Context) {
	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		bad(c, "浏览记录参数错误")
		return
	}
	query := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history").Where("user_id = ?", middleware.UID(c))
	if len(req.HistoryIDs) > 0 {
		query = query.Where("id IN ?", req.HistoryIDs)
	}
	if err := query.Delete(&record{}).Error; err != nil {
		internal(c)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (r listRow) view() gin.H {
	return gin.H{
		"history_id": r.HistoryID, "product_id": r.ProductID, "store_id": r.StoreID, "store_name": r.StoreName,
		"title": r.Title, "image": r.CoverURL, "price": r.Price, "sales": r.Sales,
		"viewed_at": r.ViewedAt.Format("2006-01-02 15:04:05"),
	}
}

func positiveID(raw string) (uint64, bool) {
	id, err := strconv.ParseUint(raw, 10, 64)
	return id, err == nil && id > 0
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	return page, limit
}

func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func internal(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "浏览记录服务异常")
}
