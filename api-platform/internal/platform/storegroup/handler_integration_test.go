package storegroup

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

func TestStoreGroupIntegrationTreeAndStatus(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_STORE_GROUP_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_STORE_GROUP_TEST_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DELETE FROM qixi_crm_a_store_group WHERE name LIKE ?", "中文分组验收%")
	db.Exec("DELETE FROM qixi_crm_a_store_group WHERE name LIKE ?", "中文分组验收%")
	gin.SetMode(gin.TestMode)
	h := NewHandler(db)
	r := gin.New()
	r.POST("/groups", h.Create)
	r.PUT("/groups/:id", h.Update)
	r.POST("/groups/:id/status", h.SetStatus)
	r.POST("/groups/:id/template", h.SetTemplate)
	create := func(parent uint, name string) uint {
		raw, _ := json.Marshal(gin.H{"parent_id": parent, "name": name, "address": "虚构中文验收地址"})
		req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("create %s=%d %s", name, w.Code, w.Body.String())
		}
		var out struct {
			Data struct {
				ID uint `json:"id"`
			} `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &out)
		return out.Data.ID
	}
	root := create(0, "中文分组验收根")
	setTemplate := func(groupID, templateID uint) int {
		raw, _ := json.Marshal(gin.H{"diy_page_id": templateID})
		req := httptest.NewRequest(http.MethodPost, "/groups/"+strconv.FormatUint(uint64(groupID), 10)+"/template", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	if setTemplate(root, 999999999) != http.StatusBadRequest {
		t.Fatal("missing template should fail")
	}
	if setTemplate(root, 4001) != http.StatusOK {
		t.Fatal("existing template should bind")
	}
	if setTemplate(999999999, 0) != http.StatusNotFound {
		t.Fatal("missing group template update should fail")
	}
	if w := func() *httptest.ResponseRecorder {
		raw, _ := json.Marshal(gin.H{"parent_id": 0, "name": "中文分组验收根", "address": "虚构中文验收地址"})
		req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}(); w.Code != http.StatusConflict {
		t.Fatalf("duplicate sibling=%d %s", w.Code, w.Body.String())
	}
	child := create(root, "中文分组验收子")
	grand := create(child, "中文分组验收孙")
	move := func(id, parent uint) int {
		raw, _ := json.Marshal(gin.H{"parent_id": parent, "name": "中文分组验收根", "address": "虚构中文验收地址"})
		req := httptest.NewRequest(http.MethodPut, "/groups/"+strconv.FormatUint(uint64(id), 10), bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}
	if move(root, grand) != 400 {
		t.Fatal("cycle move should fail")
	}
	raw, _ := json.Marshal(gin.H{"enabled": false})
	req := httptest.NewRequest(http.MethodPost, "/groups/"+strconv.FormatUint(uint64(root), 10)+"/status", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	var enabled int64
	if err := db.Table(groupTable).Where("id IN ? AND status=1", []uint{root, child, grand}).Count(&enabled).Error; err != nil || enabled != 0 {
		t.Fatalf("cascade=%d err=%v", enabled, err)
	}
	raw, _ = json.Marshal(gin.H{"parent_id": 0, "name": "中文分组验收缺失商户", "address": "虚构中文验收地址", "merchant_ids": []uint{999999999}})
	req = httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("missing merchant=%d", w.Code)
	}
}
