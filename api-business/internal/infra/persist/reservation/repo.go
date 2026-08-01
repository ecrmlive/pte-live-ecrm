package reservationpersist

import (
	"context"
	"fmt"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/reservation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repo 只访问 qixi_crm_b_ 预约活动、时段和占位账本。
// 商品/商户展示信息来自业务读模型 qixi_crm_b_product_view。
type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListProducts(ctx context.Context, merID *uint, page, limit int) ([]reservation.ProductView, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Table("qixi_crm_b_reservation_activity AS a").
		Joins("INNER JOIN qixi_crm_b_product_view p ON p.product_id = a.product_id").
		Where("a.status = 1 AND p.sale_status = 1")
	if merID != nil {
		q = q.Where("a.merchant_id = ?", *merID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	type row struct {
		ProductID uint
		MerID     uint
		StoreID   uint
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
	err := q.Select(`p.product_id, p.merchant_id AS mer_id, p.store_id, p.title AS store_name,
		p.cover_url AS image, p.price, COALESCE(p.original_price, 0) AS ot_price, p.stock,
		p.merchant_name AS mer_name, a.show_reservation_days AS show_days, a.reservation_type AS res_type`).
		Order("a.product_reservation_id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]reservation.ProductView, 0, len(rows))
	for _, x := range rows {
		out = append(out, reservation.ProductView{
			ProductID: x.ProductID, MerID: x.MerID, StoreID: x.StoreID, MerName: x.MerName,
			StoreName: x.StoreName, Image: x.Image, Price: x.Price, OtPrice: x.OtPrice, Stock: x.Stock,
			ShowDays: x.ShowDays, ReserveType: x.ResType,
		})
	}
	return out, total, nil
}

func (r *Repo) GetProduct(ctx context.Context, productID uint) (*reservation.ProductView, error) {
	type row struct {
		ProductID uint
		MerID     uint
		StoreID   uint
		StoreName string
		Image     string
		Price     float64
		OtPrice   float64
		Stock     uint
		MerName   string
		ShowDays  int
		ResType   int8
	}
	var x row
	err := r.db.WithContext(ctx).Table("qixi_crm_b_reservation_activity AS a").
		Joins("INNER JOIN qixi_crm_b_product_view p ON p.product_id = a.product_id").
		Where("a.product_id = ? AND a.status = 1 AND p.sale_status = 1", productID).
		Select(`p.product_id, p.merchant_id AS mer_id, p.store_id, p.title AS store_name, p.cover_url AS image,
		p.price, COALESCE(p.original_price, 0) AS ot_price, p.stock, p.merchant_name AS mer_name,
		a.show_reservation_days AS show_days, a.reservation_type AS res_type`).
		First(&x).Error
	if err != nil {
		return nil, err
	}
	return &reservation.ProductView{
		ProductID: x.ProductID, MerID: x.MerID, StoreID: x.StoreID, MerName: x.MerName, StoreName: x.StoreName,
		Image: x.Image, Price: x.Price, OtPrice: x.OtPrice, Stock: x.Stock,
		ShowDays: x.ShowDays, ReserveType: x.ResType,
	}, nil
}

func (r *Repo) GetConfig(ctx context.Context, productID uint) (*reservation.Config, error) {
	var c reservation.Config
	err := r.db.WithContext(ctx).Where("product_id = ? AND status = 1", productID).First(&c).Error
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
			"merchant_id":           c.MerID,
			"store_id":              c.StoreID,
			"reservation_type":      c.ReservationType,
			"show_reservation_days": c.ShowReservationDays,
			"is_cancel_reservation": c.IsCancelReservation,
			"time_period":           c.TimePeriod,
			"status":                1,
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
	err := r.db.WithContext(ctx).Model(&reservation.Booking{}).
		Where("product_id = ? AND slot_id = ? AND booking_date = ? AND status = 1", productID, slotID, date).
		Count(&n).Error
	return n, err
}

// Book 以时段行锁串行化同一时段的并发创建；只有账本写入成功时，外层订单事务才会提交。
func (r *Repo) Book(ctx context.Context, productID, slotID uint, date string, orderID, uid uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var slot reservation.Slot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("attr_reservation_id = ? AND product_id = ?", slotID, productID).First(&slot).Error; err != nil {
			return err
		}
		var booked int64
		if err := tx.Model(&reservation.Booking{}).
			Where("product_id = ? AND slot_id = ? AND booking_date = ? AND status = 1", productID, slotID, date).
			Count(&booked).Error; err != nil {
			return err
		}
		if booked >= int64(slot.Stock) {
			return reservation.ErrFull
		}
		row := &reservation.Booking{ProductID: productID, SlotID: slotID, Date: date, OrderID: orderID, UID: uid, Status: 1}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		return tx.Model(&reservation.Slot{}).Where("attr_reservation_id = ?", slotID).
			UpdateColumn("use_num", gorm.Expr("use_num + 1")).Error
	})
}

func errorsIsNotFound(err error) bool { return err == gorm.ErrRecordNotFound }

var _ reservation.Store = (*Repo)(nil)
