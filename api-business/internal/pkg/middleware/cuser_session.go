package middleware

import (
	"errors"
	"net/http"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CUserSessionRequired verifies the database-backed account state after JWT
// signature validation. It makes auth_version a real session-revocation
// boundary for password changes and platform-triggered user disabling.
func CUserSessionRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if claims == nil || claims.UID == 0 || claims.IdentityVersion == 0 {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		var user struct {
			Status      int8   `gorm:"column:status"`
			AuthVersion uint64 `gorm:"column:auth_version"`
		}
		err := db.WithContext(c.Request.Context()).Table("qixi_crm_b_user").
			Select("status,auth_version").Where("id=?", claims.UID).Take(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) || err == nil && (user.Status != 1 || user.AuthVersion != claims.IdentityVersion) {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "会话校验失败")
			c.Abort()
			return
		}
		c.Next()
	}
}
