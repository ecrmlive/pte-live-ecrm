// Package nativestoreconfig holds per-store JSON config and pickup addresses.
package nativestoreconfig

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	merchantDB *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(merchantDB, businessDB *gorm.DB) *Handler {
	return &Handler{merchantDB: merchantDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/product/units", h.listUnits)
	r.POST("/product/units", h.saveUnits)
	r.GET("/store-user-labels", h.listUserLabels)
	r.POST("/store-user-labels", h.saveUserLabels)
	r.GET("/store-auto-label-rules", h.listAutoLabels)
	r.POST("/store-auto-label-rules", h.saveAutoLabels)
	r.GET("/store-pickup-points", h.listPickupPoints)
	r.POST("/store-pickup-points", h.createPickupPoint)
	r.PUT("/store-pickup-points/:id", h.updatePickupPoint)
	r.DELETE("/store-pickup-points/:id", h.deletePickupPoint)
	r.GET("/user/search-records", h.listSearchRecords)
}

type unitItem struct {
	UnitID uint64 `json:"unit_id"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
}

type labelItem struct {
	LabelID uint64 `json:"label_id"`
	Name    string `json:"name"`
	Sort    int    `json:"sort"`
	Status  int8   `json:"status"`
}

type autoLabelRule struct {
	RuleID   uint64 `json:"rule_id"`
	Name     string `json:"name"`
	RuleType string `json:"rule_type"`
	Status   int8   `json:"status"`
}

type pickupPoint struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	StoreID     uint64 `gorm:"column:store_id"`
	AddressType string `gorm:"column:address_type"`
	ContactName string `gorm:"column:contact_name"`
	Mobile      string `gorm:"column:mobile"`
	RegionCode  string `gorm:"column:region_code"`
	Detail      string `gorm:"column:detail"`
	IsDefault   int8   `gorm:"column:is_default"`
}

func (pickupPoint) TableName() string { return "qixi_crm_m_store_address" }

func (h *Handler) listUnits(c *gin.Context) {
	items, err := readConfigSlice[unitItem](h, c, "product_units")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取商品单位失败")
		return
	}
	if len(items) == 0 {
		items = []unitItem{{UnitID: 1, Name: "件", Sort: 0}, {UnitID: 2, Name: "盒", Sort: 1}}
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) saveUnits(c *gin.Context) {
	var req struct {
		List []unitItem `json:"list"`
	}
	if c.ShouldBindJSON(&req) != nil || len(req.List) == 0 {
		response.Fail(c, http.StatusBadRequest, "商品单位列表不能为空")
		return
	}
	for i := range req.List {
		req.List[i].Name = strings.TrimSpace(req.List[i].Name)
		if req.List[i].Name == "" {
			response.Fail(c, http.StatusBadRequest, "商品单位名称不能为空")
			return
		}
		if req.List[i].UnitID == 0 {
			req.List[i].UnitID = uint64(i + 1)
		}
	}
	if err := h.writeConfig(c, "product_units", req.List); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存商品单位失败")
		return
	}
	response.OK(c, gin.H{"list": req.List})
}

func (h *Handler) listUserLabels(c *gin.Context) {
	items, err := readConfigSlice[labelItem](h, c, "store_user_labels")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取用户标签失败")
		return
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) saveUserLabels(c *gin.Context) {
	var req struct {
		List []labelItem `json:"list"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "用户标签参数不正确")
		return
	}
	for i := range req.List {
		req.List[i].Name = strings.TrimSpace(req.List[i].Name)
		if req.List[i].LabelID == 0 {
			req.List[i].LabelID = uint64(i + 1)
		}
	}
	if err := h.writeConfig(c, "store_user_labels", req.List); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存用户标签失败")
		return
	}
	response.OK(c, gin.H{"list": req.List})
}

func (h *Handler) listAutoLabels(c *gin.Context) {
	items, err := readConfigSlice[autoLabelRule](h, c, "store_auto_label_rules")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取自动标签规则失败")
		return
	}
	response.OK(c, gin.H{"list": items})
}

func (h *Handler) saveAutoLabels(c *gin.Context) {
	var req struct {
		List []autoLabelRule `json:"list"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "自动标签规则参数不正确")
		return
	}
	for i := range req.List {
		if req.List[i].RuleID == 0 {
			req.List[i].RuleID = uint64(i + 1)
		}
	}
	if err := h.writeConfig(c, "store_auto_label_rules", req.List); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存自动标签规则失败")
		return
	}
	response.OK(c, gin.H{"list": req.List})
}

func (h *Handler) listPickupPoints(c *gin.Context) {
	page, limit := pagination(c)
	q := h.merchantDB.WithContext(c.Request.Context()).Model(&pickupPoint{}).
		Where("store_id = ? AND address_type = 'pickup'", middleware.StoreID(c))
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(contact_name LIKE ? OR mobile LIKE ? OR detail LIKE ?)", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询自提点失败")
		return
	}
	var rows []pickupPoint
	if err := q.Order("is_default DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询自提点失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, pickupJSON(row))
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func (h *Handler) createPickupPoint(c *gin.Context) {
	row, err := h.bindPickup(c, 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	row.StoreID = uint64(middleware.StoreID(c))
	row.AddressType = "pickup"
	if err := h.merchantDB.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "新增自提点失败")
		return
	}
	response.OK(c, pickupJSON(row))
}

func (h *Handler) updatePickupPoint(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.bindPickup(c, id)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	res := h.merchantDB.WithContext(c.Request.Context()).Model(&pickupPoint{}).Where("id = ? AND store_id = ? AND address_type = 'pickup'", id, middleware.StoreID(c)).Updates(map[string]any{
		"contact_name": row.ContactName, "mobile": row.Mobile, "region_code": row.RegionCode, "detail": row.Detail, "is_default": row.IsDefault,
	})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新自提点失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "自提点不存在")
		return
	}
	row.ID = id
	response.OK(c, pickupJSON(row))
}

func (h *Handler) deletePickupPoint(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res := h.merchantDB.WithContext(c.Request.Context()).Where("id = ? AND store_id = ? AND address_type = 'pickup'", id, middleware.StoreID(c)).Delete(&pickupPoint{})
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除自提点失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "自提点不存在")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) listSearchRecords(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user_browse_history AS h").
		Select("h.id,h.user_id,h.product_id,h.viewed_at,COALESCE(p.title,'') AS product_title").
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id=h.product_id").
		Where("h.store_id = ?", middleware.StoreID(c))
	if userID, err := strconv.ParseUint(strings.TrimSpace(c.Query("user_id")), 10, 64); err == nil && userID > 0 {
		q = q.Where("h.user_id = ?", userID)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("p.title LIKE ?", "%"+keyword+"%")
	}
	if from := strings.TrimSpace(c.Query("date_from")); from != "" {
		if t, err := time.ParseInLocation("2006-01-02", from, time.Local); err == nil {
			q = q.Where("h.viewed_at >= ?", t)
		}
	}
	if to := strings.TrimSpace(c.Query("date_to")); to != "" {
		if t, err := time.ParseInLocation("2006-01-02", to, time.Local); err == nil {
			q = q.Where("h.viewed_at < ?", t.AddDate(0, 0, 1))
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询浏览记录失败")
		return
	}
	type browseRow struct {
		ID           uint64    `gorm:"column:id"`
		UserID       uint64    `gorm:"column:user_id"`
		ProductID    uint64    `gorm:"column:product_id"`
		ProductTitle string    `gorm:"column:product_title"`
		ViewedAt     time.Time `gorm:"column:viewed_at"`
	}
	var rows []browseRow
	if err := q.Order("h.viewed_at DESC,h.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询浏览记录失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id": row.ID, "user_id": row.UserID, "product_id": row.ProductID,
			"product_title": row.ProductTitle, "viewed_at": row.ViewedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "limit": limit})
}

func readConfigSlice[T any](h *Handler, c *gin.Context, key string) ([]T, error) {
	var raw json.RawMessage
	err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_config").
		Select("config_value").Where("store_id = ? AND config_key = ?", middleware.StoreID(c), key).Scan(&raw).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || len(raw) == 0 || string(raw) == "null" {
		return []T{}, nil
	}
	if err != nil {
		return nil, err
	}
	var items []T
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (h *Handler) writeConfig(c *gin.Context, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return h.merchantDB.WithContext(c.Request.Context()).Exec(
		"INSERT INTO qixi_crm_m_config (store_id, config_key, config_value) VALUES (?, ?, CAST(? AS JSON)) ON DUPLICATE KEY UPDATE config_value=VALUES(config_value)",
		middleware.StoreID(c), key, string(payload),
	).Error
}

func (h *Handler) bindPickup(c *gin.Context, id uint64) (pickupPoint, error) {
	var req struct {
		ContactName string `json:"contact_name"`
		Mobile      string `json:"mobile"`
		RegionCode  string `json:"region_code"`
		Detail      string `json:"detail"`
		IsDefault   *int8  `json:"is_default"`
	}
	if c.ShouldBindJSON(&req) != nil {
		return pickupPoint{}, errors.New("自提点参数不正确")
	}
	name := strings.TrimSpace(req.ContactName)
	mobile := strings.TrimSpace(req.Mobile)
	detail := strings.TrimSpace(req.Detail)
	if name == "" || mobile == "" || detail == "" {
		return pickupPoint{}, errors.New("自提点联系人、手机号和地址不能为空")
	}
	isDefault := int8(0)
	if req.IsDefault != nil && *req.IsDefault == 1 {
		isDefault = 1
	}
	return pickupPoint{ID: id, ContactName: name, Mobile: mobile, RegionCode: strings.TrimSpace(req.RegionCode), Detail: detail, IsDefault: isDefault}, nil
}

func pickupJSON(row pickupPoint) gin.H {
	return gin.H{"id": row.ID, "contact_name": row.ContactName, "mobile": row.Mobile, "region_code": row.RegionCode, "detail": row.Detail, "is_default": row.IsDefault}
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
