// Package commentmoderation owns application of platform product-comment
// moderation commands inside the business database.
package commentmoderation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CommandSubject = "qixi.platform.product-comment-moderation-command.v1"

type command struct {
	CommentID         uint64   `json:"comment_id"`
	Action            string   `json:"action"`
	OperatorID        uint64   `json:"operator_id"`
	IdempotencyKey    string   `json:"idempotency_key"`
	Note              string   `json:"note,omitempty"`
	ProductID         uint64   `json:"product_id,omitempty"`
	Score             int      `json:"score,omitempty"`
	Content           string   `json:"content,omitempty"`
	VirtualAuthorName string   `json:"virtual_author_name,omitempty"`
	Sort              int      `json:"sort,omitempty"`
	Media             []string `json:"media,omitempty"`
	MediaSet          bool     `json:"media_set,omitempty"`
}
type result struct {
	CommentID uint64 `json:"comment_id"`
	Status    string `json:"status,omitempty"`
	Code      string `json:"code,omitempty"`
}
type comment struct {
	ID                uint64     `gorm:"column:id"`
	ProductID         uint64     `gorm:"column:product_id"`
	StoreID           uint64     `gorm:"column:store_id"`
	Score             int        `gorm:"column:score"`
	Content           string     `gorm:"column:content"`
	Media             string     `gorm:"column:media"`
	VirtualAuthorName string     `gorm:"column:virtual_author_name"`
	Sort              int        `gorm:"column:sort"`
	Source            string     `gorm:"column:source"`
	Status            string     `gorm:"column:status"`
	DeletedAt         *time.Time `gorm:"column:deleted_at"`
}

var (
	errNotFound = errors.New("comment not found")
	errConflict = errors.New("comment moderation conflict")
	errInvalid  = errors.New("comment moderation invalid")
)

func StartCommandSubscriber(ctx context.Context, businessDB *gorm.DB, natsURL string) (*nats.Conn, error) {
	if businessDB == nil || strings.TrimSpace(natsURL) == "" {
		return nil, nil
	}
	nc, err := nats.Connect(natsURL, nats.Name("pte_live_ecrm_api_business_product_comment_moderation"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.QueueSubscribe(CommandSubject, "pte_live_ecrm_business_product_comment_moderation", func(msg *nats.Msg) {
		out, applyErr := Apply(ctx, businessDB, msg.Data)
		if applyErr != nil && out.Code == "" {
			out.Code = "failed"
		}
		body, _ := json.Marshal(out)
		_ = msg.Respond(body)
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nil
}
func Apply(ctx context.Context, db *gorm.DB, raw []byte) (out result, err error) {
	var in command
	if json.Unmarshal(raw, &in) != nil || !valid(in) {
		return result{CommentID: in.CommentID, Code: "invalid"}, errInvalid
	}
	in.IdempotencyKey, in.Note, in.Content, in.VirtualAuthorName = strings.TrimSpace(in.IdempotencyKey), strings.TrimSpace(in.Note), strings.TrimSpace(in.Content), strings.TrimSpace(in.VirtualAuthorName)
	for i := range in.Media {
		in.Media[i] = strings.TrimSpace(in.Media[i])
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var prior struct {
			CommentID uint64 `gorm:"column:comment_id"`
			Status    string `gorm:"column:to_status"`
		}
		if e := tx.Table("qixi_crm_b_product_comment_moderation_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&prior).Error; e == nil {
			if in.CommentID > 0 && prior.CommentID != in.CommentID {
				return errConflict
			}
			out = result{CommentID: prior.CommentID, Status: prior.Status}
			return nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		if in.Action == "create_virtual" {
			return createVirtual(tx, &in, &out)
		}
		var row comment
		if e := tx.Table("qixi_crm_b_product_comment").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id=? AND deleted_at IS NULL", in.CommentID).Take(&row).Error; e != nil {
			return e
		}
		return mutateVirtualOrStatus(tx, row, &in, &out)
	})
	if err != nil && isDuplicate(err) {
		var prior struct {
			CommentID uint64 `gorm:"column:comment_id"`
			Status    string `gorm:"column:to_status"`
		}
		e := db.WithContext(ctx).Table("qixi_crm_b_product_comment_moderation_audit").Where("idempotency_key=?", in.IdempotencyKey).Take(&prior).Error
		if e == nil && (in.CommentID == 0 || prior.CommentID == in.CommentID) {
			return result{CommentID: prior.CommentID, Status: prior.Status}, nil
		}
		if e == nil {
			return result{CommentID: in.CommentID, Code: "conflict"}, errConflict
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result{CommentID: in.CommentID, Code: "not_found"}, errNotFound
	}
	if errors.Is(err, errConflict) {
		return result{CommentID: in.CommentID, Code: "conflict"}, errConflict
	}
	if err != nil {
		return result{CommentID: in.CommentID, Code: "failed"}, err
	}
	return out, nil
}
func createVirtual(tx *gorm.DB, in *command, out *result) error {
	var product struct {
		StoreID uint64 `gorm:"column:store_id"`
	}
	if err := tx.Table("qixi_crm_b_product_view").Select("store_id").Where("product_id=?", in.ProductID).Take(&product).Error; err != nil {
		return err
	}
	media, err := json.Marshal(in.Media)
	if err != nil {
		return err
	}
	row := comment{ProductID: in.ProductID, StoreID: product.StoreID, Score: in.Score, Content: in.Content, Media: string(media), VirtualAuthorName: in.VirtualAuthorName, Sort: in.Sort, Source: "virtual", Status: "published"}
	if err := tx.Table("qixi_crm_b_product_comment").Create(&row).Error; err != nil {
		return err
	}
	if err := audit(tx, row.ID, "", row.Status, in); err != nil {
		return err
	}
	*out = result{CommentID: row.ID, Status: row.Status}
	return nil
}
func mutateVirtualOrStatus(tx *gorm.DB, row comment, in *command, out *result) error {
	from, to := row.Status, row.Status
	changes := map[string]any{}
	switch in.Action {
	case "publish", "hide":
		next, ok := nextStatus(row.Status, in.Action)
		if !ok {
			return errConflict
		}
		to = next
		changes["status"] = next
	case "update_virtual":
		if row.Source != "virtual" {
			return errConflict
		}
		changes["score"], changes["content"], changes["virtual_author_name"] = in.Score, in.Content, in.VirtualAuthorName
		if in.MediaSet {
			media, err := json.Marshal(in.Media)
			if err != nil {
				return err
			}
			changes["media"] = string(media)
		}
	case "sort_virtual":
		if row.Source != "virtual" {
			return errConflict
		}
		changes["sort"] = in.Sort
	case "delete_virtual":
		if row.Source != "virtual" {
			return errConflict
		}
		now := time.Now()
		changes["deleted_at"] = now
		to = "deleted"
	default:
		return errInvalid
	}
	if res := tx.Table("qixi_crm_b_product_comment").Where("id=? AND deleted_at IS NULL", row.ID).Updates(changes); res.Error != nil {
		return res.Error
	} else if res.RowsAffected != 1 {
		return errConflict
	}
	if err := audit(tx, row.ID, from, to, in); err != nil {
		return err
	}
	*out = result{CommentID: row.ID, Status: to}
	return nil
}
func audit(tx *gorm.DB, commentID uint64, from, to string, in *command) error {
	return tx.Table("qixi_crm_b_product_comment_moderation_audit").Create(map[string]any{"comment_id": commentID, "from_status": from, "to_status": to, "action": in.Action, "note": in.Note, "operator_admin_id": in.OperatorID, "idempotency_key": in.IdempotencyKey, "created_at": time.Now()}).Error
}
func valid(in command) bool {
	key, note := strings.TrimSpace(in.IdempotencyKey), strings.TrimSpace(in.Note)
	if in.OperatorID == 0 || len([]rune(key)) < 8 || len([]rune(key)) > 128 || len([]rune(note)) > 500 {
		return false
	}
	switch in.Action {
	case "publish", "hide", "delete_virtual":
		return in.CommentID > 0
	case "create_virtual":
		return in.CommentID == 0 && in.ProductID > 0 && validVirtual(in)
	case "update_virtual":
		return in.CommentID > 0 && validVirtual(in)
	case "sort_virtual":
		return in.CommentID > 0 && in.Sort >= 0 && in.Sort <= 999999
	default:
		return false
	}
}
func validVirtual(in command) bool {
	if in.Score < 1 || in.Score > 5 || len([]rune(strings.TrimSpace(in.Content))) == 0 || len([]rune(strings.TrimSpace(in.Content))) > 2000 || len([]rune(strings.TrimSpace(in.VirtualAuthorName))) == 0 || len([]rune(strings.TrimSpace(in.VirtualAuthorName))) > 64 || in.Sort < 0 || in.Sort > 999999 || len(in.Media) > 9 {
		return false
	}
	for _, media := range in.Media {
		if len([]rune(strings.TrimSpace(media))) == 0 || len([]rune(strings.TrimSpace(media))) > 1024 {
			return false
		}
	}
	return true
}
func nextStatus(current, action string) (string, bool) {
	switch action {
	case "publish":
		return "published", current == "pending" || current == "hidden"
	case "hide":
		return "hidden", current == "pending" || current == "published"
	default:
		return "", false
	}
}
func isDuplicate(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
