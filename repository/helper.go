package repository

import (
	"context"

	"gorm.io/gorm"
)

func (r *Registry) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx := Extract(ctx); tx != nil {
		return fn(ctx)
	}

	return r.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			return fn(Inject(ctx, tx))
		},
	)
}
