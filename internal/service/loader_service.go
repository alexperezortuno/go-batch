package service

import (
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/repository"
)

type LoaderService struct {
	Repo *repository.LoaderRepo
}

func (s *LoaderService) InsertUserBatch(users []domain.User) error {
	return s.Repo.BulkUserInsert(users)
}
