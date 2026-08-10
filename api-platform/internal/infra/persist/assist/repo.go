package assistpersist

import (
	"context"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/assist"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]assist.ProductAssist, int64, error) {
	q := r.db.WithContext(ctx).Model(&assist.ProductAssist{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyOn {
		q = q.Where("status = 1 AND is_show = 1 AND product_status = 1 AND action_status = 1")
		q = q.Where("start_time <= NOW() AND end_time >= NOW()")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []assist.ProductAssist
	err := q.Order("product_assist_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*assist.ProductAssist, error) {
	var row assist.ProductAssist
	err := r.db.WithContext(ctx).Where("product_assist_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, a *assist.ProductAssist) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repo) Update(ctx context.Context, a *assist.ProductAssist) error {
	return r.db.WithContext(ctx).Model(a).Where("product_assist_id = ?", a.ProductAssistID).Updates(map[string]interface{}{
		"start_time": a.StartTime, "end_time": a.EndTime, "status": a.Status,
		"assist_count": a.AssistCount, "assist_user_count": a.AssistUserCount,
		"assist_price": a.AssistPrice, "stock": a.Stock, "is_show": a.IsShow,
		"store_name": a.StoreName, "store_info": a.StoreInfo,
		"product_status": a.ProductStatus, "action_status": a.ActionStatus, "refusal": a.Refusal,
	}).Error
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&assist.ProductAssist{}).Where("product_assist_id = ?", id).
		Updates(map[string]interface{}{"is_del": 1, "is_show": 0, "action_status": -1}).Error
}

func (r *Repo) DecStock(ctx context.Context, id uint, num int) error {
	res := r.db.WithContext(ctx).Model(&assist.ProductAssist{}).
		Where("product_assist_id = ? AND is_del = 0 AND stock >= ?", id, num).
		Update("stock", gorm.Expr("stock - ?", num))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repo) IncStock(ctx context.Context, id uint, num int) error {
	if num <= 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&assist.ProductAssist{}).
		Where("product_assist_id = ? AND is_del = 0", id).
		Update("stock", gorm.Expr("stock + ?", num)).Error
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price, cost float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		Cost      float64 `gorm:"column:cost"`
		MerID     uint    `gorm:"column:mer_id"`
		MerName   string  `gorm:"column:mer_name"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view AS p").
		Select("p.store_name, p.cover_url AS image, p.price, 0 AS cost, p.merchant_id AS mer_id, p.merchant_name AS mer_name").
		Where("p.product_id = ?", productID).
		Limit(1).Scan(&row).Error
	if err != nil {
		return "", "", "", 0, 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", "", 0, 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.MerName, row.Price, row.Cost, row.MerID, nil
}

func (r *Repo) LoadNickname(ctx context.Context, uid uint) (string, error) {
	var nick string
	err := r.db.WithContext(ctx).Table("qixi_crm_b_user").Select("nickname").Where("id = ?", uid).Scan(&nick).Error
	if err != nil {
		return "", err
	}
	if nick == "" {
		nick = "用户"
	}
	return nick, nil
}

func (r *Repo) CreateSet(ctx context.Context, s *assist.AssistSet) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *Repo) GetSet(ctx context.Context, id uint) (*assist.AssistSet, error) {
	var row assist.AssistSet
	err := r.db.WithContext(ctx).Where("product_assist_set_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateSet(ctx context.Context, s *assist.AssistSet) error {
	return r.db.WithContext(ctx).Model(s).Where("product_assist_set_id = ?", s.ProductAssistSetID).Updates(map[string]interface{}{
		"status": s.Status, "yet_assist_count": s.YetAssistCount,
	}).Error
}

func (r *Repo) ListSetsByAssist(ctx context.Context, assistID uint, onlyOpen bool, limit int) ([]assist.AssistSet, error) {
	q := r.db.WithContext(ctx).Model(&assist.AssistSet{}).
		Where("product_assist_id = ? AND is_del = 0", assistID)
	if onlyOpen {
		q = q.Where("status = ?", assist.SetStatusRunning)
	}
	var rows []assist.AssistSet
	err := q.Order("product_assist_set_id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *Repo) ListSetsAdmin(ctx context.Context, q assist.AdminSetQuery) ([]assist.AssistSet, int64, error) {
	db := r.db.WithContext(ctx).Model(&assist.AssistSet{}).Where("is_del = 0")
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
		db = db.Where("uid IN (?)",
			r.db.WithContext(ctx).Table("qixi_crm_b_user").Select("id").Where("nickname LIKE ?", like),
		)
	}
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		assistIDs := r.db.WithContext(ctx).Model(&assist.ProductAssist{}).
			Select("product_assist_id").Where("is_del = 0 AND store_name LIKE ?", like)
		productIDs := r.db.WithContext(ctx).Table("qixi_crm_b_product_view").
			Select("product_id").Where("store_name LIKE ? OR title LIKE ?", like, like)
		conds := []string{
			"product_assist_id IN (?)",
			"product_id IN (?)",
			"CAST(product_assist_set_id AS CHAR) = ?",
			"CAST(product_id AS CHAR) = ?",
		}
		args := []any{assistIDs, productIDs, kw, kw}
		if id, err := strconv.ParseUint(kw, 10, 64); err == nil && id > 0 {
			conds = append(conds, "product_assist_set_id = ?", "product_id = ?")
			args = append(args, id, id)
		}
		db = db.Where("("+strings.Join(conds, " OR ")+")", args...)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []assist.AssistSet
	err := db.Order("product_assist_set_id DESC").
		Offset((q.Page - 1) * q.Limit).Limit(q.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) FindOpenSetByUID(ctx context.Context, assistID, uid uint) (*assist.AssistSet, error) {
	var row assist.AssistSet
	err := r.db.WithContext(ctx).
		Where("product_assist_id = ? AND uid = ? AND is_del = 0 AND status IN ?", assistID, uid,
			[]int{assist.SetStatusRunning, assist.SetStatusDone}).
		Order("product_assist_set_id DESC").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateHelper(ctx context.Context, u *assist.AssistUser) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repo) HasHelped(ctx context.Context, setID, uid uint) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&assist.AssistUser{}).
		Where("product_assist_set_id = ? AND uid = ?", setID, uid).Count(&n).Error
	return n > 0, err
}

func (r *Repo) CountHelpsByUID(ctx context.Context, assistID, uid uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&assist.AssistUser{}).
		Where("product_assist_id = ? AND uid = ?", assistID, uid).Count(&n).Error
	return n, err
}

func (r *Repo) ListHelpers(ctx context.Context, setID uint) ([]assist.AssistUser, error) {
	var rows []assist.AssistUser
	err := r.db.WithContext(ctx).Where("product_assist_set_id = ?", setID).
		Order("product_assist_user_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) CountAssistStats(ctx context.Context, assistID uint) (success, pay, all int, err error) {
	if assistID == 0 {
		return 0, 0, 0, nil
	}
	var successN, payN, allN int64
	if err = r.db.WithContext(ctx).Model(&assist.AssistSet{}).
		Where("product_assist_id = ? AND is_del = 0 AND status IN ?", assistID,
			[]int{assist.SetStatusDone, assist.SetStatusPaid}).
		Count(&successN).Error; err != nil {
		return 0, 0, 0, err
	}
	if err = r.db.WithContext(ctx).Model(&assist.AssistSet{}).
		Where("product_assist_id = ? AND is_del = 0 AND status = ?", assistID, assist.SetStatusPaid).
		Count(&payN).Error; err != nil {
		return 0, 0, 0, err
	}
	if err = r.db.WithContext(ctx).Model(&assist.AssistUser{}).
		Where("product_assist_id = ?", assistID).
		Count(&allN).Error; err != nil {
		return 0, 0, 0, err
	}
	return int(successN), int(payN), int(allN), nil
}

var _ assist.Store = (*Repo)(nil)
