// Package wechatreply 公众号自动回复（关注欢迎 / 关键字 / 默认）。
package wechatreply

import (
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
	menuRead   = "app.wechat_reply.read"
	menuManage = "app.wechat_reply.manage"
	tableName  = "qixi_crm_a_wechat_reply"

	keySubscribe = "subscribe"
	keyDefault   = "default"
)

var reservedKeys = map[string]struct{}{
	keySubscribe: {},
	keyDefault:   {},
}

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

	r.GET("/wechat/replies", access, read, h.List)
	r.GET("/wechat/replies/special/:key", access, read, h.GetSpecial)
	r.PUT("/wechat/replies/special/:key", access, write, h.SaveSpecial)
	r.GET("/wechat/replies/match", access, read, h.Match)
	r.GET("/wechat/replies/:id", access, read, h.Detail)
	r.POST("/wechat/replies", access, write, h.Create)
	r.PUT("/wechat/replies/:id", access, write, h.Update)
	r.PUT("/wechat/replies/:id/status", access, write, h.SetStatus)
	r.DELETE("/wechat/replies/:id", access, write, h.Delete)
}

type replyRow struct {
	WechatReplyID uint64    `json:"wechat_reply_id" gorm:"column:wechat_reply_id"`
	ReplyKey      string    `json:"key" gorm:"column:reply_key"`
	ReplyType     string    `json:"type" gorm:"column:reply_type"`
	Content       string    `json:"content" gorm:"column:content"`
	Status        int8      `json:"status" gorm:"column:status"`
	Sort          int       `json:"sort" gorm:"column:sort"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time"`
}

func (replyRow) TableName() string { return tableName }

type saveReq struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Status  *int8  `json:"status"`
	Sort    *int   `json:"sort"`
}

type specialSaveReq struct {
	Content string `json:"content"`
	Status  *int8  `json:"status"`
	Type    string `json:"type"`
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
	keyword := strings.TrimSpace(c.Query("keyword"))
	kind := strings.TrimSpace(c.DefaultQuery("kind", "keyword"))

	q := h.adminDB.WithContext(c.Request.Context()).Table(tableName)
	switch kind {
	case "all":
		// no filter
	case "special":
		q = q.Where("reply_key IN ?", []string{keySubscribe, keyDefault})
	default:
		q = q.Where("reply_key NOT IN ?", []string{keySubscribe, keyDefault})
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("reply_key LIKE ? OR content LIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	var rows []replyRow
	if err := q.Order("sort ASC, wechat_reply_id DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "count": total})
}

func (h *Handler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var row replyRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "数据不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) GetSpecial(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	if _, ok := reservedKeys[key]; !ok {
		response.Fail(c, http.StatusBadRequest, "仅支持 subscribe / default")
		return
	}
	row, err := h.getByKey(c, key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.OK(c, replyRow{
				ReplyKey:  key,
				ReplyType: "text",
				Content:   "",
				Status:    0,
			})
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) SaveSpecial(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	if _, ok := reservedKeys[key]; !ok {
		response.Fail(c, http.StatusBadRequest, "仅支持 subscribe / default")
		return
	}
	var req specialSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		response.Fail(c, http.StatusBadRequest, "请填写回复内容")
		return
	}
	if utf8.RuneCountInString(content) > 2000 {
		response.Fail(c, http.StatusBadRequest, "回复内容不能超过 2000 字")
		return
	}
	replyType := strings.TrimSpace(req.Type)
	if replyType == "" {
		replyType = "text"
	}
	if replyType != "text" {
		response.Fail(c, http.StatusBadRequest, "当前仅支持文本回复")
		return
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	existing, err := h.getByKey(c, key)
	now := time.Now()
	if err == gorm.ErrRecordNotFound {
		row := replyRow{
			ReplyKey:   key,
			ReplyType:  replyType,
			Content:    content,
			Status:     status,
			Sort:       0,
			CreateTime: now,
			UpdateTime: now,
		}
		if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
			Select("reply_key", "reply_type", "content", "status", "sort", "create_time", "update_time").
			Create(&row).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "保存失败")
			return
		}
		response.OK(c, row)
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	updates := map[string]any{
		"reply_type":  replyType,
		"content":     content,
		"status":      status,
		"update_time": now,
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", existing.WechatReplyID).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存失败")
		return
	}
	existing.ReplyType = replyType
	existing.Content = content
	existing.Status = status
	existing.UpdateTime = now
	response.OK(c, existing)
}

func (h *Handler) Create(c *gin.Context) {
	var req saveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	key, replyType, content, status, sort, errMsg := normalizeKeyword(req)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	if _, ok := reservedKeys[key]; ok {
		response.Fail(c, http.StatusBadRequest, "关键字不能使用保留词 subscribe / default")
		return
	}
	var count int64
	_ = h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("reply_key = ?", key).Count(&count).Error
	if count > 0 {
		response.Fail(c, http.StatusBadRequest, "关键字已存在")
		return
	}
	now := time.Now()
	row := replyRow{
		ReplyKey:   key,
		ReplyType:  replyType,
		Content:    content,
		Status:     status,
		Sort:       sort,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Select("reply_key", "reply_type", "content", "status", "sort", "create_time", "update_time").
		Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var existing replyRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).Take(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "数据不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	if _, ok := reservedKeys[existing.ReplyKey]; ok {
		response.Fail(c, http.StatusBadRequest, "请使用关注/默认回复接口编辑")
		return
	}

	var req saveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	key, replyType, content, status, sort, errMsg := normalizeKeyword(req)
	if errMsg != "" {
		response.Fail(c, http.StatusBadRequest, errMsg)
		return
	}
	if _, ok := reservedKeys[key]; ok {
		response.Fail(c, http.StatusBadRequest, "关键字不能使用保留词 subscribe / default")
		return
	}
	var dup int64
	_ = h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("reply_key = ? AND wechat_reply_id <> ?", key, id).Count(&dup).Error
	if dup > 0 {
		response.Fail(c, http.StatusBadRequest, "关键字已存在")
		return
	}
	updates := map[string]any{
		"reply_key":   key,
		"reply_type":  replyType,
		"content":     content,
		"status":      status,
		"sort":        sort,
		"update_time": time.Now(),
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	status := int8(0)
	if req.Status == 1 {
		status = 1
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).
		Updates(map[string]any{"status": status, "update_time": time.Now()})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
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
	var existing replyRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).Take(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "数据不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	if _, ok := reservedKeys[existing.ReplyKey]; ok {
		response.Fail(c, http.StatusBadRequest, "关注/默认回复不可删除，可关闭启用状态")
		return
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("wechat_reply_id = ?", id).Delete(&replyRow{})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// Match 供本地预览 / 后续公众号回调复用：优先精确关键字，否则回落 default。
func (h *Handler) Match(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		response.Fail(c, http.StatusBadRequest, "请输入关键字")
		return
	}
	row, err := h.getEnabledByKey(c, key)
	if err == gorm.ErrRecordNotFound && key != keyDefault {
		row, err = h.getEnabledByKey(c, keyDefault)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.OK(c, gin.H{"matched": false, "reply": nil})
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"matched": true, "reply": row})
}

func (h *Handler) getByKey(c *gin.Context, key string) (replyRow, error) {
	var row replyRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("reply_key = ?", key).Take(&row).Error
	return row, err
}

func (h *Handler) getEnabledByKey(c *gin.Context, key string) (replyRow, error) {
	var row replyRow
	err := h.adminDB.WithContext(c.Request.Context()).Table(tableName).
		Where("reply_key = ? AND status = 1", key).Take(&row).Error
	return row, err
}

func normalizeKeyword(req saveReq) (key, replyType, content string, status int8, sort int, errMsg string) {
	key = strings.TrimSpace(req.Key)
	content = strings.TrimSpace(req.Content)
	replyType = strings.TrimSpace(req.Type)
	if replyType == "" {
		replyType = "text"
	}
	status = 1
	if req.Status != nil {
		status = *req.Status
	}
	if req.Sort != nil {
		sort = *req.Sort
	}
	if key == "" {
		return "", "", "", 0, 0, "请填写关键字"
	}
	if utf8.RuneCountInString(key) > 64 {
		return "", "", "", 0, 0, "关键字不能超过 64 个字"
	}
	if replyType != "text" {
		return "", "", "", 0, 0, "当前仅支持文本回复"
	}
	if content == "" {
		return "", "", "", 0, 0, "请填写回复内容"
	}
	if utf8.RuneCountInString(content) > 2000 {
		return "", "", "", 0, 0, "回复内容不能超过 2000 字"
	}
	return key, replyType, content, status, sort, ""
}
