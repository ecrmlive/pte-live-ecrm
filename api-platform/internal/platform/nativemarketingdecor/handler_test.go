package nativemarketingdecor

import "testing"

func TestNormalizeRequiresName(t *testing.T) {
	_, _, _, _, _, _, _, _, _, err := normalize(upsertInput{Name: ""}, 1, 0)
	if err == nil {
		t.Fatal("empty name must be rejected")
	}
	name, code, _, remark, payload, status, sort, _, _, err := normalize(upsertInput{
		Name: "夏日焕新氛围图", Code: "summer-atmosphere", Remark: "中文演示", Status: intPtr(1), Sort: intPtr(10),
	}, 0, 0)
	if err != nil || name != "夏日焕新氛围图" || code != "summer-atmosphere" || remark != "中文演示" || payload == "" || status != 1 || sort != 10 {
		t.Fatalf("valid input rejected: %#v err=%v", name, err)
	}
}

func intPtr(v int) *int { return &v }
