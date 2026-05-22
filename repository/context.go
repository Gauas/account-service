package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionKey struct{}

func Inject(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, TransactionKey{}, tx)
}

func Extract(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(TransactionKey{}).(*gorm.DB)

	if !ok {
		return nil
	}

	return tx
}

func (r Repository[T]) Resolve(ctx context.Context) *gorm.DB {
	if tx := Extract(ctx); tx != nil {
		return tx.WithContext(ctx)
	}

	return r.db.WithContext(ctx)
}
