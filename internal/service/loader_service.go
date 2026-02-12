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

	// Creamos el batcher para model.User
	b := batch.New[domain.User](
		batch.WithSize(batchSize),
		batch.WithTimeout(2*time.Second),
		batch.WithConcurrency(4),
	)

	// Aquí defines qué hacer con cada batch
	b.Start(ctx, func(ctx context.Context, items []domain.User) error {
		// Aquí llamas a tu repositorio
		return s.Repo.InsertBatchHeavy(ctx, items)
	})

	// En vez de insertar uno por uno:
	for _, user := range users {
		b.TrySubmit(user)
	}

	return nil
}
