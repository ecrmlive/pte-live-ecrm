package txctx

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

// With attaches a *gorm.DB (typically a transaction) to ctx for cross-domain writes.
func With(ctx context.Context, db *gorm.DB) context.Context {
	if db == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, db)
}

// DB returns the transactional DB from ctx, or fallback.
func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if v, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok && v != nil {
		return v
	}
	return fallback
}
