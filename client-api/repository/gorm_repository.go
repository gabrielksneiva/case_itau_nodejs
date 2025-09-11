package repository

import (
	"errors"

	"gorm.io/gorm"
)


type GormRepository[T any, ID comparable] struct {
	db *gorm.DB
}

func NewGormRepository[T any, ID comparable](db *gorm.DB) *GormRepository[T, ID] {
	return &GormRepository[T, ID]{db: db}
}

func (r *GormRepository[T, ID]) DB() *gorm.DB {
	return r.db
}

func (r *GormRepository[T, ID]) Migrate() error {
	var model T
	return r.db.AutoMigrate(&model)
}

func (r *GormRepository[T, ID]) Create(t *T) error {
	return r.db.Create(t).Error
}

func (r *GormRepository[T, ID]) GetAll() ([]T, error) {
	var list []T
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormRepository[T, ID]) GetByID(id ID) (T, error) {
	var e T
	if err := r.db.First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e, ErrRepoNotFound
		}
		return e, err
	}
	return e, nil
}

func (r *GormRepository[T, ID]) Update(t *T) error {
	return r.db.Save(t).Error
}

func (r *GormRepository[T, ID]) Delete(id ID) error {
	res := r.db.Delete(new(T), id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRepoNotFound
	}
	return nil
}
