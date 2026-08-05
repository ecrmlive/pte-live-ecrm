package userlist

import (
	"encoding/json"
	"testing"
)

func TestMaskAndPagingFailSafe(t *testing.T) {
	if got := mask("13500000001"); got != "135****0001" {
		t.Fatalf("masked mobile=%q", got)
	}
	if got := mask("1234"); got != "" {
		t.Fatalf("short mobile must not be returned: %q", got)
	}
}

func TestPositiveID(t *testing.T) {
	if id, ok := positiveID(" 9101 "); !ok || id != 9101 {
		t.Fatalf("valid id=%d ok=%v", id, ok)
	}
	for _, raw := range []string{"", "0", "-1", "abc"} {
		if _, ok := positiveID(raw); ok {
			t.Fatalf("invalid user id %q accepted", raw)
		}
	}
}

func TestOptionalPositiveQueryParser(t *testing.T) {
	// The HTTP handler only accepts an optional strictly-positive ID; keep this
	// parsing rule separate from the DB query to prevent accidental broad queries.
	if id, ok := positiveID("3002"); !ok || id != 3002 {
		t.Fatalf("coupon id=%d ok=%v", id, ok)
	}
	for _, raw := range []string{"0", "-3", "优惠券"} {
		if _, ok := positiveID(raw); ok {
			t.Fatalf("invalid coupon filter %q accepted", raw)
		}
	}
}

func TestValidAdjustment(t *testing.T) {
	valid := assetAdjustmentInput{AssetType: "balance", Amount: json.Number("-0.01"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"}
	if !validAdjustment(valid) {
		t.Fatal("valid balance adjustment rejected")
	}
	valid.AssetType, valid.Amount = "points", json.Number("5")
	if !validAdjustment(valid) {
		t.Fatal("valid point adjustment rejected")
	}
	for _, invalid := range []assetAdjustmentInput{
		{AssetType: "commission", Amount: json.Number("1"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "points", Amount: json.Number("1.5"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "balance", Amount: json.Number("0"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "balance", Amount: json.Number("1.001"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "balance", Amount: json.Number("1e2"), Reason: "虚构中文工单", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "balance", Amount: json.Number("1"), Reason: "x", IdempotencyKey: "adjust-9101-001"},
		{AssetType: "balance", Amount: json.Number("1"), Reason: "虚构中文工单", IdempotencyKey: "short"},
	} {
		if validAdjustment(invalid) {
			t.Fatalf("invalid adjustment accepted: %#v", invalid)
		}
	}
}

func TestAdjustmentMinorUsesFixedCents(t *testing.T) {
	for _, item := range []struct {
		input assetAdjustmentInput
		want  int64
	}{
		{assetAdjustmentInput{AssetType: "balance", Amount: json.Number("0.10")}, 10},
		{assetAdjustmentInput{AssetType: "balance", Amount: json.Number("0.2")}, 20},
		{assetAdjustmentInput{AssetType: "balance", Amount: json.Number("-36.50")}, -3650},
		{assetAdjustmentInput{AssetType: "points", Amount: json.Number("268")}, 268},
	} {
		if got, ok := adjustmentMinor(item.input); !ok || got != item.want {
			t.Fatalf("adjustmentMinor(%+v)=%d,%v want %d", item.input, got, ok, item.want)
		}
	}
	if got := formatMinor(30); got != "0.30" {
		t.Fatalf("formatMinor=%q", got)
	}
}

func TestNullableMemberLevel(t *testing.T) {
	if got := nullableLevel(0); got != nil {
		t.Fatalf("ordinary member must be stored as NULL, got %#v", got)
	}
	if got := nullableLevel(8102); got != uint64(8102) {
		t.Fatalf("member level mapping=%#v", got)
	}
}

func TestNormalizeUserIDs(t *testing.T) {
	got, ok := normalizeUserIDs([]uint64{9102, 9101, 9102})
	if !ok || len(got) != 2 || got[0] != 9101 || got[1] != 9102 {
		t.Fatalf("normalized ids=%v ok=%v", got, ok)
	}
	for _, ids := range [][]uint64{nil, {0}, make([]uint64, 101)} {
		if _, ok := normalizeUserIDs(ids); ok {
			t.Fatalf("invalid ids accepted: %v", ids)
		}
	}
}

func TestSameUserIDsJSON(t *testing.T) {
	if !sameUserIDsJSON("[9102, 9101, 9102]", []uint64{9101, 9102}) {
		t.Fatal("JSON user ids should compare by normalized meaning")
	}
	if sameUserIDsJSON("not-json", []uint64{9101}) || sameUserIDsJSON("[9101]", []uint64{9102}) {
		t.Fatal("invalid or different JSON user ids accepted")
	}
}

func TestNormalizeLabelIDs(t *testing.T) {
	got, ok := normalizeLabelIDs([]uint64{9402, 9401, 9402})
	if !ok || len(got) != 2 || got[0] != 9401 || got[1] != 9402 {
		t.Fatalf("normalized label ids=%v ok=%v", got, ok)
	}
	if got, ok := normalizeLabelIDs(nil); !ok || len(got) != 0 {
		t.Fatalf("empty label list must mean clear, got=%v ok=%v", got, ok)
	}
	if _, ok := normalizeLabelIDs([]uint64{0}); ok {
		t.Fatal("zero label id accepted")
	}
}

func TestUserAdminCommandValidation(t *testing.T) {
	if !validUserCommand("虚构中文工单", "user-command-9101") || validUserCommand("x", "user-command-9101") || validUserCommand("虚构中文工单", "short") {
		t.Fatal("user admin command validation mismatch")
	}
	for _, value := range []string{"", "/demo/avatar.png", "https://example.invalid/avatar.png"} {
		if !validAvatarURL(value) {
			t.Fatalf("valid avatar URL rejected: %q", value)
		}
	}
	if validAvatarURL("http://example.invalid/avatar.png") || validAvatarURL("javascript:alert(1)") {
		t.Fatal("unsafe avatar URL accepted")
	}
	if fingerprint("create", "DEMO-PC-USER", "七禧演示") != fingerprint("create", "DEMO-PC-USER", "七禧演示") {
		t.Fatal("request fingerprint must be stable")
	}
}

func TestCSVCellPreventsSpreadsheetFormulaInjection(t *testing.T) {
	for raw, want := range map[string]string{
		"=SUM(1,1)": "'=SUM(1,1)",
		"+123":      "'+123",
		"-123":      "'-123",
		"@DDE":      "'@DDE",
		"七禧体验用户":    "七禧体验用户",
	} {
		if got := csvCell(raw); got != want {
			t.Fatalf("csv cell %q=%q, want %q", raw, got, want)
		}
	}
	if exportStatusValue(nil) != -1 {
		t.Fatal("empty export status must have a stable fingerprint marker")
	}
}
