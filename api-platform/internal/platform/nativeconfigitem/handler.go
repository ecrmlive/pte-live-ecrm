// Package nativeconfigitem manages platform config list items
// (hot_search / group_data / system_form / backup) in qixi_crm_a_config_item.
package nativeconfigitem

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
	platformOnly := middleware.RequireAdminRoles("platform")
	platformOrOps := middleware.RequireAdminRoles("platform", "operations")
	for _, kind := range []itemKind{
		{path: "/maintain/hot-search", typ: "hot_search", read: "maintain.hot_search", manage: "maintain.hot_search.manage", roles: platformOnly},
		{path: "/maintain/group-data", typ: "group_data", read: "maintain.group_data", manage: "maintain.group_data.manage", roles: platformOnly},
		{path: "/maintain/backups", typ: "backup", read: "maintain.backup", manage: "maintain.backup.manage", roles: platformOnly},
		{path: "/diy/system-forms", typ: "system_form", read: "operations.system_form", manage: "operations.system_form.manage", roles: platformOrOps},
	} {
		k := kind
		read := middleware.RequireAdminMenu(h.adminDB, k.read)
		manage := middleware.RequireAdminMenu(h.adminDB, k.manage)
		r.GET(k.path, k.roles, read, h.list(k.typ))
		r.GET(k.path+"/:id", k.roles, read, h.detail(k.typ))
		r.POST(k.path, k.roles, manage, h.create(k.typ))
		r.PUT(k.path+"/:id", k.roles, manage, h.update(k.typ))
		r.PUT(k.path+"/:id/status", k.roles, manage, h.setStatus(k.typ))
		r.DELETE(k.path+"/:id", k.roles, manage, h.remove(k.typ))
	}
}

type itemKind struct {
	path, typ, read, manage string
	roles                   gin.HandlerFunc
}

type itemRow struct {
	ID        uint64    `gorm:"column:id"`
	ItemType  string    `gorm:"column:item_type"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	Remark    string    `gorm:"column:remark"`
	Payload   string    `gorm:"column:payload"`
	Status    int       `gorm:"column:status"`
	Sort      int       `gorm:"column:sort"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type upsertInput struct {
	Name    string         `json:"name"`
	Code    string         `json:"code"`
	Remark  string         `json:"remark"`
	Payload map[string]any `json:"payload"`
	Status  *int           `json:"status"`
	Sort    *int           `json:"sort"`
}

func (h *Handler) list(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pageLimit(c)
		q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_item").
			Where("item_type = ? AND is_del = 0", itemType)
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
		rows := make([]itemRow, 0)
		if err := q.Select("id,item_type,name,code,remark,payload,status,sort,updated_at").
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

func (h *Handler) detail(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, itemType, c.Param("id"))
		if !ok {
			return
		}
		response.OK(c, toItem(*row))
	}
}

func (h *Handler) create(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in upsertInput
		if err := c.ShouldBindJSON(&in); err != nil {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		name, code, remark, payload, status, sort, err := normalize(in, 1, 0)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		created := struct {
			ID       uint64 `gorm:"column:id;primaryKey"`
			ItemType string `gorm:"column:item_type"`
			Name     string `gorm:"column:name"`
			Code     string `gorm:"column:code"`
			Remark   string `gorm:"column:remark"`
			Payload  string `gorm:"column:payload"`
			Status   int    `gorm:"column:status"`
			Sort     int    `gorm:"column:sort"`
			IsDel    int    `gorm:"column:is_del"`
		}{ItemType: itemType, Name: name, Code: code, Remark: remark, Payload: payload, Status: status, Sort: sort}
		if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_item").Create(&created).Error; err != nil {
			fail(c)
			return
		}
		row, ok := h.load(c, itemType, strconv.FormatUint(created.ID, 10))
		if !ok {
			return
		}
		response.OK(c, toItem(*row))
	}
}

func (h *Handler) update(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, itemType, c.Param("id"))
		if !ok {
			return
		}
		var in upsertInput
		if err := c.ShouldBindJSON(&in); err != nil {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		name, code, remark, payload, status, sort, err := normalize(in, row.Status, row.Sort)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		res := h.adminDB.WithContext(c.Request.Context()).Exec(`
			UPDATE qixi_crm_a_config_item
			SET name=?, code=?, remark=?, payload=?, status=?, sort=?
			WHERE id=? AND item_type=? AND is_del=0`,
			name, code, remark, payload, status, sort, row.ID, itemType,
		)
		if res.Error != nil {
			fail(c)
			return
		}
		if res.RowsAffected == 0 {
			response.Fail(c, http.StatusNotFound, "记录不存在")
			return
		}
		updated, ok := h.load(c, itemType, strconv.FormatUint(row.ID, 10))
		if !ok {
			return
		}
		response.OK(c, toItem(*updated))
	}
}

func (h *Handler) setStatus(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, itemType, c.Param("id"))
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
			`UPDATE qixi_crm_a_config_item SET status=? WHERE id=? AND item_type=? AND is_del=0`,
			*in.Status, row.ID, itemType,
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

func (h *Handler) remove(itemType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, ok := h.load(c, itemType, c.Param("id"))
		if !ok {
			return
		}
		res := h.adminDB.WithContext(c.Request.Context()).Exec(
			`UPDATE qixi_crm_a_config_item SET is_del=1, status=0 WHERE id=? AND item_type=? AND is_del=0`,
			row.ID, itemType,
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

func (h *Handler) load(c *gin.Context, itemType, rawID string) (*itemRow, bool) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "ID 错误")
		return nil, false
	}
	var row itemRow
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config_item").
		Select("id,item_type,name,code,remark,payload,status,sort,updated_at").
		Where("id = ? AND item_type = ? AND is_del = 0", id, itemType).Take(&row).Error
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

func normalize(in upsertInput, defaultStatus, defaultSort int) (string, string, string, string, int, int, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", "", 0, 0, errors.New("名称必填且不超过 128 字")
	}
	code := strings.TrimSpace(in.Code)
	if utf8.RuneCountInString(code) > 64 {
		return "", "", "", "", 0, 0, errors.New("标识不超过 64 字")
	}
	remark := strings.TrimSpace(in.Remark)
	if utf8.RuneCountInString(remark) > 500 {
		return "", "", "", "", 0, 0, errors.New("备注不超过 500 字")
	}
	payloadMap := in.Payload
	if payloadMap == nil {
		payloadMap = map[string]any{}
	}
	raw, err := json.Marshal(payloadMap)
	if err != nil {
		return "", "", "", "", 0, 0, errors.New("扩展字段编码失败")
	}
	status := defaultStatus
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return "", "", "", "", 0, 0, errors.New("状态错误")
		}
		status = *in.Status
	}
	sort := defaultSort
	if in.Sort != nil {
		sort = *in.Sort
	}
	return name, code, remark, string(raw), status, sort, nil
}

func toItem(row itemRow) gin.H {
	payload := map[string]any{}
	_ = json.Unmarshal([]byte(row.Payload), &payload)
	return gin.H{
		"id":         row.ID,
		"item_type":  row.ItemType,
		"name":       row.Name,
		"code":       row.Code,
		"remark":     row.Remark,
		"payload":    payload,
		"status":     row.Status,
		"enabled":    row.Status == 1,
		"sort":       row.Sort,
		"updated_at": row.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
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
	response.Fail(c, http.StatusInternalServerError, "配置条目操作失败")
}
