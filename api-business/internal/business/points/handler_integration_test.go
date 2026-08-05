package points

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestQuoteRejectsOffSaleProduct proves the C-end creation path consults the
// same business projection that the platform supervision page updates.
func TestQuoteRejectsOffSaleProduct(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_POINTS_BUSINESS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_POINTS_BUSINESS_TEST_DSN to run points integration test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	const productID = 987675001
	_ = db.Table("qixi_crm_b_points_product_view").Where("product_id=?", productID).Delete(nil).Error
	t.Cleanup(func() { _ = db.Table("qixi_crm_b_points_product_view").Where("product_id=?", productID).Delete(nil).Error })
	if err := db.Table("qixi_crm_b_points_product_view").Create(map[string]any{"product_id": productID, "merchant_id": 1, "store_id": 1, "merchant_name": "中文积分验收商户", "store_name": "中文积分验收店", "title": "中文积分验收商品", "original_price": "88.00", "points_required": 120, "stock": 5, "sale_status": 1, "version": 1}).Error; err != nil {
		t.Fatal(err)
	}
	h := NewHandler(db)
	if _, _, got, err := h.quote(context.Background(), input{ProductID: productID, CartNum: 2}); err != nil || got != 240 {
		t.Fatalf("on-sale quote points=%d err=%v", got, err)
	}
	if err := db.Table("qixi_crm_b_points_product_view").Where("product_id=?", productID).Update("sale_status", 0).Error; err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := h.quote(context.Background(), input{ProductID: productID, CartNum: 1}); !errors.Is(err, errNotFound) {
		t.Fatalf("off-sale quote error=%v, want errNotFound", err)
	}
}
