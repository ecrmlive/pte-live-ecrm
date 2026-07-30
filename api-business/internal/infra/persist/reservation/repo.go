package reservationpersist

import (
	"context"
	"fmt"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/reservation"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListProducts(ctx context.Context, merID *uint, page, limit int) ([]reservation.ProductView, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Table("qixi_m_admin_store_product AS p").
		Where("p.is_del = 0 AND p.type = ? AND p.status = 1 AND p.is_show = 1 AND p.mer_status = 1",
			reservation.ProductTypeReservation)
	if merID != nil {
		q = q.Where("p.mer_id = ?", *merID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	type row struct {
		ProductID uint
		MerID     uint
		StoreName string
		Image     string
		Price     float64
		OtPrice   float64
		Stock     uint
		MerName   string
		ShowDays  int
		ResType   int8
	}
	var rows []row
	err := q.Select(`p.product_id, p.mer_id, p.store_name, p.image, p.price, p.ot_price, p.stock,
		IFNULL(m.mer_name,'') AS mer_name,
		IFNULL(r.show_reservation_days,7) AS show_days,
		IFNULL(r.reservation_type,1) AS res_type`).
		Joins("LEFT JOIN qixi_m_admin_merchant m ON m.mer_id = p.mer_id").
		Joins("LEFT JOIN qixi_m_admin_store_product_reservation r ON r.product_id = p.product_id AND r.is_del = 0").
		Order("p.product_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]reservation.ProductView, 0, len(rows))
	for _, x := range rows {
		out = append(out, reservation.ProductView{
			ProductID: x.ProductID, MerID: x.MerID, MerName: x.MerName,
			StoreName: x.StoreName, Image: x.Image, Price: x.Price, OtPrice: x.OtPrice, Stock: x.Stock,
			ShowDays: x.ShowDays, ReserveType: x.ResType,
		})
	}
	return out, total, nil
}

func (r *Repo) GetProduct(ctx context.Context, productID uint) (*reservation.ProductView, error) {
	var p catalog.Product
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_del = 0 AND type = ?", productID, reservation.ProductTypeReservation).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	v := &reservation.ProductView{
		ProductID: p.ProductID, MerID: p.MerID, StoreName: p.StoreName,
		Image: p.Image, Price: p.Price, OtPrice: p.OtPrice, Stock: p.Stock, ShowDays: 7, ReserveType: 1,
	}
	_ = r.db.WithContext(ctx).Table("qixi_m_admin_merchant").Select("mer_name").
		Where("mer_id = ?", p.MerID).Scan(&v.MerName).Error
	var cfg reservation.Config
	if err := r.db.WithContext(ctx).Where("product_id = ? AND is_del = 0", productID).First(&cfg).Error; err == nil {
		v.ShowDays = cfg.ShowReservationDays
		v.ReserveType = cfg.ReservationType
	}
	return v, nil
}

func (r *Repo) GetConfig(ctx context.Context, productID uint) (*reservation.Config, error) {
	var c reservation.Config
	err := r.db.WithContext(ctx).Where("product_id = ? AND is_del = 0", productID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) UpsertConfig(ctx context.Context, c *reservation.Config) error {
	var old reservation.Config
	err := r.db.WithContext(ctx).Where("product_id = ?", c.ProductID).First(&old).Error
	if errorsIsNotFound(err) {
		return r.db.WithContext(ctx).Create(c).Error
	}
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&reservation.Config{}).
		Where("product_id = ?", c.ProductID).
		Updates(map[string]interface{}{
			"reservation_type":      c.ReservationType,
			"show_reservation_days": c.ShowReservationDays,
			"is_cancel_reservation": c.IsCancelReservation,
			"time_period":           c.TimePeriod,
			"is_del":                0,
		}).Error
}

func (r *Repo) ListSlots(ctx context.Context, productID uint) ([]reservation.Slot, error) {
	var rows []reservation.Slot
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).
		Order("start_time ASC, attr_reservation_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) ReplaceSlots(ctx context.Context, productID uint, slots []reservation.Slot) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&reservation.Slot{}).Error; err != nil {
			return err
		}
		for i := range slots {
			slots[i].ProductID = productID
			if slots[i].Unique == "" {
				slots[i].Unique = fmt.Sprintf("rsv%08d", i+1)
			}
			if err := tx.Create(&slots[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) GetSlot(ctx context.Context, slotID uint) (*reservation.Slot, error) {
	var s reservation.Slot
	err := r.db.WithContext(ctx).Where("attr_reservation_id = ?", slotID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repo) CountBooked(ctx context.Context, productID, slotID uint, date string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("qixi_m_app_store_order AS o").
		Joins("INNER JOIN qixi_m_app_store_order_product op ON op.order_id = o.order_id").
		Where("o.is_del = 0 AND o.reservation_id = ? AND o.reservation_date = ? AND op.product_id = ?",
			slotID, date, productID).
		Count(&n).Error
	return n, err
}

func (r *Repo) BumpSlotUse(ctx context.Context, slotID uint) error {
	return r.db.WithContext(ctx).Model(&reservation.Slot{}).
		Where("attr_reservation_id = ?", slotID).
		UpdateColumn("use_num", gorm.Expr("use_num + 1")).Error
}

func errorsIsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

var _ reservation.Store = (*Repo)(nil)
