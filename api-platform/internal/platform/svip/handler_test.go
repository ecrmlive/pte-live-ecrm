package svip

import "testing"

func TestMaskPhone(t *testing.T) {
	if got := maskPhone("13912340000"); got != "139****0000" {
		t.Fatalf("mask=%q", got)
	}
	if got := maskPhone("123"); got != "" {
		t.Fatalf("short phone=%q", got)
	}
}

func TestValidPlan(t *testing.T) {
	if !validPlan(&planInput{Name: "虚构演示会员", Price: 29, PlanType: "period", DurationDays: 30, Benefits: []string{"会员专享价"}, Status: 1, Sort: 10}) {
		t.Fatal("valid period plan rejected")
	}
	if !validPlan(&planInput{Name: "虚构体验会员", Price: 0, PlanType: "trial", DurationDays: 7, Benefits: []string{"体验权益"}, Status: 1}) {
		t.Fatal("valid trial plan rejected")
	}
	if validPlan(&planInput{Name: "错误试用", Price: 1, PlanType: "trial", DurationDays: 7, Benefits: []string{"体验权益"}, Status: 1}) || validPlan(&planInput{Name: "错误永久", Price: 99, PlanType: "lifetime", DurationDays: 30, Benefits: []string{"权益"}, Status: 1}) {
		t.Fatal("invalid plan accepted")
	}
}

func TestNormalizedBenefitsRejectsDuplicates(t *testing.T) {
	in := &planInput{Name: "重复权益", Price: 29, PlanType: "period", DurationDays: 30, Benefits: []string{"会员专享价", "会员专享价"}, Status: 1}
	if validPlan(in) {
		t.Fatal("duplicate benefits accepted")
	}
}
