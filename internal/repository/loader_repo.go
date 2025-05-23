package repository

import (
	"github.com/alexperezortuno/go-batch/internal/domain"
	"gorm.io/gorm"
)

type LoaderRepo struct {
	DB *gorm.DB
}

func (r *LoaderRepo) BulkUserInsert(users []domain.User, batchSize int) error {
	if len(users) == 0 {
		return nil // Nada que insertar
	}

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		if len(batch) == 0 {
			continue
		}

		if err := r.DB.Create(&batch).Error; err != nil {
			return err
		}
	}
	return nil
}
