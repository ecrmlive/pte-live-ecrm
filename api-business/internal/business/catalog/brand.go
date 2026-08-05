package catalog

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

var errInvalidBrand = errors.New("品牌名称错误")

type brandView struct {
	Name         string `gorm:"column:brand_name"`
	ProductCount int64  `gorm:"column:product_count"`
}

// Brands is sourced from the published product view. Brand labels are never
// accepted from the browser as a product mutation and the API remains inside
// the business database boundary.
func (h *Handler) Brands(c *gin.Context) {
	scope, err := h.resolveStoreScope(c)
	if err != nil {
		writeScopeError(c, err)
		return
	}
	query := sellableProductQuery(h.db.WithContext(c.Request.Context()).Model(&productView{}).
		Where("sale_status = 1 AND brand_name <> ''"))
	if scope.MerchantID != 0 {
		query = query.Where("merchant_id = ?", scope.MerchantID)
	}
	rows := make([]brandView, 0)
	if err := query.Select("brand_name, COUNT(*) AS product_count").Group("brand_name").Order("product_count DESC, brand_name ASC").Scan(&rows).Error; err != nil {
		fail(c, err)
		return
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, map[string]any{"name": row.Name, "product_count": row.ProductCount})
	}
	response.OK(c, map[string]any{"list": items})
}

func brandFilter(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if len([]rune(value)) > 64 || strings.ContainsAny(value, "\x00\r\n") {
		return "", errInvalidBrand
	}
	return value, nil
}

func writeBrandError(c *gin.Context, err error) bool {
	if errors.Is(err, errInvalidBrand) {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return true
	}
	return false
}
