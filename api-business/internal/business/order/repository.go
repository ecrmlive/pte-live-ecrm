package order

import (
	"context"
	"errors"
	"math"

	"gorm.io/gorm"
)

var ErrCartOwnership = errors.New("购物车商品不存在或无权访问")

// LoadCheckout reads only business-owned cart data and the merchant-published
// consumption view. It never joins a qixi_crm_m_* table.
func LoadCheckout(ctx context.Context, db *gorm.DB, userID uint64, cartIDs []uint64) (Checkout, error) {
	ids, err := uniqueCartIDs(cartIDs)
	if err != nil {
		return Checkout{}, err
	}
	var rows []checkoutRow
	err = db.WithContext(ctx).Table("qixi_crm_b_cart AS c").
		Select("c.id AS cart_id,c.product_id,c.sku_key,c.quantity,ps.merchant_sku_id,ps.spec_snapshot,p.merchant_id,p.merchant_name,p.store_id,p.store_name,p.title,p.cover_url,ps.price,ps.stock,ps.sale_status,p.product_type,s.integral_enabled,s.integral_points_per_yuan,s.integral_max_deduction_bps,p.svip_price_type,p.svip_price").
		Joins("INNER JOIN qixi_crm_b_product_view AS p ON p.product_id = c.product_id").
		Joins("INNER JOIN qixi_crm_b_product_sku_view AS ps ON ps.product_id = c.product_id AND ps.sku_key = c.sku_key").
		Joins("INNER JOIN qixi_crm_b_store_view AS s ON s.store_id = p.store_id AND s.status = 1").
		Where("c.user_id = ? AND c.id IN ?", userID, ids).Order("c.id ASC").Scan(&rows).Error
	if err != nil {
		return Checkout{}, err
	}
	if len(rows) != len(ids) {
		return Checkout{}, ErrCartOwnership
	}
	svipActive, err := activeSVIP(ctx, db, userID)
	if err != nil {
		return Checkout{}, err
	}
	lines := make([]CartLine, 0, len(rows))
	for _, row := range rows {
		cents, err := cents(row.Price)
		if err != nil {
			return Checkout{}, err
		}
		if row.MerchantSKUID == 0 {
			return Checkout{}, ErrUnavailableCart
		}
		memberCents := ResolveSVIPUnitCents(cents, row.ProductType, svipActive, row.SVIPPriceType, row.SVIPPrice)
		lines = append(lines, CartLine{CartID: row.CartID, ProductID: row.ProductID, MerchantSKUID: row.MerchantSKUID, SKUKey: row.SKUKey, SpecSnapshot: normalizeSpecSnapshot(row.SpecSnapshot), MerchantID: row.MerchantID, MerchantName: row.MerchantName, StoreID: row.StoreID, StoreName: row.StoreName, Title: row.Title, CoverURL: row.CoverURL, ListCents: cents, UnitCents: memberCents, SVIPApplied: memberCents < cents, IntegralEnabled: row.IntegralEnabled, IntegralPointsPerYuan: row.IntegralPointsPerYuan, IntegralMaxDeductionBps: row.IntegralMaxDeductionBps, Quantity: row.Quantity, Stock: row.Stock, SaleStatus: row.SaleStatus, ProductType: row.ProductType})
	}
	return BuildCheckout(lines)
}

type checkoutRow struct {
	CartID                  uint64  `gorm:"column:cart_id"`
	ProductID               uint64  `gorm:"column:product_id"`
	MerchantSKUID           uint64  `gorm:"column:merchant_sku_id"`
	SKUKey                  string  `gorm:"column:sku_key"`
	SpecSnapshot            string  `gorm:"column:spec_snapshot"`
	Quantity                int     `gorm:"column:quantity"`
	MerchantID              uint64  `gorm:"column:merchant_id"`
	MerchantName            string  `gorm:"column:merchant_name"`
	StoreID                 uint64  `gorm:"column:store_id"`
	StoreName               string  `gorm:"column:store_name"`
	Title                   string  `gorm:"column:title"`
	CoverURL                string  `gorm:"column:cover_url"`
	Price                   float64 `gorm:"column:price"`
	Stock                   int     `gorm:"column:stock"`
	SaleStatus              int8    `gorm:"column:sale_status"`
	ProductType             int8    `gorm:"column:product_type"`
	SVIPPriceType           int8    `gorm:"column:svip_price_type"`
	SVIPPrice               float64 `gorm:"column:svip_price"`
	IntegralEnabled         bool    `gorm:"column:integral_enabled"`
	IntegralPointsPerYuan   int64   `gorm:"column:integral_points_per_yuan"`
	IntegralMaxDeductionBps int64   `gorm:"column:integral_max_deduction_bps"`
}

func uniqueCartIDs(ids []uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return nil, ErrEmptyCart
	}
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, ErrCartOwnership
		}
		if _, ok := seen[id]; ok {
			return nil, ErrCartOwnership
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

func cents(amount float64) (int64, error) {
	if amount < 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, ErrUnavailableCart
	}
	return int64(math.Round(amount * 100)), nil
}
