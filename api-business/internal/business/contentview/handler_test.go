package contentview

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAgreementLabel(t *testing.T) {
	label, ok := agreementLabel("sys_user_agree")
	if !ok || label != "用户协议" {
		t.Fatalf("agreementLabel(sys_user_agree) = (%q, %v)", label, ok)
	}
	label, ok = agreementLabel("sys_coupon_agree")
	if !ok || label != "优惠券使用说明" {
		t.Fatalf("agreementLabel(sys_coupon_agree) = (%q, %v)", label, ok)
	}
	label, ok = agreementLabel("sys_extension_agree")
	if !ok || label != "佣金说明" {
		t.Fatalf("agreementLabel(sys_extension_agree) = (%q, %v)", label, ok)
	}
	label, ok = agreementLabel("sys_brokerage")
	if !ok || label != "分销等级规则" {
		t.Fatalf("agreementLabel(sys_brokerage) = (%q, %v)", label, ok)
	}
	if _, ok := agreementLabel("not-a-contract"); ok {
		t.Fatal("unknown agreement key must be rejected")
	}
}

func TestArticleResponse(t *testing.T) {
	row := contentView{ContentID: 2201, Title: "夏日焕新购物指南"}
	response := contentResponse(row, "article")
	if response["article_id"] != uint64(2201) || response["title"] != "夏日焕新购物指南" {
		t.Fatalf("unexpected article response: %#v", response)
	}
}

func TestNoticeResponse(t *testing.T) {
	row := contentView{ContentID: 2001, Title: "七禧商城服务公告"}
	response := noticeResponse(row)
	if response["notice_id"] != uint64(2001) || response["title"] != "七禧商城服务公告" {
		t.Fatalf("unexpected notice response: %#v", response)
	}
}

func TestPageParamsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		query string
		page  int
		limit int
	}{
		{query: "?page=-1&limit=101", page: 1, limit: 20},
		{query: "?page=2&limit=15", page: 2, limit: 15},
	} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodGet, "/notices"+tc.query, nil)
		page, limit := pageParams(ctx)
		if page != tc.page || limit != tc.limit {
			t.Fatalf("query %s: page=%d limit=%d, want %d,%d", tc.query, page, limit, tc.page, tc.limit)
		}
	}
}
