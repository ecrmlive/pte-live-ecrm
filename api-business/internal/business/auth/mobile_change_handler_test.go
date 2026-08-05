package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestChangeMobileRouteIsAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	public := gin.New()
	authed := gin.New()
	NewHandler(nil, nil, nil).Register(public, authed)
	for _, route := range public.Routes() {
		if route.Method == "POST" && route.Path == "/auth/mobile/change" {
			t.Fatal("mobile change must not be registered on public routes")
		}
	}
	for _, route := range authed.Routes() {
		if route.Method == "POST" && route.Path == "/auth/mobile/change" {
			return
		}
	}
	t.Fatal("authenticated mobile change route not registered")
}
