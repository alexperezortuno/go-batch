package repository

import (
	"github.com/alexperezortuno/go-batch/internal/domain"
	"gorm.io/gorm"
)

type LoaderRepo struct {
	DB *gorm.DB
}

func (r *LoaderRepo) BulkUserInsert(users []domain.User) error {
	return r.DB.Create(&users).Error
}
