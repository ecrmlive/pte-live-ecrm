package feedbackmoderation

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const Subject = "qixi.platform.feedback-moderation-command.v1"

type command struct {
	FeedbackID     uint64 `json:"feedback_id"`
	CategoryID     uint64 `json:"category_id,omitempty"`
	Action         string `json:"action"`
	Reply          string `json:"reply"`
	Name           string `json:"name,omitempty"`
	PID            uint64 `json:"pid,omitempty"`
	Sort           int    `json:"sort,omitempty"`
	Status         int    `json:"status,omitempty"`
	OperatorID     uint64 `json:"operator_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
type result struct {
	FeedbackID uint64 `json:"feedback_id"`
	CategoryID uint64 `json:"category_id,omitempty"`
	Status     string `json:"status,omitempty"`
	Code       string `json:"code,omitempty"`
}
type feedback struct {
	ID     uint64 `gorm:"column:id"`
	Status string `gorm:"column:status"`
}

func Start(ctx context.Context, db *gorm.DB, url string) (*nats.Conn, error) {
	if db == nil || strings.TrimSpace(url) == "" {
		return nil, nil
	}
	nc, e := nats.Connect(url, nats.Name("pte_live_ecrm_api_business_feedback_moderation"))
	if e != nil {
		return nil, e
	}
	_, e = nc.QueueSubscribe(Subject, "pte_live_ecrm_business_feedback_moderation", func(m *nats.Msg) { out, _ := Apply(ctx, db, m.Data); b, _ := json.Marshal(out); _ = m.Respond(b) })
	if e != nil {
		nc.Close()
		return nil, e
	}
	return nc, nil
}
func Apply(ctx context.Context, db *gorm.DB, raw []byte) (out result, err error) {
	var in command
	if json.Unmarshal(raw, &in) != nil || !valid(in) {
		return result{FeedbackID: in.FeedbackID, Code: "invalid"}, errors.New("invalid")
	}
	in.Reply = strings.TrimSpace(in.Reply)
	in.Name = strings.TrimSpace(in.Name)
	if strings.HasPrefix(in.Action, "category_") {
		return applyCategory(ctx, db, in)
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old struct {
			FeedbackID uint64 `gorm:"column:feedback_id"`
			Status     string `gorm:"column:to_status"`
			Action     string `gorm:"column:action"`
		}
		if e := tx.Table("qixi_crm_b_user_feedback_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&old).Error; e == nil {
			if old.FeedbackID != in.FeedbackID || old.Action != in.Action {
				return errors.New("conflict")
			}
			out = result{FeedbackID: old.FeedbackID, Status: old.Status}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		var row feedback
		if e := tx.Table("qixi_crm_b_user_feedback").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id=? AND deleted_at IS NULL", in.FeedbackID).Take(&row).Error; e != nil {
			return e
		}
		next, ok := next(row.Status, in.Action)
		if !ok {
			return errors.New("conflict")
		}
		changes := map[string]any{"status": next}
		where := "id=? AND status=?"
		args := []any{row.ID, row.Status}
		if in.Action == "delete" {
			changes = map[string]any{"deleted_at": time.Now()}
			where = "id=? AND deleted_at IS NULL"
			args = []any{row.ID}
		}
		if in.Reply != "" {
			changes["reply"] = in.Reply
		}
		if e := tx.Table("qixi_crm_b_user_feedback").Where(where, args...).Updates(changes).Error; e != nil {
			return e
		}
		if e := tx.Table("qixi_crm_b_user_feedback_audit").Create(map[string]any{"feedback_id": row.ID, "from_status": row.Status, "to_status": next, "action": in.Action, "reply": in.Reply, "operator_admin_id": in.OperatorID, "idempotency_key": in.IdempotencyKey, "created_at": time.Now()}).Error; e != nil {
			return e
		}
		out = result{FeedbackID: row.ID, Status: next}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result{FeedbackID: in.FeedbackID, Code: "not_found"}, err
	}
	if err != nil {
		return result{FeedbackID: in.FeedbackID, Code: "conflict"}, err
	}
	return out, nil
}
func valid(x command) bool {
	k, r, n := strings.TrimSpace(x.IdempotencyKey), strings.TrimSpace(x.Reply), strings.TrimSpace(x.Name)
	if x.OperatorID == 0 || len([]rune(k)) < 8 || len([]rune(k)) > 128 || len([]rune(r)) > 1000 || x.Sort < 0 || x.Sort > 9999 {
		return false
	}
	if (x.Action == "reply" && r != "") || x.Action == "close" || x.Action == "delete" {
		return x.FeedbackID > 0
	}
	if x.Action == "category_create" {
		return x.CategoryID == 0 && len([]rune(n)) > 0 && len([]rune(n)) <= 32 && (x.Status == 0 || x.Status == 1)
	}
	if x.Action == "category_update" {
		return x.CategoryID > 0 && len([]rune(n)) > 0 && len([]rune(n)) <= 32 && (x.Status == 0 || x.Status == 1)
	}
	if x.Action == "category_status" {
		return x.CategoryID > 0 && (x.Status == 0 || x.Status == 1)
	}
	return x.Action == "category_delete" && x.CategoryID > 0
}
func next(s, a string) (string, bool) {
	if a == "reply" && s == "pending" {
		return "replied", true
	}
	if a == "close" && (s == "pending" || s == "replied") {
		return "closed", true
	}
	if a == "delete" && (s == "pending" || s == "replied" || s == "closed") {
		return "deleted", true
	}
	return "", false
}
