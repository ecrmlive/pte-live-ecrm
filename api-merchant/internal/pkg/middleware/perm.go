package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
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

// RestrictRegionConsole 在区域数据范围尚未迁移到 qixi_crm_a_data_scope 前，阻止
// 区域角色调用遗留业务 handler，避免其绕过区域隔离读取全量数据。认证与菜单接口保留，
// 业务接口在完成 repository 迁移后再按数据范围逐项开放。
func RestrictRegionConsole() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || !hasRole(claims.Roles, "region") {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/auth/") ||
			(strings.HasPrefix(c.Request.URL.Path, "/api/platform/v1/customer-service/") && hasRole(claims.Roles, "customer_service")) {
			c.Next()
			return
		}
		response.Fail(c, http.StatusForbidden, "区域数据范围迁移中，当前接口暂未开放")
		c.Abort()
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
				strings.HasPrefix(path, "/api/platform/v1/finance/withdraws"))
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
