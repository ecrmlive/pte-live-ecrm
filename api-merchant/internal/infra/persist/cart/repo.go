package cartpersist

import (
	"context"
	"errors"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/cart"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/catalog"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) FindActive(ctx context.Context, uid uint, productID uint, unique string) (*cart.Cart, error) {
	var row cart.Cart
	err := r.db.WithContext(ctx).
		Where("uid = ? AND product_id = ? AND product_attr_unique = ? AND is_pay = 0 AND is_del = 0 AND is_new = 0",
			uid, productID, unique).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, c *cart.Cart) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateNum(ctx context.Context, cartID uint64, uid uint, num uint) error {
	return r.db.WithContext(ctx).Model(&cart.Cart{}).
		Where("cart_id = ? AND uid = ? AND is_del = 0 AND is_pay = 0", cartID, uid).
		Update("cart_num", num).Error
}

func (r *Repo) SoftDelete(ctx context.Context, cartID uint64, uid uint) error {
	return r.db.WithContext(ctx).Model(&cart.Cart{}).
		Where("cart_id = ? AND uid = ?", cartID, uid).
		Update("is_del", 1).Error
}

func (r *Repo) ListByUID(ctx context.Context, uid uint) ([]cart.Cart, error) {
	var rows []cart.Cart
	err := r.db.WithContext(ctx).
		Where("uid = ? AND is_pay = 0 AND is_del = 0", uid).
		Order("cart_id DESC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetByID(ctx context.Context, cartID uint64, uid uint) (*cart.Cart, error) {
	var row cart.Cart
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND uid = ? AND is_del = 0", cartID, uid).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListByIDs(ctx context.Context, uid uint, ids []uint64) ([]cart.Cart, error) {
	var rows []cart.Cart
	err := r.db.WithContext(ctx).
		Where("uid = ? AND is_pay = 0 AND is_del = 0 AND cart_id IN ?", uid, ids).
		Find(&rows).Error
	return rows, err
}

func (r *Repo) MarkPaid(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&cart.Cart{}).
		Where("cart_id IN ?", ids).
		Update("is_pay", 1).Error
}

func (r *Repo) ResolveProduct(ctx context.Context, productID uint, unique string) (*cart.ProductView, error) {
	var p catalog.Product
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_del = 0", productID).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	goodsType := p.Type
	if p.ProductType == catalog.ProductTypePoints {
		goodsType = 1
	}
	pv := &cart.ProductView{
		ProductID:     p.ProductID,
		MerID:         p.MerID,
		StoreName:     p.StoreName,
		Image:         p.Image,
		Price:         p.Price,
		OtPrice:       p.OtPrice,
		Cost:          0,
		Stock:         p.Stock,
		Status:        p.Status,
		IsShow:        p.IsShow,
		MerStatus:     p.MerStatus,
		ProductType:   p.ProductType,
		Integral:      p.Integral,
		GoodsType:     goodsType,
		SvipPriceType: p.SvipPriceType,
		SvipPrice:     p.SvipPrice,
		MerSvipStatus: p.MerSvipStatus,
	}
	_ = r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Select("mer_name").Where("mer_id = ?", p.MerID).Scan(&pv.MerName).Error

	var sku catalog.AttrValue
	q := r.db.WithContext(ctx).Where("product_id = ?", productID)
	if unique != "" {
		q = q.Where("`unique` = ?", unique)
	}
	if err := q.Order("value_id ASC").First(&sku).Error; err == nil {
		pv.Unique = sku.Unique
		pv.Price = sku.Price
		pv.OtPrice = sku.OtPrice
		pv.Stock = sku.Stock
		pv.Cost = sku.Cost
		if sku.SvipPrice > 0 {
			pv.SvipPrice = sku.SvipPrice
		}
		if sku.Image != "" {
			pv.Image = sku.Image
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return pv, nil
}

func (r *Repo) ListAddresses(ctx context.Context, uid uint) ([]cart.Address, error) {
	var rows []cart.Address
	err := r.db.WithContext(ctx).
		Where("uid = ? AND is_del = 0", uid).
		Order("is_default DESC, address_id DESC").
		Find(&rows).Error
	return rows, err
}

func (r *Repo) GetAddress(ctx context.Context, id, uid uint) (*cart.Address, error) {
	var row cart.Address
	err := r.db.WithContext(ctx).
		Where("address_id = ? AND uid = ? AND is_del = 0", id, uid).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateAddress(ctx context.Context, a *cart.Address) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) UpdateAddress(ctx context.Context, a *cart.Address) error {
	return r.db.WithContext(ctx).Model(&cart.Address{}).
		Where("address_id = ? AND uid = ? AND is_del = 0", a.AddressID, a.UID).
		Updates(map[string]interface{}{
			"real_name":  a.RealName,
			"phone":      a.Phone,
			"province":   a.Province,
			"city":       a.City,
			"district":   a.District,
			"detail":     a.Detail,
			"post_code":  a.PostCode,
			"is_default": a.IsDefault,
		}).Error
}

func (r *Repo) SoftDeleteAddress(ctx context.Context, id, uid uint) error {
	return r.db.WithContext(ctx).Model(&cart.Address{}).
		Where("address_id = ? AND uid = ?", id, uid).
		Update("is_del", 1).Error
}

func (r *Repo) ClearDefaultAddress(ctx context.Context, uid uint) error {
	return r.db.WithContext(ctx).Model(&cart.Address{}).
		Where("uid = ? AND is_del = 0", uid).
		Update("is_default", 0).Error
}
