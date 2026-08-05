package merchantapplication

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const ReviewSubject = "qixi.platform.merchant-application.v1"
const Reviewed = "platform.merchant_application.reviewed"

type reviewEnvelope struct {
	EventID   uint64          `json:"event_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type reviewPayload struct {
	SourceApplicationID uint   `json:"source_application_id"`
	Status              string `json:"status"`
	ReviewNote          string `json:"review_note"`
}

// StartReviewProjection applies platform-owned review decisions only to the
// matching business application. The consumer is idempotent: replaying an
// already terminal decision changes no row.
func StartReviewProjection(ctx context.Context, businessDB *gorm.DB, url string) (*nats.Conn, error) {
	if url == "" {
		return nil, nil
	}
	nc, err := nats.Connect(url, nats.Name("pte_live_ecrm_api_business_merchant_application_review"))
	if err != nil {
		return nil, err
	}
	if _, err = nc.Subscribe(ReviewSubject, func(msg *nats.Msg) {
		if applyErr := ApplyReview(ctx, businessDB, msg.Data); applyErr != nil {
			log.Printf("merchant application review projection failed: %v", applyErr)
		}
	}); err != nil {
		nc.Close()
		return nil, err
	}
	return nc, nc.Flush()
}

func ApplyReview(ctx context.Context, db *gorm.DB, raw []byte) error {
	var env reviewEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	if env.EventType != Reviewed {
		return nil
	}
	var payload reviewPayload
	if err := json.Unmarshal(env.Payload, &payload); err != nil {
		return err
	}
	payload.Status = strings.TrimSpace(payload.Status)
	if payload.SourceApplicationID == 0 || (payload.Status != "approved" && payload.Status != "rejected") {
		return nil
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table("qixi_crm_b_merchant_application").Where("id = ? AND status = 'pending'", payload.SourceApplicationID).Updates(map[string]any{
			"status": payload.Status, "review_note": strings.TrimSpace(payload.ReviewNote),
		})
		return result.Error
	})
}
