// Package storegroup implements the platform-owned merchant store grouping view.
// It deliberately uses the unified admin merchant projection rather than a direct
// merchant database join, so it remains safe while cross-database projections lag.
package storegroup

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	groupTable       = "qixi_crm_a_store_group"
	groupMerchantTbl = "qixi_crm_a_store_group_merchant"
	maxGroupLevel    = 2 // root is level 0, which gives the CRMEB baseline three levels.
)

var (
	errGroupNotFound  = errors.New("店铺分组不存在")
	errGroupBadInput  = errors.New("店铺分组参数不合法")
	errGroupHasChild  = errors.New("分组下仍有子分组，不能删除")
	errGroupCycle     = errors.New("不能将分组移动到自身或其子分组")
	errGroupTooDeep   = errors.New("店铺分组最多支持三级")
	errMerchantAbsent = errors.New("关联店铺不存在于统一后台商户投影")
	errGroupConflict  = errors.New("同一上级分组下名称不能重复")
)

type Handler struct{ adminDB *gorm.DB }

func NewHandler(adminDB *gorm.DB) *Handler { return &Handler{adminDB: adminDB} }

func (h *Handler) Register(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	manage := middleware.RequireAdminMenu(h.adminDB, "merchant.group.manage")
	r.GET("/store-groups", platform, h.List)
	r.POST("/store-groups", platform, manage, h.Create)
	r.GET("/store-groups/:id", platform, manage, h.Get)
	r.PUT("/store-groups/:id", platform, manage, h.Update)
	r.DELETE("/store-groups/:id", platform, manage, h.Delete)
	r.POST("/store-groups/:id/status", platform, manage, h.SetStatus)
	r.POST("/store-groups/:id/template", platform, manage, h.SetTemplate)
	r.GET("/store-groups/:id/merchants", platform, manage, h.ListMerchants)
}

type group struct {
	ID                uint     `gorm:"column:id;primaryKey" json:"id"`
	ParentID          uint     `gorm:"column:parent_id" json:"parent_id"`
	Path              string   `gorm:"column:path" json:"path"`
	Level             uint8    `gorm:"column:level" json:"level"`
	Name              string   `gorm:"column:name" json:"name"`
	Sort              int      `gorm:"column:sort" json:"sort"`
	Status            int8     `gorm:"column:status" json:"status"`
	DiyPageID         uint     `gorm:"column:diy_page_id" json:"diy_page_id"`
	PositioningStatus int8     `gorm:"column:positioning_status" json:"positioning_status"`
	Longitude         *float64 `gorm:"column:longitude" json:"longitude"`
	Latitude          *float64 `gorm:"column:latitude" json:"latitude"`
	Address           string   `gorm:"column:address" json:"address"`
	MerchantCount     int64    `gorm:"-" json:"merchant_count"`
	MerchantIDs       []uint   `gorm:"-" json:"merchant_ids"`
	Children          []group  `gorm:"-" json:"children,omitempty"`
}

func (group) TableName() string { return groupTable }

type saveRequest struct {
	ParentID          uint     `json:"parent_id"`
	Name              string   `json:"name"`
	Sort              int      `json:"sort"`
	Status            *bool    `json:"status"`
	DiyPageID         uint     `json:"diy_page_id"`
	PositioningStatus *bool    `json:"positioning_status"`
	Longitude         *float64 `json:"longitude"`
	Latitude          *float64 `json:"latitude"`
	Address           string   `json:"address"`
	MerchantIDs       []uint   `json:"merchant_ids"`
}

func (h *Handler) List(c *gin.Context) {
	var rows []group
	q := h.adminDB.WithContext(c.Request.Context()).Table(groupTable).Order("sort DESC, id ASC")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := q.Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺分组失败")
		return
	}
	if err := h.fillGroupMembers(c, rows); err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询关联店铺失败")
		return
	}
	response.OK(c, gin.H{"list": buildTree(rows)})
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	row, err := h.findGroup(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := h.fillGroupMembers(c, []group{*row}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询关联店铺失败")
		return
	}
	var ids []uint
	if err := h.adminDB.WithContext(c.Request.Context()).Table(groupMerchantTbl).Where("store_group_id = ?", id).Pluck("merchant_id", &ids).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询关联店铺失败")
		return
	}
	row.MerchantIDs = ids
	row.MerchantCount = int64(len(ids))
	response.OK(c, row)
}

func (h *Handler) Create(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "店铺分组参数不合法")
		return
	}
	if err := validateRequest(&req); err != nil {
		writeError(c, err)
		return
	}
	var created group
	err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		level, path, err := groupParent(tx, req.ParentID, 0, "")
		if err != nil {
			return err
		}
		created = group{ParentID: req.ParentID, Level: level, Path: path, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: boolInt(req.Status, true), DiyPageID: req.DiyPageID, PositioningStatus: boolInt(req.PositioningStatus, false), Longitude: req.Longitude, Latitude: req.Latitude, Address: strings.TrimSpace(req.Address)}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		created.Path = path + strconv.FormatUint(uint64(created.ID), 10) + "/"
		if err := tx.Model(&created).Update("path", created.Path).Error; err != nil {
			return err
		}
		return replaceMembers(tx, created.ID, req.MerchantIDs)
	})
	if err != nil {
		writeError(c, err)
		return
	}
	created.MerchantIDs = uniqueIDs(req.MerchantIDs)
	created.MerchantCount = int64(len(created.MerchantIDs))
	response.OK(c, created)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "店铺分组参数不合法")
		return
	}
	if err := validateRequest(&req); err != nil {
		writeError(c, err)
		return
	}
	var updated group
	err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var current group
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, id).Error; err != nil {
			return normalizeNotFound(err)
		}
		if req.ParentID == current.ID || (req.ParentID != 0 && strings.HasPrefix(parentPath(tx, req.ParentID), current.Path)) {
			return errGroupCycle
		}
		level, prefix, err := groupParent(tx, req.ParentID, current.ID, current.Path)
		if err != nil {
			return err
		}
		var deepest uint8
		if err := tx.Table(groupTable).Where("path LIKE ?", current.Path+"%").Select("COALESCE(MAX(level), 0)").Scan(&deepest).Error; err != nil {
			return err
		}
		delta := int(level) - int(current.Level)
		if int(deepest)+delta > maxGroupLevel {
			return errGroupTooDeep
		}
		newPath := prefix + strconv.FormatUint(uint64(current.ID), 10) + "/"
		updates := map[string]any{"parent_id": req.ParentID, "path": newPath, "level": level, "name": strings.TrimSpace(req.Name), "sort": req.Sort, "diy_page_id": req.DiyPageID, "positioning_status": boolInt(req.PositioningStatus, false), "longitude": req.Longitude, "latitude": req.Latitude, "address": strings.TrimSpace(req.Address)}
		if req.Status != nil {
			updates["status"] = boolInt(req.Status, current.Status == 1)
		}
		if err := tx.Model(&current).Updates(updates).Error; err != nil {
			return err
		}
		var descendants []group
		if err := tx.Where("path LIKE ? AND id <> ?", current.Path+"%", current.ID).Find(&descendants).Error; err != nil {
			return err
		}
		for _, child := range descendants {
			if !strings.HasPrefix(child.Path, current.Path) {
				return fmt.Errorf("%w: 分组路径不一致", errGroupBadInput)
			}
			if err := tx.Model(&group{}).Where("id = ?", child.ID).Updates(map[string]any{"path": newPath + strings.TrimPrefix(child.Path, current.Path), "level": int(child.Level) + delta}).Error; err != nil {
				return err
			}
		}
		if err := replaceMembers(tx, current.ID, req.MerchantIDs); err != nil {
			return err
		}
		return tx.First(&updated, current.ID).Error
	})
	if err != nil {
		writeError(c, err)
		return
	}
	updated.MerchantIDs = uniqueIDs(req.MerchantIDs)
	updated.MerchantCount = int64(len(updated.MerchantIDs))
	response.OK(c, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var total int64
		if err := tx.Table(groupTable).Where("parent_id = ?", id).Count(&total).Error; err != nil {
			return err
		}
		if total > 0 {
			return errGroupHasChild
		}
		res := tx.Where("id = ?", id).Delete(&group{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errGroupNotFound
		}
		return tx.Table(groupMerchantTbl).Where("store_group_id = ?", id).Delete(nil).Error
	})
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		Enabled *bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Enabled == nil {
		response.Fail(c, http.StatusBadRequest, "状态参数不合法")
		return
	}
	err := h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var current group
		if err := tx.First(&current, id).Error; err != nil {
			return normalizeNotFound(err)
		}
		return tx.Table(groupTable).Where("path LIKE ?", current.Path+"%").Update("status", boolInt(req.Enabled, false)).Error
	})
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) SetTemplate(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req struct {
		DiyPageID uint `json:"diy_page_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "装修模板参数不合法")
		return
	}
	// 0 表示清空绑定；非零模板必须来自平台 DIY 页面，避免存入跨域或不存在的 ID。
	if req.DiyPageID != 0 {
		var total int64
		if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_diy_page").Where("id = ?", req.DiyPageID).Count(&total).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "校验装修模板失败")
			return
		}
		if total == 0 {
			response.Fail(c, http.StatusBadRequest, "装修模板不存在")
			return
		}
	}
	res := h.adminDB.WithContext(c.Request.Context()).Table(groupTable).Where("id = ?", id).Update("diy_page_id", req.DiyPageID)
	if res.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "设置装修模板失败")
		return
	}
	if res.RowsAffected == 0 {
		writeError(c, errGroupNotFound)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListMerchants(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if _, err := h.findGroup(c, id); err != nil {
		writeError(c, err)
		return
	}
	var rows []struct {
		MerchantID   uint   `json:"merchant_id"`
		MerchantName string `json:"merchant_name"`
		RegionID     uint   `json:"region_id"`
		Status       int8   `json:"status"`
	}
	err := h.adminDB.WithContext(c.Request.Context()).Table(groupMerchantTbl+" AS gm").
		Select("m.merchant_id, m.merchant_name, m.region_id, m.status").
		Joins("INNER JOIN qixi_crm_a_merchant_view AS m ON m.merchant_id = gm.merchant_id").
		Where("gm.store_group_id = ?", id).Order("m.merchant_id ASC").Find(&rows).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询关联店铺失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) findGroup(c *gin.Context, id uint) (*group, error) {
	var row group
	err := h.adminDB.WithContext(c.Request.Context()).First(&row, id).Error
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return &row, nil
}

func (h *Handler) fillGroupMembers(c *gin.Context, rows []group) error {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	var counts []struct {
		StoreGroupID uint
		Total        int64
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table(groupMerchantTbl).Select("store_group_id, COUNT(*) AS total").Where("store_group_id IN ?", ids).Group("store_group_id").Find(&counts).Error; err != nil {
		return err
	}
	byID := make(map[uint]int64, len(counts))
	for _, item := range counts {
		byID[item.StoreGroupID] = item.Total
	}
	for i := range rows {
		rows[i].MerchantCount = byID[rows[i].ID]
	}
	return nil
}

func validateRequest(req *saveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.TrimSpace(req.Address)
	if req.Name == "" || len([]rune(req.Name)) > 128 || len([]rune(req.Address)) > 255 {
		return errGroupBadInput
	}
	if (req.Longitude == nil) != (req.Latitude == nil) {
		return errGroupBadInput
	}
	if req.Longitude != nil && (*req.Longitude < -180 || *req.Longitude > 180 || *req.Latitude < -90 || *req.Latitude > 90) {
		return errGroupBadInput
	}
	for _, id := range req.MerchantIDs {
		if id == 0 {
			return errGroupBadInput
		}
	}
	return nil
}

func groupParent(tx *gorm.DB, parentID, selfID uint, selfPath string) (uint8, string, error) {
	if parentID == 0 {
		return 0, "/", nil
	}
	var parent group
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&parent, parentID).Error; err != nil {
		return 0, "", normalizeNotFound(err)
	}
	if parent.ID == selfID || (selfPath != "" && strings.HasPrefix(parent.Path, selfPath)) {
		return 0, "", errGroupCycle
	}
	if parent.Level >= maxGroupLevel {
		return 0, "", errGroupTooDeep
	}
	return parent.Level + 1, parent.Path, nil
}

func parentPath(tx *gorm.DB, id uint) string {
	var path string
	if id == 0 || tx.Table(groupTable).Where("id = ?", id).Pluck("path", &path).Error != nil {
		return ""
	}
	return path
}

func replaceMembers(tx *gorm.DB, groupID uint, merchantIDs []uint) error {
	ids := uniqueIDs(merchantIDs)
	if len(ids) > 0 {
		var total int64
		if err := tx.Table("qixi_crm_a_merchant_view").Where("merchant_id IN ?", ids).Count(&total).Error; err != nil {
			return err
		}
		if total != int64(len(ids)) {
			return errMerchantAbsent
		}
	}
	if err := tx.Table(groupMerchantTbl).Where("store_group_id = ?", groupID).Delete(nil).Error; err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(ids))
	for _, merchantID := range ids {
		rows = append(rows, map[string]any{"store_group_id": groupID, "merchant_id": merchantID})
	}
	return tx.Table(groupMerchantTbl).Create(rows).Error
}

func uniqueIDs(in []uint) []uint {
	seen := make(map[uint]struct{}, len(in))
	out := make([]uint, 0, len(in))
	for _, id := range in {
		if id != 0 {
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				out = append(out, id)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func buildTree(rows []group) []group {
	byParent := make(map[uint][]group)
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], row)
	}
	var build func(uint) []group
	build = func(parentID uint) []group {
		list := byParent[parentID]
		for i := range list {
			list[i].Children = build(list[i].ID)
		}
		return list
	}
	return build(0)
}

func boolInt(value *bool, fallback bool) int8 {
	if value == nil {
		if fallback {
			return 1
		}
		return 0
	}
	if *value {
		return 1
	}
	return 0
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺分组 ID 不合法")
		return 0, false
	}
	return uint(id), true
}

func normalizeNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errGroupNotFound
	}
	return err
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errGroupNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, errGroupConflict), isDuplicateEntry(err):
		response.Fail(c, http.StatusConflict, errGroupConflict.Error())
	case errors.Is(err, errGroupBadInput), errors.Is(err, errGroupCycle), errors.Is(err, errGroupTooDeep), errors.Is(err, errMerchantAbsent), errors.Is(err, errGroupHasChild):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "店铺分组操作失败")
	}
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
