package merchanttype

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMerchantTypeIntegrationCRUDAndMargin(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_MERCHANT_TYPE_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_MERCHANT_TYPE_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		var ids []uint
		db.Table(typeTable).Where("name LIKE ?", "中文类型验收%").Pluck("id", &ids)
		if len(ids) > 0 {
			db.Table(typeMenuTable).Where("merchant_type_id IN ?", ids).Delete(nil)
		}
		db.Where("name LIKE ?", "中文类型验收%").Delete(&record{})
	}()
	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/types", h.Create)
	r.GET("/types/:id", h.Get)
	r.PUT("/types/:id", h.Update)
	r.PUT("/types/:id/remark", h.UpdateRemark)
	r.PUT("/types/:id/status", h.UpdateStatus)
	r.DELETE("/types/:id", h.Delete)
	call := func(method, path string, v any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(v)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := call(http.MethodPost, "/types", gin.H{"name": "中文类型验收非法", "type_info": "虚构", "is_margin": true, "margin": 0}); w.Code != 400 {
		t.Fatalf("zero margin=%d", w.Code)
	}
	if w := call(http.MethodPost, "/types", gin.H{"name": "中文类型验收分以下金额", "type_info": "虚构", "is_margin": true, "margin": 0.001}); w.Code != http.StatusBadRequest {
		t.Fatalf("fractional margin=%d", w.Code)
	}
	w := call(http.MethodPost, "/types", gin.H{"name": "中文类型验收直营网", "type_info": "虚构中文类型", "is_margin": true, "margin": 199.5, "description": "虚构中文保证金", "remark": "中文备注", "menu_codes": []string{"order.read", "order.read", "product.read"}})
	if w.Code != 200 {
		t.Fatalf("create=%d %s", w.Code, w.Body.String())
	}
	var out struct {
		Data struct {
			ID        uint     `json:"id"`
			MenuCodes []string `json:"menu_codes"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &out)
	if out.Data.ID == 0 || len(out.Data.MenuCodes) != 2 {
		t.Fatal("menu normalization failed")
	}
	if w = call(http.MethodPost, "/types", gin.H{"name": "中文类型验收直营网", "type_info": "重复", "description": "重复"}); w.Code != http.StatusConflict {
		t.Fatalf("duplicate=%d %s", w.Code, w.Body.String())
	}
	idPath := "/types/" + strconv.FormatUint(uint64(out.Data.ID), 10)
	if w = call(http.MethodGet, idPath, nil); w.Code != http.StatusOK {
		t.Fatalf("get=%d", w.Code)
	}
	if w = call(http.MethodPut, idPath+"/remark", gin.H{"remark": "中文类型验收内部备注"}); w.Code != http.StatusOK {
		t.Fatalf("remark=%d", w.Code)
	}
	if w = call(http.MethodPut, idPath+"/status", gin.H{"enabled": false}); w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}
	if w = call(http.MethodPut, "/types/999999999", gin.H{"name": "中文类型验收不存在", "description": "不存在"}); w.Code != http.StatusNotFound {
		t.Fatalf("missing update=%d", w.Code)
	}
	if w = call(http.MethodPut, idPath, gin.H{"name": "中文类型验收旗舰", "type_info": "虚构中文类型", "is_margin": false, "margin": 0, "description": "更新", "remark": "更新备注", "menu_codes": []string{"marketing.read"}}); w.Code != 200 {
		t.Fatalf("update=%d", w.Code)
	}
	if w = call(http.MethodDelete, "/types/999999999", nil); w.Code != http.StatusNotFound {
		t.Fatalf("missing delete=%d", w.Code)
	}
	if w = call(http.MethodDelete, idPath, nil); w.Code != 200 {
		t.Fatalf("delete=%d", w.Code)
	}
	var count int64
	db.Table(typeMenuTable).Where("merchant_type_id=?", out.Data.ID).Count(&count)
	if count != 0 {
		t.Fatal("menu links remain")
	}
}
