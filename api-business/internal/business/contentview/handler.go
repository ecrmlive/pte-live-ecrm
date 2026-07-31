// Package contentview serves the C-end content projection. It deliberately
// reads qixi_crm_b_content_view only; the public PC, H5 and uni-app clients
// must not read platform-owned legacy content tables.
package contentview

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/notices", h.ListNotices)
	r.GET("/notices/:id", h.GetNotice)
	r.GET("/agreements/:key", h.GetAgreement)
}

type contentView struct {
	ContentID   uint64     `gorm:"column:content_id"`
	ContentType string     `gorm:"column:content_type"`
	Title       string     `gorm:"column:title"`
	CoverURL    string     `gorm:"column:cover_url"`
	Body        string     `gorm:"column:body"`
	Status      int8       `gorm:"column:status"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (contentView) TableName() string { return "qixi_crm_b_content_view" }

func (h *Handler) ListNotices(c *gin.Context) {
	page, limit := pageParams(c)
	query := h.db.WithContext(c.Request.Context()).Model(&contentView{}).
		Where("content_type = ? AND status = ?", "notice", 1)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		internalError(c)
		return
	}
	rows := make([]contentView, 0)
	if err := query.Order("published_at DESC, content_id DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		internalError(c)
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, noticeResponse(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) GetNotice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "公告 ID 错误")
		return
	}
	var row contentView
	err = h.db.WithContext(c.Request.Context()).
		Where("content_id = ? AND content_type = ? AND status = ?", id, "notice", 1).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "公告不存在")
		return
	}
	if err != nil {
		internalError(c)
		return
	}
	response.OK(c, noticeResponse(row))
}

func (h *Handler) GetAgreement(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	label, ok := agreementLabel(key)
	if !ok {
		response.Fail(c, http.StatusNotFound, "协议不存在")
		return
	}
	var row contentView
	err := h.db.WithContext(c.Request.Context()).
		Where("content_type = ? AND title = ? AND status = ?", "agreement", key, 1).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Keep the historical C-end contract: an enabled agreement key without
		// a published projection is readable as an empty document, not an error.
		response.OK(c, gin.H{"key": key, "label": label, "content": ""})
		return
	}
	if err != nil {
		internalError(c)
		return
	}
	response.OK(c, gin.H{"key": key, "label": label, "content": row.Body})
}

func noticeResponse(row contentView) gin.H {
	createdAt := row.UpdatedAt
	if row.PublishedAt != nil {
		createdAt = *row.PublishedAt
	}
	return gin.H{
		"notice_id":   row.ContentID,
		"title":       row.Title,
		"cover_url":   row.CoverURL,
		"content":     row.Body,
		"create_time": createdAt.Format(time.DateTime),
	}
}

func pageParams(c *gin.Context) (int, int) {
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

func agreementLabel(key string) (string, bool) {
	labels := map[string]string{
		"sys_user_agree":            "用户协议",
		"sys_userr_privacy":         "隐私政策",
		"sys_svip":                  "付费会员协议",
		"sys_product_presell_agree": "预售协议",
		"business_entry_agree":      "商户入驻协议",
		"promoter_explain":          "分销说明",
		"sys_about_us":              "关于我们",
		"sys_refund_agree":          "退款协议",
		"sys_cancel_agree":          "取消订单说明",
		"sys_recharge_agree":        "充值协议",
		"sys_integral_agree":        "积分规则",
		"mer_settle_agree":          "商户结算说明",
		"sys_lottery_agree":         "抽奖活动说明",
		"sys_deposit_agree":         "保证金说明",
	}
	label, ok := labels[key]
	return label, ok
}

func internalError(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "内容消费视图查询失败")
}
