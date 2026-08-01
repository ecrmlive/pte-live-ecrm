package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
)

const (
	CtxClaimsKey  = "auth_claims"
	CtxAdminID    = "admin_id"
	CtxMerID      = "mer_id"
	CtxStoreID    = "store_id"
	CtxStoreAppID = "store_app_id"
	CtxUID        = "uid"
)

func JWTRequired(mgr *authjwt.Manager, portal string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authenticate(c, mgr, portal, true) {
			return
		}
		c.Next()
	}
}

// JWTOptional 为可匿名浏览的接口补充已登录用户的只读上下文。只有客户端
// 传入唯一允许的 Authori-zation: Bearer 头时才校验；不传头保持匿名访问。
// 无效令牌不能降级为匿名，避免前端误把已失效会话当成新访客继续操作。
func JWTOptional(mgr *authjwt.Manager, portal string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authori-zation") == "" {
			c.Next()
			return
		}
		if !authenticate(c, mgr, portal, false) {
			return
		}
		c.Next()
	}
}

func authenticate(c *gin.Context, mgr *authjwt.Manager, portal string, required bool) bool {
	raw := c.GetHeader("Authori-zation")
	if raw == "" {
		if required {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return false
		}
		return true
	}
	token, err := authjwt.BearerToken(raw)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "认证格式错误")
		c.Abort()
		return false
	}
	claims, err := mgr.ParseExpect(token, portal, authjwt.TokenAccess)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "登录已失效")
		c.Abort()
		return false
	}
	if (portal == authjwt.PortalMerchant || portal == authjwt.PortalOpen || portal == authjwt.PortalManager) && (claims.MerchantID == 0 || claims.MerID != claims.MerchantID) {
		response.Fail(c, http.StatusForbidden, "缺少商户上下文")
		c.Abort()
		return false
	}
	if portal == authjwt.PortalApp && (claims.UID == 0 || claims.Scope != authjwt.ScopeCUser || claims.PrincipalType != authjwt.PrincipalCUser || claims.PrincipalID != claims.UID || claims.ClientPlatform == "") {
		response.Fail(c, http.StatusForbidden, "缺少用户上下文")
		c.Abort()
		return false
	}
	c.Set(CtxClaimsKey, claims)
	c.Set(CtxAdminID, claims.AdminID)
	c.Set(CtxMerID, claims.MerchantID)
	c.Set(CtxStoreID, claims.StoreID)
	c.Set(CtxStoreAppID, claims.MerchantAppID)
	c.Set(CtxUID, claims.UID)
	return true
}

func ClaimsFrom(c *gin.Context) *authjwt.Claims {
	v, ok := c.Get(CtxClaimsKey)
	if !ok {
		return nil
	}
	claims, _ := v.(*authjwt.Claims)
	return claims
}

func AdminID(c *gin.Context) uint {
	v, _ := c.Get(CtxAdminID)
	id, _ := v.(uint)
	return id
}

func MerID(c *gin.Context) uint {
	v, _ := c.Get(CtxMerID)
	id, _ := v.(uint)
	return id
}

func StoreID(c *gin.Context) uint {
	v, _ := c.Get(CtxStoreID)
	id, _ := v.(uint)
	return id
}

func StoreAppID(c *gin.Context) string {
	v, _ := c.Get(CtxStoreAppID)
	appID, _ := v.(string)
	return appID
}

func UID(c *gin.Context) uint {
	v, _ := c.Get(CtxUID)
	id, _ := v.(uint)
	return id
}
