package auth

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errInvalidCredentials = errors.New("账号或密码错误")
	errAccountNotFound    = errors.New("账号不存在")
	errNoRole             = errors.New("账号未分配后台角色")
	errInvalidMenuScope   = errors.New("菜单身份仅支持平台、商户或区域")
	errMenuScopeForbidden = errors.New("当前账号无此身份菜单权限")
	errDisplayNameExists  = errors.New("昵称已存在")
	errDisplayNameEmpty   = errors.New("昵称不能为空")
)

type adminUser struct {
	ID               uint64         `gorm:"column:id;primaryKey"`
	Username         string         `gorm:"column:username"`
	PasswordHash     string         `gorm:"column:password_hash"`
	DisplayName      string         `gorm:"column:display_name"`
	LinkedUserID     uint64         `gorm:"column:linked_user_id"`
	AvatarURL        string         `gorm:"column:avatar_url"`
	Phone            string         `gorm:"column:phone"`
	Status           int8           `gorm:"column:status"`
	DataScopeVersion uint64         `gorm:"column:data_scope_version"`
	AuthVersion      uint64         `gorm:"column:auth_version"`
	CircleAgentID    *uint64        `gorm:"column:circle_agent_id"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (adminUser) TableName() string { return "qixi_crm_a_admin_user" }

type menu struct {
	ID        uint64 `gorm:"column:id" json:"id"`
	ParentID  uint64 `gorm:"column:parent_id" json:"parent_id"`
	Code      string `gorm:"column:code" json:"code"`
	Title     string `gorm:"column:title" json:"title"`
	Icon      string `gorm:"column:icon" json:"icon"`
	RoutePath string `gorm:"column:route_path" json:"route_path"`
	Kind      string `gorm:"column:kind" json:"kind"`
	Sort      int    `gorm:"column:sort" json:"sort"`
}

func (menu) TableName() string { return "qixi_crm_a_menu" }

type profile struct {
	ID               uint64   `json:"id"`
	AdminID          uint64   `json:"admin_id"`
	Username         string   `json:"username"`
	Account          string   `json:"account"`
	DisplayName      string   `json:"display_name"`
	RealName         string   `json:"real_name"`
	Roles            []string `json:"roles"`
	DataScopeVersion uint64   `json:"data_scope_version"`
	// is_agent 是旧 Vben 页面在角色菜单迁移期间使用的兼容字段；区域角色为 1。
	IsAgent int8 `json:"is_agent"`
}

type Handler struct {
	db  *gorm.DB
	jwt *authjwt.Manager
}

func NewHandler(db *gorm.DB, jwt *authjwt.Manager) *Handler { return &Handler{db: db, jwt: jwt} }

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	// Vben 登录页启动配置。当前不启用图形验证码时仍需返回成功，避免页面在
	// 未登录状态将“无验证码”误判为接口故障。
	public.POST("/admin/index/base", h.LoginBase)
	public.POST("/auth/refresh", h.Refresh)
	authed.POST("/auth/logout", h.Logout)
	authed.GET("/auth/me", h.Me)
	authed.GET("/auth/menus", h.Menus)
	authed.GET("/auth/permissions", h.Permissions)
	authed.PUT("/auth/password", h.ChangePassword)
}

// Logout 供前端退出时调用。令牌失效由客户端清缓存完成；此处返回成功即可，
// 避免 Vben logoutApi 命中 Gin 默认 404 并弹出 axios 原始错误文案。
func (h *Handler) Logout(c *gin.Context) {
	response.OK(c, nil)
}

func (h *Handler) LoginBase(c *gin.Context) {
	response.OK(c, gin.H{
		"codeData": gin.H{},
		"settings": gin.H{
			"admin_name":   "七禧多商户·管理中心",
			"admin_bg_img": "",
		},
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Account  string `json:"account"`
	Password string `json:"password" binding:"required"`
}

type passwordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	username := req.Username
	if username == "" {
		username = req.Account
	}
	if strings.TrimSpace(username) == "" {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := h.findUserByUsername(c, username)
	if err != nil {
		h.writeLoginLog(c, 0, username, "", false)
		writeError(c, err)
		return
	}
	if user.Status != 1 {
		h.writeLoginLog(c, user.ID, username, "", false)
		response.Fail(c, http.StatusForbidden, "账号已禁用")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.writeLoginLog(c, user.ID, username, "", false)
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}
	result, err := h.issue(c, user)
	if err != nil {
		h.writeLoginLog(c, user.ID, username, "", false)
		writeError(c, err)
		return
	}
	roles := ""
	if current, ok := result["user"].(profile); ok {
		roles = strings.Join(current.Roles, ",")
	}
	h.writeLoginLog(c, user.ID, username, roles, true)
	response.OK(c, result)
}

// writeLoginLog deliberately never writes a password, bearer token or request
// body. Audit failure must not turn a successful login into a lockout.
func (h *Handler) writeLoginLog(c *gin.Context, userID uint64, username, roleCode string, success bool) {
	if h.db == nil {
		return
	}
	_ = h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_login_log").Create(map[string]any{
		"admin_user_id": nullableAdminID(userID),
		"username":      boundedAuditValue(strings.TrimSpace(username), 64),
		"role_code":     boundedAuditValue(strings.TrimSpace(roleCode), 32),
		"success":       success,
		"ip":            boundedAuditValue(c.ClientIP(), 64),
		"user_agent":    boundedAuditValue(c.GetHeader("User-Agent"), 512),
	}).Error
}

func nullableAdminID(id uint64) any {
	if id == 0 {
		return nil
	}
	return id
}

func boundedAuditValue(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := authjwt.BearerToken(c.GetHeader("Authori-zation"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	claims, err := h.jwt.ParseExpect(refreshToken, authjwt.PortalPlatform, authjwt.TokenRefresh)
	if err != nil || claims.Scope != authjwt.ScopeAdminConsole || claims.AdminID == 0 {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	var user adminUser
	if err := h.db.WithContext(c.Request.Context()).First(&user, claims.AdminID).Error; err != nil {
		writeError(c, err)
		return
	}
	if user.Status != 1 {
		response.Fail(c, http.StatusForbidden, "账号已禁用")
		return
	}
	if user.AuthVersion != claims.IdentityVersion {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	result, err := h.issue(c, &user)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"token": result["token"]})
}

func (h *Handler) Me(c *gin.Context) {
	user, err := h.findUserByID(c, uint64(middleware.AdminID(c)))
	if err != nil {
		writeError(c, err)
		return
	}
	roles, err := h.roles(c, user.ID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, h.profile(user, roles))
}

func (h *Handler) Menus(c *gin.Context) {
	rows, scope, err := h.menus(c, uint64(middleware.AdminID(c)), "")
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"menus": rows, "menu_scope": scope})
}

func (h *Handler) Permissions(c *gin.Context) {
	rows, scope, err := h.menus(c, uint64(middleware.AdminID(c)), "button")
	if err != nil {
		writeError(c, err)
		return
	}
	permissions := make([]string, 0, len(rows))
	for _, row := range rows {
		permissions = append(permissions, row.Code)
	}
	response.OK(c, gin.H{"permissions": permissions, "menu_scope": scope})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req passwordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if len(req.NewPassword) < 8 {
		response.Fail(c, http.StatusBadRequest, "新密码至少 8 位")
		return
	}
	user, err := h.findUserByID(c, uint64(middleware.AdminID(c)))
	if err != nil {
		writeError(c, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		response.Fail(c, http.StatusUnauthorized, "原密码错误")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码处理失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Model(&adminUser{}).Where("id = ?", user.ID).
		Updates(map[string]any{"password_hash": string(hash), "auth_version": gorm.Expr("auth_version + 1")}).Error; err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) issue(c *gin.Context, user *adminUser) (gin.H, error) {
	roles, err := h.roles(c, user.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errNoRole
	}
	pair, err := h.jwt.IssueAdminConsoleWithIdentityVersion(uint(user.ID), user.Username, roles, user.DataScopeVersion, user.AuthVersion)
	if err != nil {
		return nil, err
	}
	return gin.H{"token": pair, "user": h.profile(user, roles)}, nil
}

func (h *Handler) findUserByUsername(c *gin.Context, username string) (*adminUser, error) {
	var user adminUser
	err := h.db.WithContext(c.Request.Context()).Where("username = ?", strings.TrimSpace(username)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errInvalidCredentials
	}
	return &user, err
}

func (h *Handler) findUserByID(c *gin.Context, id uint64) (*adminUser, error) {
	var user adminUser
	err := h.db.WithContext(c.Request.Context()).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errAccountNotFound
	}
	return &user, err
}

func (h *Handler) roles(c *gin.Context, userID uint64) ([]string, error) {
	var rows []string
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role AS r").
		Select("r.code").Joins("INNER JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = r.id").
		Where("ur.admin_user_id = ? AND r.status = 1", userID).Order("r.code ASC").Scan(&rows).Error
	return rows, err
}

// menuScopeForUser 返回当前后台会话的单一菜单域。平台、商户、区域是三套
// 独立菜单树，绝不能因为角色菜单表的历史关联而混合返回。
func (h *Handler) menuScopeForUser(c *gin.Context, userID uint64) (string, error) {
	var roleTypes []string
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_role AS r").
		Select("DISTINCT r.role_type").
		Joins("INNER JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = r.id").
		Where("ur.admin_user_id = ? AND r.status = 1", userID).
		Order("r.role_type ASC").Scan(&roleTypes).Error
	if err != nil {
		return "", err
	}
	return resolveMenuScope(roleTypes, c.Query("scope"))
}

func resolveMenuScope(roleTypes []string, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested != "" {
		if !isMenuScope(requested) {
			return "", errInvalidMenuScope
		}
		if !containsMenuScope(roleTypes, requested) {
			return "", errMenuScopeForbidden
		}
		return requested, nil
	}

	// 没有显式指定身份时，统一后台优先展示平台菜单；其余账号只展示各自范围。
	for _, scope := range []string{"platform", "merchant", "region"} {
		if containsMenuScope(roleTypes, scope) {
			return scope, nil
		}
	}
	return "", errNoRole
}

func isMenuScope(scope string) bool {
	return scope == "platform" || scope == "merchant" || scope == "region"
}

func containsMenuScope(scopes []string, expected string) bool {
	for _, scope := range scopes {
		if scope == expected {
			return true
		}
	}
	return false
}

func (h *Handler) menus(c *gin.Context, userID uint64, kind string) ([]menu, string, error) {
	scope, err := h.menuScopeForUser(c, userID)
	if err != nil {
		return nil, "", err
	}
	query := h.db.WithContext(c.Request.Context()).Table("qixi_crm_a_menu AS m").
		Select("DISTINCT m.id,m.parent_id,m.code,m.title,m.icon,m.route_path,m.kind,m.sort").
		Joins("INNER JOIN qixi_crm_a_role_menu AS rm ON rm.menu_id = m.id").
		Joins("INNER JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = rm.role_id").
		Where("ur.admin_user_id = ? AND m.status = 1 AND m.menu_scope = ?", userID, scope)
	if kind != "" {
		query = query.Where("m.kind = ?", kind)
	}
	var rows []menu
	if err := query.Order("m.sort ASC,m.id ASC").Scan(&rows).Error; err != nil {
		return nil, "", err
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Sort < rows[j].Sort })
	return rows, scope, nil
}

func (h *Handler) profile(user *adminUser, roles []string) profile {
	isAgent := int8(0)
	if hasRole(roles, "region") {
		isAgent = 1
	}
	return profile{
		ID: user.ID, AdminID: user.ID, Username: user.Username, Account: user.Username,
		DisplayName: user.DisplayName, RealName: user.DisplayName, Roles: roles,
		DataScopeVersion: user.DataScopeVersion, IsAgent: isAgent,
	}
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

func writeError(c *gin.Context, err error) {
	if errors.Is(err, errInvalidCredentials) {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	if errors.Is(err, errAccountNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "账号不存在")
		return
	}
	if errors.Is(err, errNoRole) {
		response.Fail(c, http.StatusForbidden, err.Error())
		return
	}
	if errors.Is(err, errMenuScopeForbidden) {
		response.Fail(c, http.StatusForbidden, err.Error())
		return
	}
	if errors.Is(err, errInvalidMenuScope) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, errDisplayNameExists) || errors.Is(err, errDisplayNameEmpty) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, "服务异常")
}

// ensureAdminDisplayNameUnique 保证统一后台登录用户昵称（display_name）在未删除账号中唯一。
func ensureAdminDisplayNameUnique(c *gin.Context, tx *gorm.DB, displayName string, excludeID uint64) error {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return errDisplayNameEmpty
	}
	q := tx.WithContext(c.Request.Context()).Table((adminUser{}).TableName()).
		Where("display_name = ? AND deleted_at IS NULL", displayName)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	var exists int64
	if err := q.Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return errDisplayNameExists
	}
	return nil
}
