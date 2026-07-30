package circlepersist

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/circle"
	"gorm.io/gorm"
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
