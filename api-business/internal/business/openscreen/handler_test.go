package openscreen

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterExposesOpenScreenRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	for _, route := range r.Routes() {
		if route.Method == "GET" && route.Path == "/open-screen" {
			return
		}
	}
	t.Fatal("open-screen route is not registered")
}
