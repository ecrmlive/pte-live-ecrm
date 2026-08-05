package nativestaff

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStaffConsoleRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(nil).Register(r)
	wants := map[string]bool{"GET /setting/staff": false, "POST /setting/staff": false, "PUT /setting/staff/:id": false, "DELETE /setting/staff/:id": false}
	for _, route := range r.Routes() {
		if _, ok := wants[route.Method+" "+route.Path]; ok {
			wants[route.Method+" "+route.Path] = true
		}
	}
	for route, present := range wants {
		if !present {
			t.Errorf("missing route %s", route)
		}
	}
}

func TestChineseStaffFixtureValidation(t *testing.T) {
	fixture := saveRequest{Account: "qixi_clerk_demo", Password: "Qixi1234", Nickname: "七禧杭州店店员", Phone: "13800000001"}
	if !validCreate(fixture) {
		t.Fatal("中文员工模拟数据应通过创建校验")
	}
	shortPassword := fixture
	shortPassword.Password = "1234567"
	if validCreate(shortPassword) {
		t.Fatal("短密码不应通过创建校验")
	}
}
