package seckillpersist

import (
	"context"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/seckill"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListTimes(ctx context.Context) ([]seckill.TimeSlot, error) {
	var rows []seckill.TimeSlot
	err := r.db.WithContext(ctx).Where("status = 1").Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ListTimesAdmin(ctx context.Context, status *int8, page, limit int) ([]seckill.TimeSlot, int64, error) {
	q := r.db.WithContext(ctx).Model(&seckill.TimeSlot{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.TimeSlot
	err := q.Order("start_time ASC, seckill_time_id ASC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetTime(ctx context.Context, id uint) (*seckill.TimeSlot, error) {
	var row seckill.TimeSlot
	err := r.db.WithContext(ctx).Where("seckill_time_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateTime(ctx context.Context, t *seckill.TimeSlot) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *Repo) UpdateTime(ctx context.Context, t *seckill.TimeSlot) error {
	return r.db.WithContext(ctx).Model(&seckill.TimeSlot{}).
		Where("seckill_time_id = ?", t.SeckillTimeID).
		Updates(map[string]interface{}{
			"title":      t.Title,
			"start_time": t.StartTime,
			"end_time":   t.EndTime,
			"status":     t.Status,
			"pic":        t.Pic,
		}).Error
}

func (r *Repo) DeleteTime(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("seckill_time_id = ?", id).Delete(&seckill.TimeSlot{}).Error
}

// HasTimeOverlap 半开区间 [start,end) 与已有场次重叠则返回 true（对齐 CRMEB checkTime）。
func (r *Repo) HasTimeOverlap(ctx context.Context, start, end int, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).Model(&seckill.TimeSlot{}).
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeID > 0 {
		q = q.Where("seckill_time_id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) ListActives(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]seckill.Active, int64, error) {
	q := r.db.WithContext(ctx).Model(&seckill.Active{}).Where("delete_time IS NULL")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyOn {
		q = q.Where("status = 1 AND active_status = 1 AND product_status = 1 AND is_show = 1")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.Active
	err := q.Order("seckill_active_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) applyActiveAdminFilters(q *gorm.DB, query seckill.ActiveAdminQuery, tabType int) *gorm.DB {
	if query.MerID != nil {
		q = q.Where("mer_id = ?", *query.MerID)
	}
	if len(query.MerIDs) > 0 {
		q = q.Where("mer_id IN ?", query.MerIDs)
	}
	if name := strings.TrimSpace(query.ActiveName); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if star := query.Star; star != nil {
		q = q.Where("star = ?", *star)
	}
	if us := query.UsStatus; us != nil {
		switch *us {
		case 1:
			q = q.Where("product_status = 1 AND is_show = 1 AND delete_time IS NULL")
		case 0:
			q = q.Where("product_status = 1 AND is_show = 0 AND delete_time IS NULL")
		case -1:
			q = q.Where("product_status = -1 AND delete_time IS NULL")
		case -2:
			q = q.Where("product_status = -2")
		}
	}
	if labels := strings.TrimSpace(query.SysLabels); labels != "" {
		q = q.Where("FIND_IN_SET(?, REPLACE(sys_labels,' ',''))", labels)
	}
	switch tabType {
	case 1:
		q = q.Where("delete_time IS NULL AND product_status = 1 AND is_show = 1")
	case 2:
		q = q.Where("delete_time IS NULL AND product_status = 1 AND is_show = 0")
	case 5:
		q = q.Where("delete_time IS NOT NULL")
	case 6:
		q = q.Where("delete_time IS NULL AND product_status = 0")
	case 7:
		q = q.Where("delete_time IS NULL AND product_status = -1")
	default:
		// 不过滤 Tab
	}
	return q
}

func (r *Repo) ListActivesAdmin(ctx context.Context, query seckill.ActiveAdminQuery) ([]seckill.Active, int64, error) {
	q := r.db.WithContext(ctx).Model(&seckill.Active{})
	q = r.applyActiveAdminFilters(q, query, query.Type)
	if kw := strings.TrimSpace(query.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where(`product_id IN (
			SELECT product_id FROM qixi_crm_b_product_view
			WHERE title LIKE ? OR CAST(product_id AS CHAR) LIKE ?
		)`, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.Active
	err := q.Order("sort DESC, seckill_active_id DESC").
		Offset((query.Page - 1) * query.Limit).Limit(query.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CountActivesAdmin(ctx context.Context, query seckill.ActiveAdminQuery, tabType int) (int64, error) {
	q := r.db.WithContext(ctx).Model(&seckill.Active{})
	q = r.applyActiveAdminFilters(q, query, tabType)
	if kw := strings.TrimSpace(query.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where(`product_id IN (
			SELECT product_id FROM qixi_crm_b_product_view
			WHERE title LIKE ? OR CAST(product_id AS CHAR) LIKE ?
		)`, like, like)
	}
	var total int64
	err := q.Count(&total).Error
	return total, err
}

func (r *Repo) GetActive(ctx context.Context, id uint) (*seckill.Active, error) {
	var row seckill.Active
	err := r.db.WithContext(ctx).Where("seckill_active_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetActiveByProduct(ctx context.Context, productID uint) (*seckill.Active, error) {
	var row seckill.Active
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = 1 AND active_status = 1 AND product_status = 1 AND is_show = 1 AND delete_time IS NULL", productID).
		Order("seckill_active_id DESC").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateActive(ctx context.Context, a *seckill.Active) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) UpdateActive(ctx context.Context, a *seckill.Active) error {
	return r.db.WithContext(ctx).Model(a).Where("seckill_active_id = ?", a.SeckillActiveID).Updates(map[string]interface{}{
		"name": a.Name, "seckill_time_ids": a.SeckillTimeIDs, "start_day": a.StartDay, "end_day": a.EndDay,
		"seckill_price": a.SeckillPrice, "once_pay_count": a.OncePayCount, "all_pay_count": a.AllPayCount,
		"status": a.Status,
		"is_show": a.IsShow, "product_status": a.ProductStatus, "star": a.Star, "sort": a.Sort,
		"stock": a.Stock, "sales": a.Sales, "sys_labels": a.SysLabels, "refusal": a.Refusal,
		"update_time": a.UpdateTime,
	}).Error
}

func (r *Repo) SoftDeleteActive(ctx context.Context, id uint) error {
	now := time.Now().Unix()
	return r.db.WithContext(ctx).Model(&seckill.Active{}).Where("seckill_active_id = ?", id).
		Updates(map[string]interface{}{"delete_time": now, "status": 0}).Error
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		MerID     uint    `gorm:"column:mer_id"`
		MerName   string  `gorm:"column:mer_name"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view AS p").
		Select("p.title AS store_name, p.cover_url AS image, p.price, p.merchant_id AS mer_id, p.merchant_name AS mer_name").
		Where("p.product_id = ? AND p.sale_status = 1", productID).
		Scan(&row).Error
	if err != nil {
		return "", "", "", 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", "", 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.MerName, row.Price, row.MerID, nil
}

func (r *Repo) ListActivities(ctx context.Context, q seckill.ActivityQuery) ([]seckill.Activity, int64, error) {
	db := r.db.WithContext(ctx).Model(&seckill.Activity{}).Where("delete_time IS NULL")
	if name := strings.TrimSpace(q.Name); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.ActiveStatus != nil {
		today := time.Now().Format("2006-01-02")
		switch *q.ActiveStatus {
		case 0:
			db = db.Where("start_day > ?", today)
		case 1:
			db = db.Where("start_day <= ? AND end_day >= ?", today, today)
		case -1:
			db = db.Where("end_day < ?", today)
		}
	}
	if from := strings.TrimSpace(q.DateFrom); from != "" {
		db = db.Where("end_day >= ?", from)
	}
	if to := strings.TrimSpace(q.DateTo); to != "" {
		db = db.Where("start_day <= ?", to)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.Activity
	err := db.Order("seckill_activity_id DESC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetActivity(ctx context.Context, id uint) (*seckill.Activity, error) {
	var row seckill.Activity
	err := r.db.WithContext(ctx).Where("seckill_activity_id = ? AND delete_time IS NULL", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateActivity(ctx context.Context, a *seckill.Activity) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) UpdateActivity(ctx context.Context, a *seckill.Activity) error {
	return r.db.WithContext(ctx).Model(&seckill.Activity{}).
		Where("seckill_activity_id = ?", a.SeckillActivityID).
		Updates(map[string]interface{}{
			"name":                 a.Name,
			"seckill_time_ids":     a.SeckillTimeIDs,
			"start_day":            a.StartDay,
			"end_day":              a.EndDay,
			"once_pay_count":       a.OncePayCount,
			"all_pay_count":        a.AllPayCount,
			"product_category_ids": a.ProductCategoryIDs,
			"border_pic":           a.BorderPic,
			"status":               a.Status,
			"active_status":        a.ActiveStatus,
		}).Error
}

func (r *Repo) SoftDeleteActivity(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&seckill.Activity{}).
		Where("seckill_activity_id = ?", id).
		Updates(map[string]interface{}{"delete_time": now, "status": 0}).Error
}

func (r *Repo) RefreshActivityCounts(ctx context.Context, id uint) error {
	type agg struct {
		ProductCount  int64
		MerchantCount int64
	}
	var a agg
	_ = r.db.WithContext(ctx).Model(&seckill.Active{}).
		Select("COUNT(*) AS product_count, COUNT(DISTINCT mer_id) AS merchant_count").
		Where("activity_id = ? AND delete_time IS NULL", id).
		Scan(&a)
	var row seckill.Activity
	if err := r.db.WithContext(ctx).Where("seckill_activity_id = ?", id).First(&row).Error; err != nil {
		return err
	}
	today := time.Now().Format("2006-01-02")
	startDay := dayOnly(row.StartDay)
	endDay := dayOnly(row.EndDay)
	var activeStatus int8 = 1
	if today < startDay {
		activeStatus = 0
	} else if today > endDay {
		activeStatus = -1
	}
	return r.db.WithContext(ctx).Model(&seckill.Activity{}).
		Where("seckill_activity_id = ?", id).
		Updates(map[string]interface{}{
			"product_count":  a.ProductCount,
			"merchant_count": a.MerchantCount,
			"active_status":  activeStatus,
		}).Error
}

func dayOnly(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 10 && s[4] == '-' && s[7] == '-' {
		return s[:10]
	}
	return s
}

func (r *Repo) ActivityOpsAggregate(ctx context.Context, activityID uint) (salesTotal, stockTotal int64, payMoney float64, err error) {
	type agg struct {
		SalesTotal int64
		StockTotal int64
		PayMoney   float64
	}
	var a agg
	err = r.db.WithContext(ctx).Model(&seckill.Active{}).
		Select("COALESCE(SUM(sales),0) AS sales_total, COALESCE(SUM(stock),0) AS stock_total, COALESCE(SUM(sales * seckill_price),0) AS pay_money").
		Where("activity_id = ? AND delete_time IS NULL", activityID).
		Scan(&a).Error
	return a.SalesTotal, a.StockTotal, a.PayMoney, err
}

func (r *Repo) ActivityPanelStats(ctx context.Context, activityID uint, merID *uint) (ordersPeople, payPeople, payOrders int64, payMoney float64, err error) {
	db := r.db.WithContext(ctx).Model(&seckill.ActivityStatOrder{}).Where("activity_id = ?", activityID)
	if merID != nil && *merID > 0 {
		db = db.Where("mer_id = ?", *merID)
	}
	type panel struct {
		OrdersPeople int64
		PayPeople    int64
		PayOrders    int64
		PayMoney     float64
	}
	var p panel
	err = db.Select(`
		COUNT(DISTINCT CASE WHEN status > -1 THEN uid END) AS orders_people,
		COUNT(DISTINCT CASE WHEN paid = 1 AND status > -1 THEN uid END) AS pay_people,
		COUNT(DISTINCT CASE WHEN paid = 1 AND status > -1 THEN id END) AS pay_orders,
		COALESCE(SUM(CASE WHEN paid = 1 AND status > -1 THEN pay_price ELSE 0 END),0) AS pay_money
	`).Scan(&p).Error
	return p.OrdersPeople, p.PayPeople, p.PayOrders, p.PayMoney, err
}

func (r *Repo) ListActivityStatPeople(ctx context.Context, activityID uint, q seckill.ActivityStatQuery) ([]seckill.ActivityStatPeople, int64, error) {
	db := r.db.WithContext(ctx).Model(&seckill.ActivityStatPeople{}).Where("activity_id = ?", activityID)
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where("nickname LIKE ? OR phone LIKE ? OR CAST(uid AS CHAR) LIKE ?", like, like, like)
	}
	if from := strings.TrimSpace(q.DateFrom); from != "" {
		db = db.Where("DATE(last_join_time) >= ?", from)
	}
	if to := strings.TrimSpace(q.DateTo); to != "" {
		db = db.Where("DATE(last_join_time) <= ?", to)
	}
	if q.MerID != nil && *q.MerID > 0 {
		db = db.Where("mer_id = ?", *q.MerID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.ActivityStatPeople
	err := db.Order("last_join_time DESC, id DESC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListActivityStatOrders(ctx context.Context, activityID uint, q seckill.ActivityStatQuery) ([]seckill.ActivityStatOrder, int64, error) {
	db := r.db.WithContext(ctx).Model(&seckill.ActivityStatOrder{}).Where("activity_id = ?", activityID)
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where("nickname LIKE ? OR order_sn LIKE ? OR CAST(uid AS CHAR) LIKE ?", like, like, like)
	}
	if from := strings.TrimSpace(q.DateFrom); from != "" {
		db = db.Where("DATE(create_time) >= ?", from)
	}
	if to := strings.TrimSpace(q.DateTo); to != "" {
		db = db.Where("DATE(create_time) <= ?", to)
	}
	if q.MerID != nil && *q.MerID > 0 {
		db = db.Where("mer_id = ?", *q.MerID)
	}
	if q.Status != nil && *q.Status > 0 {
		db = db.Where("status = ?", *q.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.ActivityStatOrder
	err := db.Order("create_time DESC, id DESC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListActivityStatProducts(ctx context.Context, activityID uint, q seckill.ActivityStatQuery) ([]seckill.Active, int64, error) {
	db := r.db.WithContext(ctx).Model(&seckill.Active{}).Where("activity_id = ? AND delete_time IS NULL", activityID)
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where("name LIKE ? OR CAST(product_id AS CHAR) LIKE ? OR CAST(seckill_active_id AS CHAR) LIKE ?", like, like, like)
	}
	if q.MerID != nil && *q.MerID > 0 {
		db = db.Where("mer_id = ?", *q.MerID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []seckill.Active
	err := db.Order("sales DESC, seckill_active_id ASC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) LoadProductCategoryName(ctx context.Context, productID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(c.name, '') AS name
		FROM qixi_crm_b_product_view AS p
		LEFT JOIN qixi_crm_b_category_view AS c ON c.category_id = p.category_id
		WHERE p.product_id = ?
		LIMIT 1
	`, productID).Scan(&name).Error
	return name, err
}


func (r *Repo) ListActivityProductsPaged(ctx context.Context, q seckill.ActivityProductQuery) ([]seckill.Active, int64, error) {
	db := r.db.WithContext(ctx).Model(&seckill.Active{}).Where("activity_id = ? AND delete_time IS NULL", q.ActivityID)
	if q.ProductStatus != nil {
		db = db.Where("product_status = ?", *q.ProductStatus)
	}
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where(`(
			name LIKE ? OR CAST(product_id AS CHAR) LIKE ? OR CAST(seckill_active_id AS CHAR) LIKE ?
			OR product_id IN (
				SELECT product_id FROM qixi_crm_b_product_view
				WHERE title LIKE ? OR CAST(product_id AS CHAR) LIKE ?
			)
		)`, like, like, like, like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := q.Page, q.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	var rows []seckill.Active
	err := db.Order("sort DESC, seckill_active_id DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) LoadProductCategoryNames(ctx context.Context, productIDs []uint) (map[uint]string, error) {
	out := map[uint]string{}
	if len(productIDs) == 0 {
		return out, nil
	}
	type row struct {
		ProductID uint   `gorm:"column:product_id"`
		Name      string `gorm:"column:name"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, COALESCE(c.name, '') AS name
		FROM qixi_crm_b_product_view AS p
		LEFT JOIN qixi_crm_b_category_view AS c ON c.category_id = p.category_id
		WHERE p.product_id IN ?
	`, productIDs).Scan(&rows).Error
	if err != nil {
		return out, err
	}
	for _, row := range rows {
		out[row.ProductID] = row.Name
	}
	return out, nil
}

func (r *Repo) LoadProductStock(ctx context.Context, productID uint) (int, error) {
	if productID == 0 {
		return 0, nil
	}
	var stock int
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(stock, 0) FROM qixi_crm_b_product_view WHERE product_id = ? LIMIT 1
	`, productID).Scan(&stock).Error
	return stock, err
}

func (r *Repo) LoadProductSKUs(ctx context.Context, productID uint) ([]seckill.ProductSKURow, error) {
	if productID == 0 {
		return nil, nil
	}
	type row struct {
		MerchantSKUID uint    `gorm:"column:merchant_sku_id"`
		SKUKey        string  `gorm:"column:sku_key"`
		SpecSnapshot  string  `gorm:"column:spec_snapshot"`
		Price         float64 `gorm:"column:price"`
		Stock         int     `gorm:"column:stock"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT merchant_sku_id, sku_key, CAST(spec_snapshot AS CHAR) AS spec_snapshot, price, stock
		FROM qixi_crm_b_product_sku_view
		WHERE product_id = ? AND sale_status = 1
		ORDER BY merchant_sku_id ASC
	`, productID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]seckill.ProductSKURow, 0, len(rows))
	for _, x := range rows {
		out = append(out, seckill.ProductSKURow{
			MerchantSKUID: x.MerchantSKUID,
			SKUKey:        x.SKUKey,
			SpecSnapshot:  x.SpecSnapshot,
			Price:         x.Price,
			Stock:         x.Stock,
		})
	}
	return out, nil
}

func (r *Repo) FindActiveByActivityProduct(ctx context.Context, activityID, productID uint) (*seckill.Active, error) {
	if activityID == 0 || productID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var row seckill.Active
	err := r.db.WithContext(ctx).
		Where("activity_id = ? AND product_id = ? AND delete_time IS NULL", activityID, productID).
		Order("seckill_active_id DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

var _ seckill.Store = (*Repo)(nil)
