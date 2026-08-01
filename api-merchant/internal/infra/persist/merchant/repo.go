package merchantpersist

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/merchant"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

type ListMerchantsFilter struct {
	Keyword   string
	Status    *int8
	RegionIDs []uint
	Page      int
	Limit     int
}

func (r *Repo) ListMerchants(ctx context.Context, f ListMerchantsFilter) ([]merchant.Merchant, int64, error) {
	q := r.db.WithContext(ctx).Model(&merchant.Merchant{}).Where("is_del = 0")
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("mer_name LIKE ? OR mer_phone LIKE ? OR real_name LIKE ?", like, like, like)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.RegionIDs != nil {
		if len(f.RegionIDs) == 0 {
			return []merchant.Merchant{}, 0, nil
		}
		q = q.Where("region_id IN ?", f.RegionIDs)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := normalizePage(f.Page, f.Limit)
	var rows []merchant.Merchant
	err := q.Order("mer_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetMerchant(ctx context.Context, id uint) (*merchant.Merchant, error) {
	var row merchant.Merchant
	err := r.db.WithContext(ctx).Where("mer_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error {
	return r.db.WithContext(ctx).Model(&merchant.Merchant{}).
		Where("mer_id = ? AND is_del = 0", id).
		Updates(map[string]interface{}{"status": status, "mer_state": merState}).Error
}

func (r *Repo) UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error {
	return r.db.WithContext(ctx).Model(&merchant.Merchant{}).
		Where("mer_id = ? AND is_del = 0", merID).
		Update("svip_coupon_merge", merge).Error
}

func (r *Repo) UpdateShopProfile(ctx context.Context, merID uint, merName, realName, merPhone, merAddress, merInfo string) error {
	return r.db.WithContext(ctx).Model(&merchant.Merchant{}).
		Where("mer_id = ? AND is_del = 0", merID).
		Updates(map[string]interface{}{
			"mer_name": merName, "real_name": realName, "mer_phone": merPhone,
			"mer_address": merAddress, "mer_info": merInfo,
		}).Error
}

func (r *Repo) CreateMerchant(ctx context.Context, m *merchant.Merchant) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repo) CreateMerchantAdmin(ctx context.Context, a *identity.MerchantAdmin) error {
	return r.db.WithContext(ctx).Create(a).Error
}

type ListIntentionFilter struct {
	Status  *int8
	Keyword string
	Page    int
	Limit   int
}

func (r *Repo) ListIntentions(ctx context.Context, f ListIntentionFilter) ([]merchant.Intention, int64, error) {
	q := r.db.WithContext(ctx).Model(&merchant.Intention{}).Where("is_del = 0")
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("mer_name LIKE ? OR phone LIKE ? OR name LIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := normalizePage(f.Page, f.Limit)
	var rows []merchant.Intention
	err := q.Order("status ASC, mer_intention_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetIntention(ctx context.Context, id uint) (*merchant.Intention, error) {
	var row merchant.Intention
	err := r.db.WithContext(ctx).Where("mer_intention_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) SaveIntention(ctx context.Context, row *merchant.Intention) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) ListCategories(ctx context.Context) ([]merchant.Category, error) {
	var rows []merchant.Category
	err := r.db.WithContext(ctx).Order("merchant_category_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) CreateCategory(ctx context.Context, c *merchant.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, c *merchant.Category) error {
	return r.db.WithContext(ctx).Model(&merchant.Category{}).
		Where("merchant_category_id = ?", c.MerchantCategoryID).
		Updates(map[string]interface{}{
			"category_name":   c.CategoryName,
			"commission_rate": c.CommissionRate,
		}).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&merchant.Category{}, id).Error
}

func (r *Repo) WithTx(fn func(tx *Repo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{db: tx})
	})
}

func normalizePage(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
