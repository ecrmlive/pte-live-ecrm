package circlepersist

import (
	"context"
	"errors"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListCircles(ctx context.Context, keyword string, status *int8, page, limit int) ([]circle.Circle, int64, error) {
	q := r.db.WithContext(ctx).Model(&circle.Circle{})
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []circle.Circle
	err := q.Order("sort DESC, circle_id ASC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetCircle(ctx context.Context, id uint) (*circle.Circle, error) {
	var row circle.Circle
	if err := r.db.WithContext(ctx).Where("circle_id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateCircle(ctx context.Context, row *circle.Circle) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateCircle(ctx context.Context, row *circle.Circle) error {
	return r.db.WithContext(ctx).Model(&circle.Circle{}).Where("circle_id = ?", row.CircleID).Updates(map[string]any{
		"pid": row.PID, "path": row.Path, "name": row.Name, "circle_agent_id": row.CircleAgentID,
		"commission_type": row.CommissionType, "commission_rate": row.CommissionRate, "level": row.Level,
		"remark": row.Remark, "sort": row.Sort, "status": row.Status, "type": row.Type,
		"role_id": row.RoleID, "business_store_category": row.BusinessStoreCategory,
		"business_store_type": row.BusinessStoreType, "update_time": time.Now(),
	}).Error
}

func (r *Repo) DeleteCircle(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("circle_id = ?", id).Delete(&circle.Circle{}).Error
}

func (r *Repo) CountCircleChildren(ctx context.Context, id uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&circle.Circle{}).Where("pid = ?", id).Count(&count).Error
	return count, err
}

func (r *Repo) ListAgents(ctx context.Context, keyword string, status *int8, page, limit int) ([]circle.Agent, int64, error) {
	q := r.db.WithContext(ctx).Model(&circle.Agent{})
	if keyword != "" {
		q = q.Where("name LIKE ? OR phone LIKE ? OR business_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []circle.Agent
	err := q.Order("status ASC, circle_agent_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetAgent(ctx context.Context, id uint) (*circle.Agent, error) {
	var row circle.Agent
	if err := r.db.WithContext(ctx).Where("circle_agent_id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateAgent(ctx context.Context, row *circle.Agent) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateAgent(ctx context.Context, row *circle.Agent) error {
	return r.db.WithContext(ctx).Model(&circle.Agent{}).Where("circle_agent_id = ?", row.CircleAgentID).Updates(map[string]any{
		"uid": row.UID, "name": row.Name, "phone": row.Phone, "qualification": row.Qualification,
		"remark": row.Remark, "payment_method": row.PaymentMethod, "payment_name": row.PaymentName,
		"payment_account": row.PaymentAccount, "payment_bank": row.PaymentBank, "payment_qr_img": row.PaymentQRImg,
		"type": row.Type, "business_name": row.BusinessName, "business_store_category": row.BusinessStoreCategory,
		"business_store_type": row.BusinessStoreType, "update_time": time.Now(),
	}).Error
}

func (r *Repo) AuditAgent(ctx context.Context, id uint, status int8, reason string, adminID uint, auditTime time.Time) error {
	return r.db.WithContext(ctx).Model(&circle.Agent{}).Where("circle_agent_id = ? AND status = ?", id, circle.AgentPending).Updates(map[string]any{
		"status": status, "audit_reason": reason, "audit_admin_id": adminID, "audit_time": auditTime, "update_time": auditTime,
	}).Error
}

func (r *Repo) RevokeAgent(ctx context.Context, id uint, reason, idempotencyKey string, adminID uint, now time.Time) (bool, error) {
	replayed := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var agent circle.Agent
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("circle_agent_id=?", id).Take(&agent).Error; err != nil {
			return err
		}
		var previous struct {
			Reason     string `gorm:"column:reason"`
			OperatorID uint   `gorm:"column:operator_admin_id"`
		}
		if err := tx.Table("qixi_crm_a_business_zone_agent_command_audit").Where("circle_agent_id=? AND action='revoke' AND idempotency_key=?", id, idempotencyKey).Take(&previous).Error; err == nil {
			if previous.Reason != reason || previous.OperatorID != adminID {
				return circle.ErrCommandConflict
			}
			replayed = true
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if agent.Status == circle.AgentRevoked {
			return circle.ErrAgentRevoked
		}
		if agent.Status != circle.AgentApproved {
			return circle.ErrAgentNotApproved
		}
		if agent.Balance != 0 {
			return circle.ErrAgentBalance
		}
		var bound int64
		if err := tx.Table("qixi_crm_a_business_zone").Where("circle_agent_id=?", id).Count(&bound).Error; err != nil {
			return err
		}
		if bound > 0 {
			return circle.ErrAgentBound
		}
		var adminBound int64
		if err := tx.Table("qixi_crm_a_admin_user").Where("circle_agent_id=?", id).Count(&adminBound).Error; err != nil {
			return err
		}
		if adminBound > 0 {
			return circle.ErrAgentAdminBound
		}
		if err := tx.Model(&circle.Agent{}).Where("circle_agent_id=? AND status<>?", id, circle.AgentRevoked).Updates(map[string]any{"status": circle.AgentRevoked, "update_time": now}).Error; err != nil {
			return err
		}
		return tx.Table("qixi_crm_a_business_zone_agent_command_audit").Create(map[string]any{"circle_agent_id": id, "action": "revoke", "from_status": agent.Status, "to_status": circle.AgentRevoked, "reason": reason, "operator_admin_id": adminID, "idempotency_key": idempotencyKey}).Error
	})
	return replayed, err
}
