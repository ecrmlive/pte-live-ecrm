package nativerefund

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestConfirmReturnHTTPIntegration exercises the merchant console's complete
// route middleware: store JWT context, X-AppId binding, active session and
// RBAC. It must turn only an in-store awaiting-receipt refund into refunding;
// a provider callback remains the sole writer of the refunded terminal state.
func TestConfirmReturnHTTPIntegration(t *testing.T) {
	businessDSN := strings.TrimSpace(os.Getenv("ECRM_RETURN_CHAIN_BUSINESS_TEST_DSN"))
	merchantDSN := strings.TrimSpace(os.Getenv("ECRM_RETURN_CHAIN_MERCHANT_TEST_DSN"))
	if businessDSN == "" || merchantDSN == "" {
		t.Skip("set ECRM_RETURN_CHAIN_BUSINESS_TEST_DSN and ECRM_RETURN_CHAIN_MERCHANT_TEST_DSN to run confirm return HTTP acceptance")
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated business database: %v", err)
	}
	merchantDB, err := gorm.Open(mysql.Open(merchantDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated merchant database: %v", err)
	}

	const merchantID uint64 = 988810001
	const storeID uint64 = 988810002
	const otherStoreID uint64 = 988810003
	const accountID uint64 = 988810004
	const otherAccountID uint64 = 988810005
	const groupID uint64 = 988810006
	const orderID uint64 = 988810007
	const refundID uint64 = 988810008
	const paymentID uint64 = 988810009
	ctx := context.Background()
	cleanup := func() {
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund_event").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund_return_shipment").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund_transaction").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_refund").Where("id = ?", refundID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_payment_transaction").Where("id = ?", paymentID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_order").Where("id = ?", orderID).Delete(nil).Error
		_ = businessDB.WithContext(ctx).Table("qixi_crm_b_group_order").Where("id = ?", groupID).Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_aftersale_action").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_refund_hidden").Where("refund_id = ?", refundID).Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_role_menu").Where("role_code = ?", "owner").Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_menu").Where("code IN ?", []string{"refund.approve", "refund.remark", "refund.delete"}).Delete(nil).Error
		_ = merchantDB.WithContext(ctx).Table("qixi_crm_m_account").Where("id IN ?", []uint64{accountID, otherAccountID}).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)
	if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_account").Create([]map[string]any{
		{"id": accountID, "store_id": storeID, "username": "return_acceptance_owner", "password_hash": "not-a-real-password", "role_code": "owner", "display_name": "退货验收店长", "status": 1, "auth_version": 1},
		{"id": otherAccountID, "store_id": otherStoreID, "username": "return_acceptance_other", "password_hash": "not-a-real-password", "role_code": "owner", "display_name": "隔离验收店长", "status": 1, "auth_version": 1},
	}).Error; err != nil {
		t.Fatalf("seed merchant accounts: %v", err)
	}
	for _, menu := range []map[string]any{
		{"id": 988810010, "code": "refund.approve", "name": "退款同意", "path": "/refunds/approve", "is_menu": 2, "is_route": 0, "status": 1},
		{"id": 988810012, "code": "refund.remark", "name": "售后备注", "path": "/refunds/remark", "is_menu": 2, "is_route": 0, "status": 1},
		{"id": 988810013, "code": "refund.delete", "name": "隐藏退款", "path": "/refunds/delete", "is_menu": 2, "is_route": 0, "status": 1},
	} {
		if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_menu").Create(menu).Error; err != nil {
			t.Fatalf("seed refund permission: %v", err)
		}
		if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_role_menu").Create(map[string]any{"role_code": "owner", "menu_id": menu["id"]}).Error; err != nil {
			t.Fatalf("bind refund permission: %v", err)
		}
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_group_order").Create(map[string]any{
		"id": groupID, "order_no": "RETURN-CONFIRM-GROUP-988810006", "user_id": 988810011,
		"total_amount": 66.5, "pay_amount": 66.5, "total_quantity": 1, "recipient_snapshot": `{}`,
		"pay_channel": "mock", "pay_status": "paid", "idempotency_key": "fixture:return-confirm-group",
	}).Error; err != nil {
		t.Fatalf("seed group order: %v", err)
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_order").Create(map[string]any{
		"id": orderID, "group_order_id": groupID, "order_no": "RETURN-CONFIRM-988810007", "merchant_id": merchantID, "merchant_name_snapshot": "七禧确认收货验收商户", "store_id": storeID, "store_name_snapshot": "七禧确认收货验收店", "user_id": 988810011, "total_amount": 66.5, "pay_amount": 66.5, "total_quantity": 1, "recipient_snapshot": `{}`, "status": "aftersale",
	}).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_refund").Create(map[string]any{
		"id": refundID, "order_id": orderID, "refund_no": "RETURN-CONFIRM-988810008", "reason": "商品颜色与描述不符", "amount": 66.5, "refund_type": "return_and_refund", "order_status_before": "completed", "status": "awaiting_receipt", "idempotency_key": "fixture:return-confirm-refund",
	}).Error; err != nil {
		t.Fatalf("seed awaiting receipt refund: %v", err)
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_refund_return_shipment").Create(map[string]any{"refund_id": refundID, "carrier_name": "七禧演示快递", "tracking_no": "QX-CONFIRM-988810008", "remark": "买家已寄回，待仓库核验。", "submitted_by": 988810011}).Error; err != nil {
		t.Fatalf("seed return shipment: %v", err)
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_payment_transaction").Create(map[string]any{"id": paymentID, "group_order_id": groupID, "channel": "mock", "transaction_no": "RETURN-CONFIRM-PAYMENT-988810009", "amount": 66.5, "status": "succeeded"}).Error; err != nil {
		t.Fatalf("seed successful payment: %v", err)
	}

	gin.SetMode(gin.TestMode)
	jwt := authjwt.NewManager("local-confirm-return-acceptance-not-a-production-secret", time.Hour, time.Hour)
	r := gin.New()
	authed := r.Group("/api/merchant/v1")
	authed.Use(middleware.JWTRequired(jwt, authjwt.PortalMerchant), middleware.RequireStoreConsole(), middleware.RequireStoreAppID(), middleware.RequireActiveStoreSession(merchantDB))
	NewHandler(businessDB, merchantDB).Register(authed)
	owner, err := jwt.IssueStoreConsoleWithIdentityVersion(uint(accountID), uint(merchantID), uint(storeID), "qixi.return.acceptance.store", "1400000000", "return_acceptance_owner", "owner", 1)
	if err != nil {
		t.Fatal(err)
	}
	other, err := jwt.IssueStoreConsoleWithIdentityVersion(uint(otherAccountID), uint(merchantID), uint(otherStoreID), "qixi.return.acceptance.other", "1400000000", "return_acceptance_other", "owner", 1)
	if err != nil {
		t.Fatal(err)
	}
	request := func(method, suffix, token, appID, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, "/api/merchant/v1/refunds/"+strconv.FormatUint(refundID, 10)+suffix, strings.NewReader(body))
		req.Header.Set("Authori-zation", "Bearer "+token)
		req.Header.Set("X-AppId", appID)
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		return resp
	}
	if resp := request(http.MethodPost, "/confirm-return", owner.AccessToken, "qixi.return.acceptance.store", ""); resp.Code != http.StatusOK {
		t.Fatalf("confirm returned goods status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodPost, "/confirm-return", owner.AccessToken, "qixi.return.acceptance.store", ""); resp.Code != http.StatusOK {
		t.Fatalf("confirm returned goods replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodPost, "/confirm-return", other.AccessToken, "qixi.return.acceptance.other", ""); resp.Code != http.StatusNotFound {
		t.Fatalf("cross-store confirmation status=%d body=%s, want 404", resp.Code, resp.Body.String())
	}
	var fact struct {
		Status       string  `gorm:"column:status"`
		Channel      string  `gorm:"column:channel"`
		Amount       float64 `gorm:"column:amount"`
		Transactions int64   `gorm:"column:transactions"`
		Events       int64   `gorm:"column:events"`
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_refund AS r").Select("r.status,t.channel,t.amount,(SELECT COUNT(*) FROM qixi_crm_b_refund_transaction x WHERE x.refund_id = r.id) AS transactions,(SELECT COUNT(*) FROM qixi_crm_b_refund_event e WHERE e.refund_id = r.id) AS events").Joins("JOIN qixi_crm_b_refund_transaction AS t ON t.refund_id = r.id").Where("r.id = ?", refundID).Scan(&fact).Error; err != nil {
		t.Fatalf("read confirmation facts: %v", err)
	}
	if fact.Status != "refunding" || fact.Channel != "mock" || fact.Amount != 66.5 || fact.Transactions != 1 || fact.Events != 1 {
		t.Fatalf("unexpected confirmation facts: %+v", fact)
	}

	const remark = `{"note":"仓库已核验退货包裹外观，等待退款回调。","idempotency_key":"return-remark-988810008"}`
	if resp := request(http.MethodPost, "/remark", owner.AccessToken, "qixi.return.acceptance.store", remark); resp.Code != http.StatusOK {
		t.Fatalf("create Chinese remark status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodPost, "/remark", owner.AccessToken, "qixi.return.acceptance.store", remark); resp.Code != http.StatusOK {
		t.Fatalf("remark replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodPost, "/remark", owner.AccessToken, "qixi.return.acceptance.store", `{"note":"篡改后的备注","idempotency_key":"return-remark-988810008"}`); resp.Code != http.StatusConflict {
		t.Fatalf("changed remark key status=%d body=%s, want 409", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodPost, "/remark", other.AccessToken, "qixi.return.acceptance.other", remark); resp.Code != http.StatusNotFound {
		t.Fatalf("cross-store remark status=%d body=%s, want 404", resp.Code, resp.Body.String())
	}

	const hide = `{"reason":"本店已完成售后归档，隐藏后台视图。","idempotency_key":"return-hide-988810008"}`
	if resp := request(http.MethodDelete, "", owner.AccessToken, "qixi.return.acceptance.store", hide); resp.Code != http.StatusOK {
		t.Fatalf("hide refund status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodDelete, "", owner.AccessToken, "qixi.return.acceptance.store", hide); resp.Code != http.StatusOK {
		t.Fatalf("hide refund replay status=%d body=%s", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodDelete, "", owner.AccessToken, "qixi.return.acceptance.store", `{"reason":"篡改后的隐藏原因","idempotency_key":"return-hide-988810008"}`); resp.Code != http.StatusConflict {
		t.Fatalf("changed hide key status=%d body=%s, want 409", resp.Code, resp.Body.String())
	}
	if resp := request(http.MethodGet, "", owner.AccessToken, "qixi.return.acceptance.store", ""); resp.Code != http.StatusNotFound {
		t.Fatalf("hidden refund detail status=%d body=%s, want 404", resp.Code, resp.Body.String())
	}
	var actions, hidden, businessRefunds int64
	if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_aftersale_action").Where("refund_id = ? AND action = ?", refundID, "remark").Count(&actions).Error; err != nil {
		t.Fatalf("count immutable remarks: %v", err)
	}
	if err := merchantDB.WithContext(ctx).Table("qixi_crm_m_refund_hidden").Where("refund_id = ? AND store_id = ?", refundID, storeID).Count(&hidden).Error; err != nil {
		t.Fatalf("count hidden views: %v", err)
	}
	if err := businessDB.WithContext(ctx).Table("qixi_crm_b_refund").Where("id = ? AND status = ?", refundID, "refunding").Count(&businessRefunds).Error; err != nil {
		t.Fatalf("count unchanged business refund: %v", err)
	}
	if actions != 1 || hidden != 1 || businessRefunds != 1 {
		t.Fatalf("unexpected remark/hide facts: actions=%d hidden=%d businessRefunds=%d", actions, hidden, businessRefunds)
	}
}
