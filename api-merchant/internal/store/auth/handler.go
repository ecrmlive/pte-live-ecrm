package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type account struct {
	ID           uint64 `gorm:"column:id;primaryKey"`
	StoreID      uint64 `gorm:"column:store_id"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	RoleCode     string `gorm:"column:role_code"`
	Status       int8   `gorm:"column:status"`
	AuthVersion  uint64 `gorm:"column:auth_version"`
}

func (account) TableName() string { return "qixi_crm_m_account" }

type accountContext struct {
	AccountID    uint64 `gorm:"column:account_id"`
	StoreID      uint64 `gorm:"column:store_id"`
	MerchantID   uint64 `gorm:"column:merchant_id"`
	StoreAppID   string `gorm:"column:store_app_id"`
	IMSDKAppID   string `gorm:"column:im_sdk_app_id"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	RoleCode     string `gorm:"column:role_code"`
	Status       int8   `gorm:"column:status"`
	AuthVersion  uint64 `gorm:"column:auth_version"`
	StoreName    string `gorm:"column:store_name"`
	MerchantName string `gorm:"column:merchant_name"`
}

type Handler struct {
	db  *gorm.DB
	jwt *authjwt.Manager
}

func NewHandler(db *gorm.DB, jwt *authjwt.Manager) *Handler { return &Handler{db: db, jwt: jwt} }

func (h *Handler) Register(public, authed gin.IRoutes) {
	public.POST("/auth/login", h.Login)
	// Vben 店铺登录页启动配置。未启用图形验证码也必须返回成功，
	// 否则页面会把“无验证码”错误地当成接口不可用。
	public.POST("/shop/index/base", h.LoginBase)
	public.POST("/auth/refresh", h.Refresh)
	authed.GET("/auth/me", h.Me)
	authed.GET("/auth/menus", h.Menus)
	authed.GET("/auth/permissions", h.Permissions)
	authed.PUT("/auth/password", h.ChangePassword)
}

func (h *Handler) LoginBase(c *gin.Context) {
	response.OK(c, gin.H{
		"codeData": gin.H{},
		"settings": gin.H{
			"shop_name":     "七禧多商户·商户",
			"shop_logo_img": "",
			"shop_bg_img":   "",
		},
	})
}

type loginRequest struct {
	Account  string `json:"account" binding:"required"`
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
	ctx, err := h.findByUsername(c, req.Account)
	if err != nil {
		writeError(c, err)
		return
	}
	if ctx.Status != 1 {
		response.Fail(c, http.StatusForbidden, "账号或店铺已禁用")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ctx.PasswordHash), []byte(req.Password)); err != nil {
		response.Fail(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}
	h.writeSession(c, ctx)
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := authjwt.BearerToken(c.GetHeader("Authori-zation"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	claims, err := h.jwt.ParseExpect(refreshToken, authjwt.PortalMerchant, authjwt.TokenRefresh)
	if err != nil || claims.Scope != authjwt.ScopeStoreConsole || claims.AdminID == 0 || claims.StoreID == 0 {
		response.Fail(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}
	ctx, err := h.findByID(c, uint64(claims.AdminID))
	if err != nil {
		writeError(c, err)
		return
	}
	if ctx.Status != 1 || ctx.AuthVersion != claims.IdentityVersion || ctx.StoreID != uint64(claims.StoreID) || ctx.MerchantID != uint64(claims.MerchantID) || claims.MerID != claims.MerchantID || ctx.StoreAppID != claims.MerchantAppID || claims.StoreAppID != claims.MerchantAppID || ctx.IMSDKAppID != claims.IMSDKAppID {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		return
	}
	pair, err := h.jwt.IssueStoreConsoleWithIdentityVersion(uint(ctx.AccountID), uint(ctx.MerchantID), uint(ctx.StoreID), ctx.StoreAppID, ctx.IMSDKAppID, ctx.Username, ctx.RoleCode, ctx.AuthVersion)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair})
}

func (h *Handler) Me(c *gin.Context) {
	ctx, err := h.current(c)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, profile(ctx))
}
func (h *Handler) Menus(c *gin.Context) {
	ctx, err := h.current(c)
	if err != nil {
		writeError(c, err)
		return
	}
	rows, err := h.menuTree(c, ctx.RoleCode)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"menus": rows, "mer_id": ctx.MerchantID, "store_id": ctx.StoreID})
}
func (h *Handler) Permissions(c *gin.Context) {
	ctx, err := h.current(c)
	if err != nil {
		writeError(c, err)
		return
	}
	var codes []string
	err = h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_menu AS m").Select("m.code").Joins("INNER JOIN qixi_crm_m_role_menu AS rm ON rm.menu_id = m.id").Where("rm.role_code = ? AND m.status = 1", ctx.RoleCode).Order("m.code").Scan(&codes).Error
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"permissions": codes})
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
	ctx, err := h.current(c)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ctx.PasswordHash), []byte(req.OldPassword)); err != nil {
		response.Fail(c, http.StatusUnauthorized, "原密码错误")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码处理失败")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Model(&account{}).Where("id = ? AND store_id = ?", ctx.AccountID, ctx.StoreID).Updates(map[string]any{"password_hash": string(hash), "auth_version": gorm.Expr("auth_version + 1")}).Error; err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) writeSession(c *gin.Context, ctx *accountContext) {
	pair, err := h.jwt.IssueStoreConsoleWithIdentityVersion(uint(ctx.AccountID), uint(ctx.MerchantID), uint(ctx.StoreID), ctx.StoreAppID, ctx.IMSDKAppID, ctx.Username, ctx.RoleCode, ctx.AuthVersion)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{"token": pair, "user": profile(ctx)})
}
func (h *Handler) current(c *gin.Context) (*accountContext, error) {
	ctx, err := h.findByID(c, uint64(middleware.AdminID(c)))
	if err != nil {
		return nil, err
	}
	claims := middleware.ClaimsFrom(c)
	if claims == nil || ctx.Status != 1 || ctx.AuthVersion != claims.IdentityVersion || ctx.StoreID != uint64(middleware.StoreID(c)) || ctx.MerchantID != uint64(middleware.MerID(c)) || ctx.MerchantID != uint64(claims.MerchantID) || claims.MerID != claims.MerchantID || ctx.StoreAppID != claims.MerchantAppID || claims.StoreAppID != claims.MerchantAppID || ctx.IMSDKAppID != claims.IMSDKAppID {
		return nil, errors.New("登录已失效")
	}
	return ctx, nil
}
func (h *Handler) findByUsername(c *gin.Context, username string) (*accountContext, error) {
	return h.find(c, "a.username = ?", strings.TrimSpace(username))
}
func (h *Handler) findByID(c *gin.Context, id uint64) (*accountContext, error) {
	return h.find(c, "a.id = ?", id)
}
func (h *Handler) find(c *gin.Context, condition string, value any) (*accountContext, error) {
	var row accountContext
	// 查询用了 `a.id AS account_id` 投影，而 GORM 的 First 会按结构体主键
	// 自动追加 `ORDER BY a.account_id`；该列并不存在。这里无需排序，改用 Take。
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_account AS a").Select("a.id AS account_id,a.store_id,a.username,a.password_hash,a.role_code,a.status,a.auth_version,s.name AS store_name,s.app_id AS store_app_id,s.merchant_id,m.name AS merchant_name,COALESCE(im.sdk_app_id, '') AS im_sdk_app_id").Joins("INNER JOIN qixi_crm_m_store AS s ON s.id = a.store_id AND s.status = 1").Joins("INNER JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id AND m.status = 1").Joins("LEFT JOIN qixi_crm_m_im_sdk_app AS im ON im.merchant_id = s.merchant_id AND im.status = 'enabled' AND im.is_active = 1").Where(condition, value).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("账号或密码错误")
	}
	return &row, err
}

type menu struct {
	ID        uint64 `gorm:"column:id" json:"id"`
	ParentID  uint64 `gorm:"column:parent_id" json:"parent_id"`
	Name      string `gorm:"column:name" json:"name"`
	Path      string `gorm:"column:path" json:"path"`
	Component string `gorm:"column:component" json:"component"`
	Icon      string `gorm:"column:icon" json:"icon"`
	IsMenu    int8   `gorm:"column:is_menu" json:"is_menu"`
	IsRoute   int8   `gorm:"column:is_route" json:"is_route"`
	// Children 仅用于在内存中组装菜单树，不能让 GORM 将其识别为关联字段。
	Children []*menu `json:"children,omitempty" gorm:"-"`
}

func (h *Handler) menuTree(c *gin.Context, role string) ([]*menu, error) {
	var rows []*menu
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_m_menu AS m").Select("m.id,m.parent_id,m.name,m.path,m.component,m.icon,m.is_menu,m.is_route").Joins("INNER JOIN qixi_crm_m_role_menu AS rm ON rm.menu_id = m.id").Where("rm.role_code = ? AND m.status = 1", role).Order("m.sort,m.id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	byID := map[uint64]*menu{}
	roots := make([]*menu, 0)
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
	return roots, nil
}
func profile(ctx *accountContext) gin.H {
	return gin.H{"merchant_admin_id": ctx.AccountID, "mer_id": ctx.MerchantID, "store_id": ctx.StoreID, "store_app_id": ctx.StoreAppID, "merchant_app_id": ctx.StoreAppID, "im_sdk_app_id": ctx.IMSDKAppID, "mer_name": ctx.MerchantName, "store_name": ctx.StoreName, "account": ctx.Username, "real_name": ctx.Username, "phone": "", "roles": ctx.RoleCode}
}
func writeError(c *gin.Context, err error) {
	if strings.Contains(err.Error(), "账号") || strings.Contains(err.Error(), "登录") {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, "服务异常")
}
