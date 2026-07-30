package financepersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/finance"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) WithTx(fn func(tx finance.Store) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{db: tx})
	})
}

func (r *Repo) GetMerMoney(ctx context.Context, merID uint) (float64, error) {
	var row struct {
		MerMoney float64 `gorm:"column:mer_money"`
	}
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").Select("mer_money").
		Where("mer_id = ? AND is_del = 0", merID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Take(&row).Error
	if err != nil {
		return 0, err
	}
	return row.MerMoney, nil
}

func (r *Repo) DeductMerMoney(ctx context.Context, merID uint, amount float64) error {
	res := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Where("mer_id = ? AND is_del = 0 AND mer_money >= ?", merID, amount).
		Update("mer_money", gorm.Expr("mer_money - ?", amount))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return finance.ErrBalanceNotEnough
	}
	return nil
}

func (r *Repo) AddMerMoney(ctx context.Context, merID uint, amount float64) error {
	res := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Where("mer_id = ? AND is_del = 0", merID).
		Update("mer_money", gorm.Expr("mer_money + ?", amount))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return finance.ErrMerchantNotFound
	}
	return nil
}

func (r *Repo) CreateFinancial(ctx context.Context, f *finance.Financial) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *Repo) CreateRecord(ctx context.Context, rec *finance.FinancialRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *Repo) GetFinancial(ctx context.Context, id uint) (*finance.Financial, error) {
	var f finance.Financial
	err := r.db.WithContext(ctx).Where("financial_id = ?", id).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repo) ListFinancials(ctx context.Context, filter finance.ListFilter) ([]finance.Financial, int64, error) {
	page, limit := normalize(filter.Page, filter.Limit)
	q := r.db.WithContext(ctx).Model(&finance.Financial{}).Where("is_del = 0 AND type = ?", finance.TypeExtract)
	if filter.MerID != nil {
		q = q.Where("mer_id = ?", *filter.MerID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []finance.Financial
	err := q.Order("financial_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// ListPlatformFinancialsByRegions 在 SQL 层以商户所属区域过滤提现单，区域账号不得读取范围外资金数据。
func (r *Repo) ListPlatformFinancialsByRegions(ctx context.Context, regionIDs []uint, status *int, page, limit int) ([]finance.Financial, int64, error) {
	page, limit = normalize(page, limit)
	if len(regionIDs) == 0 {
		return []finance.Financial{}, 0, nil
	}
	q := r.db.WithContext(ctx).Table("qixi_m_admin_financial AS f").
		Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = f.mer_id").
		Where("f.is_del = 0 AND f.type = ? AND m.is_del = 0 AND m.region_id IN ?", finance.TypeExtract, regionIDs)
	if status != nil {
		q = q.Where("f.status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []finance.Financial
	err := q.Select("f.*").Order("f.financial_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

// GetPlatformFinancialByRegions 在 SQL 层按提现单所属商户区域过滤详情。
func (r *Repo) GetPlatformFinancialByRegions(ctx context.Context, id uint, regionIDs []uint) (*finance.Financial, error) {
	if len(regionIDs) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var row finance.Financial
	err := r.db.WithContext(ctx).Table("qixi_m_admin_financial AS f").
		Select("f.*").Joins("JOIN qixi_m_admin_merchant AS m ON m.mer_id = f.mer_id").
		Where("f.financial_id = ? AND f.is_del = 0 AND f.type = ? AND m.is_del = 0 AND m.region_id IN ?", id, finance.TypeExtract, regionIDs).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateFinancialStatus(ctx context.Context, id uint, fromStatus, toStatus int, refusal string, adminID uint) (bool, error) {
	now := time.Now()
	admin := int(adminID)
	updates := map[string]interface{}{
		"status":      toStatus,
		"status_time": now,
		"update_time": now,
		"admin_id":    admin,
	}
	if refusal != "" {
		updates["refusal"] = refusal
	}
	res := r.db.WithContext(ctx).Model(&finance.Financial{}).
		Where("financial_id = ? AND status = ? AND is_del = 0", id, fromStatus).
		Updates(updates)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}

var _ finance.Store = (*Repo)(nil)
