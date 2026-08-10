package combinationpersist

import (
	"context"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/combination"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListGroups(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]combination.ProductGroup, int64, error) {
	q := r.db.WithContext(ctx).Model(&combination.ProductGroup{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyOn {
		q = q.Where("is_show = 1 AND status = 1 AND action_status = 1 AND start_time <= NOW() AND end_time >= NOW()")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []combination.ProductGroup
	err := q.Order("product_group_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetGroup(ctx context.Context, id uint) (*combination.ProductGroup, error) {
	var row combination.ProductGroup
	err := r.db.WithContext(ctx).Where("product_group_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateGroup(ctx context.Context, g *combination.ProductGroup) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *Repo) UpdateGroup(ctx context.Context, g *combination.ProductGroup) error {
	return r.db.WithContext(ctx).Model(g).Where("product_group_id = ?", g.ProductGroupID).Updates(map[string]interface{}{
		"price": g.Price, "buying_count_num": g.BuyingCountNum, "time": g.Time,
		"start_time": g.StartTime, "end_time": g.EndTime, "is_show": g.IsShow, "status": g.Status,
		"action_status": g.ActionStatus, "product_status": g.ProductStatus, "refusal": g.Refusal,
	}).Error
}

func (r *Repo) SoftDeleteGroup(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&combination.ProductGroup{}).Where("product_group_id = ?", id).
		Update("is_del", 1).Error
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		MerID     uint    `gorm:"column:mer_id"`
		MerName   string  `gorm:"column:mer_name"`
		Cost      float64 `gorm:"column:cost"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view AS p").
		Select("p.title AS store_name, p.cover_url AS image, p.price, p.merchant_id AS mer_id, p.merchant_name AS mer_name, COALESCE(p.original_price, p.price) AS cost").
		Where("p.product_id = ? AND p.sale_status = 1", productID).
		Scan(&row).Error
	if err != nil {
		return "", "", "", 0, 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", "", 0, 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.MerName, row.Price, row.Cost, row.MerID, nil
}

func (r *Repo) CreateBuying(ctx context.Context, b *combination.Buying) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *Repo) GetBuying(ctx context.Context, id uint) (*combination.Buying, error) {
	var row combination.Buying
	err := r.db.WithContext(ctx).Where("group_buying_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateBuying(ctx context.Context, b *combination.Buying) error {
	return r.db.WithContext(ctx).Model(b).Where("group_buying_id = ?", b.GroupBuyingID).Updates(map[string]interface{}{
		"status": b.Status, "yet_buying_num": b.YetBuyingNum,
	}).Error
}

func (r *Repo) ListOpenBuyings(ctx context.Context, productGroupID uint, limit int) ([]combination.Buying, error) {
	var rows []combination.Buying
	err := r.db.WithContext(ctx).
		Where("product_group_id = ? AND status = 0 AND is_del = 0", productGroupID).
		Order("group_buying_id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *Repo) ListBuyingsAdmin(ctx context.Context, q combination.AdminBuyingQuery) ([]combination.Buying, int64, error) {
	db := r.db.WithContext(ctx).Model(&combination.Buying{}).Where("is_del = 0")
	if q.MerID != nil {
		db = db.Where("mer_id = ?", *q.MerID)
	}
	if len(q.MerIDs) > 0 {
		db = db.Where("mer_id IN ?", q.MerIDs)
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if from := strings.TrimSpace(q.DateFrom); from != "" {
		db = db.Where("create_time >= ?", from+" 00:00:00")
	}
	if to := strings.TrimSpace(q.DateTo); to != "" {
		db = db.Where("create_time <= ?", to+" 23:59:59")
	}
	if user := strings.TrimSpace(q.UserName); user != "" {
		like := "%" + user + "%"
		memberQ := r.db.WithContext(ctx).Model(&combination.Member{}).
			Select("group_buying_id").
			Where("is_del = 0 AND (is_initiator = 1 OR is_leader = 1)").
			Where("(nickname LIKE ? OR CAST(uid AS CHAR) = ?)", like, user)
		if uid, err := strconv.ParseUint(user, 10, 64); err == nil && uid > 0 {
			memberQ = r.db.WithContext(ctx).Model(&combination.Member{}).
				Select("group_buying_id").
				Where("is_del = 0 AND (is_initiator = 1 OR is_leader = 1)").
				Where("(nickname LIKE ? OR CAST(uid AS CHAR) = ? OR uid = ?)", like, user, uid)
		}
		db = db.Where("group_buying_id IN (?)", memberQ)
	}
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		groupByName := r.db.WithContext(ctx).Model(&combination.ProductGroup{}).
			Select("product_group_id").Where("is_del = 0").
			Where("product_id IN (?)",
				r.db.WithContext(ctx).Table("qixi_crm_b_product_view").
					Select("product_id").Where("store_name LIKE ? OR title LIKE ?", like, like),
			)
		conds := []string{"product_group_id IN (?)", "CAST(product_group_id AS CHAR) = ?", "CAST(group_buying_id AS CHAR) = ?"}
		args := []any{groupByName, kw, kw}
		if id, err := strconv.ParseUint(kw, 10, 64); err == nil && id > 0 {
			groupByProduct := r.db.WithContext(ctx).Model(&combination.ProductGroup{}).
				Select("product_group_id").Where("is_del = 0 AND product_id = ?", id)
			conds = append(conds, "product_group_id IN (?)", "group_buying_id = ?")
			args = append(args, groupByProduct, id)
		}
		db = db.Where("("+strings.Join(conds, " OR ")+")", args...)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []combination.Buying
	err := db.Order("group_buying_id DESC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) CreateMember(ctx context.Context, m *combination.Member) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *Repo) GetMemberByOrder(ctx context.Context, orderID uint) (*combination.Member, error) {
	var row combination.Member
	err := r.db.WithContext(ctx).Where("order_id = ? AND is_del = 0", orderID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) FindMember(ctx context.Context, buyingID, uid uint) (*combination.Member, error) {
	var row combination.Member
	err := r.db.WithContext(ctx).Where("group_buying_id = ? AND uid = ? AND is_del = 0", buyingID, uid).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) ListMembers(ctx context.Context, buyingID uint) ([]combination.Member, error) {
	var rows []combination.Member
	err := r.db.WithContext(ctx).Where("group_buying_id = ? AND is_del = 0", buyingID).
		Order("id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) UpdateMember(ctx context.Context, m *combination.Member) error {
	return r.db.WithContext(ctx).Model(m).Where("id = ?", m.ID).Update("status", m.Status).Error
}

func (r *Repo) SoftDeleteMember(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&combination.Member{}).
		Where("id = ? AND is_del = 0", id).Update("is_del", 1).Error
}

func (r *Repo) SoftDeleteBuying(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&combination.Buying{}).
		Where("group_buying_id = ? AND is_del = 0", id).
		Updates(map[string]interface{}{"is_del": 1, "status": -1}).Error
}

func (r *Repo) ListOrderIDsByBuying(ctx context.Context, buyingID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&combination.Member{}).
		Where("group_buying_id = ? AND is_del = 0 AND order_id > 0", buyingID).
		Pluck("order_id", &ids).Error
	return ids, err
}

func (r *Repo) BumpGroupSuccess(ctx context.Context, productGroupID uint) error {
	return r.db.WithContext(ctx).Model(&combination.ProductGroup{}).
		Where("product_group_id = ?", productGroupID).
		UpdateColumn("success_num", gorm.Expr("success_num + 1")).Error
}

var _ combination.Store = (*Repo)(nil)
