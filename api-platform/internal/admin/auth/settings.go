package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterSettings 提供统一后台自身的账号、角色和菜单管理。路由保持与 Vben 现有
// 页面一致，但所有读写仅作用于 qixi_crm_a_* 表。
func (h *Handler) RegisterSettings(authed gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	authed.GET("/setting/admins", platformOnly, h.ListAdmins)
	authed.GET("/setting/admin-regions", platformOnly, h.ListAdminRegions)
	authed.POST("/setting/admins", platformOnly, h.CreateAdmin)
	authed.PUT("/setting/admins/:id", platformOnly, h.UpdateAdmin)
	authed.DELETE("/setting/admins/:id", platformOnly, h.DeleteAdmin)
	authed.GET("/setting/roles", platformOnly, h.ListRoles)
	authed.POST("/setting/roles", platformOnly, h.CreateRole)
	authed.PUT("/setting/roles/:id", platformOnly, h.UpdateRole)
	authed.DELETE("/setting/roles/:id", platformOnly, h.DeleteRole)
	authed.GET("/setting/menus/tree", platformOnly, h.MenuTree)
	authed.GET("/setting/menus", platformOnly, h.ListMenus)
	authed.POST("/setting/menus", platformOnly, h.CreateMenu)
	authed.PUT("/setting/menus/:id", platformOnly, h.UpdateMenu)
	authed.DELETE("/setting/menus/:id", platformOnly, h.DeleteMenu)
}

type adminListRow struct {
	ID              uint64    `gorm:"column:id" json:"admin_id"`
	Username        string    `gorm:"column:username" json:"account"`
	DisplayName     string    `gorm:"column:display_name" json:"real_name"`
	LinkedUserID    uint64    `gorm:"column:linked_user_id" json:"linked_user_id"`
	AvatarURL       string    `gorm:"column:avatar_url" json:"avatar_url"`
	Phone           string    `gorm:"column:phone" json:"phone"`
	Status          int8      `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	Roles           string    `json:"roles"`
	RoleCodes       []string  `json:"role_codes"`
	RoleNames       string    `json:"role_names"`
	MerchantIDs     string    `json:"merchant_ids"`
	RegionIDs       string    `json:"region_ids"`
	RegionNames     string    `json:"region_names"`
	ServiceStoreIDs string    `json:"service_store_ids"`
	IsAgent         int8      `json:"is_agent"`
	Level           int8      `json:"level"`
	CircleID        uint64    `json:"circle_agent_id"`
}

type adminSaveRequest struct {
	Account         string   `json:"account"`
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	RealName        string   `json:"real_name"`
	LinkedUserID    uint64   `json:"linked_user_id"`
	AvatarURL       string   `json:"avatar_url"`
	Phone           string   `json:"phone"`
	Roles           string   `json:"roles"`
	RoleCodes       []string `json:"role_codes"`
	Status          int8     `json:"status"`
	MerchantIDs     string   `json:"merchant_ids"`
	RegionIDs       string   `json:"region_ids"`
	ServiceStoreIDs string   `json:"service_store_ids"`
	CircleAgentID   uint64   `json:"circle_agent_id"`
}

func (h *Handler) ListAdmins(c *gin.Context) {
	page, limit := pageParams(c, 20)
	query := h.db.WithContext(c.Request.Context()).Table((adminUser{}).TableName()).Where("deleted_at IS NULL")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR display_name LIKE ? OR phone LIKE ?", like, like, like)
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		status, err := strconv.ParseInt(raw, 10, 8)
		if err != nil || (status != 0 && status != 1) {
			response.Fail(c, http.StatusBadRequest, "账号状态参数错误")
			return
		}
		query = query.Where("status = ?", status)
	}
	if raw := strings.TrimSpace(c.Query("date_from")); raw != "" {
		value, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "开始时间格式错误")
			return
		}
		query = query.Where("created_at >= ?", value)
	}
	if raw := strings.TrimSpace(c.Query("date_to")); raw != "" {
		value, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "结束时间格式错误")
			return
		}
		query = query.Where("created_at < ?", value.AddDate(0, 0, 1))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		writeError(c, err)
		return
	}
	rows := make([]adminListRow, 0)
	if err := query.
		Select("id,username,display_name,linked_user_id,avatar_url,phone,status,created_at,COALESCE(circle_agent_id, 0) AS circle_agent_id").Order("id DESC").
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
		rows[i].RoleCodes = roles
		rows[i].RoleNames, err = h.roleNames(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].IsAgent = boolToInt8(hasRole(roles, "region"))
		rows[i].MerchantIDs, err = h.merchantIDs(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].RegionNames, err = h.regionNames(c, rows[i].RegionIDs)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].RegionIDs, err = h.regionIDs(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
		rows[i].ServiceStoreIDs, err = h.serviceStoreIDs(c, rows[i].ID)
		if err != nil {
			writeError(c, err)
			return
		}
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

// ListAdminRegions 供管理员表单选择区域；只返回有效区域名称和 ID，禁止前端手填 ID。
func (h *Handler) ListAdminRegions(c *gin.Context) {
	rows := make([]struct {
		ID   uint64 `gorm:"column:id" json:"value"`
		Name string `gorm:"column:name" json:"label"`
	}, 0)
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_region").
		Select("id,name").Where("status = 1").Order("parent_id ASC,id ASC").Scan(&rows).Error; err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
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
	displayName := strings.TrimSpace(req.RealName)
	if username == "" || len(req.Password) < 8 || len(roles) == 0 {
		response.Fail(c, http.StatusBadRequest, "账号、至少 8 位密码和角色必填")
		return
	}
	if hasRole(roles, "customer_service") && strings.TrimSpace(req.AvatarURL) == "" {
		response.Fail(c, http.StatusBadRequest, "客服头像不能为空")
		return
	}
	if displayName == "" {
		response.Fail(c, http.StatusBadRequest, "昵称不能为空")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(c, err)
		return
	}
	user := adminUser{Username: username, PasswordHash: string(hash), DisplayName: displayName, LinkedUserID: req.LinkedUserID, AvatarURL: strings.TrimSpace(req.AvatarURL), Phone: strings.TrimSpace(req.Phone), Status: req.Status, DataScopeVersion: 1}
	if user.Status == 0 {
		user.Status = 1
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := ensureAdminDisplayNameUnique(c, tx, displayName, 0); err != nil {
			return err
		}
		if err := validateCircleAgentBinding(c, tx, req.CircleAgentID, 0, roles); err != nil {
			return err
		}
		if req.CircleAgentID != 0 {
			user.CircleAgentID = &req.CircleAgentID
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return h.replaceAdminRolesAndScope(c, tx, user.ID, roles, req.MerchantIDs, req.RegionIDs, req.ServiceStoreIDs)
	}); err != nil {
		if errors.Is(err, errDisplayNameExists) || errors.Is(err, errDisplayNameEmpty) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
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
	if hasRole(roles, "customer_service") && strings.TrimSpace(req.AvatarURL) == "" {
		response.Fail(c, http.StatusBadRequest, "客服头像不能为空")
		return
	}
	displayName := strings.TrimSpace(req.RealName)
	if displayName == "" {
		response.Fail(c, http.StatusBadRequest, "昵称不能为空")
		return
	}
	if req.Password != "" && len(req.Password) < 8 {
		response.Fail(c, http.StatusBadRequest, "新密码至少 8 位")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := ensureAdminDisplayNameUnique(c, tx, displayName, user.ID); err != nil {
			return err
		}
		if err := validateCircleAgentBinding(c, tx, req.CircleAgentID, user.ID, roles); err != nil {
			return err
		}
		updates := map[string]any{"display_name": displayName, "linked_user_id": req.LinkedUserID, "avatar_url": strings.TrimSpace(req.AvatarURL), "phone": strings.TrimSpace(req.Phone), "status": req.Status, "data_scope_version": gorm.Expr("data_scope_version + 1"), "circle_agent_id": nullableCircleAgentID(req.CircleAgentID)}
		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			updates["password_hash"] = string(hash)
			updates["auth_version"] = gorm.Expr("auth_version + 1")
		}
		if err := tx.Model(&adminUser{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return err
		}
		return h.replaceAdminRolesAndScope(c, tx, user.ID, roles, req.MerchantIDs, req.RegionIDs, req.ServiceStoreIDs)
	}); err != nil {
		if errors.Is(err, errDisplayNameExists) || errors.Is(err, errDisplayNameEmpty) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
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

// DeleteAdmin is an audited logical deletion. Historical order, review and
// customer-service assignment records deliberately keep their operator ID.
// A deleted account is disabled, cannot pass the per-request session check,
// and is hidden from administrator and service-agent lists.
func (h *Handler) DeleteAdmin(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "管理员 ID 错误")
		return
	}
	operatorID := uint64(middleware.AdminID(c))
	if err := allowsAdminDeletion(operatorID, id, nil, 1); err != nil {
		response.Fail(c, http.StatusConflict, "不能删除当前登录账号")
		return
	}
	user, err := h.findUserByID(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	roles, err := h.roles(c, user.ID)
	if err != nil {
		writeError(c, err)
		return
	}
	otherActivePlatforms := int64(1)
	if hasRole(roles, "platform") {
		err = h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user AS u").
			Joins("JOIN qixi_crm_a_admin_user_role AS ur ON ur.admin_user_id = u.id").
			Joins("JOIN qixi_crm_a_role AS r ON r.id = ur.role_id").
			Where("u.id <> ? AND u.status = 1 AND u.deleted_at IS NULL AND r.code = ? AND r.status = 1", id, "platform").
			Count(&otherActivePlatforms).Error
		if err != nil {
			writeError(c, err)
			return
		}
	}
	if err := allowsAdminDeletion(operatorID, id, roles, otherActivePlatforms); err != nil {
		response.Fail(c, http.StatusConflict, err.Error())
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&adminUser{}).Where("id = ?", id).Updates(map[string]any{
			"status":             0,
			"auth_version":       gorm.Expr("auth_version + 1"),
			"data_scope_version": gorm.Expr("data_scope_version + 1"),
		}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&adminUser{}).Error
	}); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func allowsAdminDeletion(operatorID, targetID uint64, targetRoles []string, otherActivePlatforms int64) error {
	if operatorID == 0 || operatorID == targetID {
		return errors.New("不能删除当前登录账号")
	}
	if hasRole(targetRoles, "platform") && otherActivePlatforms == 0 {
		return errors.New("至少保留一个启用的平台管理员")
	}
	return nil
}

func (h *Handler) adminRow(c *gin.Context, user adminUser) adminListRow {
	roles, _ := h.roles(c, user.ID)
	roleNames, _ := h.roleNames(c, user.ID)
	merchantIDs, _ := h.merchantIDs(c, user.ID)
	regionIDs, _ := h.regionIDs(c, user.ID)
	regionNames, _ := h.regionNames(c, regionIDs)
	serviceStoreIDs, _ := h.serviceStoreIDs(c, user.ID)
	return adminListRow{ID: user.ID, Username: user.Username, DisplayName: user.DisplayName, LinkedUserID: user.LinkedUserID, AvatarURL: user.AvatarURL, Phone: user.Phone, Status: user.Status, Roles: strings.Join(roles, ","), RoleCodes: roles, RoleNames: roleNames, MerchantIDs: merchantIDs, RegionIDs: regionIDs, RegionNames: regionNames, ServiceStoreIDs: serviceStoreIDs, IsAgent: boolToInt8(hasRole(roles, "region")), CircleID: derefCircleAgentID(user.CircleAgentID)}
}

func nullableCircleAgentID(id uint64) any {
	if id == 0 {
		return nil
	}
	return id
}

func derefCircleAgentID(id *uint64) uint64 {
	if id == nil {
		return 0
	}
	return *id
}

// validateCircleAgentBinding 使区域代理和统一后台账号形成一对一、可撤销的归属关系。
// 仅审核通过的代理能绑定区域角色；其他角色不能借此字段获得区域身份。
func validateCircleAgentBinding(c *gin.Context, tx *gorm.DB, circleAgentID, excludingAdminID uint64, roles []string) error {
	isRegion := hasRole(roles, "region") && !hasRole(roles, "platform")
	if isRegion && circleAgentID == 0 {
		return errors.New("区域管理员必须关联已通过审核的区域代理")
	}
	if !isRegion && circleAgentID != 0 {
		return errors.New("仅区域管理员可以关联区域代理")
	}
	if circleAgentID == 0 {
		return nil
	}
	var agent struct {
		Status int8 `gorm:"column:status"`
	}
	if err := tx.WithContext(c.Request.Context()).Table("qixi_crm_a_business_zone_agent").Select("status").Where("circle_agent_id = ?", circleAgentID).Take(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("关联代理不存在")
		}
		return err
	}
	if agent.Status != 1 {
		return errors.New("仅已审核通过的代理可以关联区域管理员")
	}
	var existing uint64
	err := tx.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").Select("id").Where("circle_agent_id = ?", circleAgentID).Where("id <> ?", excludingAdminID).Limit(1).Scan(&existing).Error
	if err != nil {
		return err
	}
	if existing != 0 {
		return errors.New("该代理已关联其他后台账号")
	}
	return nil
}

func (h *Handler) replaceAdminRolesAndScope(c *gin.Context, tx *gorm.DB, userID uint64, codes []string, merchantIDs, regionIDs, serviceStoreIDs string) error {
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
	if err := tx.WithContext(c.Request.Context()).Where("admin_user_id = ? AND scope_type IN ?", userID, []string{"merchant", "region", "service_queue"}).Delete(&adminDataScope{}).Error; err != nil {
		return err
	}
	if hasRole(codes, "merchant") {
		values, err := parseIDs(merchantIDs, "授权商户 ID")
		if err != nil || len(values) == 0 {
			return errors.New("商户角色必须配置至少一个授权商户 ID")
		}
		payload, err := json.Marshal(merchantDataScope{MerchantIDs: values})
		if err != nil {
			return err
		}
		if err := tx.WithContext(c.Request.Context()).Create(&adminDataScope{AdminUserID: userID, ScopeType: "merchant", ScopeValue: payload, Version: 1}).Error; err != nil {
			return err
		}
	}
	if hasRole(codes, "region") {
		values, err := parseIDs(regionIDs, "区域 ID")
		if err != nil || len(values) == 0 {
			return errors.New("区域角色必须配置区域 ID")
		}
		payload, err := json.Marshal(values)
		if err != nil {
			return err
		}
		if err := tx.WithContext(c.Request.Context()).Create(&adminDataScope{AdminUserID: userID, ScopeType: "region", ScopeValue: payload, Version: 1}).Error; err != nil {
			return err
		}
	}
	if !hasRole(codes, "customer_service") {
		return nil
	}
	storeIDs, err := parseIDs(serviceStoreIDs, "客服授权店铺 ID")
	if err != nil || len(storeIDs) == 0 {
		return errors.New("客服角色必须配置至少一个授权店铺 ID")
	}
	payload, err := json.Marshal(serviceQueueScope{StoreIDs: storeIDs})
	if err != nil {
		return err
	}
	return tx.WithContext(c.Request.Context()).Create(&adminDataScope{AdminUserID: userID, ScopeType: "service_queue", ScopeValue: payload, Version: 1}).Error
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
	values, err := parseScopeIDs(row.ScopeValue, false)
	if err != nil {
		return "", err
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.FormatUint(value, 10))
	}
	return strings.Join(parts, ","), nil
}

func (h *Handler) roleNames(c *gin.Context, userID uint64) (string, error) {
	var names []string
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role AS r").
		Select("r.name").Joins("INNER JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = r.id").
		Where("ur.admin_user_id = ?", userID).Order("r.id ASC").Scan(&names).Error
	return strings.Join(names, "、"), err
}

func (h *Handler) regionNames(c *gin.Context, rawIDs string) (string, error) {
	ids, err := parseIDs(rawIDs, "区域 ID")
	if err != nil || len(ids) == 0 {
		return "", err
	}
	var rows []struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_region").Select("id,name").Where("id IN ?", ids).Scan(&rows).Error; err != nil {
		return "", err
	}
	nameByID := make(map[uint64]string, len(rows))
	for _, row := range rows {
		nameByID[row.ID] = row.Name
	}
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if name := strings.TrimSpace(nameByID[id]); name != "" {
			names = append(names, name)
		}
	}
	return strings.Join(names, "、"), nil
}

type serviceQueueScope struct {
	StoreIDs []uint64 `json:"store_ids"`
}

type merchantDataScope struct {
	MerchantIDs []uint64 `json:"merchant_ids"`
}

func (h *Handler) merchantIDs(c *gin.Context, userID uint64) (string, error) {
	var row adminDataScope
	err := h.db.WithContext(c.Request.Context()).Where("admin_user_id = ? AND scope_type = ?", userID, "merchant").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var scope merchantDataScope
	if err := json.Unmarshal(row.ScopeValue, &scope); err != nil {
		return "", err
	}
	return idsToCSV(scope.MerchantIDs), nil
}

func (h *Handler) serviceStoreIDs(c *gin.Context, userID uint64) (string, error) {
	var row adminDataScope
	err := h.db.WithContext(c.Request.Context()).Where("admin_user_id = ? AND scope_type = ?", userID, "service_queue").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	values, err := parseScopeIDs(row.ScopeValue, true)
	if err != nil {
		return "", err
	}
	return idsToCSV(values), nil
}

func parseScopeIDs(raw json.RawMessage, serviceQueue bool) ([]uint64, error) {
	if serviceQueue {
		var scope serviceQueueScope
		if err := json.Unmarshal(raw, &scope); err == nil && scope.StoreIDs != nil {
			return scope.StoreIDs, nil
		}
		// 兼容早期本地客服夹具使用的 JSON 数组；写入时统一使用对象结构。
		var legacyStoreIDs []uint64
		if err := json.Unmarshal(raw, &legacyStoreIDs); err != nil {
			return nil, err
		}
		return legacyStoreIDs, nil
	}
	var values []uint64
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	return values, nil
}

type roleListRow struct {
	ID        uint64    `gorm:"column:id" json:"role_id"`
	Code      string    `gorm:"column:code" json:"code"`
	Name      string    `gorm:"column:name" json:"role_name"`
	Status    int8      `gorm:"column:status" json:"status"`
	RoleType  string    `gorm:"column:role_type" json:"role_type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Rules     string    `json:"rules"`
}
type roleSaveRequest struct {
	Code     string   `json:"code"`
	RoleName string   `json:"role_name"`
	RoleType string   `json:"role_type"`
	Status   int8     `json:"status"`
	MenuIDs  []uint64 `json:"menu_ids"`
}

func (h *Handler) ListRoles(c *gin.Context) {
	page, limit := pageParams(c, 10)
	roleType := strings.TrimSpace(c.Query("role_type"))
	if roleType != "" && !validRoleType(roleType) {
		response.Fail(c, http.StatusBadRequest, "身份类型错误")
		return
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	status := strings.TrimSpace(c.Query("status"))
	query := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role")
	if roleType != "" {
		query = query.Where("role_type = ?", roleType)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		parsed, err := strconv.ParseInt(status, 10, 8)
		if err != nil || (parsed != 0 && parsed != 1) {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		query = query.Where("status = ?", parsed)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		writeError(c, err)
		return
	}
	rows := make([]roleListRow, 0)
	if err := query.Select("id,code,name,status,role_type,created_at,updated_at").Order("id ASC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
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
	}
	response.OK(c, gin.H{"list": rows, "total": total})
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req roleSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if strings.TrimSpace(req.RoleName) == "" || !validRoleType(req.RoleType) {
		response.Fail(c, http.StatusBadRequest, "身份名称和身份类型必填")
		return
	}
	code := strings.TrimSpace(req.Code)
	if code == "" {
		code = "custom_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if !validRoleCode(code) {
		response.Fail(c, http.StatusBadRequest, "身份代码仅支持小写字母、数字和下划线")
		return
	}
	if err := h.validateRoleMenuScope(c, req.RoleType, req.MenuIDs); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	row := roleListRow{Code: code, Name: strings.TrimSpace(req.RoleName), RoleType: req.RoleType, Status: req.Status}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_role").Create(map[string]any{"code": row.Code, "name": row.Name, "role_type": row.RoleType, "status": row.Status}).Error; err != nil {
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
	var current roleListRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,role_type").Where("id = ?", id).Take(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "身份不存在")
			return
		}
		writeError(c, err)
		return
	}
	if err := h.validateRoleMenuScope(c, current.RoleType, req.MenuIDs); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
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
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,code,name,status,role_type,created_at,updated_at").Where("id = ?", id).Scan(&row).Error; err != nil {
		writeError(c, err)
		return
	}
	row.Rules = idsToCSV(req.MenuIDs)
	response.OK(c, row)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "身份 ID 错误")
		return
	}
	var row roleListRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role").Select("id,code,name").Where("id = ?", id).Scan(&row).Error; err != nil {
		writeError(c, err)
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "身份不存在")
		return
	}
	if builtInRoleCode(row.Code) {
		response.Fail(c, http.StatusConflict, "系统预置身份不能删除")
		return
	}
	var assigned int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user_role").Where("role_id = ?", id).Count(&assigned).Error; err != nil {
		writeError(c, err)
		return
	}
	if assigned > 0 {
		response.Fail(c, http.StatusConflict, "身份已关联后台用户，不能删除")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_role_menu").Where("role_id = ?", id).Delete(&roleMenu{}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_a_role").Where("id = ?", id).Delete(&roleListRow{}).Error
	}); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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

func (h *Handler) validateRoleMenuScope(c *gin.Context, scope string, menuIDs []uint64) error {
	if len(menuIDs) == 0 {
		return nil
	}
	menuScope, err := validMenuScope(scope)
	if err != nil {
		return err
	}
	unique := make(map[uint64]struct{}, len(menuIDs))
	for _, id := range menuIDs {
		if id == 0 {
			return errors.New("菜单 ID 错误")
		}
		unique[id] = struct{}{}
	}
	ids := make([]uint64, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Where("menu_scope = ? AND id IN ?", menuScope, ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("所选菜单与身份归属不一致")
	}
	return nil
}

type roleMenu struct{ RoleID, MenuID uint64 }

func (roleMenu) TableName() string { return "qixi_crm_a_role_menu" }
func (h *Handler) menuIDs(c *gin.Context, roleID uint64) (string, error) {
	var ids []uint64
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role_menu").Where("role_id = ?", roleID).Order("menu_id ASC").Pluck("menu_id", &ids).Error
	return idsToCSV(ids), err
}

type menuTreeRow struct {
	ID       uint64 `gorm:"column:id" json:"menu_id"`
	ParentID uint64 `gorm:"column:parent_id" json:"pid"`
	Title    string `gorm:"column:title" json:"menu_name"`
	// 菜单树由查询结果在内存中组装；该字段不是数据库关联，必须排除 GORM 的字段解析。
	Children []*menuTreeRow `gorm:"-" json:"children,omitempty"`
}

func (h *Handler) MenuTree(c *gin.Context) {
	scope, err := validMenuScope(c.DefaultQuery("scope", "platform"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	var rows []*menuTreeRow
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Select("id,parent_id,title").Where("status = 1 AND menu_scope = ?", scope).Order("sort ASC,id ASC").Scan(&rows).Error; err != nil {
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

type menuListRow struct {
	CreatedAt time.Time `json:"created_at"`
	Icon      string    `json:"icon"`
	IsAgent   int8      `json:"is_agent"`
	IsMenu    int8      `json:"is_menu"`
	IsShow    int8      `json:"is_show"`
	Kind      string    `json:"kind"`
	MenuID    uint64    `json:"menu_id"`
	MenuName  string    `json:"menu_name"`
	MenuScope string    `json:"menu_scope"`
	PID       uint64    `json:"pid"`
	Path      string    `json:"path"`
	Route     string    `json:"route"`
	Sort      int       `json:"sort"`
}

type menuRecord struct {
	Code      string    `gorm:"column:code"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Icon      string    `gorm:"column:icon"`
	ID        uint64    `gorm:"column:id"`
	Kind      string    `gorm:"column:kind"`
	MenuScope string    `gorm:"column:menu_scope"`
	ParentID  uint64    `gorm:"column:parent_id"`
	RoutePath string    `gorm:"column:route_path"`
	Sort      int       `gorm:"column:sort"`
	Status    int8      `gorm:"column:status"`
	Title     string    `gorm:"column:title"`
}

type menuSaveRequest struct {
	Code      *string `json:"code"`
	Icon      *string `json:"icon"`
	IsShow    *int8   `json:"is_show"`
	Kind      *string `json:"kind"`
	MenuName  *string `json:"menu_name"`
	MenuScope *string `json:"menu_scope"`
	ParentID  *uint64 `json:"parent_id"`
	Path      *string `json:"path"`
	Sort      *int    `json:"sort"`
}

func validMenuScope(value string) (string, error) {
	switch strings.TrimSpace(value) {
	case "platform", "merchant", "region":
		return strings.TrimSpace(value), nil
	default:
		return "", errors.New("菜单归属仅支持平台、商户或区域")
	}
}

func validMenuKind(value string) bool {
	switch value {
	case "directory", "page", "button":
		return true
	default:
		return false
	}
}

func validMenuCode(value string) bool {
	if len(value) == 0 || len(value) > 128 {
		return false
	}
	for index, char := range value {
		if (char >= 'a' && char <= 'z') || (index > 0 && ((char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-')) {
			continue
		}
		return false
	}
	return true
}

func menuRow(record menuRecord) menuListRow {
	isMenu := int8(1)
	if record.Kind == "button" {
		isMenu = 2
	}
	return menuListRow{
		CreatedAt: record.CreatedAt,
		Icon:      record.Icon,
		IsAgent:   0,
		IsMenu:    isMenu,
		IsShow:    record.Status,
		Kind:      record.Kind,
		MenuID:    record.ID,
		MenuName:  record.Title,
		MenuScope: record.MenuScope,
		PID:       record.ParentID,
		Path:      record.RoutePath,
		Route:     record.Code,
		Sort:      record.Sort,
	}
}

func (h *Handler) listMenuRecords(c *gin.Context, scope string) ([]menuRecord, error) {
	rows := make([]menuRecord, 0)
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").
		Select("id,parent_id,code,title,icon,route_path,kind,sort,status,menu_scope,created_at").
		Where("menu_scope = ?", scope).
		Order("sort ASC, id ASC").Scan(&rows).Error
	return rows, err
}

func (h *Handler) ListMenus(c *gin.Context) {
	scope, err := validMenuScope(c.DefaultQuery("scope", "platform"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	rows, err := h.listMenuRecords(c, scope)
	if err != nil {
		writeError(c, err)
		return
	}
	out := make([]menuListRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, menuRow(row))
	}
	response.OK(c, gin.H{"list": out, "total": len(out)})
}

func menuRecordFromRequest(req menuSaveRequest, base menuRecord, creating bool) (menuRecord, error) {
	next := base
	if req.Code != nil {
		next.Code = strings.TrimSpace(*req.Code)
	}
	if req.Icon != nil {
		next.Icon = strings.TrimSpace(*req.Icon)
	}
	if req.IsShow != nil {
		next.Status = *req.IsShow
	}
	if req.Kind != nil {
		next.Kind = strings.TrimSpace(*req.Kind)
	}
	if req.MenuName != nil {
		next.Title = strings.TrimSpace(*req.MenuName)
	}
	if req.MenuScope != nil {
		next.MenuScope = strings.TrimSpace(*req.MenuScope)
	}
	if req.ParentID != nil {
		next.ParentID = *req.ParentID
	}
	if req.Path != nil {
		next.RoutePath = strings.TrimSpace(*req.Path)
	}
	if req.Sort != nil {
		next.Sort = *req.Sort
	}
	if creating && req.IsShow == nil {
		next.Status = 1
	}
	if next.Title == "" || len([]rune(next.Title)) > 64 {
		return menuRecord{}, errors.New("菜单名称不能为空且不超过 64 字")
	}
	if !validMenuCode(next.Code) {
		return menuRecord{}, errors.New("菜单标识仅支持小写字母开头、数字、点、下划线和连字符")
	}
	if len(next.Icon) > 96 {
		return menuRecord{}, errors.New("菜单图标不超过 96 字")
	}
	if len(next.RoutePath) > 255 {
		return menuRecord{}, errors.New("菜单地址不超过 255 字")
	}
	if !validMenuKind(next.Kind) {
		return menuRecord{}, errors.New("菜单类型错误")
	}
	scope, err := validMenuScope(next.MenuScope)
	if err != nil {
		return menuRecord{}, err
	}
	next.MenuScope = scope
	if next.Status != 0 && next.Status != 1 {
		return menuRecord{}, errors.New("显示状态仅支持 0/1")
	}
	return next, nil
}

func (h *Handler) validateMenuParent(c *gin.Context, record menuRecord) error {
	if record.ParentID == 0 {
		return nil
	}
	if record.ParentID == record.ID {
		return errors.New("上级菜单不能选择当前菜单")
	}
	rows, err := h.listMenuRecords(c, record.MenuScope)
	if err != nil {
		return err
	}
	parents := make(map[uint64]uint64, len(rows))
	found := false
	for _, row := range rows {
		parents[row.ID] = row.ParentID
		if row.ID == record.ParentID {
			found = true
		}
	}
	if !found {
		return errors.New("上级菜单不存在或不属于当前菜单归属")
	}
	visited := make(map[uint64]struct{}, len(parents))
	for current := record.ParentID; current != 0; current = parents[current] {
		if _, exists := visited[current]; exists {
			return errors.New("菜单层级存在循环引用")
		}
		visited[current] = struct{}{}
		if current == record.ID {
			return errors.New("上级菜单不能选择当前菜单的子菜单")
		}
	}
	return nil
}

func (h *Handler) menuHasChildren(c *gin.Context, id uint64) (bool, error) {
	var count int64
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (h *Handler) readMenuRecord(c *gin.Context, id uint64) (menuRecord, error) {
	var row menuRecord
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").
		Select("id,parent_id,code,title,icon,route_path,kind,sort,status,menu_scope,created_at").
		Where("id = ?", id).Take(&row).Error
	return row, err
}

func (h *Handler) CreateMenu(c *gin.Context) {
	var req menuSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := menuRecordFromRequest(req, menuRecord{}, true)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validateMenuParent(c, row); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Create(map[string]any{
		"parent_id":  row.ParentID,
		"code":       row.Code,
		"title":      row.Title,
		"icon":       row.Icon,
		"route_path": row.RoutePath,
		"kind":       row.Kind,
		"sort":       row.Sort,
		"status":     row.Status,
		"menu_scope": row.MenuScope,
	}).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "菜单标识已存在")
			return
		}
		writeError(c, err)
		return
	}
	var created menuRecord
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").
		Select("id,parent_id,code,title,icon,route_path,kind,sort,status,menu_scope,created_at").
		Where("code = ?", row.Code).Take(&created).Error; err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, menuRow(created))
}

func (h *Handler) UpdateMenu(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "菜单 ID 错误")
		return
	}
	var req menuSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	current, err := h.readMenuRecord(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "菜单不存在")
			return
		}
		writeError(c, err)
		return
	}
	next, err := menuRecordFromRequest(req, current, false)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if next.MenuScope != current.MenuScope {
		response.Fail(c, http.StatusBadRequest, "菜单归属不能通过编辑变更")
		return
	}
	if err := h.validateMenuParent(c, next); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if next.Kind == "button" {
		hasChildren, err := h.menuHasChildren(c, id)
		if err != nil {
			writeError(c, err)
			return
		}
		if hasChildren {
			response.Fail(c, http.StatusBadRequest, "含有子菜单时不能改为按钮权限")
			return
		}
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu").Where("id = ?", id).Updates(map[string]any{
		"parent_id":  next.ParentID,
		"code":       next.Code,
		"title":      next.Title,
		"icon":       next.Icon,
		"route_path": next.RoutePath,
		"kind":       next.Kind,
		"sort":       next.Sort,
		"status":     next.Status,
	}).Error; err != nil {
		if isDuplicate(err) {
			response.Fail(c, http.StatusConflict, "菜单标识已存在")
			return
		}
		writeError(c, err)
		return
	}
	updated, err := h.readMenuRecord(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, menuRow(updated))
}

func (h *Handler) DeleteMenu(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "菜单 ID 错误")
		return
	}
	row, err := h.readMenuRecord(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "菜单不存在")
			return
		}
		writeError(c, err)
		return
	}
	hasChildren, err := h.menuHasChildren(c, id)
	if err != nil {
		writeError(c, err)
		return
	}
	if hasChildren {
		response.Fail(c, http.StatusConflict, "请先删除子菜单")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_role_menu").Where("menu_id = ?", id).Delete(&roleMenu{}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_a_menu").Where("id = ?", row.ID).Delete(&menuRecord{}).Error
	}); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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
func parseIDs(raw, label string) ([]uint64, error) {
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
			return nil, errors.New(label + "格式错误")
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

func validRoleType(value string) bool {
	switch value {
	case "platform", "merchant", "region":
		return true
	default:
		return false
	}
}

func builtInRoleCode(value string) bool {
	switch value {
	case "platform", "merchant", "region", "customer_service", "operations":
		return true
	default:
		return false
	}
}

func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
