// Package nativemarketingdecor manages platform marketing decoration configs
// (atmosphere / border / topic / application) in qixi_crm_a_marketing_decor.
// It replaces setting_cache list stubs and never stores secrets.
package nativemarketingdecor

import (
	"encoding/json"
	"errors"
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

type Handler struct {
	adminDB *gorm.DB
}

func NewHandler(adminDB *gorm.DB) *Handler { return &Handler{adminDB: adminDB} }

func (h *Handler) Register(r gin.IRoutes) {
	platformOrOps := middleware.RequireAdminRoles("platform", "operations")
	for _, kind := range []decorKind{
		{path: "atmosphere", typ: "atmosphere", read: "marketing.atmosphere.read", manage: "marketing.atmosphere.manage"},
		{path: "border", typ: "border", read: "marketing.border.read", manage: "marketing.border.manage"},
		{path: "topic", typ: "topic", read: "marketing.topic.read", manage: "marketing.topic.manage"},
		{path: "applications", typ: "application", read: "marketing.application.read", manage: "marketing.application.manage"},
	} {
		k := kind
		read := middleware.RequireAdminMenu(h.adminDB, k.read)
		manage := middleware.RequireAdminMenu(h.adminDB, k.manage)
		r.GET("/marketing/"+k.path, platformOrOps, read, h.list(k.typ))
		r.GET("/marketing/"+k.path+"/:id", platformOrOps, read, h.detail(k.typ))
		r.POST("/marketing/"+k.path, platformOrOps, manage, h.create(k.typ))
		r.PUT("/marketing/"+k.path+"/:id", platformOrOps, manage, h.update(k.typ))
		r.PUT("/marketing/"+k.path+"/:id/status", platformOrOps, manage, h.setStatus(k.typ))
		r.DELETE("/marketing/"+k.path+"/:id", platformOrOps, manage, h.remove(k.typ))
	}
}

type decorKind struct {
	path, typ, read, manage string
}

type decorRow struct {
	ID        uint64     `gorm:"column:id"`
	DecorType string     `gorm:"column:decor_type"`
	Name      string     `gorm:"column:name"`
	Code      string     `gorm:"column:code"`
	CoverURL  string     `gorm:"column:cover_url"`
	Remark    string     `gorm:"column:remark"`
	Payload   string     `gorm:"column:payload"`
	Status    int        `gorm:"column:status"`
	Sort      int        `gorm:"column:sort"`
	StartsAt  *time.Time `gorm:"column:starts_at"`
	EndsAt    *time.Time `gorm:"column:ends_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

type upsertInput struct {
	Name     string         `json:"name"`
	Code     string         `json:"code"`
	CoverURL string         `json:"cover_url"`
	Remark   string         `json:"remark"`
	Payload  map[string]any `json:"payload"`
	Status   *int           `json:"status"`
	Sort     *int           `json:"sort"`
	StartsAt string         `json:"starts_at"`
	EndsAt   string         `json:"ends_at"`
}

func (h *Handler) list(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pageLimit(c)
		q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_marketing_decor").
			Where("decor_type = ? AND is_del = 0", decorType)
		if status := strings.TrimSpace(c.Query("status")); status != "" {
			switch status {
			case "1", "enabled":
				q = q.Where("status = 1")
			case "0", "disabled":
				q = q.Where("status = 0")
			default:
				response.Fail(c, http.StatusBadRequest, "状态错误")
				return
			}
		}
		if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
			q = q.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
		}
		var total int64
		if err := q.Count(&total).Error; err != nil {
			fail(c)
			return
		}
		rows := make([]decorRow, 0)
		if err := q.Select("id,decor_type,name,code,cover_url,remark,payload,status,sort,starts_at,ends_at,updated_at").
			Order("sort ASC, id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
			fail(c)
			return
		}
		list := make([]gin.H, 0, len(rows))
		for _, row := range rows {
			list = append(list, toItem(row))
		}
		response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
	}
}

func (h *Handler) detail(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, decorType, c.Param("id"))
		if !ok {
			return
		}
		response.OK(c, toItem(*row))
	}
}

func (h *Handler) create(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in upsertInput
		if err := c.ShouldBindJSON(&in); err != nil {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		name, code, cover, remark, payload, status, sort, starts, ends, err := normalize(in, 1, 0)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		created := struct {
			ID        uint64     `gorm:"column:id;primaryKey"`
			DecorType string     `gorm:"column:decor_type"`
			Name      string     `gorm:"column:name"`
			Code      string     `gorm:"column:code"`
			CoverURL  string     `gorm:"column:cover_url"`
			Remark    string     `gorm:"column:remark"`
			Payload   string     `gorm:"column:payload"`
			Status    int        `gorm:"column:status"`
			Sort      int        `gorm:"column:sort"`
			StartsAt  *time.Time `gorm:"column:starts_at"`
			EndsAt    *time.Time `gorm:"column:ends_at"`
			IsDel     int        `gorm:"column:is_del"`
		}{
			DecorType: decorType, Name: name, Code: code, CoverURL: cover, Remark: remark,
			Payload: payload, Status: status, Sort: sort, StartsAt: starts, EndsAt: ends, IsDel: 0,
		}
		if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_marketing_decor").Create(&created).Error; err != nil {
			fail(c)
			return
		}
		row, ok := h.load(c, decorType, strconv.FormatUint(created.ID, 10))
		if !ok {
			return
		}
		response.OK(c, toItem(*row))
	}
}

func (h *Handler) update(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, decorType, c.Param("id"))
		if !ok {
			return
		}
		var in upsertInput
		if err := c.ShouldBindJSON(&in); err != nil {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		name, code, cover, remark, payload, status, sort, starts, ends, err := normalize(in, row.Status, row.Sort)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		res := h.adminDB.WithContext(c.Request.Context()).Exec(`
			UPDATE qixi_crm_a_marketing_decor
			SET name=?, code=?, cover_url=?, remark=?, payload=?, status=?, sort=?, starts_at=?, ends_at=?
			WHERE id=? AND decor_type=? AND is_del=0`,
			name, code, cover, remark, payload, status, sort, starts, ends, row.ID, decorType,
		)
		if res.Error != nil {
			fail(c)
			return
		}
		if res.RowsAffected == 0 {
			response.Fail(c, http.StatusNotFound, "记录不存在")
			return
		}
		updated, ok := h.load(c, decorType, strconv.FormatUint(row.ID, 10))
		if !ok {
			return
		}
		response.OK(c, toItem(*updated))
	}
}

func (h *Handler) setStatus(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, decorType, c.Param("id"))
		if !ok {
			return
		}
		var in struct {
			Status *int `json:"status"`
		}
		if c.ShouldBindJSON(&in) != nil || in.Status == nil || (*in.Status != 0 && *in.Status != 1) {
			response.Fail(c, http.StatusBadRequest, "状态错误")
			return
		}
		res := h.adminDB.WithContext(c.Request.Context()).Exec(
			`UPDATE qixi_crm_a_marketing_decor SET status=? WHERE id=? AND decor_type=? AND is_del=0`,
			*in.Status, row.ID, decorType,
		)
		if res.Error != nil {
			fail(c)
			return
		}
		if res.RowsAffected == 0 {
			response.Fail(c, http.StatusNotFound, "记录不存在")
			return
		}
		response.OK(c, gin.H{"id": row.ID, "status": *in.Status})
	}
}

func (h *Handler) remove(decorType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, decorType, c.Param("id"))
		if !ok {
			return
		}
		res := h.adminDB.WithContext(c.Request.Context()).Exec(
			`UPDATE qixi_crm_a_marketing_decor SET is_del=1, status=0 WHERE id=? AND decor_type=? AND is_del=0`,
			row.ID, decorType,
		)
		if res.Error != nil {
			fail(c)
			return
		}
		if res.RowsAffected == 0 {
			response.Fail(c, http.StatusNotFound, "记录不存在")
			return
		}
		response.OK(c, gin.H{"ok": true})
	}
}

func (h *Handler) load(c *gin.Context, decorType, rawID string) (*decorRow, bool) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "ID 错误")
		return nil, false
	}
	var row decorRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_marketing_decor").
		Select("id,decor_type,name,code,cover_url,remark,payload,status,sort,starts_at,ends_at,updated_at").
		Where("id = ? AND decor_type = ? AND is_del = 0", id, decorType).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "记录不存在")
		return nil, false
	}
	if err != nil {
		fail(c)
		return nil, false
	}
	return &row, true
}

func normalize(in upsertInput, defaultStatus, defaultSort int) (string, string, string, string, string, int, int, *time.Time, *time.Time, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("名称必填且不超过 128 字")
	}
	code := strings.TrimSpace(in.Code)
	if utf8.RuneCountInString(code) > 64 {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("标识不超过 64 字")
	}
	cover := strings.TrimSpace(in.CoverURL)
	if utf8.RuneCountInString(cover) > 1024 {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("封面地址过长")
	}
	remark := strings.TrimSpace(in.Remark)
	if utf8.RuneCountInString(remark) > 500 {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("备注不超过 500 字")
	}
	payloadMap := in.Payload
	if payloadMap == nil {
		payloadMap = map[string]any{}
	}
	raw, err := json.Marshal(payloadMap)
	if err != nil {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("扩展字段编码失败")
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", "", "", "", "", 0, 0, nil, nil, errors.New("状态错误")
		}
		status = *in.Status
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	starts, err := parseOptionalTime(in.StartsAt)
	if err != nil {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("开始时间格式错误")
	}
	ends, err := parseOptionalTime(in.EndsAt)
	if err != nil {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("结束时间格式错误")
	}
	if starts != nil && ends != nil && ends.Before(*starts) {
		return "", "", "", "", "", 0, 0, nil, nil, errors.New("结束时间不能早于开始时间")
	}
	return name, code, cover, remark, string(raw), status, sort, starts, ends, nil
}

func toItem(row decorRow) gin.H {
	payload := map[string]any{}
	_ = json.Unmarshal([]byte(row.Payload), &payload)
	return gin.H{
		"id":         row.ID,
		"decor_type": row.DecorType,
		"name":       row.Name,
		"code":       row.Code,
		"cover_url":  row.CoverURL,
		"remark":     row.Remark,
		"payload":    payload,
		"status":     row.Status,
		"sort":       row.Sort,
		"starts_at":  formatTime(row.StartsAt),
		"ends_at":    formatTime(row.EndsAt),
		"updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func parseOptionalTime(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	for _, layout := range []string{time.DateTime, "2006-01-02 15:04", time.RFC3339, "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return &t, nil
		}
	}
	return nil, errors.New("bad time")
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func pageLimit(c *gin.Context) (int, int) {
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

func fail(c *gin.Context) {
	response.Fail(c, http.StatusInternalServerError, "营销装饰配置操作失败")
}
