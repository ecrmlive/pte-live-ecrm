package cart

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
)

func TestResolveSKUKeyPreservesCRMEBSKUUnique(t *testing.T) {
	if got := resolveSKUKey(42, "颜色-红色:尺码-L"); got != "颜色-红色:尺码-L" {
		t.Fatalf("resolveSKUKey() = %q, want original unique key", got)
	}
	if got := resolveSKUKey(42, "  "); got != "42" {
		t.Fatalf("resolveSKUKey() fallback = %q, want product id", got)
	}
}

func TestItemResponseSeparatesProductAndStoreNames(t *testing.T) {
	item := itemResponse(cartView{
		CartID: 1, ProductID: 2, SKUKey: "sku-red-l", Quantity: 3,
		MerchantID: 4, MerchantName: "七禧商户", StoreName: "七禧旗舰店",
		Title: "夏季短袖", CoverURL: "https://cdn.example.invalid/p.png",
		Price: 99.5, Stock: 3, SaleStatus: 1,
	})
	if item["title"] != "夏季短袖" || item["store_name"] != "七禧旗舰店" {
		t.Fatalf("product/store names are mixed: %#v", item)
	}
	if item["product_attr_unique"] != "sku-red-l" || item["is_fail"] != 0 {
		t.Fatalf("unexpected cart item: %#v", item)
	}

	outOfStock := itemResponse(cartView{Quantity: 2, Stock: 1, SaleStatus: 1})
	if outOfStock["is_fail"] != 1 {
		t.Fatalf("out-of-stock item must be unavailable: %#v", outOfStock)
	}
}

func TestVerifyMerchantContextRequiresSignedStoreBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	product := productRow{MerchantID: 8, StoreID: 12}

	newContext := func(appID string, claims *authjwt.Claims) *gin.Context {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("POST", "/cart", nil)
		if appID != "" {
			ctx.Request.Header.Set("X-AppId", appID)
		}
		if claims != nil {
			ctx.Set(middleware.CtxClaimsKey, claims)
		}
		return ctx
	}

	if err := verifyMerchantContext(newContext("", nil), product); err != errStoreContext {
		t.Fatalf("missing context error = %v, want %v", err, errStoreContext)
	}
	if err := verifyMerchantContext(newContext("qixi.store.8", &authjwt.Claims{
		AuthContext: authjwt.ContextStore, MerchantID: 8, StoreID: 12, MerchantAppID: "qixi.store.other",
	}), product); err != errStoreContext {
		t.Fatalf("mismatched app id error = %v, want %v", err, errStoreContext)
	}
	if err := verifyMerchantContext(newContext("qixi.store.8", &authjwt.Claims{
		AuthContext: authjwt.ContextStore, MerchantID: 8, StoreID: 12, MerchantAppID: "qixi.store.8",
	}), product); err != nil {
		t.Fatalf("matching store context rejected: %v", err)
	}
	if err := verifyMerchantContext(newContext("", nil), productRow{}); err != nil {
		t.Fatalf("platform product must not require a merchant context: %v", err)
	}
}
