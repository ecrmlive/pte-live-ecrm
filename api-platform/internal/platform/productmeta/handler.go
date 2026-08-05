// Package productmeta owns platform-wide product labels, guarantees and
// parameter templates. Store-specific metadata remains in api-merchant.
package productmeta

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	h.registerLabels(r, platform)
	h.registerGuarantees(r, platform)
	h.registerParameters(r, platform)
}

type label struct {
	ID          uint64    `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Color       string    `gorm:"column:color" json:"color"`
	Sort        int       `gorm:"column:sort" json:"sort"`
	Status      int8      `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (label) TableName() string { return "qixi_crm_a_product_label" }

type guarantee struct {
	ID        uint64    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Content   string    `gorm:"column:content" json:"content"`
	IconURL   string    `gorm:"column:icon_url" json:"icon_url"`
	Sort      int       `gorm:"column:sort" json:"sort"`
	Status    int8      `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (guarantee) TableName() string { return "qixi_crm_a_product_guarantee" }

type parameterTemplate struct {
	ID        uint64          `gorm:"column:id;primaryKey" json:"id"`
	Name      string          `gorm:"column:name" json:"name"`
	Values    json.RawMessage `gorm:"column:values_json" json:"values_json"`
	Sort      int             `gorm:"column:sort" json:"sort"`
	Status    int8            `gorm:"column:status" json:"status"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (parameterTemplate) TableName() string { return "qixi_crm_a_product_parameter_template" }

type labelInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Sort        int    `json:"sort"`
	Status      int8   `json:"status"`
}

type guaranteeInput struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	IconURL string `json:"icon_url"`
	Sort    int    `json:"sort"`
	Status  int8   `json:"status"`
}

type parameterInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	Sort   int      `json:"sort"`
	Status int8     `json:"status"`
}

func (h *Handler) registerLabels(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.label.manage")
	r.GET("/product/labels", platform, manage, h.listLabels)
	r.POST("/product/labels", platform, manage, h.createLabel)
	r.PUT("/product/labels/:id", platform, manage, h.updateLabel)
	r.DELETE("/product/labels/:id", platform, manage, h.deleteLabel)
}

func (h *Handler) registerGuarantees(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.guarantee.manage")
	r.GET("/product/guarantees", platform, manage, h.listGuarantees)
	r.POST("/product/guarantees", platform, manage, h.createGuarantee)
	r.PUT("/product/guarantees/:id", platform, manage, h.updateGuarantee)
	r.DELETE("/product/guarantees/:id", platform, manage, h.deleteGuarantee)
}

func (h *Handler) registerParameters(r gin.IRoutes, platform gin.HandlerFunc) {
	manage := middleware.RequireAdminMenu(h.db, "product.parameter.manage")
	r.GET("/product/parameter-templates", platform, manage, h.listParameters)
	r.POST("/product/parameter-templates", platform, manage, h.createParameter)
	r.PUT("/product/parameter-templates/:id", platform, manage, h.updateParameter)
	r.DELETE("/product/parameter-templates/:id", platform, manage, h.deleteParameter)
}

func (h *Handler) listLabels(c *gin.Context) {
	var rows []label
	if err := h.db.WithContext(c.Request.Context()).Order("sort DESC,id DESC").Find(&rows).Error; err != nil {
		internal(c, "查询商品标签失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) createLabel(c *gin.Context) {
	var in labelInput
	if c.ShouldBindJSON(&in) != nil || !validLabel(in) {
		bad(c, "商品标签参数错误")
		return
	}
	row := label{Name: strings.TrimSpace(in.Name), Description: strings.TrimSpace(in.Description), Color: strings.TrimSpace(in.Color), Sort: in.Sort, Status: in.Status}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "商品标签名称已存在", "创建商品标签失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateLabel(c *gin.Context) {
	var in labelInput
	if id(c) == 0 || c.ShouldBindJSON(&in) != nil || !validLabel(in) {
		bad(c, "商品标签参数错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Model(&label{}).Where("id=?", id(c)).Updates(map[string]any{"name": strings.TrimSpace(in.Name), "description": strings.TrimSpace(in.Description), "color": strings.TrimSpace(in.Color), "sort": in.Sort, "status": in.Status})
	if res.Error != nil {
		conflictOrInternal(c, res.Error, "商品标签名称已存在", "更新商品标签失败")
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, "商品标签不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) deleteLabel(c *gin.Context) {
	h.delete(c, &label{}, "商品标签不存在", "删除商品标签失败")
}

func (h *Handler) listGuarantees(c *gin.Context) {
	var rows []guarantee
	if err := h.db.WithContext(c.Request.Context()).Order("sort DESC,id DESC").Find(&rows).Error; err != nil {
		internal(c, "查询保障服务失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) createGuarantee(c *gin.Context) {
	var in guaranteeInput
	if c.ShouldBindJSON(&in) != nil || !validGuarantee(in) {
		bad(c, "保障服务参数错误")
		return
	}
	if !h.validGuaranteeIcon(c, in.IconURL) {
		return
	}
	row := guarantee{Name: strings.TrimSpace(in.Name), Content: strings.TrimSpace(in.Content), IconURL: strings.TrimSpace(in.IconURL), Sort: in.Sort, Status: in.Status}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "保障服务名称已存在", "创建保障服务失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateGuarantee(c *gin.Context) {
	var in guaranteeInput
	if id(c) == 0 || c.ShouldBindJSON(&in) != nil || !validGuarantee(in) {
		bad(c, "保障服务参数错误")
		return
	}
	if !h.validGuaranteeIcon(c, in.IconURL) {
		return
	}
	res := h.db.WithContext(c.Request.Context()).Model(&guarantee{}).Where("id=?", id(c)).Updates(map[string]any{"name": strings.TrimSpace(in.Name), "content": strings.TrimSpace(in.Content), "icon_url": strings.TrimSpace(in.IconURL), "sort": in.Sort, "status": in.Status})
	if res.Error != nil {
		conflictOrInternal(c, res.Error, "保障服务名称已存在", "更新保障服务失败")
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, "保障服务不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// validGuaranteeIcon keeps guarantee icons in the controlled platform asset
// library. Empty is permitted for text-only guarantees; arbitrary external
// URLs are not accepted.
func (h *Handler) validGuaranteeIcon(c *gin.Context, raw string) bool {
	icon := strings.TrimSpace(raw)
	if icon == "" {
		return true
	}
	var total int64
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_attachment_asset").Where("user_type=0 AND attachment_type=0 AND attachment_src=?", icon).Count(&total).Error
	if err != nil {
		internal(c, "校验保障图标失败")
		return false
	}
	if total != 1 {
		bad(c, "保障图标必须选择平台图片素材")
		return false
	}
	return true
}

func (h *Handler) deleteGuarantee(c *gin.Context) {
	h.delete(c, &guarantee{}, "保障服务不存在", "删除保障服务失败")
}

func (h *Handler) listParameters(c *gin.Context) {
	var rows []parameterTemplate
	if err := h.db.WithContext(c.Request.Context()).Order("sort DESC,id DESC").Find(&rows).Error; err != nil {
		internal(c, "查询商品参数失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) createParameter(c *gin.Context) {
	var in parameterInput
	if c.ShouldBindJSON(&in) != nil || !validParameter(&in) {
		bad(c, "商品参数错误")
		return
	}
	values, _ := json.Marshal(normalizeParameterValues(in.Values))
	row := parameterTemplate{Name: strings.TrimSpace(in.Name), Values: values, Sort: in.Sort, Status: in.Status}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		conflictOrInternal(c, err, "商品参数名称已存在", "创建商品参数失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) updateParameter(c *gin.Context) {
	var in parameterInput
	if id(c) == 0 || c.ShouldBindJSON(&in) != nil || !validParameter(&in) {
		bad(c, "商品参数错误")
		return
	}
	values, _ := json.Marshal(normalizeParameterValues(in.Values))
	res := h.db.WithContext(c.Request.Context()).Model(&parameterTemplate{}).Where("id=?", id(c)).Updates(map[string]any{"name": strings.TrimSpace(in.Name), "values_json": values, "sort": in.Sort, "status": in.Status})
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

func (h *Handler) delete(c *gin.Context, model any, missing, failed string) {
	if id(c) == 0 {
		bad(c, "ID 错误")
		return
	}
	res := h.db.WithContext(c.Request.Context()).Delete(model, id(c))
	if res.Error != nil {
		internal(c, failed)
		return
	}
	if res.RowsAffected != 1 {
		notFound(c, missing)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func validLabel(in labelInput) bool {
	return validName(in.Name, 64) && len([]rune(strings.TrimSpace(in.Description))) <= 255 && len([]rune(strings.TrimSpace(in.Color))) <= 32 && validStatus(in.Status)
}
func validGuarantee(in guaranteeInput) bool {
	return validName(in.Name, 64) && len([]rune(strings.TrimSpace(in.Content))) <= 1000 && len([]rune(strings.TrimSpace(in.IconURL))) <= 1024 && validStatus(in.Status)
}
func validParameter(in *parameterInput) bool {
	if !validName(in.Name, 64) || !validStatus(in.Status) || len(in.Values) == 0 || len(in.Values) > 50 {
		return false
	}
	seen := map[string]struct{}{}
	for _, value := range in.Values {
		value = strings.TrimSpace(value)
		if value == "" || len([]rune(value)) > 64 {
			return false
		}
		if _, ok := seen[value]; ok {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
func normalizeParameterValues(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strings.TrimSpace(value))
	}
	return result
}
func validName(value string, limit int) bool {
	length := len([]rune(strings.TrimSpace(value)))
	return length > 0 && length <= limit
}
func validStatus(status int8) bool            { return status == 0 || status == 1 }
func id(c *gin.Context) uint64                { value, _ := strconv.ParseUint(c.Param("id"), 10, 64); return value }
func bad(c *gin.Context, message string)      { response.Fail(c, http.StatusBadRequest, message) }
func notFound(c *gin.Context, message string) { response.Fail(c, http.StatusNotFound, message) }
func internal(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
func conflictOrInternal(c *gin.Context, err error, conflict, failed string) {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		response.Fail(c, http.StatusConflict, conflict)
		return
	}
	internal(c, failed)
}
