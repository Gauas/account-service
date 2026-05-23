package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository[T interface{}] struct {
	db *gorm.DB
}

func (r *Repository[T]) Take(ctx context.Context, args ...interface{}) (*T, error) {
	record := new(T)

	err := r.Resolve(ctx).Take(record, args...).Error

	return record, err
}

func (r *Repository[T]) GetAll(ctx context.Context, args ...interface{}) ([]*T, error) {
	var records []*T

	err := r.Resolve(ctx).Find(&records, args...).Error

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.Resolve(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	if err := r.Resolve(ctx).Updates(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

//func (r Repository[T]) UpdateWhere(ctx context.Context, values interface{}, args ...interface{}) error {
//	tx := r.Resolve(ctx)
//	tx = applyArgs(tx, args...)
//
//	return tx.Model(new(T)).Updates(values).Error
//}

func (r *Repository[T]) Delete(ctx context.Context, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("delete conditions required")
	}

	return r.Resolve(ctx).Delete(new(T), args...).Error
}

func (r *Repository[T]) Pluck(ctx context.Context, column string, dest any) error {
	return r.Resolve(ctx).Model(new(T)).Pluck(column, dest).Error
}

func (r *Repository[T]) Exists(ctx context.Context, args ...interface{}) bool {
	record := new(T)

	err := r.Resolve(ctx).Select("1").Take(record, args...).Error
	if err != nil {
		return false
	}

	return true
}

func (r Repository[T]) WithContext(ctx context.Context) *gorm.DB {
	return r.Resolve(ctx)
}
