package address

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetRejectsInvalidAddressID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(nil).Register(router)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/address/not-a-number", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestToResponseKeepsChineseAddressFields(t *testing.T) {
	got := toResponse(row{
		ID: 7, Recipient: "李明", Mobile: "13900000000", Province: "浙江省", City: "杭州市",
		District: "西湖区", Detail: "文三路 88 号", IsDefault: 1,
	})
	if got["address_id"] != uint64(7) || got["real_name"] != "李明" || got["city"] != "杭州市" || got["is_default"] != int8(1) {
		t.Fatalf("unexpected address response: %#v", got)
	}
}
