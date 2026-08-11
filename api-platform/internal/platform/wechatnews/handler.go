// Package wechatnews 公众号图文管理（对齐 CRMEB WechatNews）。
package wechatnews

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	menuRead   = "app.wechat_news.read"
	menuManage = "app.wechat_news.manage"
	tableName  = "qixi_crm_a_wechat_news"
	maxItems   = 8
)

type Handler struct {
	adminDB *gorm.DB
}

func NewHandler(adminDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenuAny(h.adminDB, menuRead, menuManage)
	write := middleware.RequireAdminMenu(h.adminDB, menuManage)
	r.GET("/wechat/news", access, read, h.List)
	r.GET("/wechat/news/:id", access, read, h.Detail)
	r.POST("/wechat/news", access, write, h.Create)
	r.PUT("/wechat/news/:id", access, write, h.Update)
	r.DELETE("/wechat/news/:id", access, write, h.Delete)
}

type newsItem struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Synopsis string `json:"synopsis"`
	Image    string `json:"image"`
	Content  string `json:"content"`
}

type newsRow struct {
	WechatNewsID uint64    `json:"wechat_news_id" gorm:"column:wechat_news_id"`
	Status       int8      `json:"status" gorm:"column:status"`
	ItemsRaw     string    `json:"-" gorm:"column:items"`
	Items        []newsItem `json:"article" gorm:"-"`
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime   time.Time `json:"update_time" gorm:"column:update_time"`
}

func (newsRow) TableName() string { return tableName }

type saveReq struct {
	Status *int8      `json:"status"`
	Data   []newsItem `json:"data"`
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	keyword := strings.TrimSpace(c.Query("cate_name"))
	q := h.adminDB.WithContext(c.Request.Context()).Table(tableName)
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("CAST(items AS CHAR) LIKE ?", like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	var rows []newsRow
	if err := q.Order("wechat_news_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	for i := range rows {
		rows[i].Items = decodeItems(rows[i].ItemsRaw)
	}
	response.OK(c, gin.H{"list": rows, "count": total})
}

func (h *Handler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var row newsRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_news_id = ?", id).Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "数据不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	row.Items = decodeItems(row.ItemsRaw)
	response.OK(c, row)
}

func (h *Handler) Create(c *gin.Context) {
	var req saveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	items, errMsg := normalizeItems(req.Data)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	raw, _ := json.Marshal(items)
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	row := newsRow{
		Status:     status,
		ItemsRaw:   string(raw),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Select("status", "items", "create_time", "update_time").
		Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存失败")
		return
	}
	response.OK(c, gin.H{"ok": true, "wechat_news_id": row.WechatNewsID})
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var req saveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	items, errMsg := normalizeItems(req.Data)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	raw, _ := json.Marshal(items)
	updates := map[string]any{"items": string(raw)}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_news_id = ?", id).Updates(updates)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "保存失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "数据不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_news_id = ?", id).Delete(&newsRow{})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "数据不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func decodeItems(raw string) []newsItem {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return []newsItem{}
	}
	var items []newsItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []newsItem{}
	}
	return items
}

func normalizeItems(in []newsItem) ([]newsItem, string) {
	if len(in) == 0 {
		return nil, "请至少添加一条图文"
	}
	if len(in) > maxItems {
		return nil, "单次最多添加 8 条图文"
	}
	out := make([]newsItem, 0, len(in))
	for i, item := range in {
		title := strings.TrimSpace(item.Title)
		author := strings.TrimSpace(item.Author)
		synopsis := strings.TrimSpace(item.Synopsis)
		image := strings.TrimSpace(item.Image)
		content := strings.TrimSpace(item.Content)
		if title == "" {
			return nil, "请填写第 " + strconv.Itoa(i+1) + " 条图文标题"
		}
		if utf8.RuneCountInString(title) > 64 {
			return nil, "标题不能超过 64 个字"
		}
		if author == "" {
			return nil, "请填写第 " + strconv.Itoa(i+1) + " 条图文作者"
		}
		if utf8.RuneCountInString(author) > 32 {
			return nil, "作者不能超过 32 个字"
		}
		if synopsis == "" {
			return nil, "请填写第 " + strconv.Itoa(i+1) + " 条图文摘要"
		}
		if utf8.RuneCountInString(synopsis) > 128 {
			return nil, "摘要不能超过 128 个字"
		}
		if image == "" {
			return nil, "请上传第 " + strconv.Itoa(i+1) + " 条图文封面"
		}
		if content == "" || content == "<p></p>" {
			return nil, "请填写第 " + strconv.Itoa(i+1) + " 条图文正文"
		}
		out = append(out, newsItem{
			Title:    title,
			Author:   author,
			Synopsis: synopsis,
			Image:    image,
			Content:  content,
		})
	}
	return out, ""
}
