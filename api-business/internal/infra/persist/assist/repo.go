package assistpersist

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/assist"
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
		Image     string  `gorm:"column:cover_url"`
		Price     float64 `gorm:"column:price"`
		Cost      float64 `gorm:"column:cost"`
		MerID     uint    `gorm:"column:merchant_id"`
		MerName   string  `gorm:"column:merchant_name"`
	}
	err = r.db.WithContext(ctx).Table("qixi_crm_b_product_view AS p").
		Select("p.store_name, p.cover_url, p.price, p.price AS cost, p.merchant_id, p.merchant_name").
		Where("p.product_id = ? AND p.sale_status = 1", productID).
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

var _ assist.Store = (*Repo)(nil)
