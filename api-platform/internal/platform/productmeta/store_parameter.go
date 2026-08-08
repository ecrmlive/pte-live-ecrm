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

// storeParameterTemplate is merchant-DB scoped (CRMEB eb_parameter_template, mer_id>0).
type storeParameterTemplate struct {
	ID           uint64          `gorm:"column:id;primaryKey" json:"id"`
	StoreID      uint64          `gorm:"column:store_id" json:"store_id"`
	TemplateName string          `gorm:"column:template_name" json:"template_name"`
	CateID       uint64          `gorm:"column:cate_id" json:"cate_id"`
	IsRequired   int8            `gorm:"column:is_required" json:"is_required"`
	ParamsJSON   json.RawMessage `gorm:"column:params_json" json:"params_json"`
	Sort         int             `gorm:"column:sort" json:"sort"`
	IsDel        int8            `gorm:"column:is_del" json:"-"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (storeParameterTemplate) TableName() string {
	return "qixi_crm_m_product_parameter_template"
}

type storeParameterItem struct {
	Name     string   `json:"name"`
	Values   []string `json:"values"`
	Required int8     `json:"required"`
	Sort     int      `json:"sort"`
}

type storeParameterListRow struct {
	ID           uint64    `json:"id"`
	StoreID      uint64    `json:"store_id"`
	MerID        uint64    `json:"mer_id"`
	MerName      string    `json:"mer_name"`
	TemplateName string    `json:"template_name"`
	Sort         int       `json:"sort"`
	CreatedAt    time.Time `json:"created_at"`
}

type storeParameterDetail struct {
	ID           uint64               `json:"id"`
	StoreID      uint64               `json:"store_id"`
	MerID        uint64               `json:"mer_id"`
	MerName      string               `json:"mer_name"`
	TemplateName string               `json:"template_name"`
	CateID       uint64               `json:"cate_id"`
	IsRequired   int8                 `json:"is_required"`
	Params       []storeParameterItem `json:"params"`
	Sort         int                  `json:"sort"`
	CreatedAt    time.Time            `json:"created_at"`
}

type storeParameterCopyInput struct {
	TemplateName string               `json:"template_name"`
	CateIDs      []uint64             `json:"cate_ids"`
	Params       []storeParameterItem `json:"params"`
	Sort         *int                 `json:"sort"`
	Status       *int8                `json:"status"`
}

type storeParameterCreateInput struct {
	MerID        uint64               `json:"mer_id"`
	StoreID      uint64               `json:"store_id"`
	TemplateName string               `json:"template_name"`
	CateID       uint64               `json:"cate_id"`
	IsRequired   *int8                `json:"is_required"`
	Params       []storeParameterItem `json:"params"`
	Sort         int                  `json:"sort"`
}

func (h *Handler) registerStoreParameters(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.parameter.store.manage")
	r.GET("/product/store-parameter-templates", platform, manage, h.listStoreParameters)
	r.POST("/product/store-parameter-templates", platform, manage, h.createStoreParameter)
	r.GET("/product/store-parameter-templates/:id", platform, manage, h.detailStoreParameter)
	r.POST("/product/store-parameter-templates/:id/copy", platform, manage, h.copyStoreParameter)
}

func (h *Handler) listStoreParameters(c *gin.Context) {
	if h.merchantDB == nil {
		internal(c, "商户库未配置")
		return
	}
	page, limit := pageLimit(c)
	q := h.merchantDB.WithContext(c.Request.Context()).
		Table("qixi_crm_m_product_parameter_template AS t").
		Joins("JOIN qixi_crm_m_store AS s ON s.id = t.store_id").
		Where("t.is_del = 0")

	if merID := queryUint(c, "mer_id"); merID > 0 {
		q = q.Where("s.merchant_id = ?", merID)
	}
	if storeID := queryUint(c, "store_id"); storeID > 0 {
		q = q.Where("t.store_id = ?", storeID)
	}
	if name := strings.TrimSpace(c.Query("template_name")); name != "" {
		q = q.Where("t.template_name LIKE ?", "%"+name+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		internal(c, "查询店铺商品参数失败")
		return
	}

	type scanRow struct {
		ID           uint64    `gorm:"column:id"`
		StoreID      uint64    `gorm:"column:store_id"`
		MerID        uint64    `gorm:"column:mer_id"`
		MerName      string    `gorm:"column:mer_name"`
		TemplateName string    `gorm:"column:template_name"`
		Sort         int       `gorm:"column:sort"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}
	var rows []scanRow
	err := q.Select(`t.id, t.store_id, s.merchant_id AS mer_id, s.name AS mer_name,
		t.template_name, t.sort, t.created_at`).
		Order("t.sort DESC, t.created_at DESC, t.id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error
	if err != nil {
		internal(c, "查询店铺商品参数失败")
		return
	}

	list := make([]storeParameterListRow, 0, len(rows))
	for _, row := range rows {
		list = append(list, storeParameterListRow{
			ID:           row.ID,
			StoreID:      row.StoreID,
			MerID:        row.MerID,
			MerName:      row.MerName,
			TemplateName: row.TemplateName,
			Sort:         row.Sort,
			CreatedAt:    row.CreatedAt,
		})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detailStoreParameter(c *gin.Context) {
	if h.merchantDB == nil {
		internal(c, "商户库未配置")
		return
	}
	templateID := id(c)
	if templateID == 0 {
		bad(c, "ID 错误")
		return
	}
	detail, err := h.loadStoreParameterDetail(c, templateID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			notFound(c, "参数模板不存在")
			return
		}
		internal(c, "查询参数模板详情失败")
		return
	}
	response.OK(c, detail)
}

func (h *Handler) createStoreParameter(c *gin.Context) {
	if h.merchantDB == nil {
		internal(c, "商户库未配置")
		return
	}
	var in storeParameterCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, "参数错误")
		return
	}
	name := strings.TrimSpace(in.TemplateName)
	if !validName(name, 64) {
		bad(c, "请填写参数模板名称")
		return
	}
	if in.CateID == 0 {
		bad(c, "请选择平台分类")
		return
	}
	if !h.validParameterCateIDs(c, []uint64{in.CateID}) {
		return
	}
	normalized, ok := normalizeStoreParams(in.Params)
	if !ok {
		bad(c, "参数项格式错误")
		return
	}
	if in.Sort < 0 {
		bad(c, "排序参数错误")
		return
	}
	isRequired := int8(0)
	if in.IsRequired != nil && *in.IsRequired == 1 {
		isRequired = 1
	}
	storeID, err := h.resolveStoreIDForParameter(c, in.StoreID, in.MerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			bad(c, "请选择有效店铺")
			return
		}
		internal(c, "查询店铺失败")
		return
	}
	paramsJSON, err := json.Marshal(normalized)
	if err != nil {
		bad(c, "参数项格式错误")
		return
	}
	row := storeParameterTemplate{
		StoreID:      storeID,
		TemplateName: name,
		CateID:       in.CateID,
		IsRequired:   isRequired,
		ParamsJSON:   paramsJSON,
		Sort:         in.Sort,
		IsDel:        0,
	}
	if err := h.merchantDB.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		internal(c, "新增店铺商品参数失败")
		return
	}
	detail, err := h.loadStoreParameterDetail(c, row.ID)
	if err != nil {
		response.OK(c, gin.H{"id": row.ID, "ok": true})
		return
	}
	response.OK(c, detail)
}

func (h *Handler) resolveStoreIDForParameter(c *gin.Context, storeID, merID uint64) (uint64, error) {
	if storeID > 0 {
		var found uint64
		err := h.merchantDB.WithContext(c.Request.Context()).
			Table("qixi_crm_m_store").
			Select("id").
			Where("id = ? AND status = 1", storeID).
			Limit(1).
			Scan(&found).Error
		if err != nil {
			return 0, err
		}
		if found == 0 {
			return 0, gorm.ErrRecordNotFound
		}
		return found, nil
	}
	if merID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	var found uint64
	err := h.merchantDB.WithContext(c.Request.Context()).
		Table("qixi_crm_m_store").
		Select("id").
		Where("merchant_id = ? AND status = 1", merID).
		Order("id ASC").
		Limit(1).
		Scan(&found).Error
	if err != nil {
		return 0, err
	}
	if found == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return found, nil
}

func (h *Handler) copyStoreParameter(c *gin.Context) {
	if h.merchantDB == nil {
		internal(c, "商户库未配置")
		return
	}
	templateID := id(c)
	if templateID == 0 {
		bad(c, "ID 错误")
		return
	}
	src, err := h.loadStoreParameterDetail(c, templateID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			notFound(c, "参数模板不存在")
			return
		}
		internal(c, "查询参数模板失败")
		return
	}

	var in storeParameterCopyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		bad(c, "复制参数错误")
		return
	}
	name := strings.TrimSpace(in.TemplateName)
	if name == "" {
		name = strings.TrimSpace(src.TemplateName)
	}
	if !validName(name, 64) {
		bad(c, "请填写参数模板名称")
		return
	}
	params := in.Params
	if len(params) == 0 {
		params = src.Params
	}
	normalized, ok := normalizeStoreParams(params)
	if !ok {
		bad(c, "参数项格式错误")
		return
	}
	cateIDs := uniqueUint64(in.CateIDs)
	if src.CateID > 0 && len(cateIDs) == 0 {
		cateIDs = []uint64{src.CateID}
	}
	if !h.validParameterCateIDs(c, cateIDs) {
		return
	}
	status := int8(1)
	if in.Status != nil {
		if !validStatus(*in.Status) {
			bad(c, "状态参数错误")
			return
		}
		status = *in.Status
	}
	sortVal := src.Sort
	if in.Sort != nil {
		sortVal = *in.Sort
	}

	platformParams := make([]parameterItem, 0, len(normalized))
	for _, item := range normalized {
		platformParams = append(platformParams, parameterItem{
			Name:     item.Name,
			Values:   item.Values,
			Required: item.Required,
			Sort:     item.Sort,
		})
	}
	row, err := buildParameterRow(parameterInput{
		Name:    name,
		CateIDs: cateIDs,
		Params:  platformParams,
		Sort:    sortVal,
		Status:  status,
	})
	if err != nil {
		bad(c, "复制参数错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "平台参数模板名称已存在", "复制到平台参数模板失败")
		return
	}
	response.OK(c, gin.H{
		"ok":                   true,
		"platform_template_id": row.ID,
		"name":                 row.Name,
	})
}

func (h *Handler) loadStoreParameterDetail(c *gin.Context, templateID uint64) (*storeParameterDetail, error) {
	type scanRow struct {
		ID           uint64          `gorm:"column:id"`
		StoreID      uint64          `gorm:"column:store_id"`
		MerID        uint64          `gorm:"column:mer_id"`
		MerName      string          `gorm:"column:mer_name"`
		TemplateName string          `gorm:"column:template_name"`
		CateID       uint64          `gorm:"column:cate_id"`
		IsRequired   int8            `gorm:"column:is_required"`
		ParamsJSON   json.RawMessage `gorm:"column:params_json"`
		Sort         int             `gorm:"column:sort"`
		CreatedAt    time.Time       `gorm:"column:created_at"`
	}
	var row scanRow
	err := h.merchantDB.WithContext(c.Request.Context()).
		Table("qixi_crm_m_product_parameter_template AS t").
		Joins("JOIN qixi_crm_m_store AS s ON s.id = t.store_id").
		Select(`t.id, t.store_id, s.merchant_id AS mer_id, s.name AS mer_name,
			t.template_name, t.cate_id, t.is_required, t.params_json, t.sort, t.created_at`).
		Where("t.id = ? AND t.is_del = 0", templateID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	params, err := parseStoreParams(row.ParamsJSON)
	if err != nil {
		params = []storeParameterItem{}
	}
	return &storeParameterDetail{
		ID:           row.ID,
		StoreID:      row.StoreID,
		MerID:        row.MerID,
		MerName:      row.MerName,
		TemplateName: row.TemplateName,
		CateID:       row.CateID,
		IsRequired:   row.IsRequired,
		Params:       params,
		Sort:         row.Sort,
		CreatedAt:    row.CreatedAt,
	}, nil
}

func parseStoreParams(raw json.RawMessage) ([]storeParameterItem, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []storeParameterItem{}, nil
	}
	var items []storeParameterItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func normalizeStoreParams(items []storeParameterItem) ([]storeParameterItem, bool) {
	if len(items) == 0 || len(items) > 50 {
		return nil, false
	}
	out := make([]storeParameterItem, 0, len(items))
	seenName := map[string]struct{}{}
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" || len([]rune(name)) > 32 {
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
			if v == "" || len([]rune(v)) > 64 {
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
		out = append(out, storeParameterItem{
			Name:     name,
			Values:   values,
			Required: req,
			Sort:     item.Sort,
		})
	}
	return out, true
}

func queryUint(c *gin.Context, key string) uint64 {
	v, _ := strconv.ParseUint(strings.TrimSpace(c.Query(key)), 10, 64)
	return v
}
