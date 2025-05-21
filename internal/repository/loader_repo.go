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

func (r *LoaderRepo) BulkInsert(users []domain.User, batchSize int) error {
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		if err := r.DB.Create(users[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}
