package repository

import (
    "traza-go/pkg/entity"
    "traza-go/pkg/util"
    "gorm.io/gorm"
)

type LoginRepository interface {
    Save(login *entity.Login) error
}

type LoginRepositoryImpl struct {
    base util.BaseRepository[entity.Login]
}

func NewLoginRepository(db *gorm.DB) LoginRepository {
    return &LoginRepositoryImpl{
        base: util.BaseRepository[entity.Login]{DB: db},
    }
}

func (r *LoginRepositoryImpl) Save(login *entity.Login) error {
    return r.base.Save(login)
}
