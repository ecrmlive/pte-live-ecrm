package merchantpersist

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
	merchantapplicationevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantapplication"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repo struct {
	adminDB *gorm.DB
}

// NewRepo 只读取后台库中的商户监管投影。跨库业务事实通过 api-merchant/NATS
// 同步到该投影，统一后台不直连店铺数据库。
func NewRepo(adminDB *gorm.DB) *Repo { return &Repo{adminDB: adminDB} }

type ListMerchantsFilter struct {
	Keyword     string
	Status      *int8
	CategoryID  *uint
	TypeID      *uint
	RegionID    *uint
	IsBest      *int8
	OfflinePay  *int8
	DateFrom    string
	DateTo      string
	MerchantIDs []uint
	RegionIDs   []uint
	Page        int
	Limit       int
}

const merchantListSelect = `v.merchant_id AS mer_id,
            COALESCE(v.category_id, 0) AS category_id,
            COALESCE(v.type_id, 0) AS type_id,
            COALESCE(v.business_id, 0) AS business_id,
            v.merchant_name AS mer_name,
            COALESCE(b.name, v.owner_name, '') AS owner_name,
            v.contact_name AS real_name,
            v.contact_mobile AS mer_phone,
            COALESCE(v.address, '') AS mer_address,
            COALESCE(v.mer_info, '') AS mer_info,
            COALESCE(v.mer_keyword, '') AS mer_keyword,
            COALESCE(v.mark, '') AS mark,
            v.status,
            v.status AS mer_state,
            COALESCE(v.is_audit, 1) AS is_audit,
            COALESCE(v.is_best, 0) AS is_best,
            COALESCE(v.offline_pay, 0) AS offline_pay,
            COALESCE(v.is_trader, 0) AS is_trader,
            COALESCE(v.is_bro_room, 0) AS is_bro_room,
            COALESCE(v.is_bro_goods, 0) AS is_bro_goods,
            COALESCE(v.commission_switch, 0) AS commission_switch,
            COALESCE(v.commission_rate, 0) AS commission_rate,
            COALESCE(v.mer_account, '') AS mer_account,
            COALESCE(v.sub_mchid, '') AS sub_mchid,
            COALESCE(v.applyment_id, '') AS applyment_id,
            COALESCE(v.care_count, 0) AS care_count,
            COALESCE(v.care_ficti, 0) AS care_ficti,
            COALESCE(v.goods_type, '') AS goods_type,
            COALESCE(v.platform_category_ids, '') AS platform_category_ids,
            COALESCE(v.mer_star, 5) AS mer_star,
            COALESCE(v.sort, 0) AS sort,
            0 AS svip_coupon_merge,
            COALESCE(v.region_id, 0) AS region_id,
            v.created_at AS create_time,
            COALESCE(c.name, '') AS category_name,
            COALESCE(t.name, '') AS type_name,
            COALESCE(z.name, '') AS region_name,
            COALESCE(d.state, 'not_required') AS deposit_state,
            COALESCE(d.required_amount, 0) AS deposit_required,
            COALESCE(d.available_amount, 0) AS deposit_available,
            COALESCE(t.margin, 0) AS type_margin,
            COALESCE(t.is_margin, 0) AS type_is_margin`

func (r *Repo) merchantListBase(ctx context.Context) *gorm.DB {
	return r.adminDB.WithContext(ctx).
		Table("qixi_crm_a_merchant_view AS v").
		Joins("LEFT JOIN qixi_crm_a_merchant_category AS c ON c.id = v.category_id").
		Joins("LEFT JOIN qixi_crm_a_merchant_type AS t ON t.id = v.type_id").
		Joins("LEFT JOIN qixi_crm_a_business_zone AS z ON z.circle_id = v.region_id").
		Joins("LEFT JOIN qixi_crm_a_business_zone AS b ON b.circle_id = v.business_id AND b.type = 1").
		Joins("LEFT JOIN qixi_crm_a_merchant_deposit_account AS d ON d.merchant_id = v.merchant_id").
		Select(merchantListSelect)
}

func (r *Repo) applyMerchantListFilter(q *gorm.DB, f ListMerchantsFilter) *gorm.DB {
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("v.merchant_name LIKE ? OR v.contact_name LIKE ? OR v.contact_mobile LIKE ? OR v.owner_name LIKE ?", like, like, like, like)
	}
	if f.Status != nil {
		q = q.Where("v.status = ?", *f.Status)
	}
	if f.CategoryID != nil {
		q = q.Where("v.category_id = ?", *f.CategoryID)
	}
	if f.TypeID != nil {
		q = q.Where("v.type_id = ?", *f.TypeID)
	}
	if f.RegionID != nil {
		q = q.Where("v.region_id = ?", *f.RegionID)
	}
	if f.IsBest != nil {
		q = q.Where("v.is_best = ?", *f.IsBest)
	}
	if f.OfflinePay != nil {
		q = q.Where("v.offline_pay = ?", *f.OfflinePay)
	}
	if f.DateFrom != "" {
		q = q.Where("v.created_at >= ?", f.DateFrom)
	}
	if f.DateTo != "" {
		q = q.Where("v.created_at < DATE_ADD(?, INTERVAL 1 DAY)", f.DateTo)
	}
	if f.MerchantIDs != nil || f.RegionIDs != nil {
		if len(f.MerchantIDs) == 0 && len(f.RegionIDs) == 0 {
			return q.Where("1 = 0")
		}
		switch {
		case len(f.MerchantIDs) > 0 && len(f.RegionIDs) > 0:
			q = q.Where("(v.merchant_id IN ? OR v.region_id IN ?)", f.MerchantIDs, f.RegionIDs)
		case len(f.MerchantIDs) > 0:
			q = q.Where("v.merchant_id IN ?", f.MerchantIDs)
		default:
			q = q.Where("v.region_id IN ?", f.RegionIDs)
		}
	}
	return q
}

func (r *Repo) ListMerchants(ctx context.Context, f ListMerchantsFilter) ([]merchant.Merchant, int64, error) {
	q := r.applyMerchantListFilter(r.merchantListBase(ctx), f)
	if f.MerchantIDs != nil || f.RegionIDs != nil {
		if len(f.MerchantIDs) == 0 && len(f.RegionIDs) == 0 {
			return []merchant.Merchant{}, 0, nil
		}
	}
	var total int64
	countQ := r.applyMerchantListFilter(
		r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view AS v"),
		f,
	)
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := normalizePage(f.Page, f.Limit)
	var rows []merchant.Merchant
	err := q.Order("v.sort ASC, v.merchant_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	return rows, total, err
}

func (r *Repo) GetMerchant(ctx context.Context, id uint) (*merchant.Merchant, error) {
	var row merchant.Merchant
	err := r.merchantListBase(ctx).Where("v.merchant_id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	ids, err := r.listStoreGroupIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	row.StoreGroupIDs = ids
	row.GoodsTypes = parseCSVInts(row.GoodsType)
	row.PlatformCategoryIDList = parseCSVUints(row.PlatformCategoryIDs)
	return &row, nil
}

func (r *Repo) UpdateMerchantStatus(ctx context.Context, id uint, status, merState int8) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Where("merchant_id = ?", id).Update("status", status).Error
}

func (r *Repo) UpdateSvipCouponMerge(ctx context.Context, merID uint, merge int8) error {
	return fmt.Errorf("商户会员配置应由 api-merchant 处理，merchant_id=%d", merID)
}

func (r *Repo) UpdateMerchant(ctx context.Context, m *merchant.Merchant) error {
	ownerName, err := r.resolveBusinessOwnerName(ctx, m.BusinessID)
	if err != nil {
		return err
	}
	m.OwnerName = ownerName
	if err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Where("merchant_id = ?", m.MerID).
		Updates(map[string]any{
			"merchant_name":     m.MerName,
			"owner_name":        m.OwnerName,
			"contact_name":      m.RealName,
			"contact_mobile":    m.MerPhone,
			"address":           m.MerAddress,
			"mer_info":          m.MerInfo,
			"mer_keyword":       m.MerKeyword,
			"mark":              m.Mark,
			"category_id":       m.CategoryID,
			"type_id":           m.TypeID,
			"business_id":       m.BusinessID,
			"region_id":         nullableUint(m.RegionID),
			"is_best":           m.IsBest,
			"offline_pay":       m.OfflinePay,
			"is_trader":         m.IsTrader,
			"is_audit":          m.IsAudit,
			"is_bro_room":       m.IsBroRoom,
			"is_bro_goods":      m.IsBroGoods,
			"commission_switch": m.CommissionSwitch,
			"commission_rate":   m.CommissionRate,
			"mer_account":       m.MerAccount,
			"sub_mchid":         m.SubMchid,
			"applyment_id":      m.ApplymentID,
			"care_count":             m.CareCount,
			"care_ficti":             m.CareFicti,
			"goods_type":             m.GoodsType,
			"platform_category_ids":  m.PlatformCategoryIDs,
			"mer_star":               m.MerStar,
			"sort":                   m.Sort,
			"status":                 m.Status,
		}).Error; err != nil {
		return err
	}
	if m.StoreGroupIDs != nil {
		return r.replaceStoreGroups(ctx, m.MerID, m.StoreGroupIDs)
	}
	return nil
}

func (r *Repo) CreateMerchant(ctx context.Context, m *merchant.Merchant) error {
	ownerName, err := r.resolveBusinessOwnerName(ctx, m.BusinessID)
	if err != nil {
		return err
	}
	m.OwnerName = ownerName
	var id uint
	if err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Select("COALESCE(MAX(merchant_id),0)+1").Scan(&id).Error; err != nil {
		return err
	}
	if err := r.adminDB.WithContext(ctx).Exec(
		`INSERT INTO qixi_crm_a_merchant_view
		 (merchant_id, merchant_name, owner_name, contact_name, contact_mobile, address, category_id, type_id, business_id, region_id, status,
		  is_best, offline_pay, is_trader, is_audit, is_bro_room, is_bro_goods, commission_switch, commission_rate,
		  mer_keyword, mer_info, mer_account, sub_mchid, applyment_id, care_count, care_ficti, goods_type, platform_category_ids, mer_star, sort, mark, created_at)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())`,
		id, m.MerName, m.OwnerName, m.RealName, m.MerPhone, m.MerAddress, m.CategoryID, m.TypeID, m.BusinessID, nullableUint(m.RegionID),
		m.Status, m.IsBest, m.OfflinePay, m.IsTrader, m.IsAudit, m.IsBroRoom, m.IsBroGoods, m.CommissionSwitch, m.CommissionRate,
		m.MerKeyword, m.MerInfo, m.MerAccount, m.SubMchid, m.ApplymentID, m.CareCount, m.CareFicti, m.GoodsType, m.PlatformCategoryIDs, m.MerStar, m.Sort, m.Mark,
	).Error; err != nil {
		return err
	}
	m.MerID = id
	return r.replaceStoreGroups(ctx, m.MerID, m.StoreGroupIDs)
}

func nullableUint(v uint) any {
	if v == 0 {
		return nil
	}
	return v
}

func (r *Repo) resolveBusinessOwnerName(ctx context.Context, businessID uint) (string, error) {
	if businessID == 0 {
		return "", nil
	}
	var name string
	err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_business_zone").
		Select("name").
		Where("circle_id = ? AND type = 1 AND status = 1", businessID).
		Scan(&name).Error
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", merchant.ErrBadParam
	}
	return name, nil
}

func parseCSVInts(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []int{}
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func parseCSVUints(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint{}
	}
	parts := strings.Split(raw, ",")
	out := make([]uint, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64)
		if err == nil {
			out = append(out, uint(n))
		}
	}
	return out
}

func (r *Repo) listStoreGroupIDs(ctx context.Context, merchantID uint) ([]uint, error) {
	var ids []uint
	err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_store_group_merchant").
		Where("merchant_id = ?", merchantID).
		Pluck("store_group_id", &ids).Error
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

func (r *Repo) replaceStoreGroups(ctx context.Context, merchantID uint, groupIDs []uint) error {
	db := r.adminDB.WithContext(ctx)
	if err := db.Exec(`DELETE FROM qixi_crm_a_store_group_merchant WHERE merchant_id = ?`, merchantID).Error; err != nil {
		return err
	}
	seen := map[uint]struct{}{}
	for _, gid := range groupIDs {
		if gid == 0 {
			continue
		}
		if _, ok := seen[gid]; ok {
			continue
		}
		seen[gid] = struct{}{}
		if err := db.Exec(
			`INSERT INTO qixi_crm_a_store_group_merchant (store_group_id, merchant_id) VALUES (?,?)`,
			gid, merchantID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) UpsertMerchantView(ctx context.Context, m *merchant.Merchant) error {
	return r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_view").
		Exec(`INSERT INTO qixi_crm_a_merchant_view
		 (merchant_id,merchant_name,owner_name,contact_name,contact_mobile,address,category_id,type_id,business_id,region_id,status,
		  is_best,offline_pay,is_trader,is_audit,is_bro_room,is_bro_goods,commission_switch,commission_rate,
		  mer_keyword,mer_info,mer_account,sub_mchid,applyment_id,care_count,care_ficti,sort,mark,created_at)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())
		 ON DUPLICATE KEY UPDATE
		 merchant_name=VALUES(merchant_name),owner_name=VALUES(owner_name),contact_name=VALUES(contact_name),
		 contact_mobile=VALUES(contact_mobile),address=VALUES(address),category_id=VALUES(category_id),
		 type_id=VALUES(type_id),business_id=VALUES(business_id),region_id=VALUES(region_id),status=VALUES(status),
		 is_best=VALUES(is_best),offline_pay=VALUES(offline_pay),is_trader=VALUES(is_trader),
		 is_audit=VALUES(is_audit),is_bro_room=VALUES(is_bro_room),is_bro_goods=VALUES(is_bro_goods),
		 commission_switch=VALUES(commission_switch),commission_rate=VALUES(commission_rate),
		 mer_keyword=VALUES(mer_keyword),mer_info=VALUES(mer_info),mer_account=VALUES(mer_account),
		 sub_mchid=VALUES(sub_mchid),applyment_id=VALUES(applyment_id),care_count=VALUES(care_count),
		 care_ficti=VALUES(care_ficti),sort=VALUES(sort),mark=VALUES(mark)`,
			m.MerID, m.MerName, m.OwnerName, m.RealName, m.MerPhone, m.MerAddress, m.CategoryID, m.TypeID, m.BusinessID, nullableUint(m.RegionID),
			m.Status, m.IsBest, m.OfflinePay, m.IsTrader, m.IsAudit, m.IsBroRoom, m.IsBroGoods, m.CommissionSwitch, m.CommissionRate,
			m.MerKeyword, m.MerInfo, m.MerAccount, m.SubMchid, m.ApplymentID, m.CareCount, m.CareFicti, m.Sort, m.Mark).Error
}

type ListIntentionFilter struct {
	Status    *int8
	Keyword   string
	RegionIDs []uint
	Page      int
	Limit     int
}

func (r *Repo) ListIntentions(ctx context.Context, f ListIntentionFilter) ([]merchant.Intention, int64, error) {
	q := r.adminDB.WithContext(ctx).
		Table("qixi_crm_a_merchant_application").
		Select(`id AS mer_intention_id, COALESCE(source_application_id, 0) AS source_application_id, COALESCE(applicant_user_id, 0) AS uid,
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
	if f.RegionIDs != nil {
		if len(f.RegionIDs) == 0 {
			return []merchant.Intention{}, 0, nil
		}
		q = q.Where("region_id IN ?", f.RegionIDs)
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

func (r *Repo) GetIntention(ctx context.Context, id uint, regionIDs []uint) (*merchant.Intention, error) {
	var row merchant.Intention
	q := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_application").
		Select(`id AS mer_intention_id, COALESCE(source_application_id, 0) AS source_application_id, COALESCE(applicant_user_id, 0) AS uid,
            contact_mobile AS phone, merchant_name AS mer_name, contact_name AS name,
            created_at AS create_time,
            CASE status WHEN 'pending' THEN 0 WHEN 'approved' THEN 1 WHEN 'rejected' THEN 2 ELSE 0 END AS status,
            review_note AS fail_msg, review_note AS mark,
            COALESCE(region_id, 0) AS circle_id,
            0 AS mer_id, 0 AS merchant_category_id, 0 AS mer_type_id, '' AS images`).
		Where("id = ?", id)
	if regionIDs != nil {
		if len(regionIDs) == 0 {
			return nil, gorm.ErrRecordNotFound
		}
		q = q.Where("region_id IN ?", regionIDs)
	}
	err := q.Take(&row).Error
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
	updated := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_application").
		Where("id = ? AND status = 'pending'", row.MerIntentionID).
		Updates(map[string]any{"status": status, "review_note": note, "region_id": row.CircleID, "reviewed_at": gorm.Expr("NOW()")})
	if updated.Error != nil {
		return updated.Error
	}
	if updated.RowsAffected != 1 {
		return merchant.ErrAlreadyAudited
	}
	return merchantapplicationevent.EnqueueReview(r.adminDB.WithContext(ctx), merchantapplicationevent.ReviewPayload{SourceApplicationID: row.SourceApplicationID, Status: status, ReviewNote: note})
}

func (r *Repo) AssignIntentionRegion(ctx context.Context, id, regionID uint) (bool, error) {
	result := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_application").
		Where("id = ? AND status = 'pending'", id).Update("region_id", regionID)
	return result.RowsAffected == 1, result.Error
}

func (r *Repo) ListCategories(ctx context.Context) ([]merchant.Category, error) {
	var rows []merchant.Category
	err := r.adminDB.WithContext(ctx).Order("id ASC").
		Table("qixi_crm_a_merchant_category").
		Select("id AS merchant_category_id, name AS category_name, commission_rate").Scan(&rows).Error
	return rows, err
}

func (r *Repo) CreateCategory(ctx context.Context, c *merchant.Category) error {
	row := struct {
		ID             uint    `gorm:"column:id"`
		Name           string  `gorm:"column:name"`
		CommissionRate float64 `gorm:"column:commission_rate"`
	}{Name: c.CategoryName, CommissionRate: c.CommissionRate}
	if err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_category").Create(&row).Error; err != nil {
		return normalizeCategoryError(err)
	}
	c.MerchantCategoryID = row.ID
	return nil
}

func (r *Repo) UpdateCategory(ctx context.Context, c *merchant.Category) error {
	result := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_category").
		Where("id = ?", c.MerchantCategoryID).Updates(map[string]any{"name": c.CategoryName, "commission_rate": c.CommissionRate})
	if result.Error != nil {
		return normalizeCategoryError(result.Error)
	}
	if result.RowsAffected == 0 {
		var total int64
		if err := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_category").Where("id = ?", c.MerchantCategoryID).Count(&total).Error; err != nil {
			return err
		}
		if total == 0 {
			return merchant.ErrNotFound
		}
	}
	return nil
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	result := r.adminDB.WithContext(ctx).Table("qixi_crm_a_merchant_category").Where("id = ?", id).Delete(nil)
	if result.Error != nil {
		return normalizeCategoryError(result.Error)
	}
	if result.RowsAffected == 0 {
		return merchant.ErrNotFound
	}
	return nil
}

func normalizeCategoryError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return merchant.ErrConflict
	}
	return err
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
