package memberlevel

import "testing"

func TestMembershipLevelValidation(t *testing.T) {
	good := &input{
		Name:        "白银会员",
		Rank:        2,
		IconURL:     "https://example.com/silver.png",
		GrowthValue: 100,
		BgImage:     "https://example.com/bg.png",
		Status:      1,
	}
	if !valid(good, false) {
		t.Fatal("valid member level rejected")
	}
	if valid(&input{Name: "", Rank: 1, GrowthValue: 0, Status: 1}, false) {
		t.Fatal("empty name accepted")
	}
	if valid(&input{Name: "更新", Rank: 1, GrowthValue: 0, Status: 1}, true) {
		t.Fatal("zero version update accepted")
	}
	if valid(&input{Name: "负成长", Rank: 1, GrowthValue: -1, Status: 1}, false) {
		t.Fatal("negative growth accepted")
	}
}

func TestParseRulesGrowthValue(t *testing.T) {
	got := parseRules(`{"value":500,"image":"/bg.png","description":"x"}`)
	if got.Value != 500 || got.Image != "/bg.png" {
		t.Fatalf("parseRules=%+v", got)
	}
}
