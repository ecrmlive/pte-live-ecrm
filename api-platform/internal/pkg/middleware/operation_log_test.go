package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShouldAuditAdminMutation(t *testing.T) {
	for _, tc := range []struct {
		method string
		status int
		want   bool
	}{
		{http.MethodPost, http.StatusOK, true},
		{http.MethodPut, http.StatusNoContent, true},
		{http.MethodDelete, http.StatusCreated, true},
		{http.MethodGet, http.StatusOK, false},
		{http.MethodPost, http.StatusBadRequest, false},
		{http.MethodPatch, http.StatusInternalServerError, false},
	} {
		if got := shouldAuditAdminMutation(tc.method, tc.status); got != tc.want {
			t.Fatalf("shouldAuditAdminMutation(%s,%d)=%v, want %v", tc.method, tc.status, got, tc.want)
		}
	}
}

func TestOperationResourceUsesRouteAndParam(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	c.Params = gin.Params{{Key: "coupon_id", Value: "3002"}}
	resource, id := operationResource("/api/platform/v1/user-list/:id/coupons/:coupon_id/issue", c)
	if resource != "user-list" || id != "3002" {
		t.Fatalf("resource=%q id=%q", resource, id)
	}
	if got := operationAction(http.MethodPost, "/api/platform/v1/coupons"); got != "POST /api/platform/v1/coupons" {
		t.Fatalf("action=%q", got)
	}
}

func TestBoundedOperationLogValue(t *testing.T) {
	if got := boundedOperationLogValue("  platform,operations  ", 8); got != "platform" {
		t.Fatalf("bounded value=%q", got)
	}
}
