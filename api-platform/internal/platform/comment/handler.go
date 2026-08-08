// Package comment exposes platform-only moderation of business-owned product comments.
package comment

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	commentcommand "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/commentmoderation"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/queryfilter"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	businessDB, adminDB *gorm.DB
	commands            *commentcommand.Client
}

func NewHandler(businessDB, adminDB *gorm.DB, commands *commentcommand.Client) *Handler {
	return &Handler{businessDB: businessDB, adminDB: adminDB, commands: commands}
}
func (h *Handler) Register(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	read := middleware.RequireAdminMenu(h.adminDB, "product.comment.review")
	virtual := middleware.RequireAdminMenu(h.adminDB, "product.comment.virtual.manage")
	sort := middleware.RequireAdminMenu(h.adminDB, "product.comment.sort")
	remove := middleware.RequireAdminMenu(h.adminDB, "product.comment.delete")
	r.GET("/product/comments", platform, read, h.List)
	r.GET("/product/comments/:id", platform, read, h.Get)
	r.POST("/product/comments/:id/moderate", platform, read, h.Moderate)
	r.POST("/product/comments/virtual", platform, virtual, h.CreateVirtual)
	r.PUT("/product/comments/:id/virtual", platform, virtual, h.UpdateVirtual)
	r.PUT("/product/comments/:id/sort", platform, sort, h.SortVirtual)
	r.DELETE("/product/comments/:id", platform, remove, h.DeleteVirtual)
}

type row struct {
	ID                  uint64    `gorm:"column:id" json:"id"`
	ProductID           uint64    `gorm:"column:product_id" json:"product_id"`
	StoreID             uint64    `gorm:"column:store_id" json:"store_id"`
	Score               int       `gorm:"column:score" json:"score"`
	Content             string    `gorm:"column:content" json:"content"`
	Media               string    `gorm:"column:media" json:"media"`
	ReplyContent        string    `gorm:"column:reply_content" json:"reply_content"`
	Source              string    `gorm:"column:source" json:"source"`
	VirtualAuthorName   string    `gorm:"column:virtual_author_name" json:"virtual_author_name"`
	VirtualAuthorAvatar string    `gorm:"column:virtual_author_avatar" json:"virtual_author_avatar"`
	Sort                int       `gorm:"column:sort" json:"sort"`
	Status              string    `gorm:"column:status" json:"status"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	ProductTitle        string    `gorm:"column:product_title" json:"product_title"`
	ProductCover        string    `gorm:"column:product_cover" json:"product_cover"`
	UserName            string    `gorm:"column:user_name" json:"user_name"`
}
type moderateInput struct {
	Action         string `json:"action"`
	Note           string `json:"note"`
	IdempotencyKey string `json:"idempotency_key"`
}
type virtualInput struct {
	ProductID           uint64   `json:"product_id"`
	Score               int      `json:"score"`
	Content             string   `json:"content"`
	VirtualAuthorName   string   `json:"virtual_author_name"`
	VirtualAuthorAvatar string   `json:"virtual_author_avatar"`
	Sort                int      `json:"sort"`
	AttachmentIDs       []uint64 `json:"attachment_ids"`
	IdempotencyKey      string   `json:"idempotency_key"`
}
type sortInput struct {
	Sort           int    `json:"sort"`
	IdempotencyKey string `json:"idempotency_key"`
}
type deleteInput struct {
	IdempotencyKey string `json:"idempotency_key"`
	Note           string `json:"note"`
}

const columns = "pc.id,pc.product_id,pc.store_id,pc.score,pc.content,pc.media,pc.reply_content,pc.source,pc.virtual_author_name,pc.virtual_author_avatar,pc.sort,pc.status,pc.created_at,COALESCE(p.title,'') AS product_title,COALESCE(p.cover_url,'') AS product_cover,CASE WHEN pc.source='virtual' THEN pc.virtual_author_name ELSE COALESCE(u.nickname,'') END AS user_name"

func (h *Handler) List(c *gin.Context) {
	page, limit := paging(c)
	q := h.base(c)
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if !validStatus(status) {
			bad(c, "评论状态错误")
			return
		}
		q = q.Where("pc.status=?", status)
	}
	if productID, raw := uintParam(c.Query("product_id")); raw {
		if productID == 0 {
			bad(c, "商品 ID 错误")
			return
		}
		q = q.Where("pc.product_id=?", productID)
	}
	if userName := strings.TrimSpace(c.Query("user_name")); userName != "" {
		like := "%" + userName + "%"
		q = q.Where(
			"(pc.source='virtual' AND pc.virtual_author_name LIKE ?) OR (pc.source<>'virtual' AND u.nickname LIKE ?)",
			like, like,
		)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("pc.content LIKE ? OR p.title LIKE ? OR CAST(pc.product_id AS CHAR) = ?", like, like, keyword)
	}
	q = queryfilter.ApplyCreatedAtRange(q, c, "pc.created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		failed(c)
		return
	}
	orderSQL := "pc.sort DESC,pc.id DESC"
	if sortField := strings.TrimSpace(c.Query("sort_field")); sortField == "score" {
		if strings.EqualFold(strings.TrimSpace(c.Query("sort_order")), "asc") {
			orderSQL = "pc.score ASC,pc.id DESC"
		} else {
			orderSQL = "pc.score DESC,pc.id DESC"
		}
	}
	var rows []row
	if err := q.Select(columns).Order(orderSQL).Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		failed(c)
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}
func (h *Handler) Get(c *gin.Context) {
	commentID := id(c)
	if commentID == 0 {
		bad(c, "评论 ID 错误")
		return
	}
	var item row
	err := h.base(c).Select(columns).Where("pc.id=?", commentID).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "评论不存在")
		return
	}
	if err != nil {
		failed(c)
		return
	}
	response.OK(c, item)
}
func (h *Handler) Moderate(c *gin.Context) {
	commentID := id(c)
	var in moderateInput
	if commentID == 0 || c.ShouldBindJSON(&in) != nil {
		bad(c, "评论审核参数错误")
		return
	}
	h.dispatch(c, commentcommand.Command{CommentID: commentID, Action: strings.TrimSpace(in.Action), OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey), Note: strings.TrimSpace(in.Note)}, "评论审核参数错误")
}
func (h *Handler) CreateVirtual(c *gin.Context) {
	var in virtualInput
	if c.ShouldBindJSON(&in) != nil {
		bad(c, "虚拟评论参数错误")
		return
	}
	media, ok := h.media(c, in.AttachmentIDs)
	if !ok {
		return
	}
	h.dispatch(c, commentcommand.Command{Action: "create_virtual", OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey), ProductID: in.ProductID, Score: in.Score, Content: strings.TrimSpace(in.Content), VirtualAuthorName: strings.TrimSpace(in.VirtualAuthorName), VirtualAuthorAvatar: strings.TrimSpace(in.VirtualAuthorAvatar), Sort: in.Sort, Media: media}, "虚拟评论参数错误")
}
func (h *Handler) UpdateVirtual(c *gin.Context) {
	commentID := id(c)
	var in virtualInput
	if commentID == 0 || c.ShouldBindJSON(&in) != nil {
		bad(c, "虚拟评论参数错误")
		return
	}
	media, ok := h.media(c, in.AttachmentIDs)
	if !ok {
		return
	}
	h.dispatch(c, commentcommand.Command{CommentID: commentID, Action: "update_virtual", OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey), Score: in.Score, Content: strings.TrimSpace(in.Content), VirtualAuthorName: strings.TrimSpace(in.VirtualAuthorName), VirtualAuthorAvatar: strings.TrimSpace(in.VirtualAuthorAvatar), Sort: in.Sort, Media: media, MediaSet: in.AttachmentIDs != nil}, "虚拟评论参数错误")
}
func (h *Handler) SortVirtual(c *gin.Context) {
	commentID := id(c)
	var in sortInput
	if commentID == 0 || c.ShouldBindJSON(&in) != nil {
		bad(c, "评论排序参数错误")
		return
	}
	h.dispatch(c, commentcommand.Command{CommentID: commentID, Action: "sort_virtual", OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey), Sort: in.Sort}, "评论排序参数错误")
}
func (h *Handler) DeleteVirtual(c *gin.Context) {
	commentID := id(c)
	var in deleteInput
	if commentID == 0 || c.ShouldBindJSON(&in) != nil {
		bad(c, "删除虚拟评论参数错误")
		return
	}
	h.dispatch(c, commentcommand.Command{CommentID: commentID, Action: "delete_virtual", OperatorID: uint64(middleware.AdminID(c)), IdempotencyKey: strings.TrimSpace(in.IdempotencyKey), Note: strings.TrimSpace(in.Note)}, "删除虚拟评论参数错误")
}
func (h *Handler) dispatch(c *gin.Context, command commentcommand.Command, invalidMessage string) {
	if !commentcommand.Valid(command) {
		bad(c, invalidMessage)
		return
	}
	out, err := h.commands.Dispatch(c.Request.Context(), command)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "评论命令服务不可用")
		return
	}
	switch out.Code {
	case "":
		response.OK(c, gin.H{"comment_id": out.CommentID, "status": out.Status})
	case "not_found":
		response.Fail(c, http.StatusNotFound, "评论不存在")
	case "conflict":
		response.Fail(c, http.StatusConflict, "评论状态已变化、类型不匹配或幂等键冲突")
	case "schema":
		response.Fail(c, http.StatusServiceUnavailable, "业务库评论表缺少字段，请先执行 sql/business/patch_product_comment_virtual_avatar.sql")
	case "invalid":
		bad(c, invalidMessage)
	default:
		response.Fail(c, http.StatusServiceUnavailable, "评论命令处理失败")
	}
}
func (h *Handler) media(c *gin.Context, ids []uint64) ([]string, bool) {
	if len(ids) > 9 {
		bad(c, "评论图片最多 9 张")
		return nil, false
	}
	if len(ids) == 0 {
		return []string{}, true
	}
	seen := map[uint64]bool{}
	for _, id := range ids {
		if id == 0 || seen[id] {
			bad(c, "评论图片素材错误")
			return nil, false
		}
		seen[id] = true
	}
	var rows []struct {
		ID  uint64 `gorm:"column:attachment_id"`
		Src string `gorm:"column:attachment_src"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_attachment_asset").Select("attachment_id,attachment_src").Where("user_type=0 AND attachment_type=0 AND attachment_id IN ?", ids).Scan(&rows).Error; err != nil || len(rows) != len(ids) {
		bad(c, "评论图片素材不存在或不可用")
		return nil, false
	}
	src := map[uint64]string{}
	for _, row := range rows {
		src[row.ID] = row.Src
	}
	media := make([]string, 0, len(ids))
	for _, id := range ids {
		media = append(media, src[id])
	}
	return media, true
}
func (h *Handler) base(c *gin.Context) *gorm.DB {
	return h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_product_comment AS pc").
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id=pc.product_id").
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id=pc.user_id").
		Where("pc.deleted_at IS NULL")
}
func validStatus(value string) bool {
	return value == "pending" || value == "published" || value == "hidden"
}
func id(c *gin.Context) uint64 { value, _ := strconv.ParseUint(c.Param("id"), 10, 64); return value }
func uintParam(raw string) (uint64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	return value, err == nil
}
func paging(c *gin.Context) (int, int) {
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
func bad(c *gin.Context, message string) { response.Fail(c, http.StatusBadRequest, message) }
func failed(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "评论监管查询失败")
}
