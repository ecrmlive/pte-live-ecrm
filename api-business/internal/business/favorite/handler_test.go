package favorite

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductFavoriteViewKeepsChineseFields(t *testing.T) {
	got := (productFavorite{
		ProductID: 12, MerchantID: 3, StoreID: 8,
		MerchantName: "星河商贸", StoreName: "星河优选店", Title: "夏日轻薄防晒衣",
		CoverURL: "https://example.invalid/demo/sunscreen.png", Price: 89.9, Sales: 26, Stock: 9,
	}).view()
	if got["id"] != uint64(12) || got["title"] != "夏日轻薄防晒衣" || got["store_name"] != "星河优选店" || got["sales"] != 26 {
		t.Fatalf("unexpected product favorite view: %#v", got)
	}
}

func TestStoreFavoriteViewKeepsMerchantContext(t *testing.T) {
	got := (storeFavorite{StoreID: 8, MerchantID: 3, StoreName: "星河优选店", StoreAppID: "crm.store.3", FollowerCount: 15}).view()
	if got["store_id"] != uint64(8) || got["mer_id"] != uint64(3) || got["merchant_app_id"] != "crm.store.3" || got["follower_count"] != int64(15) {
		t.Fatalf("unexpected store favorite view: %#v", got)
	}
}

func TestPathIDRejectsInvalidOrZero(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, raw := range []string{"", "0", "-1", "abc"} {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = gin.Params{{Key: "id", Value: raw}}
		if _, ok := pathID(ctx); ok {
			t.Fatalf("path id %q must be rejected", raw)
		}
	}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Params = gin.Params{{Key: "id", Value: "72"}}
	if id, ok := pathID(ctx); !ok || id != 72 {
		t.Fatalf("valid id parsed as %d, %v", id, ok)
	}
}

func TestFavoritesPaginationBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/?page=-1&limit=99", nil)
	page, limit := favoritesPagination(ctx)
	if page != 1 || limit != 50 {
		t.Fatalf("page=%d limit=%d", page, limit)
	}
}
