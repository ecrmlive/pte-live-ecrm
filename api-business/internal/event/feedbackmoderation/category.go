package feedbackmoderation

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type category struct {
	ID        uint64 `gorm:"column:id"`
	PID       uint64 `gorm:"column:pid"`
	Name      string `gorm:"column:name"`
	Sort      int    `gorm:"column:sort"`
	Status    int    `gorm:"column:status"`
	DeletedAt any    `gorm:"column:deleted_at"`
}

// applyCategory owns category mutations in the business service.  The platform
// process can read projections but cannot write business facts directly.
func applyCategory(ctx context.Context, db *gorm.DB, in command) (out result, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old struct {
			CategoryID uint64 `gorm:"column:category_id"`
			Action     string `gorm:"column:action"`
		}
		if e := tx.Table("qixi_crm_b_user_feedback_category_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&old).Error; e == nil {
			if old.Action != in.Action || (in.CategoryID > 0 && old.CategoryID != in.CategoryID) {
				return errors.New("conflict")
			}
			out = result{CategoryID: old.CategoryID}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}

		if in.Action == "category_create" {
			if e := validateParent(tx, 0, in.PID); e != nil {
				return e
			}
			row := category{PID: in.PID, Name: in.Name, Sort: in.Sort, Status: in.Status}
			if e := tx.Table("qixi_crm_b_user_feedback_category").Create(&row).Error; e != nil {
				return e
			}
			if e := categoryAudit(tx, row, in); e != nil {
				return e
			}
			out = result{CategoryID: row.ID}
			return nil
		}

		var row category
		if e := tx.Table("qixi_crm_b_user_feedback_category").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id=? AND deleted_at IS NULL", in.CategoryID).Take(&row).Error; e != nil {
			return e
		}
		changes := map[string]any{}
		switch in.Action {
		case "category_update":
			if e := validateParent(tx, row.ID, in.PID); e != nil {
				return e
			}
			changes["pid"], changes["name"], changes["sort"], changes["status"] = in.PID, in.Name, in.Sort, in.Status
			row.PID, row.Name, row.Sort, row.Status = in.PID, in.Name, in.Sort, in.Status
		case "category_status":
			changes["status"], row.Status = in.Status, in.Status
		case "category_delete":
			var childCount int64
			if e := tx.Table("qixi_crm_b_user_feedback_category").Where("pid=? AND deleted_at IS NULL", row.ID).Count(&childCount).Error; e != nil {
				return e
			}
			if childCount > 0 {
				return errHasChild
			}
			changes["deleted_at"] = time.Now()
		default:
			return errors.New("invalid")
		}
		if e := tx.Table("qixi_crm_b_user_feedback_category").Where("id=? AND deleted_at IS NULL", row.ID).Updates(changes).Error; e != nil {
			return e
		}
		if e := categoryAudit(tx, row, in); e != nil {
			return e
		}
		out = result{CategoryID: row.ID}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result{CategoryID: in.CategoryID, Code: "not_found"}, err
	}
	if errors.Is(err, errInvalidParent) {
		return result{CategoryID: in.CategoryID, Code: "invalid_parent"}, err
	}
	if errors.Is(err, errHasChild) {
		return result{CategoryID: in.CategoryID, Code: "has_child"}, err
	}
	if err != nil {
		return result{CategoryID: in.CategoryID, Code: "conflict"}, err
	}
	return out, nil
}

var (
	errInvalidParent = errors.New("invalid_parent")
	errHasChild      = errors.New("has_child")
)

// validateParent enforces CRMEB-style two-level trees: parent must be a top-level
// category (or 0), and cannot be the node itself.
func validateParent(tx *gorm.DB, selfID, pid uint64) error {
	if pid == 0 {
		return nil
	}
	if selfID > 0 && pid == selfID {
		return errInvalidParent
	}
	var parent category
	if e := tx.Table("qixi_crm_b_user_feedback_category").Where("id=? AND deleted_at IS NULL", pid).Take(&parent).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return errInvalidParent
		}
		return e
	}
	if parent.PID != 0 {
		return errInvalidParent
	}
	return nil
}

func categoryAudit(tx *gorm.DB, row category, in command) error {
	return tx.Table("qixi_crm_b_user_feedback_category_audit").Create(map[string]any{
		"category_id": row.ID, "action": in.Action, "name": row.Name, "sort": row.Sort, "status": row.Status,
		"operator_admin_id": in.OperatorID, "idempotency_key": in.IdempotencyKey, "created_at": time.Now(),
	}).Error
}
