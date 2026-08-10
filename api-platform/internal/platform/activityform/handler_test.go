package activityform

import (
	"testing"
	"time"
)

func TestNormalizeRequiresCoreFields(t *testing.T) {
	_, _, _, _, _, _, _, _, _, _, _, err := normalize(upsertInput{Name: ""}, 1, 0)
	if err == nil || err.Error() != "请输入活动名称" {
		t.Fatalf("want 请输入活动名称 got %v", err)
	}
	_, _, _, _, _, _, _, _, _, _, _, err = normalize(upsertInput{
		Name: "新品内测报名", CoverURL: "https://example.com/c.png",
	}, 1, 0)
	if err == nil || err.Error() != "请上传活动分享海报" {
		t.Fatalf("want poster error got %v", err)
	}
	quota := 100
	status := 1
	sort := 10
	name, info, cover, poster, color, formID, q, st, so, starts, ends, err := normalize(upsertInput{
		Name: "新品内测报名", Info: "中文演示", CoverURL: "https://example.com/c.png",
		PosterURL: "https://example.com/p.png", Color: "#E8F5E9", FormID: 8105,
		Quota: &quota, Status: &status, Sort: &sort,
		StartsAt: "2026-08-01 00:00:00", EndsAt: "2026-12-31 23:59:59",
	}, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if name != "新品内测报名" || info != "中文演示" || cover == "" || poster == "" || color != "#E8F5E9" || formID != 8105 || q != 100 || st != 1 || so != 10 || starts == nil || ends == nil {
		t.Fatalf("unexpected normalize result")
	}
}

func TestTimeStatus(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	if timeStatus(&future, &future, 1) != 0 {
		t.Fatal("want 未开始")
	}
	if timeStatus(&past, &future, 1) != 1 {
		t.Fatal("want 进行中")
	}
	if timeStatus(&past, &past, 1) != -1 {
		t.Fatal("want 已结束")
	}
	if timeStatus(&past, &future, 0) != 1 {
		t.Fatal("is_show=0 still computes time status")
	}
}

func TestParseFormFields(t *testing.T) {
	fields := parseFormFields(`{"fields":[{"key":"attendee_name","label":"参会姓名","type":"text"},"legacy"]}`)
	if len(fields) != 2 || fields[0].Label != "参会姓名" || fields[1].Key != "legacy" {
		t.Fatalf("fields=%#v", fields)
	}
}

func TestSignupCountText(t *testing.T) {
	if signupCountText(2, 100) != "2/100" {
		t.Fatal(signupCountText(2, 100))
	}
	if signupCountText(3, 0) != "3/不限制" {
		t.Fatal(signupCountText(3, 0))
	}
}
