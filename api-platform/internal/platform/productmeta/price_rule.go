package productmeta

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var htmlTagStripper = regexp.MustCompile(`(?s)<[^>]*>`)

type priceRule struct {
	ID        uint64          `gorm:"column:id;primaryKey" json:"id"`
	Name      string          `gorm:"column:name" json:"name"`
	CateIDs   json.RawMessage `gorm:"column:cate_ids_json" json:"-"`
	IsDefault int8            `gorm:"column:is_default" json:"is_default"`
	Content   string          `gorm:"column:content" json:"content"`
	Sort      int             `gorm:"column:sort" json:"sort"`
	Status    int8            `gorm:"column:status" json:"status"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (priceRule) TableName() string { return "qixi_crm_a_product_price_rule" }

type priceRuleInput struct {
	Name    string   `json:"name"`
	CateIDs []uint64 `json:"cate_ids"`
	Content string   `json:"content"`
	Sort    int      `json:"sort"`
	Status  int8     `json:"status"`
}

type priceRuleListItem struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	CateIDs       []uint64  `json:"cate_ids"`
	CateNames     []string  `json:"cate_names"`
	CateNamesText string    `json:"cate_names_text"`
	IsDefault     int8      `json:"is_default"`
	Content       string    `json:"content"`
	Sort          int       `json:"sort"`
	Status        int8      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (h *Handler) registerPriceRules(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.price_description.manage")
	r.GET("/product/price-rules", platform, manage, h.listPriceRules)
	r.GET("/product/price-rules/:id", platform, manage, h.getPriceRule)
	r.POST("/product/price-rules", platform, manage, h.createPriceRule)
	r.PUT("/product/price-rules/:id", platform, manage, h.updatePriceRule)
	r.PUT("/product/price-rules/:id/status", platform, manage, h.updatePriceRuleStatus)
	r.DELETE("/product/price-rules/:id", platform, manage, h.deletePriceRule)
}

func (h *Handler) listPriceRules(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.db.WithContext(c.Request.Context()).Model(&priceRule{})
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	statusRaw := strings.TrimSpace(c.Query("status"))
	if statusRaw == "" {
		statusRaw = strings.TrimSpace(c.Query("is_show"))
	}
	if statusRaw != "" {
		status, err := strconv.ParseInt(statusRaw, 10, 8)
		if err == nil && (status == 0 || status == 1) {
			q = q.Where("status = ?", status)
		}
	}
	if cateID, _ := strconv.ParseUint(strings.TrimSpace(c.Query("cate_id")), 10, 64); cateID > 0 {
		raw, _ := json.Marshal(cateID)
		q = q.Where(
			"(is_default = 1 OR JSON_CONTAINS(IFNULL(cate_ids_json, JSON_ARRAY()), CAST(? AS JSON), '$'))",
			string(raw),
		)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		internal(c, "查询价格说明失败")
		return
	}
	var rows []priceRule
	if err := q.Order("sort DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		internal(c, "查询价格说明失败")
		return
	}
	items, err := h.toPriceRuleListItems(c, rows)
	if err != nil {
		internal(c, "查询价格说明失败")
		return
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) getPriceRule(c *gin.Context) {
	ruleID := id(c)
	if ruleID == 0 {
		bad(c, "价格说明错误")
		return
	}
	var row priceRule
	if err := h.db.WithContext(c.Request.Context()).First(&row, ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			notFound(c, "价格说明不存在")
			return
		}
		internal(c, "查询价格说明失败")
		return
	}
	items, err := h.toPriceRuleListItems(c, []priceRule{row})
	if err != nil || len(items) == 0 {
		internal(c, "查询价格说明失败")
		return
	}
	response.OK(c, items[0])
}

func (h *Handler) createPriceRule(c *gin.Context) {
	var in priceRuleInput
	if c.ShouldBindJSON(&in) != nil || !validPriceRule(&in) {
		bad(c, "价格说明参数错误")
		return
	}
	if !h.validPriceRuleCateIDs(c, in.CateIDs) {
		return
	}
	row, err := buildPriceRuleRow(in)
	if err != nil {
		bad(c, "价格说明参数错误")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "价格说明名称已存在", "创建价格说明失败")
		return
	}
	items, _ := h.toPriceRuleListItems(c, []priceRule{row})
	if len(items) == 1 {
		response.OK(c, items[0])
		return
	}
	response.OK(c, row)
}

func (h *Handler) updatePriceRule(c *gin.Context) {
	ruleID := id(c)
	var in priceRuleInput
	if ruleID == 0 || c.ShouldBindJSON(&in) != nil || !validPriceRule(&in) {
		bad(c, "价格说明参数错误")
		return
	}
	if !h.validPriceRuleCateIDs(c, in.CateIDs) {
		return
	}
	row, err := buildPriceRuleRow(in)
	if err != nil {
		bad(c, "价格说明参数错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Model(&priceRule{}).Where("id=?", ruleID).Updates(map[string]any{
		"name":          row.Name,
		"cate_ids_json": row.CateIDs,
		"is_default":    row.IsDefault,
		"content":       row.Content,
		"sort":          row.Sort,
		"status":        row.Status,
	})
	if res.Error != nil {
		conflictOrInternal(c, res.Error, "价格说明名称已存在", "更新价格说明失败")
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, "价格说明不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) updatePriceRuleStatus(c *gin.Context) {
	ruleID := id(c)
	if ruleID == 0 {
		bad(c, "ID 错误")
		return
	}
	var in struct {
		Status int8  `json:"status"`
		IsShow *int8 `json:"is_show"`
	}
	if c.ShouldBindJSON(&in) != nil {
		bad(c, "价格说明状态参数错误")
		return
	}
	status := in.Status
	if in.IsShow != nil {
		status = *in.IsShow
	}
	if !validStatus(status) {
		bad(c, "价格说明状态参数错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Model(&priceRule{}).Where("id=?", ruleID).Update("status", status)
	if res.Error != nil {
		internal(c, "更新价格说明状态失败")
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, "价格说明不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) deletePriceRule(c *gin.Context) {
	h.delete(c, &priceRule{}, "价格说明不存在", "删除价格说明失败")
}

func (h *Handler) validPriceRuleCateIDs(c *gin.Context, cateIDs []uint64) bool {
	ids := uniqueUint64(cateIDs)
	if len(ids) == 0 {
		return true
	}
	if len(ids) > 200 {
		bad(c, "关联分类过多")
		return false
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

func (h *Handler) toPriceRuleListItems(c *gin.Context, rows []priceRule) ([]priceRuleListItem, error) {
	allIDs := map[uint64]struct{}{}
	parsed := make([][]uint64, len(rows))
	for i, row := range rows {
		cateIDs := parsePriceRuleCateIDs(row.CateIDs)
		parsed[i] = cateIDs
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
	out := make([]priceRuleListItem, 0, len(rows))
	for i, row := range rows {
		names := make([]string, 0, len(parsed[i]))
		for _, cateID := range parsed[i] {
			if name := nameByID[cateID]; name != "" {
				names = append(names, name)
			}
		}
		text := strings.Join(names, ",")
		if row.IsDefault == 1 && len(names) == 0 {
			text = "全部商品"
		}
		out = append(out, priceRuleListItem{
			ID:            row.ID,
			Name:          row.Name,
			CateIDs:       parsed[i],
			CateNames:     names,
			CateNamesText: text,
			IsDefault:     row.IsDefault,
			Content:       row.Content,
			Sort:          row.Sort,
			Status:        row.Status,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}
	return out, nil
}

func buildPriceRuleRow(in priceRuleInput) (priceRule, error) {
	cateIDs := uniqueUint64(in.CateIDs)
	cateRaw, err := json.Marshal(cateIDs)
	if err != nil {
		return priceRule{}, err
	}
	isDefault := int8(1)
	if len(cateIDs) > 0 {
		isDefault = 0
	}
	status := in.Status
	if status != 0 && status != 1 {
		status = 1
	}
	return priceRule{
		Name:      strings.TrimSpace(in.Name),
		CateIDs:   cateRaw,
		IsDefault: isDefault,
		Content:   strings.TrimSpace(in.Content),
		Sort:      in.Sort,
		Status:    status,
	}, nil
}

func parsePriceRuleCateIDs(raw json.RawMessage) []uint64 {
	if len(raw) == 0 || string(raw) == "null" {
		return []uint64{}
	}
	var cateIDs []uint64
	_ = json.Unmarshal(raw, &cateIDs)
	return uniqueUint64(cateIDs)
}

func validPriceRule(in *priceRuleInput) bool {
	if !validName(in.Name, 64) {
		return false
	}
	if in.Status != 0 && in.Status != 1 {
		return false
	}
	if in.Sort < 0 {
		return false
	}
	if len(in.CateIDs) > 200 {
		return false
	}
	plain := stripHTMLText(in.Content)
	if plain == "" || utf8.RuneCountInString(in.Content) > 200000 {
		return false
	}
	return true
}

func stripHTMLText(html string) string {
	text := htmlTagStripper.ReplaceAllString(html, "")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.TrimSpace(text)
	return text
}
