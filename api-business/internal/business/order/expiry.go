package order

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errNoExpiredPendingOrder = errors.New("no expired pending order")

// ExpireUnpaid closes at most batch pending orders created before now-TTL.
// Row locks plus the pending predicate make the sweep safe when more than one
// API replica is running. It never closes a paid group or a successful payment.
func ExpireUnpaid(ctx context.Context, db *gorm.DB, now time.Time, ttl time.Duration, batch int) (int, error) {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if batch <= 0 {
		batch = 50
	}
	cutoff := now.UTC().Add(-ttl)
	expired := 0
	for expired < batch {
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var group groupRow
			err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Where("pay_status = ? AND created_at <= ?", "pending", cutoff).Order("id ASC").Limit(1).First(&group).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errNoExpiredPendingOrder
			}
			if err != nil {
				return err
			}
			return closePendingGroup(tx, group)
		})
		if errors.Is(err, errNoExpiredPendingOrder) {
			return expired, nil
		}
		if err != nil {
			return expired, err
		}
		expired++
	}
	return expired, nil
}

// StartUnpaidExpiryWorker consumes the existing job configuration. Multiple
// instances are permitted: the database row lock prevents duplicate releases.
func StartUnpaidExpiryWorker(ctx context.Context, db *gorm.DB, tick, ttl time.Duration, batch int) {
	if tick <= 0 {
		tick = 30 * time.Second
	}
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	if batch <= 0 {
		batch = 50
	}
	run := func() {
		count, err := ExpireUnpaid(ctx, db, time.Now(), ttl, batch)
		if err != nil {
			log.Printf("unpaid order expiry sweep failed: %v", err)
			return
		}
		if count > 0 {
			log.Printf("closed %d expired unpaid order groups", count)
		}
	}
	run()
	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}
