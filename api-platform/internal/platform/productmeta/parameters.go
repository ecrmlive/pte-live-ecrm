package productmeta

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type parameterTemplate struct {
	ID        uint64          `gorm:"column:id;primaryKey" json:"id"`
	Name      string          `gorm:"column:name" json:"name"`
	CateIDs   json.RawMessage `gorm:"column:cate_ids_json" json:"-"`
	Params    json.RawMessage `gorm:"column:params_json" json:"-"`
	Values    json.RawMessage `gorm:"column:values_json" json:"-"`
	Sort      int             `gorm:"column:sort" json:"sort"`
	Status    int8            `gorm:"column:status" json:"status"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (parameterTemplate) TableName() string { return "qixi_crm_a_product_parameter_template" }

type parameterItem struct {
	Name     string   `json:"name"`
	Values   []string `json:"values"`
	Required int8     `json:"required"`
	Sort     int      `json:"sort"`
}

type parameterInput struct {
	Name    string          `json:"name"`
	CateIDs []uint64        `json:"cate_ids"`
	Params  []parameterItem `json:"params"`
	Sort    int             `json:"sort"`
	Status  int8            `json:"status"`
}

type parameterListItem struct {
	ID            uint64          `json:"id"`
	Name          string          `json:"name"`
	CateIDs       []uint64        `json:"cate_ids"`
	CateNames     []string        `json:"cate_names"`
	CateNamesText string          `json:"cate_names_text"`
	Params        []parameterItem `json:"params"`
	IsRequired    int8            `json:"is_required"`
	Sort          int             `json:"sort"`
	Status        int8            `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (h *Handler) registerParameters(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.parameter.manage")
	r.GET("/product/parameter-templates", platform, manage, h.listParameters)
	r.GET("/product/parameter-templates/:id", platform, manage, h.getParameter)
	r.POST("/product/parameter-templates", platform, manage, h.createParameter)
	r.PUT("/product/parameter-templates/:id", platform, manage, h.updateParameter)
	r.DELETE("/product/parameter-templates/:id", platform, manage, h.deleteParameter)
}

func (h *Handler) listParameters(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.db.WithContext(c.Request.Context()).Model(&parameterTemplate{})
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		name = strings.TrimSpace(c.Query("template_name"))
	}
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if cateID, _ := strconv.ParseUint(strings.TrimSpace(c.Query("cate_id")), 10, 64); cateID > 0 {
		raw, _ := json.Marshal(cateID)
		q = q.Where("JSON_CONTAINS(IFNULL(cate_ids_json, JSON_ARRAY()), CAST(? AS JSON), '$')", string(raw))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		internal(c, "查询商品参数失败")
		return
	}
	var rows []parameterTemplate
	if err := q.Order("sort DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		internal(c, "查询商品参数失败")
		return
	}
	items, err := h.toParameterListItems(c, rows)
	if err != nil {
		internal(c, "查询商品参数失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) getParameter(c *gin.Context) {
	templateID := id(c)
	if templateID == 0 {
		bad(c, "商品参数错误")
		return
	}
	var row parameterTemplate
	if err := h.db.WithContext(c.Request.Context()).First(&row, templateID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			notFound(c, "商品参数不存在")
			return
		}
		internal(c, "查询商品参数失败")
		return
	}
	items, err := h.toParameterListItems(c, []parameterTemplate{row})
	if err != nil || len(items) == 0 {
		internal(c, "查询商品参数失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) createParameter(c *gin.Context) {
	var in parameterInput
	if c.ShouldBindJSON(&in) != nil || !validParameter(&in) {
		bad(c, "商品参数错误")
		return
	}
	if !h.validParameterCateIDs(c, in.CateIDs) {
		return
	}
	row, err := buildParameterRow(in)
	if err != nil {
		bad(c, "商品参数错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "商品参数名称已存在", "创建商品参数失败")
		return
	}
	items, _ := h.toParameterListItems(c, []parameterTemplate{row})
	if len(items) == 1 {
		response.OK(c, items[0])
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateParameter(c *gin.Context) {
	templateID := id(c)
	var in parameterInput
	if templateID == 0 || c.ShouldBindJSON(&in) != nil || !validParameter(&in) {
		bad(c, "商品参数错误")
		return
	}
	if !h.validParameterCateIDs(c, in.CateIDs) {
		return
	}
	row, err := buildParameterRow(in)
	if err != nil {
		bad(c, "商品参数错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Model(&parameterTemplate{}).Where("id=?", templateID).Updates(map[string]any{
		"name":          row.Name,
		"cate_ids_json": row.CateIDs,
		"params_json":   row.Params,
		"values_json":   row.Values,
		"sort":          row.Sort,
		"status":        row.Status,
	})
	if res.Error != nil {
		conflictOrInternal(c, res.Error, "商品参数名称已存在", "更新商品参数失败")
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, "商品参数不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) deleteParameter(c *gin.Context) {
	h.delete(c, &parameterTemplate{}, "商品参数不存在", "删除商品参数失败")
}

func (h *Handler) validParameterCateIDs(c *gin.Context, cateIDs []uint64) bool {
	if len(cateIDs) == 0 {
		bad(c, "请选择关联分类")
		return false
	}
	seen := map[uint64]struct{}{}
	ids := make([]uint64, 0, len(cateIDs))
	for _, cateID := range cateIDs {
		if cateID == 0 {
			bad(c, "关联分类无效")
			return false
		}
		if _, ok := seen[cateID]; ok {
			continue
		}
		seen[cateID] = struct{}{}
		ids = append(ids, cateID)
	}
	var total int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_platform_category").Where("id IN ?", ids).Count(&total).Error; err != nil {
		internal(c, "校验关联分类失败")
		return false
	}
	if total != int64(len(ids)) {
		bad(c, "关联分类不存在")
		return false
	}
	return true
}

func (h *Handler) toParameterListItems(c *gin.Context, rows []parameterTemplate) ([]parameterListItem, error) {
	allIDs := map[uint64]struct{}{}
	parsed := make([]struct {
		cateIDs []uint64
		params  []parameterItem
	}, len(rows))
	for i, row := range rows {
		cateIDs, params := parseParameterStored(row)
		parsed[i].cateIDs = cateIDs
		parsed[i].params = params
		for _, cateID := range cateIDs {
			allIDs[cateID] = struct{}{}
		}
	}
	nameByID := map[uint64]string{}
	if len(allIDs) > 0 {
		ids := make([]uint64, 0, len(allIDs))
		for cateID := range allIDs {
			ids = append(ids, cateID)
		}
		var cats []struct {
			ID   uint64 `gorm:"column:id"`
			Name string `gorm:"column:name"`
		}
		if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_platform_category").
			Select("id,name").Where("id IN ?", ids).Find(&cats).Error; err != nil {
			return nil, err
		}
		for _, cat := range cats {
			nameByID[cat.ID] = cat.Name
		}
	}
	out := make([]parameterListItem, 0, len(rows))
	for i, row := range rows {
		names := make([]string, 0, len(parsed[i].cateIDs))
		for _, cateID := range parsed[i].cateIDs {
			if name := nameByID[cateID]; name != "" {
				names = append(names, name)
			}
		}
		out = append(out, parameterListItem{
			ID:            row.ID,
			Name:          row.Name,
			CateIDs:       parsed[i].cateIDs,
			CateNames:     names,
			CateNamesText: strings.Join(names, ","),
			Params:        parsed[i].params,
			IsRequired:    parameterIsRequired(parsed[i].params),
			Sort:          row.Sort,
			Status:        row.Status,
			CreatedAt:     row.CreatedAt,
		})
	}
	return out, nil
}

func buildParameterRow(in parameterInput) (parameterTemplate, error) {
	cateIDs := uniqueUint64(in.CateIDs)
	params := normalizeParameterItems(in.Params)
	cateRaw, err := json.Marshal(cateIDs)
	if err != nil {
		return parameterTemplate{}, err
	}
	paramsRaw, err := json.Marshal(params)
	if err != nil {
		return parameterTemplate{}, err
	}
	snapshot := []string{}
	if len(params) > 0 {
		snapshot = params[0].Values
	}
	valuesRaw, err := json.Marshal(snapshot)
	if err != nil {
		return parameterTemplate{}, err
	}
	status := in.Status
	if status != 0 && status != 1 {
		status = 1
	}
	return parameterTemplate{
		Name:    strings.TrimSpace(in.Name),
		CateIDs: cateRaw,
		Params:  paramsRaw,
		Values:  valuesRaw,
		Sort:    in.Sort,
		Status:  status,
	}, nil
}

func parseParameterStored(row parameterTemplate) ([]uint64, []parameterItem) {
	cateIDs := []uint64{}
	_ = json.Unmarshal(row.CateIDs, &cateIDs)
	params := []parameterItem{}
	if len(row.Params) > 0 && string(row.Params) != "null" {
		_ = json.Unmarshal(row.Params, &params)
	}
	if len(params) == 0 && len(row.Values) > 0 && string(row.Values) != "null" {
		var values []string
		if json.Unmarshal(row.Values, &values) == nil && len(values) > 0 {
			params = []parameterItem{{Name: row.Name, Values: values, Required: 0, Sort: 0}}
		}
	}
	return cateIDs, normalizeParameterItems(params)
}

func parameterIsRequired(params []parameterItem) int8 {
	for _, item := range params {
		if item.Required == 1 {
			return 1
		}
	}
	return 0
}

func validParameter(in *parameterInput) bool {
	if !validName(in.Name, 64) || len(in.CateIDs) == 0 || len(in.CateIDs) > 100 {
		return false
	}
	if in.Status != 0 && in.Status != 1 {
		return false
	}
	if len(in.Params) == 0 || len(in.Params) > 50 {
		return false
	}
	seenName := map[string]struct{}{}
	for _, item := range in.Params {
		name := strings.TrimSpace(item.Name)
		if name == "" || len([]rune(name)) > 32 {
			return false
		}
		if _, ok := seenName[name]; ok {
			return false
		}
		seenName[name] = struct{}{}
		if item.Required != 0 && item.Required != 1 {
			return false
		}
		if len(item.Values) == 0 || len(item.Values) > 50 {
			return false
		}
		seenVal := map[string]struct{}{}
		for _, value := range item.Values {
			value = strings.TrimSpace(value)
			if value == "" || len([]rune(value)) > 64 {
				return false
			}
			if _, ok := seenVal[value]; ok {
				return false
			}
			seenVal[value] = struct{}{}
		}
	}
	return true
}

func normalizeParameterItems(items []parameterItem) []parameterItem {
	out := make([]parameterItem, 0, len(items))
	for _, item := range items {
		values := make([]string, 0, len(item.Values))
		for _, value := range item.Values {
			values = append(values, strings.TrimSpace(value))
		}
		required := item.Required
		if required != 1 {
			required = 0
		}
		out = append(out, parameterItem{
			Name:     strings.TrimSpace(item.Name),
			Values:   values,
			Required: required,
			Sort:     item.Sort,
		})
	}
	return out
}

func uniqueUint64(ids []uint64) []uint64 {
	seen := map[uint64]struct{}{}
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
