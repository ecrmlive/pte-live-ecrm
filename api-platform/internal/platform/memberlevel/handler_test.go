package memberlevel

import "testing"

func TestMembershipLevelValidation(t *testing.T) {
	good := &input{Name: "悦享会员", Rank: 2, Rules: `{"description":"成长值满 100 自动升级"}`, Benefits: `["专属优惠提醒","会员活动优先参与"]`, Status: 1}
	if !valid(good, false) {
		t.Fatal("valid member level rejected")
	}
	bad := &input{Name: "重复权益", Rank: 2, Rules: `{}`, Benefits: `["权益","权益"]`, Status: 1}
	if valid(bad, false) {
		t.Fatal("invalid member level accepted")
	}
	if valid(&input{Name: "更新", Rank: 1, Rules: `{"description":"x"}`, Benefits: `["权益"]`, Status: 1}, true) {
		t.Fatal("zero version update accepted")
	}
}
