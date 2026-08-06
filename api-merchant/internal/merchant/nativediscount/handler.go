// Package nativediscount exposes store discount packages (优惠套餐) on
// qixi_crm_m_marketing_activity (activity_type=discount) and mirrors active
// rows into qixi_crm_b_marketing_activity_view for C-end / platform reads.
package nativediscount

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/listquery"
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
	manage := middleware.RequireStorePermission(h.merchantDB, "marketing.discounts.manage")
	r.GET("/marketing/discounts", h.list)
	r.GET("/marketing/discounts/:id", h.detail)
	r.POST("/marketing/discounts", manage, h.create)
	r.PUT("/marketing/discounts/:id", manage, h.update)
	r.PUT("/marketing/discounts/:id/status", manage, h.setStatus)
	r.DELETE("/marketing/discounts/:id", manage, h.remove)
}

type discountRules struct {
	PackagePrice float64  `json:"package_price"`
	ProductIDs   []uint64 `json:"product_ids"`
	FreeShipping bool     `json:"free_shipping"`
	Remark       string   `json:"remark"`
}

type activityRow struct {
	ID        uint64     `gorm:"column:id"`
	StoreID   uint64     `gorm:"column:store_id"`
	Name      string     `gorm:"column:name"`
	Rules     string     `gorm:"column:rules"`
	Status    string     `gorm:"column:status"`
	StartsAt  *time.Time `gorm:"column:starts_at"`
	EndsAt    *time.Time `gorm:"column:ends_at"`
}

type upsertInput struct {
	Name         string   `json:"name"`
	PackagePrice float64  `json:"package_price"`
	ProductIDs   []uint64 `json:"product_ids"`
	FreeShipping bool     `json:"free_shipping"`
	Remark       string   `json:"remark"`
	Status       string   `json:"status"`
	StartsAt     string   `json:"starts_at"`
	EndsAt       string   `json:"ends_at"`
}

func (h *Handler) list(c *gin.Context) {
	page, limit := pageLimit(c)
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_marketing_activity").
		Where("store_id = ? AND activity_type = ?", middleware.StoreID(c), "discount")
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if !validStatus(status) {
			response.Fail(c, http.StatusBadRequest, "活动状态错误")
			return
		}
		q = q.Where("status = ?", status)
	}
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	q = listquery.ApplyTimeColumnDateRange(q, "starts_at", strings.TrimSpace(c.Query("date_from")), strings.TrimSpace(c.Query("date_to")))
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "查询优惠套餐失败")
		return
	}
	rows := make([]activityRow, 0)
	if err := q.Select("id,store_id,name,rules,status,starts_at,ends_at").
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "查询优惠套餐失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, toItem(row))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) detail(c *gin.Context) {
	row, ok := h.loadOwned(c, c.Param("id"))
	if !ok {
		return
	}
	response.OK(c, toItem(*row))
}

func (h *Handler) create(c *gin.Context) {
	var in upsertInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, rulesJSON, status, starts, ends, err := normalizeInput(in, "draft")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	storeID := uint64(middleware.StoreID(c))
	var id uint64
	err = h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		created := struct {
			ID           uint64     `gorm:"column:id;primaryKey"`
			StoreID      uint64     `gorm:"column:store_id"`
			ActivityType string     `gorm:"column:activity_type"`
			Name         string     `gorm:"column:name"`
			Rules        string     `gorm:"column:rules"`
			Status       string     `gorm:"column:status"`
			StartsAt     *time.Time `gorm:"column:starts_at"`
			EndsAt       *time.Time `gorm:"column:ends_at"`
		}{
			StoreID: storeID, ActivityType: "discount", Name: name,
			Rules: rulesJSON, Status: status, StartsAt: starts, EndsAt: ends,
		}
		if err := tx.Table("qixi_crm_m_marketing_activity").Create(&created).Error; err != nil {
			return err
		}
		id = created.ID
		if id == 0 {
			return errors.New("missing activity id")
		}
		return upsertView(h.businessDB.WithContext(c.Request.Context()), id, storeID, name, rulesJSON, status, starts, ends)
	})
	if err != nil {
		fail(c, "创建优惠套餐失败")
		return
	}
	row, ok := h.loadOwned(c, strconv.FormatUint(id, 10))
	if !ok {
		return
	}
	response.OK(c, toItem(*row))
}

func (h *Handler) update(c *gin.Context) {
	row, ok := h.loadOwned(c, c.Param("id"))
	if !ok {
		return
	}
	var in upsertInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	name, rulesJSON, status, starts, ends, err := normalizeInput(in, row.Status)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(
			`UPDATE qixi_crm_m_marketing_activity SET name=?, rules=?, status=?, starts_at=?, ends_at=?
			 WHERE id=? AND store_id=? AND activity_type='discount'`,
			name, rulesJSON, status, starts, ends, row.ID, row.StoreID,
		)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errNotFound
		}
		return upsertView(h.businessDB.WithContext(c.Request.Context()), row.ID, row.StoreID, name, rulesJSON, status, starts, ends)
	})
	if err != nil {
		if errors.Is(err, errNotFound) {
			response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
			return
		}
		fail(c, "更新优惠套餐失败")
		return
	}
	updated, ok := h.loadOwned(c, strconv.FormatUint(row.ID, 10))
	if !ok {
		return
	}
	response.OK(c, toItem(*updated))
}

func (h *Handler) setStatus(c *gin.Context) {
	row, ok := h.loadOwned(c, c.Param("id"))
	if !ok {
		return
	}
	var in struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil || !validStatus(strings.TrimSpace(in.Status)) {
		response.Fail(c, http.StatusBadRequest, "活动状态错误")
		return
	}
	status := strings.TrimSpace(in.Status)
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(
			`UPDATE qixi_crm_m_marketing_activity SET status=? WHERE id=? AND store_id=? AND activity_type='discount'`,
			status, row.ID, row.StoreID,
		)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errNotFound
		}
		return upsertView(h.businessDB.WithContext(c.Request.Context()), row.ID, row.StoreID, row.Name, row.Rules, status, row.StartsAt, row.EndsAt)
	})
	if err != nil {
		if errors.Is(err, errNotFound) {
			response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
			return
		}
		fail(c, "更新优惠套餐状态失败")
		return
	}
	updated, ok := h.loadOwned(c, strconv.FormatUint(row.ID, 10))
	if !ok {
		return
	}
	response.OK(c, toItem(*updated))
}

func (h *Handler) remove(c *gin.Context) {
	row, ok := h.loadOwned(c, c.Param("id"))
	if !ok {
		return
	}
	err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(
			`DELETE FROM qixi_crm_m_marketing_activity WHERE id=? AND store_id=? AND activity_type='discount'`,
			row.ID, row.StoreID,
		)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errNotFound
		}
		return h.businessDB.WithContext(c.Request.Context()).
			Exec(`DELETE FROM qixi_crm_b_marketing_activity_view WHERE activity_id=? AND activity_type='discount'`, row.ID).Error
	})
	if err != nil {
		if errors.Is(err, errNotFound) {
			response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
			return
		}
		fail(c, "删除优惠套餐失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) loadOwned(c *gin.Context, rawID string) (*activityRow, bool) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "活动 ID 错误")
		return nil, false
	}
	var row activityRow
	err = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_marketing_activity").
		Select("id,store_id,name,rules,status,starts_at,ends_at").
		Where("id = ? AND store_id = ? AND activity_type = ?", id, middleware.StoreID(c), "discount").
		Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "优惠套餐不存在")
		return nil, false
	}
	if err != nil {
		fail(c, "查询优惠套餐失败")
		return nil, false
	}
	return &row, true
}

func upsertView(businessDB *gorm.DB, id, storeID uint64, name, rulesJSON, status string, starts, ends *time.Time) error {
	viewStatus := 0
	if status == "active" {
		viewStatus = 1
	}
	return businessDB.Exec(`
		INSERT INTO qixi_crm_b_marketing_activity_view
		  (activity_id,store_id,activity_type,name,rules,status,version,starts_at,ends_at)
		VALUES (?,?,?,?,?,?,1,?,?)
		ON DUPLICATE KEY UPDATE
		  store_id=VALUES(store_id), activity_type=VALUES(activity_type), name=VALUES(name),
		  rules=VALUES(rules), status=VALUES(status), version=version+1,
		  starts_at=VALUES(starts_at), ends_at=VALUES(ends_at)`,
		id, storeID, "discount", name, rulesJSON, viewStatus, starts, ends,
	).Error
}

func normalizeInput(in upsertInput, defaultStatus string) (string, string, string, *time.Time, *time.Time, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return "", "", "", nil, nil, errors.New("活动名称必填且不超过 128 字")
	}
	if in.PackagePrice <= 0 {
		return "", "", "", nil, nil, errors.New("套餐价必须大于 0")
	}
	if len(in.ProductIDs) == 0 {
		return "", "", "", nil, nil, errors.New("请至少选择一个商品")
	}
	if utf8.RuneCountInString(strings.TrimSpace(in.Remark)) > 255 {
		return "", "", "", nil, nil, errors.New("备注不能超过 255 字")
	}
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = defaultStatus
	}
	if !validStatus(status) {
		return "", "", "", nil, nil, errors.New("活动状态错误")
	}
	starts, err := parseOptionalTime(in.StartsAt)
	if err != nil {
		return "", "", "", nil, nil, errors.New("开始时间格式错误")
	}
	ends, err := parseOptionalTime(in.EndsAt)
	if err != nil {
		return "", "", "", nil, nil, errors.New("结束时间格式错误")
	}
	if starts != nil && ends != nil && ends.Before(*starts) {
		return "", "", "", nil, nil, errors.New("结束时间不能早于开始时间")
	}
	raw, err := json.Marshal(discountRules{
		PackagePrice: in.PackagePrice,
		ProductIDs:   in.ProductIDs,
		FreeShipping: in.FreeShipping,
		Remark:       strings.TrimSpace(in.Remark),
	})
	if err != nil {
		return "", "", "", nil, nil, errors.New("规则编码失败")
	}
	return name, string(raw), status, starts, ends, nil
}

func toItem(row activityRow) gin.H {
	rules := discountRules{}
	_ = json.Unmarshal([]byte(row.Rules), &rules)
	return gin.H{
		"activity_id":   row.ID,
		"store_id":      row.StoreID,
		"name":          row.Name,
		"package_price": rules.PackagePrice,
		"product_ids":   rules.ProductIDs,
		"free_shipping": rules.FreeShipping,
		"remark":        rules.Remark,
		"status":        row.Status,
		"starts_at":     formatTime(row.StartsAt),
		"ends_at":       formatTime(row.EndsAt),
	}
}

func validStatus(status string) bool {
	switch status {
	case "draft", "pending", "active", "closed", "rejected":
		return true
	default:
		return false
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

func fail(c *gin.Context, msg string) {
	response.Fail(c, http.StatusInternalServerError, msg)
}

var errNotFound = errors.New("discount not found")
