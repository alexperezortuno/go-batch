package service

import (
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/repository"
)

type LoaderService struct {
	Repo *repository.LoaderRepo
}

func (s *LoaderService) InsertUsers(users []domain.User, batchSize int) error {
	return s.Repo.BulkInsert(users, batchSize)
}
