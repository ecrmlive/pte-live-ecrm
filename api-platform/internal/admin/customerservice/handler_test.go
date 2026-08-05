package customerservice

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestPositiveInt(t *testing.T) {
	if got := positiveInt("0", 20); got != 20 {
		t.Fatalf("zero got %d", got)
	}
	if got := positiveInt("101", 20); got != 101 {
		t.Fatalf("positive got %d", got)
	}
}

func TestHasRole(t *testing.T) {
	if !hasRole([]string{"customer_service"}, "customer_service") {
		t.Fatal("customer service role must be recognized")
	}
	if hasRole([]string{"operations"}, "customer_service") {
		t.Fatal("operations must not get customer service queue access")
	}
}

func TestQuickReplyValidationAndStoreScope(t *testing.T) {
	in := &quickReplyInput{StoreID: 1001, Title: "  发货时效  ", Content: "  您好，虚构店铺将在 48 小时内发货。  "}
	if !validQuickReply(in) {
		t.Fatal("valid Chinese quick reply must be accepted")
	}
	if in.Title != "发货时效" || in.Content != "您好，虚构店铺将在 48 小时内发货。" || in.Status != "enabled" {
		t.Fatalf("quick reply was not normalized: %#v", in)
	}
	if !includesStore([]uint64{1001, 1002}, 1002) || includesStore([]uint64{1001, 1002}, 2001) {
		t.Fatal("store scope must not expand beyond authorized stores")
	}
	if validQuickReply(&quickReplyInput{StoreID: 1001, Title: "", Content: "内容"}) {
		t.Fatal("blank title must be rejected")
	}
	if validQuickReply(&quickReplyInput{StoreID: 1001, Title: "退款", Content: "内容", Status: "deleted"}) {
		t.Fatal("unsupported status must be rejected")
	}
}

func TestUserNoteChineseLengthBoundary(t *testing.T) {
	valid := &userNoteInput{Content: "  虚构用户备注：已告知七禧演示茶铺的售后处理时效。  "}
	if !validUserNote(valid) || valid.Content != "虚构用户备注：已告知七禧演示茶铺的售后处理时效。" {
		t.Fatal("valid Chinese user note must be accepted and trimmed")
	}
	tooLong := make([]rune, 501)
	for index := range tooLong {
		tooLong[index] = '测'
	}
	if validUserNote(&userNoteInput{Content: string(tooLong)}) {
		t.Fatal("user note longer than 500 characters must be rejected")
	}
}

func TestValidateClaimRejectsConcurrentCustomerServiceOperator(t *testing.T) {
	owner := uint64(101)
	if err := validateClaim("open", nil, owner); err != nil {
		t.Fatalf("unclaimed open thread should be claimable: %v", err)
	}
	if err := validateClaim("open", &owner, owner); err != nil {
		t.Fatalf("same operator retry should be idempotent: %v", err)
	}
	if err := validateClaim("open", &owner, 102); !errors.Is(err, errThreadTaken) {
		t.Fatalf("other operator must not replace current owner, got %v", err)
	}
	if err := validateClaim("closed", nil, owner); !errors.Is(err, errThreadClosed) {
		t.Fatalf("closed thread must not be claimable, got %v", err)
	}
}

func TestValidTransferRequiresChineseReasonAndIdempotencyKey(t *testing.T) {
	in := &transferInput{TargetAdminID: 9302, Reason: "  虚构演示：转交居家商品咨询队列。  "}
	if !validTransfer(in, "fixture-cs-transfer-9900201") {
		t.Fatal("valid transfer should be accepted")
	}
	if in.Reason != "虚构演示：转交居家商品咨询队列。" {
		t.Fatalf("transfer reason was not normalized: %q", in.Reason)
	}
	if validTransfer(&transferInput{TargetAdminID: 9302, Reason: "原因"}, "") {
		t.Fatal("transfer without idempotency key must be rejected")
	}
	if validTransfer(&transferInput{TargetAdminID: 0, Reason: "原因"}, "fixture-transfer") {
		t.Fatal("transfer without a target operator must be rejected")
	}
}

func TestServiceSettingsChineseValidationAndQueueScope(t *testing.T) {
	settings := &serviceSettings{
		AutoReplyEnabled:    true,
		AutoReplyText:       "  您好，虚构演示客服将在工作时间内回复您。  ",
		QueueMode:           " round_robin ",
		MaxSessionsPerAgent: 12,
	}
	if !validServiceSettings(settings) {
		t.Fatal("valid Chinese service settings must be accepted")
	}
	if settings.AutoReplyText != "您好，虚构演示客服将在工作时间内回复您。" || settings.QueueMode != "round_robin" {
		t.Fatalf("service settings were not normalized: %#v", settings)
	}
	if validServiceSettings(&serviceSettings{AutoReplyEnabled: true, QueueMode: "manual", MaxSessionsPerAgent: 20}) {
		t.Fatal("enabled auto reply requires content")
	}
	if validServiceSettings(&serviceSettings{QueueMode: "manual", MaxSessionsPerAgent: 201}) {
		t.Fatal("session cap above 200 must be rejected")
	}
	if !sharesStore([]uint64{1001, 1002}, []uint64{2001, 1002}) || sharesStore([]uint64{1001}, []uint64{2001}) {
		t.Fatal("customer service roster must only share authorized stores")
	}
}

func TestServiceAgentUserMobileMask(t *testing.T) {
	if got := maskMobile("13900000001"); got != "139****0001" {
		t.Fatalf("masked mobile = %q", got)
	}
	if got := maskMobile("12345"); got != "" {
		t.Fatalf("short mobile must not leak: %q", got)
	}
}

func TestServiceAgentEmptyScopeUsesJSONArray(t *testing.T) {
	payload, err := json.Marshal(serviceAgent{ID: 1, ServiceStoreIDs: serviceStoreIDs(nil)})
	if err != nil || strings.Contains(string(payload), `"service_store_ids":null`) || !strings.Contains(string(payload), `"service_store_ids":[]`) {
		t.Fatalf("service queue scope JSON must be an array, payload=%s err=%v", payload, err)
	}
}
