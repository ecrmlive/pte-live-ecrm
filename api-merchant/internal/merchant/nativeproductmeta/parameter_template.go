package nativeproductmeta

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type parameterTemplate struct {
	ID           uint64          `gorm:"column:id;primaryKey"`
	StoreID      uint64          `gorm:"column:store_id"`
	TemplateName string          `gorm:"column:template_name"`
	CateID       uint64          `gorm:"column:cate_id"`
	IsRequired   int8            `gorm:"column:is_required"`
	ParamsJSON   json.RawMessage `gorm:"column:params_json"`
	Sort         int             `gorm:"column:sort"`
	IsDel        int8            `gorm:"column:is_del"`
	CreatedAt    time.Time       `gorm:"column:created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at"`
}

func (parameterTemplate) TableName() string {
	return "qixi_crm_m_product_parameter_template"
}

type parameterItem struct {
	Name     string   `json:"name"`
	Values   []string `json:"values"`
	Required int8     `json:"required"`
	Sort     int      `json:"sort"`
}

type parameterTemplateInput struct {
	TemplateName string          `json:"template_name"`
	CateID       uint64          `json:"cate_id"`
	IsRequired   *int8           `json:"is_required"`
	Params       []parameterItem `json:"params"`
	Sort         int             `json:"sort"`
}

func (h *Handler) listParameterTemplates(c *gin.Context) {
	page, limit := pagination(c)
	q := h.db.WithContext(c.Request.Context()).Model(&parameterTemplate{}).
		Where("store_id = ? AND is_del = 0", middleware.StoreID(c))
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("template_name LIKE ?", "%"+keyword+"%")
	}
	if name := strings.TrimSpace(c.Query("template_name")); name != "" {
		q = q.Where("template_name LIKE ?", "%"+name+"%")
	}
	if cateRaw := strings.TrimSpace(c.Query("cate_id")); cateRaw != "" {
		if cateID, err := strconv.ParseUint(cateRaw, 10, 64); err == nil && cateID > 0 {
			q = q.Where("cate_id = ?", cateID)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品参数模板失败")
		return
	}
	var rows []parameterTemplate
	if err := q.Order("sort DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品参数模板失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, parameterTemplateJSON(row))
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detailParameterTemplate(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "参数模板编号不正确")
		return
	}
	row, err := h.ownedParameterTemplate(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "参数模板不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询参数模板失败")
		return
	}
	response.OK(c, parameterTemplateJSON(row))
}

func (h *Handler) createParameterTemplate(c *gin.Context) {
	var req parameterTemplateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "商品参数模板参数不正确")
		return
	}
	name, cateID, isRequired, paramsJSON, sort, ok := normalizeParameterTemplate(req)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品参数模板参数不正确")
		return
	}
	row := parameterTemplate{
		StoreID:      uint64(middleware.StoreID(c)),
		TemplateName: name,
		CateID:       cateID,
		IsRequired:   isRequired,
		ParamsJSON:   paramsJSON,
		Sort:         sort,
		IsDel:        0,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "新增商品参数模板失败")
		return
	}
	response.OK(c, parameterTemplateJSON(row))
}

func (h *Handler) updateParameterTemplate(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "参数模板编号不正确")
		return
	}
	var req parameterTemplateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "商品参数模板参数不正确")
		return
	}
	name, cateID, isRequired, paramsJSON, sort, ok := normalizeParameterTemplate(req)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "商品参数模板参数不正确")
		return
	}
	row, err := h.ownedParameterTemplate(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "参数模板不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询参数模板失败")
		return
	}
	row.TemplateName = name
	row.CateID = cateID
	row.IsRequired = isRequired
	row.ParamsJSON = paramsJSON
	row.Sort = sort
	if err := h.db.WithContext(c.Request.Context()).Save(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品参数模板失败")
		return
	}
	response.OK(c, parameterTemplateJSON(row))
}

func (h *Handler) deleteParameterTemplate(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "参数模板编号不正确")
		return
	}
	row, err := h.ownedParameterTemplate(c, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "参数模板不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询参数模板失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Model(&row).Update("is_del", 1).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除商品参数模板失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ownedParameterTemplate(c *gin.Context, templateID uint64) (parameterTemplate, error) {
	var row parameterTemplate
	err := h.db.WithContext(c.Request.Context()).
		Where("id = ? AND store_id = ? AND is_del = 0", templateID, middleware.StoreID(c)).
		Take(&row).Error
	return row, err
}

func parameterTemplateJSON(row parameterTemplate) gin.H {
	params, _ := parseParameterItems(row.ParamsJSON)
	return gin.H{
		"template_id":   row.ID,
		"template_name": row.TemplateName,
		"cate_id":       row.CateID,
		"is_required":   row.IsRequired,
		"params":        params,
		"sort":          row.Sort,
		"create_time":   row.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func normalizeParameterTemplate(req parameterTemplateInput) (string, uint64, int8, json.RawMessage, int, bool) {
	name := strings.TrimSpace(req.TemplateName)
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return "", 0, 0, nil, 0, false
	}
	if req.CateID == 0 {
		return "", 0, 0, nil, 0, false
	}
	if req.Sort < 0 {
		return "", 0, 0, nil, 0, false
	}
	isRequired := int8(0)
	if req.IsRequired != nil && *req.IsRequired == 1 {
		isRequired = 1
	}
	params, ok := normalizeParameterItems(req.Params)
	if !ok {
		return "", 0, 0, nil, 0, false
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return "", 0, 0, nil, 0, false
	}
	return name, req.CateID, isRequired, raw, req.Sort, true
}

func parseParameterItems(raw json.RawMessage) ([]parameterItem, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []parameterItem{}, nil
	}
	var items []parameterItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func normalizeParameterItems(items []parameterItem) ([]parameterItem, bool) {
	if len(items) == 0 || len(items) > 50 {
		return nil, false
	}
	out := make([]parameterItem, 0, len(items))
	seenName := map[string]struct{}{}
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" || utf8.RuneCountInString(name) > 32 {
			return nil, false
		}
		key := strings.ToLower(name)
		if _, ok := seenName[key]; ok {
			return nil, false
		}
		seenName[key] = struct{}{}
		if len(item.Values) == 0 || len(item.Values) > 50 {
			return nil, false
		}
		values := make([]string, 0, len(item.Values))
		seenVal := map[string]struct{}{}
		for _, v := range item.Values {
			v = strings.TrimSpace(v)
			if v == "" || utf8.RuneCountInString(v) > 64 {
				return nil, false
			}
			if _, ok := seenVal[v]; ok {
				return nil, false
			}
			seenVal[v] = struct{}{}
			values = append(values, v)
		}
		req := item.Required
		if req != 0 && req != 1 {
			return nil, false
		}
		if item.Sort < 0 {
			return nil, false
		}
		out = append(out, parameterItem{
			Name:     name,
			Values:   values,
			Required: req,
			Sort:     item.Sort,
		})
	}
	return out, true
}
