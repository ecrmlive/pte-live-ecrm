package order

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrOrderNotArchivable = errors.New("订单尚未完成，暂不能隐藏")

// CanArchiveGroup preserves all financial and fulfilment records. It only
// controls whether the owning user can hide a completed or cancelled order.
func CanArchiveGroup(group groupRow, orders []orderRow) bool {
	if group.PayStatus == "closed" {
		return true
	}
	if group.PayStatus != "paid" || len(orders) == 0 {
		return false
	}
	for _, order := range orders {
		if order.Status != "completed" && order.Status != "cancelled" {
			return false
		}
	}
	return true
}

// ArchiveGroup performs a user-visible soft archive only. It does not delete
// group orders, payments, order items, refunds, assets, or provider evidence.
func ArchiveGroup(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var group groupRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", groupOrderID, userID).First(&group).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderOwnership
			}
			return err
		}
		if group.ArchivedAt != nil {
			return nil
		}
		var orders []orderRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("group_order_id = ? AND user_id = ?", group.ID, userID).Find(&orders).Error; err != nil {
			return err
		}
		if !CanArchiveGroup(group, orders) {
			return ErrOrderNotArchivable
		}
		now := time.Now().UTC()
		updated := tx.Model(&groupRow{}).Where("id = ? AND user_id = ? AND user_archived_at IS NULL", group.ID, userID).Update("user_archived_at", now)
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrOrderOwnership
		}
		return nil
	})
}

func (h *Handler) Archive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		bad(c, "订单 ID 错误")
		return
	}
	if err := ArchiveGroup(c.Request.Context(), h.db, uint64(middleware.UID(c)), id); err != nil {
		writeOrderError(c, err)
		return
	}
	response.OK(c, gin.H{"group_order_id": id, "archived": true})
}
