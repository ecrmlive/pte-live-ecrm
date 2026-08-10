package nativediscount

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestToItemComboAndUnlimited(t *testing.T) {
	starts := time.Date(2026, 7, 1, 10, 0, 0, 0, time.Local)
	ends := time.Date(2026, 8, 31, 23, 59, 59, 0, time.Local)
	row := viewRow{
		ActivityID: 5102,
		StoreID:    1,
		Name:       "夏日香氛搭配套餐",
		Rules: `{
			"package_price":129.00,
			"package_type":1,
			"is_limit":0,
			"limit_num":0,
			"is_time":1,
			"free_shipping":true,
			"create_time":"2026-07-01 09:30:00",
			"products":[
				{"product_id":1004,"store_name":"无火藤条香氛","image":"https://picsum.photos/seed/qixi-discount-main/120/120","type":0,"spec":"| 69.00"},
				{"product_id":1006,"store_name":"香氛扩香石","image":"https://picsum.photos/seed/qixi-discount-combo/120/120","type":1,"spec":"| 39.00"}
			]
		}`,
		Status:   1,
		Version:  1,
		StartsAt: &starts,
		EndsAt:   &ends,
	}
	item := toItem(row, "七禧演示店铺")
	if item["package_type"] != 1 {
		t.Fatalf("package_type=%v", item["package_type"])
	}
	if item["package_type_label"] != "搭配套餐" {
		t.Fatalf("label=%v", item["package_type_label"])
	}
	if item["remain_label"] != "不限量" {
		t.Fatalf("remain=%v", item["remain_label"])
	}
	if item["store_name"] != "七禧演示店铺" {
		t.Fatalf("store=%v", item["store_name"])
	}
	mp, ok := item["main_products"].([]gin.H)
	if !ok || len(mp) != 1 {
		t.Fatalf("main_products=%#v", item["main_products"])
	}
	cp, ok := item["combo_products"].([]gin.H)
	if !ok || len(cp) != 1 {
		t.Fatalf("combo_products=%#v", item["combo_products"])
	}
}

func TestToItemFixedUnlimitedTime(t *testing.T) {
	row := viewRow{
		ActivityID: 5103,
		StoreID:    2,
		Name:       "居家固定套餐",
		Rules: `{
			"package_price":88.00,
			"type":0,
			"is_limit":1,
			"limit_num":50,
			"is_time":0,
			"create_time":"2026-06-15 14:20:00",
			"products":[
				{"product_id":1101,"store_name":"棉柔毛巾三件套","image":"https://picsum.photos/seed/qixi-discount-fixed/120/120","type":0,"spec":"| 88.00"}
			]
		}`,
		Status: 1,
	}
	item := toItem(row, "七禧居家优选店")
	if item["package_type"] != 0 {
		t.Fatalf("package_type=%v", item["package_type"])
	}
	if item["time_label"] != "不限时" {
		t.Fatalf("time=%v", item["time_label"])
	}
	if item["remain_label"] != "50" {
		t.Fatalf("remain=%v", item["remain_label"])
	}
	if item["starts_at"] != "" || item["ends_at"] != "" {
		t.Fatalf("unexpected time %v %v", item["starts_at"], item["ends_at"])
	}
	cp, ok := item["combo_products"].([]gin.H)
	if !ok || len(cp) != 0 {
		t.Fatalf("fixed should have empty combo: %#v", item["combo_products"])
	}
}
