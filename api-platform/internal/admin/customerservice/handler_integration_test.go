package customerservice

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	adminauth "github.com/crmlive/pte-live-ecrm/api-platform/internal/admin/auth"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	platformoperationlog "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/operationlog"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestCustomerServiceHTTPRBACAndQueueClosure uses the production route stack
// against an isolated MySQL database. It deliberately keeps only fake Chinese
// identities and verifies that a customer-service account cannot enumerate or
// change another store's queue.
func TestCustomerServiceHTTPRBACAndQueueClosure(t *testing.T) {
	adminDSN := strings.TrimSpace(os.Getenv("ECRM_CUSTOMER_SERVICE_ADMIN_TEST_DSN"))
	businessDSN := strings.TrimSpace(os.Getenv("ECRM_CUSTOMER_SERVICE_BUSINESS_TEST_DSN"))
	if adminDSN == "" || businessDSN == "" {
		t.Skip("set ECRM_CUSTOMER_SERVICE_ADMIN_TEST_DSN and ECRM_CUSTOMER_SERVICE_BUSINESS_TEST_DSN to run customer-service integration test")
	}
	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	businessDB, err := gorm.Open(mysql.Open(businessDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	roles := []string{"platform", "merchant", "region", "operations", "customer_service"}
	type roleRow struct {
		ID   uint
		Code string
	}
	var rows []roleRow
	if err := adminDB.Table("qixi_crm_a_role").Select("id,code").Where("code IN ? AND status = 1", roles).Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	roleIDs := map[string]uint{}
	for _, row := range rows {
		roleIDs[row.Code] = row.ID
	}
	if len(roleIDs) != len(roles) {
		t.Fatalf("five role fixture missing: %#v", roleIDs)
	}

	const (
		platformID      uint = 987672001
		merchantID      uint = 987672002
		regionID        uint = 987672003
		operationsID    uint = 987672004
		serviceID       uint = 987672005
		targetServiceID uint = 987672006
		nonCurrentID    uint = 987672007
		storeA          uint = 987671001
		storeB          uint = 987671002
		userA           uint = 987671011
		userB           uint = 987671012
		bindingA        uint = 987671021
		bindingB        uint = 987671022
	)
	adminIDs := []uint{platformID, merchantID, regionID, operationsID, serviceID, targetServiceID, nonCurrentID}
	storeIDs := []uint{storeA, storeB}
	userIDs := []uint{userA, userB}
	bindingIDs := []uint{bindingA, bindingB}
	var previousSettings serviceConfigRow
	hadPreviousSettings := adminDB.Table("qixi_crm_a_config").Where("config_key = ?", serviceSettingsConfigKey).Take(&previousSettings).Error == nil
	cleanup := func() {
		_ = businessDB.Table("qixi_crm_b_customer_service_assignment_log").Where("binding_id IN ?", bindingIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_customer_service_message").Where("binding_id IN ?", bindingIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_customer_service_user_note").Where("user_id IN ?", userIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_customer_service_quick_reply").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_customer_service_binding").Where("id IN ?", bindingIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_user_profile").Where("user_id IN ?", userIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_user").Where("id IN ?", userIDs).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_merchant_im_sdk_app_view").Where("merchant_id IN ?", []uint{storeA, storeB}).Delete(nil).Error
		_ = businessDB.Table("qixi_crm_b_store_view").Where("store_id IN ?", storeIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_config").Where("config_key = ?", serviceSettingsConfigKey).Delete(nil).Error
		if hadPreviousSettings {
			_ = adminDB.Table("qixi_crm_a_config").Create(&previousSettings).Error
		}
		_ = adminDB.Table("qixi_crm_a_operation_log").Where("admin_user_id IN ?", adminIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_data_scope").Where("admin_user_id IN ?", adminIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user_role").Where("admin_user_id IN ?", adminIDs).Delete(nil).Error
		_ = adminDB.Table("qixi_crm_a_admin_user").Where("id IN ?", adminIDs).Delete(nil).Error
	}
	cleanup()
	t.Cleanup(cleanup)

	adminFixtures := []struct {
		id          uint
		role        string
		username    string
		displayName string
	}{
		{platformID, "platform", "cs-platform-acceptance", "中文客服平台验收员"},
		{merchantID, "merchant", "cs-merchant-acceptance", "中文客服商户验收员"},
		{regionID, "region", "cs-region-acceptance", "中文客服区域验收员"},
		{operationsID, "operations", "cs-operations-acceptance", "中文客服运营验收员"},
		{serviceID, "customer_service", "cs-agent-acceptance", "中文客服张敏"},
		{targetServiceID, "customer_service", "cs-target-acceptance", "中文客服李岚"},
		{nonCurrentID, "customer_service", "cs-non-current-acceptance", "中文客服周岚"},
	}
	for _, item := range adminFixtures {
		if err := adminDB.Table("qixi_crm_a_admin_user").Create(map[string]any{
			"id": item.id, "username": item.username, "password_hash": "not-used-by-integration-test",
			"display_name": item.displayName, "status": 1, "auth_version": 1, "data_scope_version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := adminDB.Table("qixi_crm_a_admin_user_role").Create(map[string]any{"admin_user_id": item.id, "role_id": roleIDs[item.role]}).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []uint{serviceID, targetServiceID, nonCurrentID} {
		if err := adminDB.Table("qixi_crm_a_data_scope").Create(map[string]any{
			"admin_user_id": id, "scope_type": "service_queue", "scope_value": `{"store_ids":[987671001]}`, "version": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, item := range []struct {
		store    uint
		merchant uint
		name     string
		appID    string
		sdkID    string
	}{
		{storeA, storeA, "七禧中文客服演示店", "cs-acceptance-store-a", "cs-demo-sdk-a"},
		{storeB, storeB, "七禧范围外演示店", "cs-acceptance-store-b", "cs-demo-sdk-b"},
	} {
		if err := businessDB.Table("qixi_crm_b_store_view").Create(map[string]any{
			"store_id": item.store, "merchant_id": item.merchant, "store_app_id": item.appID, "store_name": item.name, "status": 1,
		}).Error; err != nil {
			t.Fatal(err)
		}
		if err := businessDB.Table("qixi_crm_b_merchant_im_sdk_app_view").Create(map[string]any{
			"merchant_id": item.merchant, "sdk_app_id": item.sdkID,
		}).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, item := range []struct {
		id       uint
		nickname string
		mobile   string
	}{
		{userA, "中文客服用户小满", "13900000011"},
		{userB, "中文范围外用户小雪", "13900000012"},
	} {
		if err := businessDB.Table("qixi_crm_b_user").Create(map[string]any{"id": item.id, "nickname": item.nickname, "mobile": item.mobile, "status": 1, "auth_version": 1}).Error; err != nil {
			t.Fatal(err)
		}
		if err := businessDB.Table("qixi_crm_b_user_profile").Create(map[string]any{"user_id": item.id, "bio": "中文模拟客服资料", "source_channel": "h5"}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := businessDB.Table("qixi_crm_b_customer_service_binding").Create([]map[string]any{
		{"id": bindingA, "user_id": userA, "store_id": storeA, "im_conversation_id": "cs-acceptance-a", "status": "open", "last_msg": "咨询虚构商品的发货时效"},
		// 同一虚构用户在两家店铺分别建会话，用于验证用户备注严格按店铺隔离。
		{"id": bindingB, "user_id": userA, "store_id": storeB, "im_conversation_id": "cs-acceptance-b", "status": "open", "last_msg": "范围外会话不得泄露"},
	}).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_customer_service_message").Create(map[string]any{
		"binding_id": bindingA, "merchant_id": storeA, "sender_role": "system", "sender_id": 0, "msg_type": "system", "content": "中文模拟订单事件，不保存聊天正文",
	}).Error; err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	jwtManager := authjwt.NewManager(strings.Repeat("x", 32), time.Minute, 2*time.Minute)
	router := gin.New()
	authed := router.Group("/api/platform/v1")
	authed.Use(
		middleware.JWTRequired(jwtManager, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(adminDB),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
		middleware.AuditAdminMutation(adminDB),
	)
	NewHandler(adminDB, businessDB).Register(authed)
	adminauth.NewHandler(adminDB, jwtManager).RegisterSettings(authed)
	platformoperationlog.New(adminDB).Register(authed)

	call := func(item struct {
		id                          uint
		role, username, displayName string
	}, method, path string, body any, key string) *httptest.ResponseRecorder {
		var raw []byte
		if body != nil {
			var marshalErr error
			raw, marshalErr = json.Marshal(body)
			if marshalErr != nil {
				t.Fatal(marshalErr)
			}
		}
		pair, issueErr := jwtManager.IssueAdminConsole(item.id, item.username, []string{item.role}, 1)
		if issueErr != nil {
			t.Fatal(issueErr)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Authori-zation", "Bearer "+pair.AccessToken)
		if key != "" {
			req.Header.Set("X-Idempotency-Key", key)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	platform := adminFixtures[0]
	merchant := adminFixtures[1]
	region := adminFixtures[2]
	operations := adminFixtures[3]
	service := adminFixtures[4]
	targetService := adminFixtures[5]
	nonCurrent := adminFixtures[6]

	settings := gin.H{"auto_reply_enabled": true, "auto_reply_text": "虚构中文客服会在工作时间内回复。", "queue_mode": "round_robin", "max_sessions_per_agent": 12}
	if got := call(platform, http.MethodPut, "/api/platform/v1/customer-service/settings", settings, "").Code; got != http.StatusOK {
		t.Fatalf("platform settings update=%d, want 200", got)
	}
	for _, item := range []struct {
		id                          uint
		role, username, displayName string
	}{merchant, region, operations, service} {
		if got := call(item, http.MethodPut, "/api/platform/v1/customer-service/settings", settings, "").Code; got != http.StatusForbidden {
			t.Fatalf("%s settings update=%d, want 403", item.role, got)
		}
	}
	if got := call(platform, http.MethodGet, "/api/platform/v1/customer-service/threads", nil, "").Code; got != http.StatusOK {
		t.Fatalf("platform thread list=%d, want 200", got)
	}
	for _, item := range []struct {
		id                          uint
		role, username, displayName string
	}{merchant, region, operations} {
		if got := call(item, http.MethodGet, "/api/platform/v1/customer-service/threads", nil, "").Code; got != http.StatusForbidden {
			t.Fatalf("%s thread list=%d, want 403", item.role, got)
		}
	}
	serviceList := call(service, http.MethodGet, "/api/platform/v1/customer-service/threads", nil, "")
	if serviceList.Code != http.StatusOK || !strings.Contains(serviceList.Body.String(), "987671021") || strings.Contains(serviceList.Body.String(), "987671022") {
		t.Fatalf("service queue scope is invalid: code=%d body=%s", serviceList.Code, serviceList.Body.String())
	}
	if got := call(service, http.MethodGet, "/api/platform/v1/customer-service/threads/987671022", nil, "").Code; got != http.StatusNotFound {
		t.Fatalf("out-of-scope thread detail=%d, want 404", got)
	}
	quickReply := gin.H{"store_id": storeA, "title": "发货时效", "content": "虚构店铺将在 48 小时内发货。", "status": "enabled"}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/quick-replies", quickReply, "").Code; got != http.StatusOK {
		t.Fatalf("service quick reply create=%d, want 200", got)
	}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/quick-replies", gin.H{"store_id": storeB, "title": "越权测试", "content": "不得写入范围外店铺。", "status": "enabled"}, "").Code; got != http.StatusForbidden {
		t.Fatalf("out-of-scope quick reply create=%d, want 403", got)
	}
	if got := call(platform, http.MethodPost, "/api/platform/v1/customer-service/quick-replies", gin.H{"store_id": storeB, "title": "范围外客服模板", "content": "该模板仅用于校验客服数据范围。", "status": "enabled"}, "").Code; got != http.StatusOK {
		t.Fatalf("platform quick reply create for scope check=%d, want 200", got)
	}
	quickReplies := call(service, http.MethodGet, "/api/platform/v1/customer-service/quick-replies?store_id=987671001", nil, "")
	if quickReplies.Code != http.StatusOK || !strings.Contains(quickReplies.Body.String(), "发货时效") || strings.Contains(quickReplies.Body.String(), "范围外客服模板") {
		t.Fatalf("service quick reply list scope=%d body=%s", quickReplies.Code, quickReplies.Body.String())
	}
	if got := call(service, http.MethodGet, "/api/platform/v1/customer-service/quick-replies?store_id=987671002", nil, "").Code; got != http.StatusForbidden {
		t.Fatalf("out-of-scope quick reply list=%d, want 403", got)
	}
	var quickReplyA, quickReplyB struct{ ID uint }
	if err := businessDB.Table("qixi_crm_b_customer_service_quick_reply").Select("id").Where("store_id = ? AND title = ?", storeA, "发货时效").Take(&quickReplyA).Error; err != nil {
		t.Fatal(err)
	}
	if err := businessDB.Table("qixi_crm_b_customer_service_quick_reply").Select("id").Where("store_id = ? AND title = ?", storeB, "范围外客服模板").Take(&quickReplyB).Error; err != nil {
		t.Fatal(err)
	}
	if got := call(service, http.MethodPut, "/api/platform/v1/customer-service/quick-replies/"+strconv.FormatUint(uint64(quickReplyA.ID), 10), gin.H{"title": "售后时效说明", "content": "中文模拟编辑：虚构店铺将在 72 小时内发货。", "status": "disabled"}, "").Code; got != http.StatusOK {
		t.Fatalf("service quick reply update=%d, want 200", got)
	}
	if got := call(service, http.MethodPut, "/api/platform/v1/customer-service/quick-replies/"+strconv.FormatUint(uint64(quickReplyB.ID), 10), gin.H{"title": "越权修改", "content": "不得改写范围外模板。", "status": "disabled"}, "").Code; got != http.StatusNotFound {
		t.Fatalf("out-of-scope quick reply update=%d, want 404", got)
	}
	if got := call(service, http.MethodDelete, "/api/platform/v1/customer-service/quick-replies/"+strconv.FormatUint(uint64(quickReplyA.ID), 10), nil, "").Code; got != http.StatusOK {
		t.Fatalf("service quick reply delete=%d, want 200", got)
	}
	quickReplies = call(service, http.MethodGet, "/api/platform/v1/customer-service/quick-replies?store_id=987671001", nil, "")
	if quickReplies.Code != http.StatusOK || strings.Contains(quickReplies.Body.String(), "售后时效说明") {
		t.Fatalf("soft-deleted quick reply still listed: code=%d body=%s", quickReplies.Code, quickReplies.Body.String())
	}
	var deletedQuickReply struct {
		Status    string
		UpdatedBy uint64
		DeletedAt *time.Time
	}
	if err := businessDB.Table("qixi_crm_b_customer_service_quick_reply").Select("status,updated_by,deleted_at").Where("id = ?", quickReplyA.ID).Take(&deletedQuickReply).Error; err != nil || deletedQuickReply.Status != "disabled" || deletedQuickReply.UpdatedBy != uint64(serviceID) || deletedQuickReply.DeletedAt == nil {
		t.Fatalf("quick reply lifecycle=%#v err=%v", deletedQuickReply, err)
	}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/claim", nil, "").Code; got != http.StatusOK {
		t.Fatalf("service claim=%d, want 200", got)
	}
	transfer := gin.H{"target_admin_id": targetServiceID, "reason": "中文模拟验收：转交同店铺售后咨询。"}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", transfer, "cs-acceptance-transfer-001").Code; got != http.StatusOK {
		t.Fatalf("service transfer=%d, want 200", got)
	}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", transfer, "cs-acceptance-transfer-001").Code; got != http.StatusOK {
		t.Fatalf("service transfer replay=%d, want 200", got)
	}
	if got := call(service, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", gin.H{"target_admin_id": targetServiceID, "reason": "同一键不得覆盖原始中文原因。"}, "cs-acceptance-transfer-001").Code; got != http.StatusConflict {
		t.Fatalf("changed transfer replay=%d, want 409", got)
	}
	if got := call(nonCurrent, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", gin.H{"target_admin_id": serviceID, "reason": "非当前领取客服不得覆盖会话归属。"}, "cs-non-current-transfer-001").Code; got != http.StatusForbidden {
		t.Fatalf("non-current service transfer=%d, want 403", got)
	}
	if got := call(service, http.MethodPut, "/api/platform/v1/customer-service/threads/987671021/user-note", gin.H{"content": "中文模拟备注：已说明虚构商品售后时效。"}, "").Code; got != http.StatusOK {
		t.Fatalf("service note update=%d, want 200", got)
	}
	if got := call(service, http.MethodPut, "/api/platform/v1/customer-service/threads/987671022/user-note", gin.H{"content": "不得写入范围外用户。"}, "").Code; got != http.StatusNotFound {
		t.Fatalf("out-of-scope note update=%d, want 404", got)
	}
	if got := call(platform, http.MethodPut, "/api/platform/v1/customer-service/threads/987671022/user-note", gin.H{"content": "中文模拟备注：同一用户在范围外店铺的独立售后说明。"}, "").Code; got != http.StatusOK {
		t.Fatalf("platform cross-store note update=%d, want 200", got)
	}
	var assigned uint
	if err := businessDB.Table("qixi_crm_b_customer_service_binding").Select("assigned_admin_id").Where("id = ?", bindingA).Scan(&assigned).Error; err != nil || assigned != targetServiceID {
		t.Fatalf("transfer assignment=%d err=%v, want %d", assigned, err, targetServiceID)
	}
	var transferAuditCount int64
	if err := businessDB.Table("qixi_crm_b_customer_service_assignment_log").Where("binding_id = ?", bindingA).Count(&transferAuditCount).Error; err != nil || transferAuditCount != 1 {
		t.Fatalf("transfer audit count=%d err=%v, want 1", transferAuditCount, err)
	}
	var note string
	if err := businessDB.Table("qixi_crm_b_customer_service_user_note").Select("content").Where("user_id = ? AND store_id = ?", userA, storeA).Scan(&note).Error; err != nil || note != "中文模拟备注：已说明虚构商品售后时效。" {
		t.Fatalf("service note=%q err=%v", note, err)
	}
	var crossStoreNote string
	if err := businessDB.Table("qixi_crm_b_customer_service_user_note").Select("content").Where("user_id = ? AND store_id = ?", userA, storeB).Scan(&crossStoreNote).Error; err != nil || crossStoreNote != "中文模拟备注：同一用户在范围外店铺的独立售后说明。" {
		t.Fatalf("cross-store service note=%q err=%v", crossStoreNote, err)
	}

	// 管理员逻辑删除必须同时收口五角色写权限、会话、客服名册和转接目标。
	// 先验证四种非平台角色不能删，再删除虚构客服李岚；历史转接审计仍保留。
	for _, item := range []struct {
		id                          uint
		role, username, displayName string
	}{merchant, region, operations, service} {
		if got := call(item, http.MethodDelete, "/api/platform/v1/setting/admins/987672006", nil, "").Code; got != http.StatusForbidden {
			t.Fatalf("%s delete admin=%d, want 403", item.role, got)
		}
	}
	if got := call(platform, http.MethodDelete, "/api/platform/v1/setting/admins/987672001", nil, "").Code; got != http.StatusConflict {
		t.Fatalf("platform self-delete=%d, want 409", got)
	}
	if got := call(platform, http.MethodDelete, "/api/platform/v1/setting/admins/987672006", nil, "").Code; got != http.StatusOK {
		t.Fatalf("platform delete customer service=%d, want 200", got)
	}
	const deleteAction = "DELETE /api/platform/v1/setting/admins/:id"
	var deleteAudit struct {
		AdminUserID  uint64
		RoleCode     string
		Action       string
		ResourceType string
		ResourceID   string
		RequestID    string
	}
	if err := adminDB.Table("qixi_crm_a_operation_log").Where("action = ? AND resource_type = ? AND resource_id = ?", deleteAction, "setting", "987672006").Take(&deleteAudit).Error; err != nil || deleteAudit.AdminUserID != uint64(platformID) || deleteAudit.RoleCode != "platform" || deleteAudit.RequestID == "" {
		t.Fatalf("successful delete audit=%#v err=%v", deleteAudit, err)
	}
	var failedDeleteAudits int64
	if err := adminDB.Table("qixi_crm_a_operation_log").Where("action = ? AND resource_type = ? AND resource_id = ? AND admin_user_id <> ?", deleteAction, "setting", "987672006", platformID).Count(&failedDeleteAudits).Error; err != nil || failedDeleteAudits != 0 {
		t.Fatalf("failed delete must not be audited: count=%d err=%v", failedDeleteAudits, err)
	}
	for _, item := range []struct {
		id                          uint
		role, username, displayName string
	}{merchant, region, operations, service} {
		if got := call(item, http.MethodGet, "/api/platform/v1/operation-logs", nil, "").Code; got != http.StatusForbidden {
			t.Fatalf("%s operation-log list=%d, want 403", item.role, got)
		}
	}
	if logs := call(platform, http.MethodGet, "/api/platform/v1/operation-logs?action=DELETE%20%2Fapi%2Fplatform%2Fv1%2Fsetting%2Fadmins%2F%3Aid", nil, ""); logs.Code != http.StatusOK || !strings.Contains(logs.Body.String(), deleteAction) || strings.Contains(logs.Body.String(), "not-used-by-integration-test") {
		t.Fatalf("platform operation-log list=%d body=%s", logs.Code, logs.Body.String())
	}
	var deleted struct {
		Status           int8
		AuthVersion      uint64
		DataScopeVersion uint64
		DeletedAt        *time.Time
	}
	if err := adminDB.Table("qixi_crm_a_admin_user").Select("status,auth_version,data_scope_version,deleted_at").Where("id = ?", targetServiceID).Scan(&deleted).Error; err != nil || deleted.Status != 0 || deleted.AuthVersion != 2 || deleted.DataScopeVersion != 2 || deleted.DeletedAt == nil {
		t.Fatalf("logical delete state=%#v err=%v", deleted, err)
	}
	if got := call(targetService, http.MethodGet, "/api/platform/v1/customer-service/threads", nil, "").Code; got != http.StatusUnauthorized {
		t.Fatalf("deleted service session=%d, want 401", got)
	}
	agents := call(platform, http.MethodGet, "/api/platform/v1/customer-service/agents", nil, "")
	if agents.Code != http.StatusOK || strings.Contains(agents.Body.String(), "987672006") {
		t.Fatalf("deleted service agent remains visible: code=%d body=%s", agents.Code, agents.Body.String())
	}
	if got := call(platform, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", gin.H{"target_admin_id": targetServiceID, "reason": "已删除客服不得成为新的转接目标。"}, "cs-deleted-target-001").Code; got != http.StatusBadRequest {
		t.Fatalf("deleted service transfer target=%d, want 400", got)
	}
	if got := call(platform, http.MethodPost, "/api/platform/v1/customer-service/threads/987671021/transfer", gin.H{"target_admin_id": serviceID, "reason": "平台将遗留会话转交给仍启用的中文客服张敏。"}, "cs-reassign-after-delete-001").Code; got != http.StatusOK {
		t.Fatalf("platform reassign deleted service thread=%d, want 200", got)
	}
	var historicalTargetCount int64
	if err := businessDB.Table("qixi_crm_b_customer_service_assignment_log").Where("binding_id = ? AND target_admin_id = ?", bindingA, targetServiceID).Count(&historicalTargetCount).Error; err != nil || historicalTargetCount != 1 {
		t.Fatalf("historical deleted-agent transfer audit=%d err=%v", historicalTargetCount, err)
	}
}
