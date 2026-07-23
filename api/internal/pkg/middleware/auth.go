package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

const (
	CtxClaimsKey = "auth_claims"
	CtxAdminID   = "admin_id"
	CtxMerID     = "mer_id"
	CtxUID       = "uid"
)

func JWTRequired(mgr *authjwt.Manager, portal string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer"))
		token = strings.TrimSpace(token)
		if token == "" {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		claims, err := mgr.ParseExpect(token, portal, authjwt.TokenAccess)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if (portal == authjwt.PortalMerchant || portal == authjwt.PortalOpen || portal == authjwt.PortalManager || portal == authjwt.PortalService) && claims.MerID == 0 {
			response.Fail(c, http.StatusForbidden, "缺少商户上下文")
			c.Abort()
			return
		}
		if portal == authjwt.PortalApp && claims.UID == 0 {
			response.Fail(c, http.StatusForbidden, "缺少用户上下文")
			c.Abort()
			return
		}
		c.Set(CtxClaimsKey, claims)
		c.Set(CtxAdminID, claims.AdminID)
		c.Set(CtxMerID, claims.MerID)
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

func UID(c *gin.Context) uint {
	v, _ := c.Get(CtxUID)
	id, _ := v.(uint)
	return id
}
