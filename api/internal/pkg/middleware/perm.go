package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

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
