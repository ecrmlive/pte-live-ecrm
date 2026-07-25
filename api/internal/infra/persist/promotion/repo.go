package promotionpersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/promotion"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) WithTx(fn func(tx promotion.Store) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo{db: tx})
	})
}

func (r *Repo) CreateCoupon(ctx context.Context, c *promotion.Coupon) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateCoupon(ctx context.Context, c *promotion.Coupon) error {
	return r.db.WithContext(ctx).Model(c).Where("coupon_id = ?", c.CouponID).Updates(map[string]interface{}{
		"title":         c.Title,
		"coupon_price":  c.CouponPrice,
		"use_min_price": c.UseMinPrice,
		"coupon_time":   c.CouponTime,
		"status":        c.Status,
		"total_count":   c.TotalCount,
		"remain_count":  c.RemainCount,
		"is_limited":    c.IsLimited,
		"sort":          c.Sort,
	}).Error
}

func (r *Repo) UpdateCouponStatus(ctx context.Context, id uint, merID *uint, status int8) (bool, error) {
	q := r.db.WithContext(ctx).Model(&promotion.Coupon{}).Where("coupon_id = ? AND is_del = 0", id)
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	res := q.Update("status", status)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repo) SoftDeleteCoupon(ctx context.Context, id uint, merID *uint) (bool, error) {
	q := r.db.WithContext(ctx).Model(&promotion.Coupon{}).Where("coupon_id = ? AND is_del = 0", id)
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	res := q.Update("is_del", 1)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repo) GetCoupon(ctx context.Context, id uint) (*promotion.Coupon, error) {
	var row promotion.Coupon
	err := r.db.WithContext(ctx).Where("coupon_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListCoupons(ctx context.Context, merID *uint, typ *int, page, limit int) ([]promotion.Coupon, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.Coupon{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if typ != nil {
		q = q.Where("type = ?", *typ)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.Coupon
	err := q.Order("sort DESC, coupon_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListCenter(ctx context.Context, page, limit int) ([]promotion.Coupon, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.Coupon{}).
		Where("is_del = 0 AND status = 1 AND send_type = 0").
		Where("(is_limited = 0 OR remain_count > 0)")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.Coupon
	err := q.Order("type DESC, sort DESC, coupon_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) DecRemain(ctx context.Context, couponID uint) (bool, error) {
	res := r.db.WithContext(ctx).Model(&promotion.Coupon{}).
		Where("coupon_id = ? AND is_del = 0 AND is_limited = 1 AND remain_count > 0", couponID).
		Update("remain_count", gorm.Expr("remain_count - 1"))
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *Repo) HasReceived(ctx context.Context, uid, couponID uint) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&promotion.IssueUser{}).
		Where("uid = ? AND coupon_id = ?", uid, couponID).Count(&n).Error
	return n > 0, err
}

func (r *Repo) CreateIssueUser(ctx context.Context, row *promotion.IssueUser) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) CreateCouponUser(ctx context.Context, u *promotion.CouponUser) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repo) GetCouponUser(ctx context.Context, id uint) (*promotion.CouponUser, error) {
	var row promotion.CouponUser
	err := r.db.WithContext(ctx).Where("coupon_user_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListCouponUsers(ctx context.Context, uid uint, status *int, page, limit int) ([]promotion.CouponUser, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.CouponUser{}).Where("uid = ?", uid)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.CouponUser
	err := q.Order("coupon_user_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListCouponUsersByIDs(ctx context.Context, uid uint, ids []uint) ([]promotion.CouponUser, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []promotion.CouponUser
	err := r.db.WithContext(ctx).
		Where("uid = ? AND coupon_user_id IN ?", uid, ids).
		Find(&rows).Error
	return rows, err
}

func (r *Repo) ListUsablePlatform(ctx context.Context, uid uint, orderAmount float64) ([]promotion.CouponUser, error) {
	now := time.Now()
	var rows []promotion.CouponUser
	err := r.db.WithContext(ctx).Table("qixi_store_coupon_user AS cu").
		Select("cu.*").
		Joins("INNER JOIN qixi_store_coupon AS c ON c.coupon_id = cu.coupon_id").
		Where("cu.uid = ? AND cu.status = ? AND cu.is_fail = 0", uid, promotion.UserUnused).
		Where("c.type = ? AND c.is_del = 0", promotion.CouponTypePlatform).
		Where("cu.use_min_price <= ?", int(orderAmount)).
		Where("(cu.start_time IS NULL OR cu.start_time <= ?)", now).
		Where("(cu.end_time IS NULL OR cu.end_time >= ?)", now).
		Order("cu.coupon_price DESC, cu.coupon_user_id ASC").
		Find(&rows).Error
	for i := range rows {
		rows[i].CouponKind = promotion.CouponTypePlatform
	}
	return rows, err
}

func (r *Repo) MarkCouponUsersUsed(ctx context.Context, uid uint, ids []uint, at time.Time) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	res := r.db.WithContext(ctx).Model(&promotion.CouponUser{}).
		Where("uid = ? AND coupon_user_id IN ? AND status = ?", uid, ids, promotion.UserUnused).
		Updates(map[string]interface{}{
			"status":   promotion.UserUsed,
			"use_time": at,
		})
	return res.RowsAffected, res.Error
}

func (r *Repo) GetUserSpread(ctx context.Context, uid uint) (spreadUID uint, isPromoter int8, err error) {
	var row struct {
		SpreadUID  uint `gorm:"column:spread_uid"`
		IsPromoter int8 `gorm:"column:is_promoter"`
	}
	err = r.db.WithContext(ctx).Table("qixi_user").
		Select("spread_uid, is_promoter").
		Where("uid = ?", uid).
		Take(&row).Error
	if err != nil {
		return 0, 0, err
	}
	return row.SpreadUID, row.IsPromoter, nil
}

func (r *Repo) IsPromoter(ctx context.Context, uid uint) (bool, error) {
	var row struct {
		IsPromoter uint8 `gorm:"column:is_promoter"`
		Status     int8  `gorm:"column:status"`
	}
	err := r.db.WithContext(ctx).Table("qixi_user").
		Select("is_promoter, status").
		Where("uid = ?", uid).
		Take(&row).Error
	if err != nil {
		return false, err
	}
	return row.IsPromoter == 1 && row.Status == 1, nil
}

func (r *Repo) SetUserSpread(ctx context.Context, uid, spreadUID uint) error {
	res := r.db.WithContext(ctx).Table("qixi_user").
		Where("uid = ? AND spread_uid = 0", uid).
		Update("spread_uid", spreadUID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return promotion.ErrSpreadBound
	}
	return nil
}

func (r *Repo) CreateSpreadLog(ctx context.Context, log *promotion.SpreadLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *Repo) ListSpreadLogs(ctx context.Context, page, limit int) ([]promotion.SpreadLog, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.SpreadLog{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.SpreadLog
	err := q.Order("user_spread_log_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CountSpreadChildren(ctx context.Context, spreadUID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("qixi_user").Where("spread_uid = ?", spreadUID).Count(&n).Error
	return n, err
}

func (r *Repo) AddBrokerage(ctx context.Context, uid uint, amount float64) (float64, error) {
	res := r.db.WithContext(ctx).Table("qixi_user").
		Where("uid = ?", uid).
		Updates(map[string]interface{}{
			"brokerage_price": gorm.Expr("brokerage_price + ?", amount),
			"now_money":       gorm.Expr("now_money + ?", amount),
		})
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	var bal float64
	err := r.db.WithContext(ctx).Table("qixi_user").Select("now_money").Where("uid = ?", uid).Scan(&bal).Error
	return bal, err
}

func (r *Repo) CreateBill(ctx context.Context, b *promotion.UserBill) error {
	if b.CreateTime.IsZero() {
		b.CreateTime = time.Now()
	}
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *Repo) HasBill(ctx context.Context, uid uint, category, typ, linkID string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&promotion.UserBill{}).
		Where("uid = ? AND category = ? AND type = ? AND link_id = ?", uid, category, typ, linkID).
		Count(&n).Error
	return n > 0, err
}

func (r *Repo) ListBills(ctx context.Context, uid *uint, category string, page, limit int) ([]promotion.UserBill, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.UserBill{})
	if uid != nil {
		q = q.Where("uid = ?", *uid)
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.UserBill
	err := q.Order("bill_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
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

var _ promotion.Store = (*Repo)(nil)
