package comment

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler           { return &Handler{db: db} }
func (h *Handler) RegisterPublic(r gin.IRoutes) { r.GET("/products/:id/comments", h.publicList) }
func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/comments", h.create)
	r.GET("/comments/mine", h.mine)
}

type createRequest struct {
	OrderItemID uint64 `json:"order_item_id"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
}
type ownedItem struct {
	ID        uint64 `gorm:"column:id"`
	ProductID uint64 `gorm:"column:product_id"`
	StoreID   uint64 `gorm:"column:store_id"`
	Status    string `gorm:"column:status"`
}
type row struct {
	ID                uint64    `gorm:"column:id"`
	OrderItemID       uint64    `gorm:"column:order_item_id"`
	UserID            uint64    `gorm:"column:user_id"`
	ProductID         uint64    `gorm:"column:product_id"`
	StoreID           uint64    `gorm:"column:store_id"`
	Score             int       `gorm:"column:score"`
	Content           string    `gorm:"column:content"`
	Media             string    `gorm:"column:media"`
	Reply             string    `gorm:"column:reply_content"`
	Source            string    `gorm:"column:source"`
	VirtualAuthorName string    `gorm:"column:virtual_author_name"`
	Sort              int       `gorm:"column:sort"`
	Status            string    `gorm:"column:status"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	Title             string    `gorm:"column:title"`
}

func (h *Handler) create(c *gin.Context) {
	var req createRequest
	if c.ShouldBindJSON(&req) != nil {
		bad(c, "评价参数错误")
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.OrderItemID == 0 || len([]rune(req.Content)) == 0 || req.Score < 1 || req.Score > 5 || len([]rune(req.Content)) > 2000 {
		bad(c, "评价内容不合法")
		return
	}
	uid := uint64(middleware.UID(c))
	var created row
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var item ownedItem
		if err := tx.Table("qixi_crm_b_order_item AS oi").Select("oi.id,oi.product_id,oi.store_id,o.status").Joins("JOIN qixi_crm_b_order AS o ON o.id=oi.order_id").Where("oi.id=? AND o.user_id=?", req.OrderItemID, uid).Take(&item).Error; err != nil {
			return err
		}
		if item.Status != "completed" {
			return errNotEligible
		}
		created = row{OrderItemID: item.ID, UserID: uid, ProductID: item.ProductID, StoreID: item.StoreID, Score: req.Score, Content: req.Content, Status: "pending"}
		if err := tx.Table("qixi_crm_b_product_comment").Create(map[string]any{"order_item_id": created.OrderItemID, "user_id": created.UserID, "product_id": created.ProductID, "store_id": created.StoreID, "score": created.Score, "content": created.Content, "media": "[]", "status": created.Status}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_product_comment").Where("order_item_id=?", item.ID).Take(&created).Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单商品不存在")
			return
		}
		if err == errNotEligible {
			bad(c, "仅已完成订单可以评价")
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			bad(c, "该商品已评价")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "评价提交失败")
		return
	}
	response.OK(c, created.view())
}
func (h *Handler) mine(c *gin.Context) {
	page, limit := page(c)
	uid := uint64(middleware.UID(c))
	q := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_product_comment AS pc").Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id=pc.product_id").Where("pc.user_id=?", uid)
	var total int64
	if q.Count(&total).Error != nil {
		response.Fail(c, 500, "评价查询失败")
		return
	}
	var rows []row
	if q.Select("pc.*,p.title").Order("pc.id DESC").Offset((page-1)*limit).Limit(limit).Scan(&rows).Error != nil {
		response.Fail(c, 500, "评价查询失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, x := range rows {
		list = append(list, x.view())
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}
func (h *Handler) publicList(c *gin.Context) {
	id, ok := id(c.Param("id"))
	if !ok {
		bad(c, "商品 ID 错误")
		return
	}
	var rows []row
	if h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_product_comment").Select("id,COALESCE(order_item_id,0) AS order_item_id,user_id,product_id,store_id,score,content,media,reply_content,source,virtual_author_name,sort,status,created_at").Where("product_id=? AND status='published' AND deleted_at IS NULL", id).Order("sort DESC,id DESC").Limit(30).Scan(&rows).Error != nil {
		response.Fail(c, 500, "评价查询失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, x := range rows {
		list = append(list, x.view())
	}
	response.OK(c, gin.H{"list": list})
}
func (r row) view() gin.H {
	return gin.H{"comment_id": r.ID, "order_item_id": r.OrderItemID, "product_id": r.ProductID, "score": r.Score, "content": r.Content, "media": r.Media, "reply_content": r.Reply, "source": r.Source, "virtual_author_name": r.VirtualAuthorName, "sort": r.Sort, "status": r.Status, "title": r.Title, "create_time": r.CreatedAt.Format("2006-01-02 15:04:05")}
}

var errNotEligible = &commentError{}

type commentError struct{}

func (*commentError) Error() string { return "not eligible" }
func id(raw string) (uint64, bool) {
	v, e := strconv.ParseUint(raw, 10, 64)
	return v, e == nil && v > 0
}
func page(c *gin.Context) (int, int) {
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
func bad(c *gin.Context, msg string) { response.Fail(c, http.StatusBadRequest, msg) }
