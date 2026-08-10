package nativemarketingdecor

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNormalizeRequiresName(t *testing.T) {
	_, _, _, _, _, _, _, _, _, err := normalize(upsertInput{Name: ""}, "topic", 1, 0)
	if err == nil {
		t.Fatal("empty name must be rejected")
	}
	_, _, _, _, _, _, _, _, _, err = normalize(upsertInput{
		Name: "居家香氛专场", Code: "home-fragrance-topic", CoverURL: "https://example.com/list.png",
		Payload: map[string]any{}, Status: intPtr(1), Sort: intPtr(10),
	}, "topic", 0, 0)
	if err == nil || err.Error() != "请选择关联标签" {
		t.Fatalf("topic without label_id want 请选择关联标签 got %v", err)
	}
	name, code, cover, remark, payload, status, sort, _, _, err := normalize(upsertInput{
		Name: "居家香氛专场", Code: "home-fragrance-topic", CoverURL: "https://example.com/list.png",
		Remark: "中文演示", Status: intPtr(1), Sort: intPtr(10),
		Payload: map[string]any{
			"label_id": 7502,
			"banner":   []any{"https://example.com/b1.png", ""},
			"image":    "https://example.com/theme.png",
			"color":    "#F5E6D3",
			"type":     2,
		},
	}, "topic", 0, 0)
	if err != nil || name != "居家香氛专场" || code != "home-fragrance-topic" || cover == "" || remark != "中文演示" || payload == "" || status != 1 || sort != 10 {
		t.Fatalf("valid input rejected: %#v err=%v", name, err)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["scope_type"]; ok {
		t.Fatalf("non-scope decor payload should not force scope_type: %#v", m)
	}
	if v, _ := asInt(m["label_id"]); v != 7502 {
		t.Fatalf("label_id=%v", m["label_id"])
	}
	if v, _ := asInt(m["type"]); v != 2 {
		t.Fatalf("type=%v", m["type"])
	}
}

func TestNormalizeAtmosphereRequiresTimeAndScope(t *testing.T) {
	_, _, _, _, _, _, _, _, _, err := normalize(upsertInput{
		Name: "测试氛围", ScopeType: intPtr(0),
	}, "atmosphere", 1, 0)
	if err == nil {
		t.Fatal("atmosphere without time must fail")
	}
	_, _, _, _, _, _, _, _, _, err = normalize(upsertInput{
		Name: "测试氛围", ScopeType: intPtr(1), StartsAt: "2026-08-01 00:00:00", EndsAt: "2026-08-31 23:59:59",
	}, "atmosphere", 1, 0)
	if err == nil || err.Error() != "请选择指定商品" {
		t.Fatalf("want 请选择指定商品 got %v", err)
	}
	_, _, _, _, payload, _, _, _, _, err := normalize(upsertInput{
		Name: "测试氛围", ScopeType: intPtr(1), SpuIDs: []uint64{11, 11, 22},
		StartsAt: "2026-08-01 00:00:00", EndsAt: "2026-08-31 23:59:59", Status: intPtr(1),
	}, "atmosphere", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(payload), &m)
	ids := uniqIDs(asUint64Slice(m["spu_ids"]))
	if len(ids) != 2 || ids[0] != 11 || ids[1] != 22 {
		t.Fatalf("spu_ids=%v", ids)
	}
}

func TestNormalizeBorderRequiresTimeAndScope(t *testing.T) {
	_, _, _, _, _, _, _, _, _, err := normalize(upsertInput{
		Name: "测试边框", ScopeType: intPtr(0),
	}, "border", 1, 0)
	if err == nil {
		t.Fatal("border without time must fail")
	}
	_, _, _, _, payload, _, _, _, _, err := normalize(upsertInput{
		Name: "测试边框", ScopeType: intPtr(0),
		StartsAt: "2026-08-01 00:00:00", EndsAt: "2026-08-31 23:59:59", Status: intPtr(1),
	}, "border", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(payload), &m)
	if v, _ := asInt(m["scope_type"]); v != 0 {
		t.Fatalf("scope_type=%v", m["scope_type"])
	}
}

func TestCalcActivityStatus(t *testing.T) {
	now := time.Date(2026, 8, 10, 12, 0, 0, 0, time.Local)
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	if got := calcActivityStatus(&future, nil, now); got != 0 {
		t.Fatalf("want 未开始=0 got %d", got)
	}
	if got := calcActivityStatus(&past, &future, now); got != 1 {
		t.Fatalf("want 进行中=1 got %d", got)
	}
	if got := calcActivityStatus(nil, &past, now); got != -1 {
		t.Fatalf("want 已结束=-1 got %d", got)
	}
	if activityStatusText(1) != "进行中" {
		t.Fatal("activityStatusText")
	}
}

func intPtr(v int) *int { return &v }
