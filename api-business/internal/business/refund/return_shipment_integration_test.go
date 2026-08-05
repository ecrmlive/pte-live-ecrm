package refund

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestReturnShipmentHTTPIntegration uses the production C-end JWT and active
// session middleware. It proves that a buyer can register one UTF-8 shipment,
// retry that exact request safely, and cannot overwrite it or access another
// buyer's after-sale request.
func TestReturnShipmentHTTPIntegration(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("ECRM_RETURN_CHAIN_BUSINESS_TEST_DSN"))
	if dsn == "" {
		t.Skip("set ECRM_RETURN_CHAIN_BUSINESS_TEST_DSN to run return shipment HTTP acceptance")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated business database: %v", err)
	}

	const buyerID uint64 = 988800001
	const otherBuyerID uint64 = 988800002
	const groupID uint64 = 988800003
	const orderID uint64 = 988800004
	const refundID uint64 = 988800005
	const otherOrderID uint64 = 988800006
	const otherRefundID uint64 = 988800007
	ctx := context.Background()
	cleanup := func() {
		for _, id := range []uint64{refundID, otherRefundID} {
			_ = db.WithContext(ctx).Table("qixi_crm_b_refund_event").Where("refund_id = ?", id).Delete(nil).Error
			_ = db.WithContext(ctx).Table("qixi_crm_b_refund_return_shipment").Where("refund_id = ?", id).Delete(nil).Error
			_ = db.WithContext(ctx).Table("qixi_crm_b_refund").Where("id = ?", id).Delete(nil).Error
		}
		_ = db.WithContext(ctx).Table("qixi_crm_b_order").Where("id IN ?", []uint64{orderID, otherOrderID}).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_group_order").Where("id = ?", groupID).Delete(nil).Error
		_ = db.WithContext(ctx).Table("qixi_crm_b_user").Where("id IN ?", []uint64{buyerID, otherBuyerID}).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	if err := db.WithContext(ctx).Table("qixi_crm_b_user").Create([]map[string]any{
		{"id": buyerID, "nickname": "退货验收张三", "mobile": "13900000001", "status": 1, "auth_version": 1},
		{"id": otherBuyerID, "nickname": "隔离验收李四", "mobile": "13900000002", "status": 1, "auth_version": 1},
	}).Error; err != nil {
		t.Fatalf("seed buyers: %v", err)
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_group_order").Create(map[string]any{
		"id": groupID, "order_no": "RETURN-SHIPMENT-GROUP-988800003", "user_id": buyerID,
		"total_amount": 88.8, "pay_amount": 88.8, "total_quantity": 1, "recipient_snapshot": `{}`,
		"pay_channel": "mock", "pay_status": "paid", "idempotency_key": "fixture:return-shipment-group",
	}).Error; err != nil {
		t.Fatalf("seed group order: %v", err)
	}
	for _, row := range []map[string]any{
		{"id": orderID, "group_order_id": groupID, "order_no": "RETURN-SHIPMENT-988800004", "merchant_id": 988800010, "merchant_name_snapshot": "七禧退货验收商户", "store_id": 988800011, "store_name_snapshot": "七禧退货验收店", "user_id": buyerID, "total_amount": 88.8, "pay_amount": 88.8, "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale"},
		{"id": otherOrderID, "group_order_id": groupID, "order_no": "RETURN-SHIPMENT-988800006", "merchant_id": 988800010, "merchant_name_snapshot": "七禧退货验收商户", "store_id": 988800011, "store_name_snapshot": "七禧退货验收店", "user_id": otherBuyerID, "total_amount": 88.8, "pay_amount": 88.8, "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale"},
	} {
		if err := db.WithContext(ctx).Table("qixi_crm_b_order").Create(row).Error; err != nil {
			t.Fatalf("seed after-sale order: %v", err)
		}
	}
	for _, row := range []map[string]any{
		{"id": refundID, "order_id": orderID, "refund_no": "RETURN-SHIPMENT-988800005", "reason": "尺寸不合适，申请退货退款", "amount": 88.8, "refund_type": "return_and_refund", "order_status_before": "completed", "status": "awaiting_return", "idempotency_key": "fixture:return-shipment"},
		{"id": otherRefundID, "order_id": otherOrderID, "refund_no": "RETURN-SHIPMENT-988800007", "reason": "他人售后隔离夹具", "amount": 88.8, "refund_type": "return_and_refund", "order_status_before": "completed", "status": "awaiting_return", "idempotency_key": "fixture:return-shipment-other"},
	} {
		if err := db.WithContext(ctx).Table("qixi_crm_b_refund").Create(row).Error; err != nil {
			t.Fatalf("seed refund: %v", err)
		}
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager("local-return-shipment-acceptance-not-a-production-secret", time.Hour, time.Hour)
	r := gin.New()
	r.Use(middleware.JWTRequired(jwt, authjwt.PortalApp), middleware.CUserSessionRequired(db))
	NewHandler(db).Register(r)
	token, err := jwt.IssueCUserWithIdentityVersion(uint(buyerID), "退货验收张三", "h5", 1)
	if err != nil {
		t.Fatal(err)
	}
	request := func(refund uint64, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/refunds/"+strconv.FormatUint(refund, 10)+"/return-shipment", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authori-zation", "Bearer "+token.AccessToken)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		return resp
	}
	const shipment = `{"carrier_name":"七禧演示快递","tracking_no":"QX-RETURN-988800005","remark":"纸箱完整，内含订单商品。"}`
	if resp := request(refundID, shipment); resp.Code != http.StatusOK {
		t.Fatalf("submit return shipment status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(refundID, shipment); resp.Code != http.StatusOK {
		t.Fatalf("exact shipment replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(refundID, `{"carrier_name":"七禧演示快递","tracking_no":"QX-CHANGED-988800005","remark":"纸箱完整，内含订单商品。"}`); resp.Code != http.StatusBadRequest {
		t.Fatalf("changed shipment status=%d body=%s, want 400", resp.Code, resp.Body.String())
	}
	if resp := request(otherRefundID, shipment); resp.Code != http.StatusNotFound {
		t.Fatalf("other buyer refund status=%d body=%s, want 404", resp.Code, resp.Body.String())
	}
	var fact struct {
		Status      string `gorm:"column:status"`
		CarrierName string `gorm:"column:carrier_name"`
		TrackingNo  string `gorm:"column:tracking_no"`
		Remark      string `gorm:"column:remark"`
		Events      int64  `gorm:"column:events"`
	}
	if err := db.WithContext(ctx).Table("qixi_crm_b_refund AS r").Select("r.status,s.carrier_name,s.tracking_no,s.remark,(SELECT COUNT(*) FROM qixi_crm_b_refund_event e WHERE e.refund_id = r.id) AS events").Joins("JOIN qixi_crm_b_refund_return_shipment AS s ON s.refund_id = r.id").Where("r.id = ?", refundID).Scan(&fact).Error; err != nil {
		t.Fatalf("read shipment facts: %v", err)
	}
	if fact.Status != "awaiting_receipt" || fact.CarrierName != "七禧演示快递" || fact.TrackingNo != "QX-RETURN-988800005" || fact.Remark != "纸箱完整，内含订单商品。" || fact.Events != 1 {
		t.Fatalf("unexpected shipment facts: %+v", fact)
	}
}
