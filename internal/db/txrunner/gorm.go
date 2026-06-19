package txrunner

import (
	"context"

	"gorm.io/gorm"
)

type gormRunner struct {
	db *gorm.DB
}

func New(db *gorm.DB) Runner {
	return &gormRunner{db: db}
}

func (r *gormRunner) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(WithTx(ctx, tx))
	})
}
