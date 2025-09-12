package repositories

import (
	"context"

	"gorm.io/gorm"
)

type gormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) IRepository[T] {
	return &gormRepository[T]{db: db}
}

func (r *gormRepository[T]) WithPreload(associations ...string) *gormRepository[T] {
	db := r.db
	for _, assoc := range associations {
		db = db.Preload(assoc)
	}
	return &gormRepository[T]{db: db}
}

func (r *gormRepository[T]) Find(ctx context.Context, where any, order string, limit, offset int) ([]T, error) {
	var results []T
	tx := r.db.WithContext(ctx)
	if where != nil {
		tx = tx.Where(where)
	}
	if order != "" {
		tx = tx.Order(order)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if offset > 0 {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *gormRepository[T]) FindOne(ctx context.Context, where any) (*T, error) {
	var result T
	tx := r.db.WithContext(ctx)
	if where != nil {
		tx = tx.Where(where)
	}
	if err := tx.First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *gormRepository[T]) InsertOne(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *gormRepository[T]) UpdateOne(ctx context.Context, where any, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(new(T)).Where(where).Updates(updates).Error
}

func (r *gormRepository[T]) DeleteOne(ctx context.Context, where any) error {
	return r.db.WithContext(ctx).Where(where).Delete(new(T)).Error
}

func (r *gormRepository[T]) Count(ctx context.Context, where any) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx)
	if where != nil {
		tx = tx.Where(where)
	}
	if err := tx.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
