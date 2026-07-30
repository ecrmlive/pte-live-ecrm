package merchantpersist

import (
	"context"
	"fmt"

	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/merchant"
	"gorm.io/gorm"
)

type Repo struct {
	adminDB *gorm.DB
}

// NewRepo 只读取后台库中的商户监管投影。跨库业务事实通过 api-merchant/NATS
// 同步到该投影，统一后台不直连店铺数据库。
func NewRepo(adminDB *gorm.DB) *Repo { return &Repo{adminDB: adminDB} }

type ListMerchantsFilter struct {
	Keyword   string
	Status    *int8
	RegionIDs []uint
	Page      int
	Limit     int
}

func (r *Repo) ListMerchants(ctx context.Context, f ListMerchantsFilter) ([]merchant.Merchant, int64, error) {
	q := r.adminDB.WithContext(ctx).
		Table("qixi_crm_a_merchant_view").
		Select(`merchant_id AS mer_id, merchant_name AS mer_name, contact_name AS real_name, contact_mobile AS mer_phone,
            '' AS mer_address, '' AS mer_info, '' AS mark, status,
            status AS mer_state, 1 AS is_audit, 0 AS svip_coupon_merge,
            COALESCE(region_id, 0) AS region_id, created_at AS create_time`)
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("merchant_name LIKE ? OR contact_name LIKE ? OR contact_mobile LIKE ?", like, like, like)
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
	err := q.Order("merchant_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) GetMerchant(ctx context.Context, id uint) (*merchant.Merchant, error) {
	var row merchant.Merchant
	err := r.adminDB.WithContext(ctx).
		Table("qixi_crm_a_merchant_view").
		Select(`merchant_id AS mer_id, merchant_name AS mer_name, contact_name AS real_name, contact_mobile AS mer_phone,
            '' AS mer_address, '' AS mer_info, '' AS mark, status,
            status AS mer_state, 1 AS is_audit, 0 AS svip_coupon_merge,
            COALESCE(region_id, 0) AS region_id, created_at AS create_time`).
		Where("merchant_id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Where("merchant_id = ?", id).Update("status", status).Error
}

func (r *Repo) UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error {
	return fmt.Errorf("商户会员配置应由 api-merchant 处理，merchant_id=%d", merID)
}

func (r *Repo) UpdateShopProfile(ctx context.Context, merID uint, merName, realName, merPhone, merAddress, merInfo string) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Where("merchant_id = ?", merID).Update("merchant_name", merName).Error
}

func (r *Repo) CreateMerchant(ctx context.Context, m *merchant.Merchant) error {
	return fmt.Errorf("商户创建必须由 api-merchant 的受控命令处理")
}

func (r *Repo) CreateMerchantAdmin(ctx context.Context, a *identity.MerchantAdmin) error {
	return fmt.Errorf("店铺账号创建必须由 api-merchant 的受控命令处理")
}

type ListIntentionFilter struct {
	Status  *int8
	Keyword string
	Page    int
	Limit   int
}

func (r *Repo) ListIntentions(ctx context.Context, f ListIntentionFilter) ([]merchant.Intention, int64, error) {
	q := r.adminDB.WithContext(ctx).
		Table("qixi_crm_a_merchant_application").
		Select(`id AS mer_intention_id, COALESCE(applicant_user_id, 0) AS uid,
            contact_mobile AS phone, merchant_name AS mer_name, contact_name AS name,
            created_at AS create_time,
            CASE status WHEN 'pending' THEN 0 WHEN 'approved' THEN 1 WHEN 'rejected' THEN 2 ELSE 0 END AS status,
            review_note AS fail_msg, review_note AS mark,
            COALESCE(region_id, 0) AS circle_id,
            0 AS mer_id, 0 AS merchant_category_id, 0 AS mer_type_id, '' AS images`)
	if f.Status != nil {
		status := map[int8]string{
			merchant.IntentionPending:  "pending",
			merchant.IntentionApproved: "approved",
			merchant.IntentionRejected: "rejected",
		}[*f.Status]
		if status == "" {
			return []merchant.Intention{}, 0, nil
		}
		q = q.Where("status = ?", status)
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("merchant_name LIKE ? OR contact_mobile LIKE ? OR contact_name LIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := normalizePage(f.Page, f.Limit)
	var rows []merchant.Intention
	err := q.Order("FIELD(status, 'pending', 'approved', 'rejected'), id DESC").
		Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) GetIntention(ctx context.Context, id uint) (*merchant.Intention, error) {
	var row merchant.Intention
	err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_application").
		Select(`id AS mer_intention_id, COALESCE(applicant_user_id, 0) AS uid,
            contact_mobile AS phone, merchant_name AS mer_name, contact_name AS name,
            created_at AS create_time,
            CASE status WHEN 'pending' THEN 0 WHEN 'approved' THEN 1 WHEN 'rejected' THEN 2 ELSE 0 END AS status,
            review_note AS fail_msg, review_note AS mark,
            COALESCE(region_id, 0) AS circle_id,
            0 AS mer_id, 0 AS merchant_category_id, 0 AS mer_type_id, '' AS images`).
		Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) SaveIntention(ctx context.Context, row *merchant.Intention) error {
	status := map[int8]string{
		merchant.IntentionPending:  "pending",
		merchant.IntentionApproved: "approved",
		merchant.IntentionRejected: "rejected",
	}[row.Status]
	if status == "" {
		return fmt.Errorf("unsupported merchant application status %d", row.Status)
	}
	note := row.Mark
	if row.Status == merchant.IntentionRejected {
		note = row.FailMsg
	}
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_application").
		Where("id = ?", row.MerIntentionID).
		Updates(map[string]any{"status": status, "review_note": note, "reviewed_at": gorm.Expr("NOW()")}).Error
}

func (r *Repo) ListCategories(ctx context.Context) ([]merchant.Category, error) {
	var rows []merchant.Category
	err := r.adminDB.WithContext(ctx).Order("id ASC").
		Table("qixi_crm_a_platform_category").
		Select("id AS merchant_category_id, name AS category_name, 0 AS commission_rate").Scan(&rows).Error
	return rows, err
}

func (r *Repo) CreateCategory(ctx context.Context, c *merchant.Category) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_platform_category").Create(map[string]any{
		"name": c.CategoryName,
	}).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, c *merchant.Category) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_platform_category").
		Where("id = ?", c.MerchantCategoryID).Update("name", c.CategoryName).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_platform_category").Where("id = ?", id).Delete(nil).Error
}

func (r *Repo) WithTx(fn func(tx *Repo) error) error {
	return r.adminDB.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{adminDB: tx})
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
