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
		batch.WithRetry(3, 100*time.Second),
	)

	// Aquí defines qué hacer con cada batch
	b.Start(ctx, func(ctx context.Context, items []domain.User) error {
		// Aquí llamas a tu repositorio
		return s.Repo.InsertBatchHeavy(ctx, items)
	})
	defer b.Stop()

	// En vez de insertar uno por uno:
	for _, user := range users {
		err := b.Submit(user)
		if err != nil {
			return err
		}
	}

	return nil
}

//func (s *LoaderService) ProcessFile(reader io.Reader) error {
//
//	ctx := context.Background()
//
//	b := batch.New[model.User](
//		batch.WithSize(1000),
//		batch.WithConcurrency(4),
//	)
//
//	b.Start(ctx, func(ctx context.Context, items []model.User) error {
//		return s.repo.InsertBatchHeavy(ctx, items)
//	})
//	defer b.Stop()
//
//	scanner := csv.NewReader(reader)
//
//	for {
//		record, err := scanner.Read()
//		if err == io.EOF {
//			break
//		}
//
//		user := parse(record)
//
//		if err := b.Submit(user); err != nil {
//			return err
//		}
//	}
//
//	return nil
//
//}
