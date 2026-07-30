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
		Select("c.id AS cart_id,c.product_id,c.sku_key,c.quantity,p.merchant_id,p.merchant_name,p.store_id,p.store_name,p.title,p.cover_url,p.price,p.stock,p.sale_status,p.product_type").
		Joins("INNER JOIN qixi_crm_b_product_view AS p ON p.product_id = c.product_id").
		Where("c.user_id = ? AND c.id IN ?", userID, ids).Order("c.id ASC").Scan(&rows).Error
	if err != nil {
		return Checkout{}, err
	}
	if len(rows) != len(ids) {
		return Checkout{}, ErrCartOwnership
	}
	lines := make([]CartLine, 0, len(rows))
	for _, row := range rows {
		cents, err := cents(row.Price)
		if err != nil {
			return Checkout{}, err
		}
		lines = append(lines, CartLine{CartID: row.CartID, ProductID: row.ProductID, SKUKey: row.SKUKey, MerchantID: row.MerchantID, MerchantName: row.MerchantName, StoreID: row.StoreID, StoreName: row.StoreName, Title: row.Title, CoverURL: row.CoverURL, UnitCents: cents, Quantity: row.Quantity, Stock: row.Stock, SaleStatus: row.SaleStatus, ProductType: row.ProductType})
	}
	return BuildCheckout(lines)
}

type checkoutRow struct {
	CartID       uint64  `gorm:"column:cart_id"`
	ProductID    uint64  `gorm:"column:product_id"`
	SKUKey       string  `gorm:"column:sku_key"`
	Quantity     int     `gorm:"column:quantity"`
	MerchantID   uint64  `gorm:"column:merchant_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	StoreID      uint64  `gorm:"column:store_id"`
	StoreName    string  `gorm:"column:store_name"`
	Title        string  `gorm:"column:title"`
	CoverURL     string  `gorm:"column:cover_url"`
	Price        float64 `gorm:"column:price"`
	Stock        int     `gorm:"column:stock"`
	SaleStatus   int8    `gorm:"column:sale_status"`
	ProductType  int8    `gorm:"column:product_type"`
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
