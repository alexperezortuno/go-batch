package handler

import (
	"bytes"
	"encoding/csv"
	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/repository"
	"github.com/alexperezortuno/go-batch/internal/service"
	"github.com/alexperezortuno/go-batch/internal/utils/logger"
	"io"
	"os"
	"strconv"
)

func ProcessUserCSV(cfg config.Config, db *repository.Database, logger *logger.Logger) error {
	data, err := os.ReadFile(cfg.FileProcessing.Paths.UserFile)
	if err != nil {
		logger.Error("Error reading file", err)
		return err
	}

	repo := &repository.LoaderRepo{DB: db.Db}
	svc := &service.LoaderService{Repo: repo}

	reader := csv.NewReader(bytes.NewReader(data))
	var batch []domain.User
	batchSize := cfg.FileProcessing.BatchSize

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		age, _ := strconv.Atoi(record[2])
		user := domain.User{
			Name:  record[0],
			Email: record[1],
			Age:   age,
		}
		batch = append(batch, user)

		if len(batch) >= batchSize {
			if err := svc.InsertUserBatch(batch); err != nil {
				return err
			}
			batch = []domain.User{}
		}
	}

	if len(batch) > 0 {
		return svc.InsertUserBatch(batch)
	}
	return nil
}
