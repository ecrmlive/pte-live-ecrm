package distributionpersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/distribution"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/txctx"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) dbx(ctx context.Context) *gorm.DB {
	return txctx.DB(ctx, r.db).WithContext(ctx)
}

func (r *Repo) WithTx(ctx context.Context, fn func(ctx context.Context, tx distribution.Store) error) error {
	return r.dbx(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := txctx.With(ctx, tx)
		return fn(txCtx, &Repo{db: tx})
	})
}

func (r *Repo) GetUser(ctx context.Context, uid uint) (*distribution.UserSpread, error) {
	var u distribution.UserSpread
	err := r.dbx(ctx).Table("qixi_user").
		Select("uid, spread_uid, is_promoter, status").
		Where("uid = ?", uid).
		Take(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) SetSpreadUID(ctx context.Context, uid, spreadUID uint) (bool, error) {
	res := r.dbx(ctx).Table("qixi_user").
		Where("uid = ? AND spread_uid = 0", uid).
		Update("spread_uid", spreadUID)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repo) CreateLog(ctx context.Context, log *distribution.SpreadLog) error {
	if log.CreateTime.IsZero() {
		log.CreateTime = time.Now()
	}
	return r.dbx(ctx).Create(log).Error
}

func (r *Repo) CountChildren(ctx context.Context, spreadUID uint) (int64, error) {
	var n int64
	err := r.dbx(ctx).Table("qixi_user").Where("spread_uid = ?", spreadUID).Count(&n).Error
	return n, err
}

func (r *Repo) ListLogs(ctx context.Context, page, limit int) ([]distribution.SpreadLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.dbx(ctx).Model(&distribution.SpreadLog{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []distribution.SpreadLog
	err := q.Order("user_spread_log_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}
