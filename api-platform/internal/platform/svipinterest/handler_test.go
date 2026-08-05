package svipinterest

import "testing"

func TestInterestValidation(t *testing.T) {
	if !valid(&input{Name: "会员专享价", Description: "购买指定商品享受会员专享价格", IconURL: "/demo/svip-price.png", Status: 1, Sort: 10}, false) {
		t.Fatal("valid interest rejected")
	}
	if valid(&input{Name: "错误图标", IconURL: "http://unsafe.example/icon.png", Status: 1}, false) {
		t.Fatal("unsafe icon accepted")
	}
	if valid(&input{Name: "更新", Status: 1}, true) {
		t.Fatal("zero version update accepted")
	}
}
