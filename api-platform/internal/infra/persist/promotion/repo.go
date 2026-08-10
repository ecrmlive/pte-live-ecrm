package promotionpersist

import (
	"context"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/promotion"
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
		"title":          c.Title,
		"coupon_price":   c.CouponPrice,
		"use_min_price":  c.UseMinPrice,
		"coupon_type":    c.CouponType,
		"coupon_time":    c.CouponTime,
		"use_start_time": c.UseStartTime,
		"use_end_time":   c.UseEndTime,
		"is_timeout":     c.IsTimeout,
		"start_time":     c.StartTime,
		"end_time":       c.EndTime,
		"send_type":      c.SendType,
		"full_reduction": c.FullReduction,
		"status":         c.Status,
		"total_count":    c.TotalCount,
		"remain_count":   c.RemainCount,
		"is_limited":     c.IsLimited,
		"sort":           c.Sort,
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

func (r *Repo) ListCoupons(ctx context.Context, merID *uint, typ *int, page, limit int, filter promotion.CouponListFilter) ([]promotion.Coupon, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.Coupon{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if filter.StoreOnly {
		q = q.Where("mer_id > 0")
	}
	if len(filter.MerIDs) > 0 {
		q = q.Where("mer_id IN ?", filter.MerIDs)
	}
	if typ != nil {
		q = q.Where("type = ?", *typ)
	}
	if kw := strings.TrimSpace(filter.Keyword); kw != "" {
		q = q.Where("title LIKE ?", "%"+kw+"%")
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.SendType != nil {
		q = q.Where("send_type = ?", *filter.SendType)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.Coupon
	err := q.Order("sort DESC, coupon_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) UpsertCouponTemplateView(ctx context.Context, c *promotion.Coupon) error {
	if c == nil || c.CouponID == 0 {
		return nil
	}
	starts, ends := c.UseStartTime, c.UseEndTime
	if c.CouponType == int8(promotion.TemplateDays) {
		starts, ends = nil, nil
	}
	var existing struct {
		Version uint64 `gorm:"column:version"`
	}
	err := r.db.WithContext(ctx).Table("qixi_crm_b_coupon_template_view").
		Select("version").Where("coupon_id = ?", c.CouponID).Take(&existing).Error
	row := map[string]interface{}{
		"coupon_id":      c.CouponID,
		"store_id":       c.MerID,
		"name":           c.Title,
		"discount_type":  "amount",
		"discount_value": c.CouponPrice,
		"min_amount":     float64(c.UseMinPrice),
		"starts_at":      starts,
		"ends_at":        ends,
		"status":         c.Status,
		"version":        uint64(1),
	}
	if err == nil {
		row["version"] = existing.Version + 1
		return r.db.WithContext(ctx).Table("qixi_crm_b_coupon_template_view").
			Where("coupon_id = ?", c.CouponID).Updates(row).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.WithContext(ctx).Table("qixi_crm_b_coupon_template_view").Create(row).Error
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

// CountCouponUsage 统计业务库用户券领取/使用量（qixi_crm_b_coupon_user）。
func (r *Repo) CountCouponUsage(ctx context.Context, couponID uint) (received, used int64, err error) {
	if err = r.db.WithContext(ctx).Table("qixi_crm_b_coupon_user").
		Where("coupon_id = ?", couponID).Count(&received).Error; err != nil {
		return 0, 0, err
	}
	if err = r.db.WithContext(ctx).Table("qixi_crm_b_coupon_user").
		Where("coupon_id = ? AND status = ?", couponID, "used").Count(&used).Error; err != nil {
		return 0, 0, err
	}
	return received, used, nil
}

func (r *Repo) CountCouponUsageBatch(ctx context.Context, couponIDs []uint) (map[uint][2]int64, error) {
	out := make(map[uint][2]int64, len(couponIDs))
	if len(couponIDs) == 0 {
		return out, nil
	}
	type row struct {
		CouponID uint  `gorm:"column:coupon_id"`
		Received int64 `gorm:"column:received"`
		Used     int64 `gorm:"column:used"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Table("qixi_crm_b_coupon_user").
		Select("coupon_id, COUNT(1) AS received, SUM(CASE WHEN status = 'used' THEN 1 ELSE 0 END) AS used").
		Where("coupon_id IN ?", couponIDs).
		Group("coupon_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, item := range rows {
		out[item.CouponID] = [2]int64{item.Received, item.Used}
	}
	return out, nil
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

func (r *Repo) CreateCouponSend(ctx context.Context, row *promotion.CouponSend) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) MarkCouponSendDone(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&promotion.CouponSend{}).
		Where("coupon_send_id = ?", id).Update("status", 1).Error
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

func (r *Repo) HasMerchantCouponUser(ctx context.Context, merID, uid, couponID uint) (bool, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&promotion.CouponUser{}).
		Where("mer_id = ? AND uid = ? AND coupon_id = ?", merID, uid, couponID).
		Count(&total).Error
	return total > 0, err
}

func (r *Repo) ListMerchantCouponUsers(ctx context.Context, merID uint, couponID *uint, page, limit int) ([]promotion.CouponUser, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.CouponUser{}).Where("mer_id = ?", merID)
	if couponID != nil {
		q = q.Where("coupon_id = ?", *couponID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.CouponUser
	err := q.Order("coupon_user_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListMerchantCouponSends(ctx context.Context, merID uint, page, limit int) ([]promotion.CouponSend, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Model(&promotion.CouponSend{}).Where("mer_id = ?", merID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.CouponSend
	err := q.Order("coupon_send_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListPlatformCouponSends(ctx context.Context, page, limit int, filter promotion.CouponSendListFilter) ([]promotion.CouponSendListItem, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Table("qixi_crm_b_store_coupon_send AS s").
		Joins("INNER JOIN qixi_crm_b_store_coupon AS c ON c.coupon_id = s.coupon_id").
		Where("s.mer_id = 0")
	if filter.DateFrom != "" {
		q = q.Where("s.create_time >= ?", filter.DateFrom+" 00:00:00")
	}
	if filter.DateTo != "" {
		q = q.Where("s.create_time <= ?", filter.DateTo+" 23:59:59")
	}
	if filter.CouponType != nil {
		q = q.Where("c.type = ?", *filter.CouponType)
	}
	if name := strings.TrimSpace(filter.CouponName); name != "" {
		q = q.Where("c.title LIKE ?", "%"+name+"%")
	}
	if filter.SendStatus != nil {
		q = q.Where("s.status = ?", *filter.SendStatus)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	type row struct {
		CouponSendID uint       `gorm:"column:coupon_send_id"`
		CouponID     uint       `gorm:"column:coupon_id"`
		Title        string     `gorm:"column:title"`
		Type         int        `gorm:"column:type"`
		CreateTime   time.Time  `gorm:"column:create_time"`
		CouponType   int8       `gorm:"column:coupon_type"`
		CouponTime   uint       `gorm:"column:coupon_time"`
		UseStartTime *time.Time `gorm:"column:use_start_time"`
		UseEndTime   *time.Time `gorm:"column:use_end_time"`
		Mark         *string    `gorm:"column:mark"`
		CouponNum    uint       `gorm:"column:coupon_num"`
		SendStatus   int8       `gorm:"column:send_status"`
		UseCount     int64      `gorm:"column:use_count"`
	}
	var raw []row
	err := q.Select(`
s.coupon_send_id, s.coupon_id, c.title, c.type, s.create_time,
c.coupon_type, c.coupon_time, c.use_start_time, c.use_end_time,
CAST(s.mark AS CHAR) AS mark, s.coupon_num, s.status AS send_status,
(SELECT COUNT(1) FROM qixi_crm_b_coupon_user cu
  WHERE cu.send_id = s.coupon_send_id AND cu.status = 'used') AS use_count
`).
		Order("s.coupon_send_id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&raw).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]promotion.CouponSendListItem, 0, len(raw))
	for _, item := range raw {
		mark := ""
		if item.Mark != nil {
			mark = *item.Mark
		}
		out = append(out, promotion.CouponSendListItem{
			CouponSendID: item.CouponSendID,
			CouponID:     item.CouponID,
			Title:        item.Title,
			Type:         item.Type,
			CreateTime:   item.CreateTime,
			CouponType:   item.CouponType,
			CouponTime:   item.CouponTime,
			UseStartTime: item.UseStartTime,
			UseEndTime:   item.UseEndTime,
			Mark:         mark,
			CouponNum:    item.CouponNum,
			UseCount:     item.UseCount,
			SendStatus:   item.SendStatus,
		})
	}
	return out, total, nil
}

func (r *Repo) GetPlatformCouponSend(ctx context.Context, sendID uint) (*promotion.CouponSendDetail, error) {
	type row struct {
		CouponSendID uint       `gorm:"column:coupon_send_id"`
		CouponID     uint       `gorm:"column:coupon_id"`
		Title        string     `gorm:"column:title"`
		Type         int        `gorm:"column:type"`
		CreateTime   time.Time  `gorm:"column:create_time"`
		CouponType   int8       `gorm:"column:coupon_type"`
		CouponTime   uint       `gorm:"column:coupon_time"`
		UseStartTime *time.Time `gorm:"column:use_start_time"`
		UseEndTime   *time.Time `gorm:"column:use_end_time"`
		Mark         *string    `gorm:"column:mark"`
		CouponNum    uint       `gorm:"column:coupon_num"`
		SendStatus   int8       `gorm:"column:send_status"`
		UseCount     int64      `gorm:"column:use_count"`
		CouponPrice  float64    `gorm:"column:coupon_price"`
		UseMinPrice  int        `gorm:"column:use_min_price"`
		IsTimeout    int8       `gorm:"column:is_timeout"`
		StartTime    *time.Time `gorm:"column:start_time"`
		EndTime      *time.Time `gorm:"column:end_time"`
		SendType     int8       `gorm:"column:send_type"`
		IsLimited    int8       `gorm:"column:is_limited"`
		TotalCount   uint       `gorm:"column:total_count"`
		RemainCount  uint       `gorm:"column:remain_count"`
		Status       int8       `gorm:"column:status"`
		Sort         uint       `gorm:"column:sort"`
	}
	var item row
	err := r.db.WithContext(ctx).Table("qixi_crm_b_store_coupon_send AS s").
		Joins("INNER JOIN qixi_crm_b_store_coupon AS c ON c.coupon_id = s.coupon_id").
		Where("s.mer_id = 0 AND s.coupon_send_id = ?", sendID).
		Select(`
s.coupon_send_id, s.coupon_id, c.title, c.type, s.create_time,
c.coupon_type, c.coupon_time, c.use_start_time, c.use_end_time,
CAST(s.mark AS CHAR) AS mark, s.coupon_num, s.status AS send_status,
(SELECT COUNT(1) FROM qixi_crm_b_coupon_user cu
  WHERE cu.send_id = s.coupon_send_id AND cu.status = 'used') AS use_count,
c.coupon_price, c.use_min_price, c.is_timeout, c.start_time, c.end_time,
c.send_type, c.is_limited, c.total_count, c.remain_count, c.status, c.sort
`).
		Take(&item).Error
	if err != nil {
		return nil, err
	}
	mark := ""
	if item.Mark != nil {
		mark = *item.Mark
	}
	return &promotion.CouponSendDetail{
		CouponSendListItem: promotion.CouponSendListItem{
			CouponSendID: item.CouponSendID,
			CouponID:     item.CouponID,
			Title:        item.Title,
			Type:         item.Type,
			CreateTime:   item.CreateTime,
			CouponType:   item.CouponType,
			CouponTime:   item.CouponTime,
			UseStartTime: item.UseStartTime,
			UseEndTime:   item.UseEndTime,
			Mark:         mark,
			CouponNum:    item.CouponNum,
			UseCount:     item.UseCount,
			SendStatus:   item.SendStatus,
		},
		CouponPrice: item.CouponPrice,
		UseMinPrice: item.UseMinPrice,
		IsTimeout:   item.IsTimeout,
		StartTime:   item.StartTime,
		EndTime:     item.EndTime,
		SendType:    item.SendType,
		IsLimited:   item.IsLimited,
		TotalCount:  item.TotalCount,
		RemainCount: item.RemainCount,
		Status:      item.Status,
		Sort:        item.Sort,
	}, nil
}

func (r *Repo) ListPlatformCouponSendUsers(ctx context.Context, sendID uint, page, limit int) ([]promotion.CouponSendUserRow, int64, error) {
	page, limit = normalize(page, limit)
	q := r.db.WithContext(ctx).Table("qixi_crm_b_coupon_user AS cu").
		Where("cu.send_id = ?", sendID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []promotion.CouponSendUserRow
	err := q.Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = cu.user_id").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id = cu.user_id").
		Select(`
cu.user_id AS user_id,
COALESCE(NULLIF(u.nickname,''), CONCAT('用户', cu.user_id)) AS nickname,
COALESCE(p.avatar_url,'') AS avatar_url,
cu.source AS source,
cu.status AS status
`).
		Order("cu.id DESC").
		Offset((page - 1) * limit).Limit(limit).
		Scan(&rows).Error
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

func (r *Repo) HasMerchantPaidOrderUser(ctx context.Context, merID, uid uint) (bool, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_order").
		Where("mer_id = ? AND uid = ? AND paid = 1 AND is_del = 0", merID, uid).
		Count(&total).Error
	return total > 0, err
}

func (r *Repo) ListUsablePlatform(ctx context.Context, uid uint, orderAmount float64) ([]promotion.CouponUser, error) {
	now := time.Now()
	var rows []promotion.CouponUser
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_coupon_user AS cu").
		Select("cu.*").
		Joins("INNER JOIN qixi_m_admin_store_coupon AS c ON c.coupon_id = cu.coupon_id").
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
	err = r.db.WithContext(ctx).Table("qixi_m_app_user").
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
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").
		Select("is_promoter, status").
		Where("uid = ?", uid).
		Take(&row).Error
	if err != nil {
		return false, err
	}
	return row.IsPromoter == 1 && row.Status == 1, nil
}

func (r *Repo) SetUserSpread(ctx context.Context, uid, spreadUID uint) error {
	res := r.db.WithContext(ctx).Table("qixi_m_app_user").
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
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").Where("spread_uid = ?", spreadUID).Count(&n).Error
	return n, err
}

func (r *Repo) AddBrokerage(ctx context.Context, uid uint, amount float64) (float64, error) {
	res := r.db.WithContext(ctx).Table("qixi_m_app_user").
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
	err := r.db.WithContext(ctx).Table("qixi_m_app_user").Select("now_money").Where("uid = ?", uid).Scan(&bal).Error
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
