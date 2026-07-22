package broadcastpersist

import (
	"context"

	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/broadcast"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context, merID *uint, onlyPublic bool, page, limit int) ([]broadcast.Room, int64, error) {
	q := r.db.WithContext(ctx).Model(&broadcast.Room{}).Where("is_del = 0")
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if onlyPublic {
		q = q.Where("is_show = 1 AND status = ?", broadcast.AuditApproved)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []broadcast.Room
	err := q.Order("sort DESC, star DESC, broadcast_room_id DESC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*broadcast.Room, error) {
	var row broadcast.Room
	err := r.db.WithContext(ctx).Where("broadcast_room_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, room *broadcast.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *Repo) Update(ctx context.Context, room *broadcast.Room) error {
	return r.db.WithContext(ctx).Model(room).Where("broadcast_room_id = ?", room.BroadcastRoomID).
		Updates(map[string]interface{}{
			"name":        room.Name,
			"cover_img":   room.CoverImg,
			"feeds_img":   room.FeedsImg,
			"start_time":  room.StartTime,
			"end_time":    room.EndTime,
			"anchor_name": room.AnchorName,
			"phone":       room.Phone,
			"status":      room.Status,
			"live_status": room.LiveStatus,
			"is_show":     room.IsShow,
			"sort":        room.Sort,
			"star":        room.Star,
			"mark":        room.Mark,
			"refusal":     room.Refusal,
		}).Error
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&broadcast.Room{}).Where("broadcast_room_id = ?", id).
		Update("is_del", 1).Error
}

func (r *Repo) ReplaceGoods(ctx context.Context, roomID uint, productIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("broadcast_room_id = ?", roomID).Delete(&broadcast.RoomGoods{}).Error; err != nil {
			return err
		}
		seen := map[uint]struct{}{}
		for i, pid := range productIDs {
			if pid == 0 {
				continue
			}
			if _, ok := seen[pid]; ok {
				continue
			}
			seen[pid] = struct{}{}
			row := broadcast.RoomGoods{
				BroadcastRoomID: roomID,
				ProductID:       pid,
				OnSale:          1,
				Sort:            i + 1,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) ListGoods(ctx context.Context, roomID uint) ([]broadcast.RoomGoods, error) {
	var rows []broadcast.RoomGoods
	err := r.db.WithContext(ctx).Where("broadcast_room_id = ?", roomID).
		Order("sort ASC, id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) LoadMerName(ctx context.Context, merID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_merchant").
		Select("mer_name").Where("mer_id = ?", merID).Scan(&name).Error
	return name, err
}

func (r *Repo) LoadProductMeta(ctx context.Context, productID uint) (storeName, image string, price float64, merID uint, err error) {
	var row struct {
		StoreName string  `gorm:"column:store_name"`
		Image     string  `gorm:"column:image"`
		Price     float64 `gorm:"column:price"`
		MerID     uint    `gorm:"column:mer_id"`
	}
	err = r.db.WithContext(ctx).Table("qixi_store_product").
		Select("store_name, image, price, mer_id").
		Where("product_id = ? AND is_del = 0", productID).
		Scan(&row).Error
	if err != nil {
		return "", "", 0, 0, err
	}
	if row.MerID == 0 {
		return "", "", 0, 0, gorm.ErrRecordNotFound
	}
	return row.StoreName, row.Image, row.Price, row.MerID, nil
}

var _ broadcast.Store = (*Repo)(nil)
