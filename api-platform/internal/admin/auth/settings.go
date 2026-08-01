package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterSettings 提供统一后台自身的账号、角色和菜单管理。路由保持与 Vben 现有
// 页面一致，但所有读写仅作用于 qixi_crm_a_* 表。
func (h *Handler) RegisterSettings(authed gin.IRoutes) {
	authed.GET("/setting/admins", h.ListAdmins)
	authed.POST("/setting/admins", h.CreateAdmin)
	authed.PUT("/setting/admins/:id", h.UpdateAdmin)
	authed.GET("/setting/roles", h.ListRoles)
	authed.POST("/setting/roles", h.CreateRole)
	authed.PUT("/setting/roles/:id", h.UpdateRole)
	authed.GET("/setting/menus/tree", h.MenuTree)
}

type adminListRow struct {
	ID          uint64    `gorm:"column:id" json:"admin_id"`
	Username    string    `gorm:"column:username" json:"account"`
	DisplayName string    `gorm:"column:display_name" json:"real_name"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	Status      int8      `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	Roles       string    `json:"roles"`
	RegionIDs   string    `json:"region_ids"`
	IsAgent     int8      `json:"is_agent"`
	Level       int8      `json:"level"`
	CircleID    uint64    `json:"circle_agent_id"`
}

type adminSaveRequest struct {
	Account   string   `json:"account"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	RealName  string   `json:"real_name"`
	Phone     string   `json:"phone"`
	Roles     string   `json:"roles"`
	RoleCodes []string `json:"role_codes"`
	Status    int8     `json:"status"`
	RegionIDs string   `json:"region_ids"`
}

func (h *Handler) ListAdmins(c *gin.Context) {
	page, limit := pageParams(c, 20)
	var total int64
	if err := h.db.WithContext(c.Request.Context()).Model(&adminUser{}).Count(&total).Error; err != nil {
		writeError(c, err)
		return
	}
	rows := make([]adminListRow, 0)
	if err := h.db.WithContext(c.Request.Context()).Table((adminUser{}).TableName()).
		Select("id,username,display_name,phone,status,created_at").Order("id DESC").
		Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		writeError(c, err)
		return
	}
	for i := range rows {
		roles, err := h.roles(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].Roles = strings.Join(roles, ",")
		rows[i].IsAgent = boolToInt8(hasRole(roles, "region"))
		rows[i].RegionIDs, err = h.regionIDs(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (h *Handler) CreateAdmin(c *gin.Context) {
	var req adminSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	username := strings.TrimSpace(req.Account)
	if username == "" {
		username = strings.TrimSpace(req.Username)
	}
	roles := normalizedRoleCodes(req)
	if username == "" || len(req.Password) < 8 || len(roles) == 0 {
		response.Fail(c, http.StatusBadRequest, "账号、至少 8 位密码和角色必填")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(c, err)
		return
	}
	user := adminUser{Username: username, PasswordHash: string(hash), DisplayName: strings.TrimSpace(req.RealName), Phone: strings.TrimSpace(req.Phone), Status: req.Status, DataScopeVersion: 1}
	if user.Status == 0 {
		user.Status = 1
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return h.replaceAdminRolesAndScope(c, tx, user.ID, roles, req.RegionIDs)
	}); err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "账号已存在")
			return
		}
		writeError(c, err)
		return
	}
	response.OK(c, h.adminRow(c, user))
}

func (h *Handler) UpdateAdmin(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "管理员 ID 错误")
		return
	}
	var req adminSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := h.findUserByID(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	roles := normalizedRoleCodes(req)
	if len(roles) == 0 {
		response.Fail(c, http.StatusBadRequest, "至少保留一个角色")
		return
	}
	if req.Password != "" && len(req.Password) < 8 {
		response.Fail(c, http.StatusBadRequest, "新密码至少 8 位")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{"display_name": strings.TrimSpace(req.RealName), "phone": strings.TrimSpace(req.Phone), "status": req.Status, "data_scope_version": gorm.Expr("data_scope_version + 1")}
		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			updates["password_hash"] = string(hash)
		}
		if err := tx.Model(&adminUser{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return err
		}
		return h.replaceAdminRolesAndScope(c, tx, user.ID, roles, req.RegionIDs)
	}); err != nil {
		writeError(c, err)
		return
	}
	updated, err := h.findUserByID(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, h.adminRow(c, *updated))
}

func (h *Handler) adminRow(c *gin.Context, user adminUser) adminListRow {
	roles, _ := h.roles(c, user.ID)
	regionIDs, _ := h.regionIDs(c, user.ID)
	return adminListRow{ID: user.ID, Username: user.Username, DisplayName: user.DisplayName, Phone: user.Phone, Status: user.Status, Roles: strings.Join(roles, ","), RegionIDs: regionIDs, IsAgent: boolToInt8(hasRole(roles, "region"))}
}

func (h *Handler) replaceAdminRolesAndScope(c *gin.Context, tx *gorm.DB, userID uint64, codes []string, regionIDs string) error {
	var roleRows []struct {
		ID   uint64
		Code string
	}
	if err := tx.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status = 1", codes).Scan(&roleRows).Error; err != nil {
		return err
	}
	if len(roleRows) != len(codes) {
		return errors.New("角色不存在或已禁用")
	}
	if err := tx.WithContext(c.Request.Context()).Where("admin_user_id = ?", userID).Delete(&adminUserRole{}).Error; err != nil {
		return err
	}
	bindings := make([]adminUserRole, 0, len(roleRows))
	for _, role := range roleRows {
		bindings = append(bindings, adminUserRole{AdminUserID: userID, RoleID: role.ID})
	}
	if err := tx.WithContext(c.Request.Context()).Create(&bindings).Error; err != nil {
		return err
	}
	if err := tx.WithContext(c.Request.Context()).Where("admin_user_id = ? AND scope_type = ?", userID, "region").Delete(&adminDataScope{}).Error; err != nil {
		return err
	}
	if !hasRole(codes, "region") {
		return nil
	}
	values, err := parseIDs(regionIDs)
	if err != nil || len(values) == 0 {
		return errors.New("区域角色必须配置区域 ID")
	}
	payload, err := json.Marshal(values)
	if err != nil {
		return err
	}
	return tx.WithContext(c.Request.Context()).Create(&adminDataScope{AdminUserID: userID, ScopeType: "region", ScopeValue: payload, Version: 1}).Error
}

type adminUserRole struct{ AdminUserID, RoleID uint64 }

func (adminUserRole) TableName() string { return "qixi_crm_a_admin_user_role" }

type adminDataScope struct {
	AdminUserID uint64          `gorm:"column:admin_user_id"`
	ScopeType   string          `gorm:"column:scope_type"`
	ScopeValue  json.RawMessage `gorm:"column:scope_value"`
	Version     uint64          `gorm:"column:version"`
}

func (adminDataScope) TableName() string { return "qixi_crm_a_data_scope" }

func (h *Handler) regionIDs(c *gin.Context, userID uint64) (string, error) {
	var row adminDataScope
	err := h.db.WithContext(c.Request.Context()).Where("admin_user_id = ? AND scope_type = ?", userID, "region").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var values []uint64
	if err := json.Unmarshal(row.ScopeValue, &values); err != nil {
		return "", err
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.FormatUint(value, 10))
	}
	return strings.Join(parts, ","), nil
}

type roleListRow struct {
	ID       uint64 `gorm:"column:id" json:"role_id"`
	Code     string `gorm:"column:code" json:"code"`
	Name     string `gorm:"column:name" json:"role_name"`
	Status   int8   `gorm:"column:status" json:"status"`
	Rules    string `json:"rules"`
	IsAgent  int8   `json:"is_agent"`
	CircleID uint64 `json:"circle_id"`
	MerID    uint64 `json:"mer_id"`
}
type roleSaveRequest struct {
	Code     string   `json:"code"`
	RoleName string   `json:"role_name"`
	Status   int8     `json:"status"`
	MenuIDs  []uint64 `json:"menu_ids"`
}

func (h *Handler) ListRoles(c *gin.Context) {
	page, limit := pageParams(c, 50)
	var total int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Count(&total).Error; err != nil {
		writeError(c, err)
		return
	}
	rows := make([]roleListRow, 0)
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,code,name,status").Order("id ASC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		writeError(c, err)
		return
	}
	for i := range rows {
		rules, err := h.menuIDs(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].Rules = rules
		rows[i].IsAgent = boolToInt8(strings.HasPrefix(rows[i].Code, "region"))
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req roleSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if !validRoleCode(req.Code) || strings.TrimSpace(req.RoleName) == "" {
		response.Fail(c, http.StatusBadRequest, "角色代码和名称必填，代码仅支持小写字母、数字和下划线")
		return
	}
	row := roleListRow{Code: strings.TrimSpace(req.Code), Name: strings.TrimSpace(req.RoleName), Status: req.Status}
	if row.Status == 0 {
		row.Status = 1
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_role").Create(map[string]any{"code": row.Code, "name": row.Name, "status": row.Status}).Error; err != nil {
			return err
		}
		var saved struct{ ID uint64 }
		if err := tx.Table("qixi_crm_a_role").Select("id").Where("code = ?", row.Code).Scan(&saved).Error; err != nil {
			return err
		}
		row.ID = saved.ID
		return h.replaceRoleMenus(c, tx, row.ID, req.MenuIDs)
	}); err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "角色代码已存在")
			return
		}
		writeError(c, err)
		return
	}
	row.IsAgent = boolToInt8(strings.HasPrefix(row.Code, "region"))
	row.Rules = idsToCSV(req.MenuIDs)
	response.OK(c, row)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "角色 ID 错误")
		return
	}
	var req roleSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if strings.TrimSpace(req.RoleName) == "" {
		response.Fail(c, http.StatusBadRequest, "角色名称必填")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_role").Where("id = ?", id).Updates(map[string]any{"name": strings.TrimSpace(req.RoleName), "status": req.Status}).Error; err != nil {
			return err
		}
		return h.replaceRoleMenus(c, tx, id, req.MenuIDs)
	}); err != nil {
		writeError(c, err)
		return
	}
	var row roleListRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,code,name,status").Where("id = ?", id).Scan(&row).Error; err != nil {
		writeError(c, err)
		return
	}
	row.Rules = idsToCSV(req.MenuIDs)
	row.IsAgent = boolToInt8(strings.HasPrefix(row.Code, "region"))
	response.OK(c, row)
}

func (h *Handler) replaceRoleMenus(c *gin.Context, tx *gorm.DB, roleID uint64, menuIDs []uint64) error {
	if err := tx.WithContext(c.Request.Context()).Where("role_id = ?", roleID).Delete(&roleMenu{}).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}
	rows := make([]roleMenu, 0, len(menuIDs))
	for _, id := range menuIDs {
		rows = append(rows, roleMenu{RoleID: roleID, MenuID: id})
	}
	return tx.WithContext(c.Request.Context()).Create(&rows).Error
}

type roleMenu struct{ RoleID, MenuID uint64 }

func (roleMenu) TableName() string { return "qixi_crm_a_role_menu" }
func (h *Handler) menuIDs(c *gin.Context, roleID uint64) (string, error) {
	var ids []uint64
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role_menu").Where("role_id = ?", roleID).Order("menu_id ASC").Pluck("menu_id", &ids).Error
	return idsToCSV(ids), err
}

type menuTreeRow struct {
	ID       uint64         `gorm:"column:id" json:"menu_id"`
	ParentID uint64         `gorm:"column:parent_id" json:"pid"`
	Title    string         `gorm:"column:title" json:"menu_name"`
	Children []*menuTreeRow `json:"children,omitempty"`
}

func (h *Handler) MenuTree(c *gin.Context) {
	var rows []*menuTreeRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Select("id,parent_id,title").Where("status = 1").Order("sort ASC,id ASC").Scan(&rows).Error; err != nil {
		writeError(c, err)
		return
	}
	byID := make(map[uint64]*menuTreeRow, len(rows))
	roots := make([]*menuTreeRow, 0)
	for _, row := range rows {
		byID[row.ID] = row
	}
	for _, row := range rows {
		if row.ParentID == 0 {
			roots = append(roots, row)
		} else if parent := byID[row.ParentID]; parent != nil {
			parent.Children = append(parent.Children, row)
		}
	}
	response.OK(c, gin.H{"list": roots})
}

func normalizedRoleCodes(req adminSaveRequest) []string {
	if len(req.RoleCodes) > 0 {
		return uniqueStrings(req.RoleCodes)
	}
	return uniqueStrings(strings.Split(req.Roles, ","))
}
func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, raw := range values {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
func parseIDs(raw string) ([]uint64, error) {
	parts := strings.Split(raw, ",")
	result := make([]uint64, 0, len(parts))
	seen := map[uint64]struct{}{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil || value == 0 {
			return nil, errors.New("区域 ID 格式错误")
		}
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result, nil
}
func pageParams(c *gin.Context, defaultLimit int) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func parseID(raw string) (uint64, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
func boolToInt8(value bool) int8 {
	if value {
		return 1
	}
	return 0
}
func idsToCSV(ids []uint64) string {
	values := make([]string, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			values = append(values, strconv.FormatUint(id, 10))
		}
	}
	return strings.Join(values, ",")
}
func validRoleCode(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 32 {
		return false
	}
	for _, char := range value {
		if !(char == '_' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}
func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
