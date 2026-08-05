package merchanttype

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const typeTable = "qixi_crm_a_merchant_type"
const typeMenuTable = "qixi_crm_a_merchant_type_menu"

var errNotFound = errors.New("店铺类型不存在")
var errInvalid = errors.New("店铺类型参数不合法")
var errConflict = errors.New("店铺类型名称已存在")

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) Register(r gin.IRoutes) {
	p := middleware.RequireAdminRoles("platform")
	m := middleware.RequireAdminMenu(h.db, "merchant.type.manage")
	r.GET("/merchant-types", p, h.List)
	r.GET("/merchant-types/:id", p, m, h.Get)
	r.POST("/merchant-types", p, m, h.Create)
	r.PUT("/merchant-types/:id", p, m, h.Update)
	r.PUT("/merchant-types/:id/remark", p, m, h.UpdateRemark)
	r.PUT("/merchant-types/:id/status", p, m, h.UpdateStatus)
	r.DELETE("/merchant-types/:id", p, m, h.Delete)
}

type record struct {
	ID          uint     `gorm:"column:id;primaryKey" json:"id"`
	Name        string   `gorm:"column:name" json:"name"`
	TypeInfo    string   `gorm:"column:type_info" json:"type_info"`
	IsMargin    int8     `gorm:"column:is_margin" json:"is_margin"`
	Margin      float64  `gorm:"column:margin" json:"margin"`
	Description string   `gorm:"column:description" json:"description"`
	Remark      string   `gorm:"column:remark" json:"remark"`
	Status      int8     `gorm:"column:status" json:"status"`
	MenuCodes   []string `gorm:"-" json:"menu_codes"`
}

func (record) TableName() string { return typeTable }

type saveReq struct {
	Name        string   `json:"name"`
	TypeInfo    string   `json:"type_info"`
	IsMargin    bool     `json:"is_margin"`
	Margin      float64  `json:"margin"`
	Description string   `json:"description"`
	Remark      string   `json:"remark"`
	Status      *bool    `json:"status"`
	MenuCodes   []string `json:"menu_codes"`
}

func (h *Handler) List(c *gin.Context) {
	var rows []record
	q := h.db.WithContext(c.Request.Context()).Order("id DESC")
	if k := strings.TrimSpace(c.Query("keyword")); k != "" {
		q = q.Where("name LIKE ?", "%"+k+"%")
	}
	if err := q.Find(&rows).Error; err != nil {
		fail(c, err)
		return
	}
	if err := h.fill(c, rows); err != nil {
		fail(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) Get(c *gin.Context) {
	id, ok := idOf(c)
	if !ok {
		return
	}
	var row record
	if err := h.db.WithContext(c.Request.Context()).First(&row, id).Error; err != nil {
		fail(c, err)
		return
	}
	if err := h.fill(c, []record{row}); err != nil {
		fail(c, err)
		return
	}
	var codes []string
	_ = h.db.WithContext(c.Request.Context()).Table(typeMenuTable).Where("merchant_type_id = ?", id).Pluck("menu_code", &codes).Error
	row.MenuCodes = codes
	response.OK(c, row)
}
func (h *Handler) Create(c *gin.Context) {
	var req saveReq
	if c.ShouldBindJSON(&req) != nil {
		fail(c, errInvalid)
		return
	}
	row, err := h.save(c, 0, req)
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) Update(c *gin.Context) {
	id, ok := idOf(c)
	if !ok {
		return
	}
	var req saveReq
	if c.ShouldBindJSON(&req) != nil {
		fail(c, errInvalid)
		return
	}
	row, err := h.save(c, id, req)
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, row)
}
func (h *Handler) UpdateRemark(c *gin.Context) {
	id, ok := idOf(c)
	if !ok {
		return
	}
	var req struct {
		Remark string `json:"remark"`
	}
	if c.ShouldBindJSON(&req) != nil || len([]rune(strings.TrimSpace(req.Remark))) > 500 {
		fail(c, errInvalid)
		return
	}
	res := h.db.WithContext(c.Request.Context()).Table(typeTable).Where("id = ?", id).Update("remark", strings.TrimSpace(req.Remark))
	if res.Error != nil || res.RowsAffected == 0 {
		fail(c, first(res.Error, errNotFound))
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, ok := idOf(c)
	if !ok {
		return
	}
	var req struct {
		Enabled *bool `json:"enabled"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Enabled == nil {
		fail(c, errInvalid)
		return
	}
	res := h.db.WithContext(c.Request.Context()).Table(typeTable).Where("id = ?", id).Update("status", boolNum(*req.Enabled))
	if res.Error != nil || res.RowsAffected == 0 {
		fail(c, first(res.Error, errNotFound))
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) Delete(c *gin.Context) {
	id, ok := idOf(c)
	if !ok {
		return
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ?", id).Delete(&record{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errNotFound
		}
		return tx.Table(typeMenuTable).Where("merchant_type_id = ?", id).Delete(nil).Error
	})
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}
func (h *Handler) save(c *gin.Context, id uint, req saveReq) (record, error) {
	if err := validate(&req); err != nil {
		return record{}, err
	}
	var out record
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		values := map[string]any{"name": strings.TrimSpace(req.Name), "type_info": strings.TrimSpace(req.TypeInfo), "is_margin": boolNum(req.IsMargin), "margin": req.Margin, "description": strings.TrimSpace(req.Description), "remark": strings.TrimSpace(req.Remark), "status": boolPtrNum(req.Status, true)}
		if id == 0 {
			out = record{Name: values["name"].(string), TypeInfo: values["type_info"].(string), IsMargin: values["is_margin"].(int8), Margin: req.Margin, Description: values["description"].(string), Remark: values["remark"].(string), Status: values["status"].(int8)}
			if err := tx.Create(&out).Error; err != nil {
				return err
			}
		} else {
			res := tx.Table(typeTable).Where("id = ?", id).Updates(values)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errNotFound
			}
			if err := tx.First(&out, id).Error; err != nil {
				return err
			}
		}
		if err := tx.Table(typeMenuTable).Where("merchant_type_id = ?", out.ID).Delete(nil).Error; err != nil {
			return err
		}
		for _, code := range unique(req.MenuCodes) {
			if err := tx.Table(typeMenuTable).Create(map[string]any{"merchant_type_id": out.ID, "menu_code": code}).Error; err != nil {
				return err
			}
		}
		out.MenuCodes = unique(req.MenuCodes)
		return nil
	})
	return out, err
}
func (h *Handler) fill(c *gin.Context, rows []record) error {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	var links []struct {
		MerchantTypeID uint   `gorm:"column:merchant_type_id"`
		MenuCode       string `gorm:"column:menu_code"`
	}
	if err := h.db.WithContext(c.Request.Context()).Table(typeMenuTable).Where("merchant_type_id IN ?", ids).Find(&links).Error; err != nil {
		return err
	}
	by := map[uint][]string{}
	for _, l := range links {
		by[l.MerchantTypeID] = append(by[l.MerchantTypeID], l.MenuCode)
	}
	for i := range rows {
		rows[i].MenuCodes = by[rows[i].ID]
	}
	return nil
}
func validate(r *saveReq) error {
	r.Name = strings.TrimSpace(r.Name)
	r.TypeInfo = strings.TrimSpace(r.TypeInfo)
	r.Description = strings.TrimSpace(r.Description)
	r.Remark = strings.TrimSpace(r.Remark)
	if r.Name == "" || len([]rune(r.Name)) > 128 || len([]rune(r.TypeInfo)) > 500 || len([]rune(r.Description)) > 65535 || len([]rune(r.Remark)) > 500 || !validMargin(r.Margin) {
		return errInvalid
	}
	if r.IsMargin && r.Margin <= 0 {
		return errors.New("保证金必须大于 0")
	}
	if !r.IsMargin {
		r.Margin = 0
	}
	for _, v := range r.MenuCodes {
		if len(strings.TrimSpace(v)) == 0 || len(v) > 128 {
			return errInvalid
		}
	}
	return nil
}
func unique(in []string) []string {
	m := map[string]struct{}{}
	out := []string{}
	for _, v := range in {
		v = strings.TrimSpace(v)
		if _, ok := m[v]; v != "" && !ok {
			m[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
func idOf(c *gin.Context) (uint, bool) {
	n, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil || n == 0 {
		fail(c, errInvalid)
		return 0, false
	}
	return uint(n), true
}
func boolNum(v bool) int8 {
	if v {
		return 1
	}
	return 0
}
func boolPtrNum(v *bool, fallback bool) int8 {
	if v == nil {
		return boolNum(fallback)
	}
	return boolNum(*v)
}
func first(a, b error) error {
	if a != nil {
		return a
	}
	return b
}
func fail(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, errNotFound) {
		response.Fail(c, http.StatusNotFound, errNotFound.Error())
		return
	}
	if errors.Is(err, errInvalid) || strings.Contains(err.Error(), "保证金") {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, errConflict) || isDuplicateEntry(err) {
		response.Fail(c, http.StatusConflict, errConflict.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, "店铺类型操作失败")
}

func validMargin(value float64) bool {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 9_999_999_999.99 {
		return false
	}
	cents := math.Round(value * 100)
	return math.Abs(value*100-cents) < 1e-7
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
