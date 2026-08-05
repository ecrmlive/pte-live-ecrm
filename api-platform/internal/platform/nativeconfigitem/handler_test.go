package nativeconfigitem

import "testing"

func TestNormalizeRequiresName(t *testing.T) {
	_, _, _, _, _, _, err := normalize(upsertInput{Name: ""}, 1, 0)
	if err == nil {
		t.Fatal("empty name must be rejected")
	}
	name, code, remark, payload, status, sort, err := normalize(upsertInput{
		Name: "夏日香氛", Code: "summer-fragrance", Remark: "中文演示热搜", Status: intPtr(1), Sort: intPtr(10),
	}, 0, 0)
	if err != nil || name != "夏日香氛" || code != "summer-fragrance" || remark != "中文演示热搜" || payload == "" || status != 1 || sort != 10 {
		t.Fatalf("valid input rejected: name=%q err=%v", name, err)
	}
}

func intPtr(v int) *int { return &v }
