package util

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func (r *BaseRepository[T]) Save(entity *T) error {
	return r.DB.Create(entity).Error
}