package nativedistribution

import "testing"

func TestDistributionFiltersFailClosed(t *testing.T) {
	for _, value := range []string{"", "0", "1"} {
		if _, ok := promoterStatus(value); !ok {
			t.Fatalf("promoter status %q should be allowed", value)
		}
	}
	if _, ok := promoterStatus("启用"); ok {
		t.Fatal("arbitrary promoter status must be rejected")
	}
	for _, value := range []string{"", "pending", "available", "settled", "voided"} {
		if _, ok := commissionStatus(value); !ok {
			t.Fatalf("commission status %q should be allowed", value)
		}
	}
	if _, ok := commissionStatus("已结算"); ok {
		t.Fatal("arbitrary commission status must be rejected")
	}
}

func TestNormalizeLevelSaveRequiresTask(t *testing.T) {
	in := levelSaveInput{Name: "测试", Rank: 2, ExtensionOne: 0, ExtensionTwo: 0}
	if _, msg := normalizeLevelSave(in); msg == "" {
		t.Fatal("level without tasks must be rejected")
	}
	in.TaskRule.SpreadUser = levelTaskItem{Name: "邀请", Num: 3}
	if _, msg := normalizeLevelSave(in); msg != "" {
		t.Fatalf("valid level save rejected: %s", msg)
	}
	in.TaskRule.PayMoney = levelTaskItem{Name: "有名无量", Num: 0}
	if _, msg := normalizeLevelSave(in); msg == "" {
		t.Fatal("mismatched task name/num must be rejected")
	}
}

func TestParseLevelTaskRule(t *testing.T) {
	rule := parseLevelTaskRule(`{"spread_user":{"name":"邀请新人","num":3,"info":""},"pay_money":{"name":"","num":0,"info":""}}`)
	if rule.SpreadUser.Name != "邀请新人" || rule.SpreadUser.Num != 3 {
		t.Fatalf("unexpected parsed rule: %+v", rule.SpreadUser)
	}
}

func TestNormalizePrivilegeSave(t *testing.T) {
	in := privilegeSaveInput{Title: " 高佣金 ", ImgURL: "https://example.com/a.png"}
	title, img, status, sort, msg := normalizePrivilegeSave(in, 1, 0)
	if msg != "" || title != "高佣金" || img == "" || status != 1 || sort != 0 {
		t.Fatalf("unexpected privilege save: title=%q img=%q status=%d sort=%d msg=%q", title, img, status, sort, msg)
	}
	in.Title = ""
	if _, _, _, _, msg = normalizePrivilegeSave(in, 1, 0); msg == "" {
		t.Fatal("empty title must be rejected")
	}
	in.Title = "高佣金"
	in.ImgURL = ""
	if _, _, _, _, msg = normalizePrivilegeSave(in, 1, 0); msg == "" {
		t.Fatal("empty image must be rejected")
	}
}

func TestNormalizePosterSave(t *testing.T) {
	in := posterSaveInput{Name: " 618 ", PicURL: "https://example.com/p.png"}
	name, pic, status, sort, msg := normalizePosterSave(in, 1, 0)
	if msg != "" || name != "618" || pic == "" || status != 1 || sort != 0 {
		t.Fatalf("unexpected poster save: name=%q pic=%q status=%d sort=%d msg=%q", name, pic, status, sort, msg)
	}
	in.Name = ""
	if _, _, _, _, msg = normalizePosterSave(in, 1, 0); msg == "" {
		t.Fatal("empty name must be rejected")
	}
	in.Name = "618"
	in.PicURL = ""
	if _, _, _, _, msg = normalizePosterSave(in, 1, 0); msg == "" {
		t.Fatal("empty pic must be rejected")
	}
}
