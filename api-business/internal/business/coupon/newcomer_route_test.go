package coupon

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterExposesAuthenticatedNewcomerRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)
	for _, route := range router.Routes() {
		if route.Method == "GET" && route.Path == "/coupons/newcomer" {
			return
		}
	}
	t.Fatal("newcomer route is not registered")
}
