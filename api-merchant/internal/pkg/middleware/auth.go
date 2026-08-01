package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
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
		raw := c.GetHeader("Authori-zation")
		if raw == "" {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		token, err := authjwt.BearerToken(raw)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "认证格式错误")
			c.Abort()
			return
		}
		claims, err := mgr.ParseExpect(token, portal, authjwt.TokenAccess)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if (portal == authjwt.PortalMerchant || portal == authjwt.PortalOpen || portal == authjwt.PortalManager) && (claims.MerchantID == 0 || claims.MerID != claims.MerchantID) {
			response.Fail(c, http.StatusForbidden, "缺少商户上下文")
			c.Abort()
			return
		}
		if portal == authjwt.PortalApp && (claims.UID == 0 || claims.Scope != authjwt.ScopeCUser || claims.PrincipalType != authjwt.PrincipalCUser || claims.PrincipalID != claims.UID || claims.ClientPlatform == "") {
			response.Fail(c, http.StatusForbidden, "缺少用户上下文")
			c.Abort()
			return
		}
		c.Set(CtxClaimsKey, claims)
		c.Set(CtxAdminID, claims.AdminID)
		c.Set(CtxMerID, claims.MerchantID)
		c.Set(CtxStoreID, claims.StoreID)
		c.Set(CtxStoreAppID, claims.MerchantAppID)
		c.Set(CtxUID, claims.UID)
		c.Next()
	}
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
