package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireAdminConsole 拒绝旧平台令牌进入新统一后台。新后台的角色、菜单与数据范围
// 全部来自 qixi_crm_a_*，不能再用只带 portal 的旧身份令牌混用。
func RequireAdminConsole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || claims.Scope != authjwt.ScopeAdminConsole || claims.PrincipalType != authjwt.PrincipalAdminUser || claims.AdminID == 0 || claims.PrincipalID != claims.AdminID || len(claims.Roles) == 0 {
			response.Fail(c, http.StatusUnauthorized, "统一后台登录已失效")
			c.Abort()
			return
		}
		c.Next()
	}
}

type adminSessionVersion struct {
	Status           int8   `gorm:"column:status"`
	AuthVersion      uint64 `gorm:"column:auth_version"`
	DataScopeVersion uint64 `gorm:"column:data_scope_version"`
}

func matchesAdminSession(claims *authjwt.Claims, session adminSessionVersion) bool {
	return claims != nil && session.Status == 1 && session.AuthVersion == claims.IdentityVersion && session.DataScopeVersion == claims.DataScopeVersion
}

// RequireAdminSession makes password resets, account disablement and data-scope
// changes effective on every unified-console request rather than only when a
// browser happens to call /auth/me or refresh its token.
func RequireAdminSession(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || claims.AdminID == 0 || adminDB == nil {
			response.Fail(c, http.StatusUnauthorized, "统一后台登录已失效")
			c.Abort()
			return
		}
		var session adminSessionVersion
		err := adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
			Select("status,auth_version,data_scope_version").Where("id = ?", claims.AdminID).Take(&session).Error
		if err != nil || !matchesAdminSession(claims, session) {
			response.Fail(c, http.StatusUnauthorized, "统一后台登录已失效")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdminRoles enforces a unified-admin role at the HTTP boundary. Hidden
// Vben menus are not an authorization mechanism for sensitive settings.
func RequireAdminRoles(expected ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil {
			response.Fail(c, http.StatusUnauthorized, "统一后台登录已失效")
			c.Abort()
			return
		}
		for _, role := range expected {
			if hasRole(claims.Roles, role) {
				c.Next()
				return
			}
		}
		response.Fail(c, http.StatusForbidden, "无权执行该后台操作")
		c.Abort()
	}
}

// RequireAdminMenu verifies a button code against the unified qixi_crm_a_ RBAC
// tables. It must be used by new unified-console handlers instead of legacy
// qixi_m_admin_system_menu checks.
func RequireAdminMenu(adminDB *gorm.DB, code string) gin.HandlerFunc {
	return RequireAdminMenuAny(adminDB, code)
}

// RequireAdminMenuAny allows the request when the admin holds any listed button code.
func RequireAdminMenuAny(adminDB *gorm.DB, codes ...string) gin.HandlerFunc {
	normalized := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code != "" {
			normalized = append(normalized, code)
		}
	}
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || claims.AdminID == 0 {
			response.Fail(c, http.StatusUnauthorized, "统一后台登录已失效")
			c.Abort()
			return
		}
		if adminDB == nil || len(normalized) == 0 {
			response.Fail(c, http.StatusInternalServerError, "统一后台按钮权限配置错误")
			c.Abort()
			return
		}
		var total int64
		err := adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_menu AS m").
			Joins("INNER JOIN qixi_crm_a_role_menu AS rm ON rm.menu_id = m.id").
			Joins("INNER JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = rm.role_id").
			Where("ur.admin_user_id = ? AND m.code IN ? AND m.kind = 'button' AND m.status = 1", claims.AdminID, normalized).
			Count(&total).Error
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "统一后台按钮权限校验失败")
			c.Abort()
			return
		}
		if total == 0 {
			response.Fail(c, http.StatusForbidden, "无权执行该后台操作")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireStoreConsole 仅接受由 qixi_crm_m_account 签发的店铺后台令牌。
func RequireStoreConsole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || claims.Scope != authjwt.ScopeStoreConsole || claims.PrincipalType != authjwt.PrincipalStoreAccount || claims.AdminID == 0 || claims.PrincipalID != claims.AdminID || claims.MerchantID == 0 || claims.MerID != claims.MerchantID || claims.StoreID == 0 || claims.MerchantAppID == "" || claims.StoreAppID != claims.MerchantAppID || len(claims.Roles) == 0 {
			response.Fail(c, http.StatusUnauthorized, "店铺后台登录已失效")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireStoreAppID 将店铺后台的 X-AppId 与 JWT 中的店铺应用标识强绑定。
func RequireStoreAppID() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		appID := strings.TrimSpace(c.GetHeader("X-AppId"))
		if claims == nil || claims.MerchantAppID == "" || appID == "" || appID != claims.MerchantAppID {
			response.Fail(c, http.StatusForbidden, "店铺应用标识不匹配")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RestrictRegionConsole blocks legacy handlers which do not enforce the
// qixi_crm_m_merchant.region_id scope.  Routes are opened one by one only
// after their current-table handler applies that scope server-side.
func RestrictRegionConsole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || hasRole(claims.Roles, "platform") || !hasRole(claims.Roles, "region") {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/auth/") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/dashboard/") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/analytics/") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/merchants") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/merchant-intentions") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/products") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/orders") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/refunds") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/finance/merchant-settlements") ||
			(strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/customer-service/") && hasRole(claims.Roles, "customer_service")) {
			c.Next()
			return
		}
		response.Fail(c, http.StatusForbidden, "区域数据范围迁移中，当前接口暂未开放")
		c.Abort()
	}
}

// RestrictRoleConsole is a deny-by-default boundary for non-platform unified
// console accounts.  Vben menus are presentation only: every non-platform
// role must be confined here to routes which have a corresponding server-side
// role and data-scope implementation.  New role routes are deliberately added
// one by one alongside their scope tests.
func RestrictRoleConsole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || hasRole(claims.Roles, "platform") {
			c.Next()
			return
		}
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/platform/v1/auth/") {
			c.Next()
			return
		}
		if hasRole(claims.Roles, "customer_service") && (strings.HasPrefix(path, "/api/platform/v1/customer-service/") || strings.HasPrefix(path, "/api/platform/v1/user-feedback") || strings.HasPrefix(path, "/api/platform/v1/dashboard/")) {
			c.Next()
			return
		}
		if hasRole(claims.Roles, "merchant") && (strings.HasPrefix(path, "/api/platform/v1/dashboard/") || strings.HasPrefix(path, "/api/platform/v1/analytics/") || strings.HasPrefix(path, "/api/platform/v1/merchants") || strings.HasPrefix(path, "/api/platform/v1/products") || strings.HasPrefix(path, "/api/platform/v1/orders") || strings.HasPrefix(path, "/api/platform/v1/refunds")) {
			c.Next()
			return
		}
		if hasRole(claims.Roles, "operations") && operationRoute(path) {
			c.Next()
			return
		}
		// Region is further constrained by RestrictRegionConsole, which follows
		// this middleware and lists only routes with region-aware filtering.
		if hasRole(claims.Roles, "region") {
			c.Next()
			return
		}
		response.Fail(c, http.StatusForbidden, "当前角色无权访问此统一后台接口")
		c.Abort()
	}
}

func operationRoute(path string) bool {
	return strings.HasPrefix(path, "/api/platform/v1/article/") ||
		strings.HasPrefix(path, "/api/platform/v1/articles") ||
		strings.HasPrefix(path, "/api/platform/v1/seckill/") ||
		strings.HasPrefix(path, "/api/platform/v1/combination/") ||
		strings.HasPrefix(path, "/api/platform/v1/presell/") ||
		strings.HasPrefix(path, "/api/platform/v1/assist/") ||
		strings.HasPrefix(path, "/api/platform/v1/points/") ||
		strings.HasPrefix(path, "/api/platform/v1/integral/") ||
		strings.HasPrefix(path, "/api/platform/v1/setting/balance") ||
		strings.HasPrefix(path, "/api/platform/v1/setting/user-setup") ||
		strings.HasPrefix(path, "/api/platform/v1/recharge/") ||
		strings.HasPrefix(path, "/api/platform/v1/svip/") ||
		strings.HasPrefix(path, "/api/platform/v1/coupons") ||
		strings.HasPrefix(path, "/api/platform/v1/distribution/") ||
		strings.HasPrefix(path, "/api/platform/v1/broadcast/") ||
		strings.HasPrefix(path, "/api/platform/v1/notices") ||
		strings.HasPrefix(path, "/api/platform/v1/community/") ||
		strings.HasPrefix(path, "/api/platform/v1/attachments") ||
		strings.HasPrefix(path, "/api/platform/v1/agreements") ||
		strings.HasPrefix(path, "/api/platform/v1/diy/")
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

// RestrictPlatformAgent 将区域管理员限制在已经完成服务端数据隔离的接口集合中。
// 新区域业务接入前必须先扩展该白名单，不能仅依赖前端菜单隐藏。
func RestrictPlatformAgent(svc *identity.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc == nil {
			c.Next()
			return
		}
		profile, err := svc.PlatformProfile(c.Request.Context(), AdminID(c))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if profile.IsAgent != 1 {
			c.Next()
			return
		}
		path := c.Request.URL.Path
		isScopedRead := c.Request.Method == http.MethodGet &&
			(strings.HasPrefix(path, "/api/platform/v1/orders") ||
				strings.HasPrefix(path, "/api/platform/v1/refunds") ||
				strings.HasPrefix(path, "/api/platform/v1/finance/withdraws") ||
				strings.HasPrefix(path, "/api/platform/v1/finance/merchant-settlements"))
		allowed := strings.HasPrefix(path, "/api/platform/v1/auth/") ||
			strings.HasPrefix(path, "/api/platform/v1/merchants") ||
			isScopedRead
		if !allowed {
			response.Fail(c, http.StatusForbidden, "区域账号无权访问此平台接口")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireMerchantMenu 校验商户端按钮/菜单权限（menu_id）。
func RequireMerchantMenu(svc *identity.Service, menuID uint) gin.HandlerFunc {
	return requireMenu(svc, menuID, true)
}

// RequirePlatformMenu 校验平台端按钮/菜单权限（menu_id）。
func RequirePlatformMenu(svc *identity.Service, menuID uint) gin.HandlerFunc {
	return requireMenu(svc, menuID, false)
}

func requireMenu(svc *identity.Service, menuID uint, merchant bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc == nil || menuID == 0 {
			c.Next()
			return
		}
		var err error
		if merchant {
			err = svc.RequireMerchantMenu(c.Request.Context(), AdminID(c), menuID)
		} else {
			err = svc.RequirePlatformMenu(c.Request.Context(), AdminID(c), menuID)
		}
		if err != nil {
			if errors.Is(err, identity.ErrNoPerm) {
				response.Fail(c, http.StatusForbidden, err.Error())
			} else if errors.Is(err, identity.ErrNotFound) {
				response.Fail(c, http.StatusUnauthorized, "未登录")
			} else {
				response.Fail(c, http.StatusInternalServerError, "权限校验失败")
			}
			c.Abort()
			return
		}
		c.Next()
	}
}
