package middleware

import (
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireActiveStoreSession rejects a token after its account has been disabled
// or its auth version has changed (for example after a password change).
func RequireActiveStoreSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if db == nil || claims == nil || claims.Scope != authjwt.ScopeStoreConsole || claims.AdminID == 0 || claims.StoreID == 0 {
			response.Fail(c, http.StatusUnauthorized, "店铺登录已失效")
			c.Abort()
			return
		}
		var row struct {
			ID          uint64 `gorm:"column:id"`
			AuthVersion uint64 `gorm:"column:auth_version"`
		}
		err := db.WithContext(c.Request.Context()).Table("qixi_crm_m_account").
			Select("id, auth_version").
			Where("id = ? AND store_id = ? AND status = 1", claims.AdminID, claims.StoreID).
			Take(&row).Error
		if err != nil || row.ID == 0 || row.AuthVersion != claims.IdentityVersion {
			response.Fail(c, http.StatusUnauthorized, "店铺登录已失效")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireStorePermission uses the new qixi_crm_m_ menu/role mapping. It never
// falls back to legacy qixi_m_* permission data.
func RequireStorePermission(db *gorm.DB, code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFrom(c)
		if db == nil || claims == nil || claims.Scope != authjwt.ScopeStoreConsole || strings.TrimSpace(code) == "" || len(claims.Roles) == 0 {
			response.Fail(c, http.StatusForbidden, "没有该操作权限")
			c.Abort()
			return
		}
		roles := make([]string, 0, len(claims.Roles))
		for _, role := range claims.Roles {
			if role = strings.TrimSpace(role); role != "" {
				roles = append(roles, role)
			}
		}
		var count int64
		err := db.WithContext(c.Request.Context()).Table("qixi_crm_m_role_menu AS rm").
			Joins("INNER JOIN qixi_crm_m_menu AS m ON m.id = rm.menu_id").
			Where("rm.role_code IN ? AND m.code = ? AND m.is_menu = 2 AND m.status = 1", roles, code).
			Count(&count).Error
		if err != nil || count == 0 {
			response.Fail(c, http.StatusForbidden, "没有该操作权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
