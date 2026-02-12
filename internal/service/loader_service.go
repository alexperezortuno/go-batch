package service

import (
	"context"
	"time"

	"github.com/alexperezortuno/go-batch/internal/batch"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/repository"
)

type LoaderService struct {
	Repo *repository.Database
}

func (s *LoaderService) ProcessUsers(users []domain.User, batchSize int) error {
	ctx := context.Background()

	// Create batcher for model.User
	b := batch.New[domain.User](
		batch.WithSize(batchSize),
		batch.WithTimeout(2*time.Second),
		batch.WithConcurrency(4),
		batch.WithRetry(3, 100*time.Second),
	)

	// Here you define what to do with each batch
	b.Start(ctx, func(ctx context.Context, items []domain.User) error {
		// Aquí llamas a tu repositorio
		return s.Repo.InsertBatchHeavy(ctx, items)
	})
	defer b.Stop()

	for _, user := range users {
		err := b.Submit(user)
		if err != nil {
			return err
		}
	}

	return nil
}
