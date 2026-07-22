package financepersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/finance"
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
	err := r.db.WithContext(ctx).Table("qixi_merchant").Select("mer_money").
		Where("mer_id = ? AND is_del = 0", merID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Take(&row).Error
	if err != nil {
		return 0, err
	}
	return row.MerMoney, nil
}

func (r *Repo) DeductMerMoney(ctx context.Context, merID uint, amount float64) error {
	res := r.db.WithContext(ctx).Table("qixi_merchant").
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
	res := r.db.WithContext(ctx).Table("qixi_merchant").
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
