package middleware

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditAdminMutation records successful unified-console mutations without
// serializing request bodies, tokens, accounts, or other sensitive payloads.
// Domain audit tables remain the source of truth for state transitions; this
// log supplies the cross-module operational trace required by the console.
func AuditAdminMutation(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if adminDB == nil || !shouldAuditAdminMutation(c.Request.Method, c.Writer.Status()) {
			return
		}
		claims := ClaimsFrom(c)
		if claims == nil || claims.AdminID == 0 {
			return
		}
		roles := append([]string(nil), claims.Roles...)
		sort.Strings(roles)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		resourceType, resourceID := operationResource(path, c)
		permissionName := operationPermissionName(resourceType)
		requestID := strings.TrimSpace(c.GetHeader("X-Request-Id"))
		if len(requestID) == 0 || len(requestID) > 64 {
			requestID = fmt.Sprintf("admin-op-%d-%d", claims.AdminID, time.Now().UnixNano())
		}
		if err := adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_operation_log").Create(map[string]any{
			"admin_user_id":   claims.AdminID,
			"role_code":       boundedOperationLogValue(strings.Join(roles, ","), 32),
			"action":          boundedOperationLogValue(operationAction(c.Request.Method, path), 128),
			"resource_type":   boundedOperationLogValue(resourceType, 64),
			"resource_id":     boundedOperationLogValue(resourceID, 64),
			"request_id":      requestID,
			"request_method":  boundedOperationLogValue(c.Request.Method, 16),
			"request_path":    boundedOperationLogValue(path, 512),
			"request_ip":      boundedOperationLogValue(c.ClientIP(), 64),
			"permission_name": boundedOperationLogValue(permissionName, 128),
		}).Error; err != nil {
			log.Printf("admin operation log write failed: %v", err)
		}
	}
}

func operationPermissionName(resource string) string {
	switch resource {
	case "setting", "maintain":
		return "系统设置"
	case "customer-service":
		return "客服管理"
	case "product", "products":
		return "商品管理"
	case "merchants":
		return "商户管理"
	case "stores":
		return "店铺管理"
	case "users", "user-list":
		return "用户管理"
	default:
		return "平台管理"
	}
}

func shouldAuditAdminMutation(method string, status int) bool {
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return false
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func operationAction(method, path string) string {
	return strings.TrimSpace(method + " " + path)
}

func operationResource(path string, c *gin.Context) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	resource := "console"
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "platform" && parts[2] == "v1" {
		resource = parts[3]
	}
	for _, key := range []string{"id", "coupon_id", "user_id", "store_id", "merchant_id"} {
		if value := strings.TrimSpace(c.Param(key)); value != "" {
			return resource, value
		}
	}
	return resource, ""
}

func boundedOperationLogValue(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}
