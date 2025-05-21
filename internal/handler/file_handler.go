package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/repository"
	"github.com/alexperezortuno/go-batch/internal/service"
	"github.com/alexperezortuno/go-batch/internal/utils/logger"
	"os"
	"strconv"
)

func ProcessUserCSV(cfg config.Config, db *repository.Database, logger *logger.Logger) error {
	// Check if the file exists
	if _, err := os.Stat(cfg.FileProcessing.Paths.UserFile); os.IsNotExist(err) {
		logger.Error("File does not exist", err)
		return fmt.Errorf("file does not exist: %w", err)
	}

	// read the file an load it into memory
	records, err := readCSVToMemory(cfg.FileProcessing.Paths.UserFile, logger)
	if err != nil {
		logger.Error("Error reading file", err)
		return err
	}

	repo := &repository.LoaderRepo{DB: db.Db}
	svc := &service.LoaderService{Repo: repo}

	var batch []domain.User
	batchSize := cfg.FileProcessing.BatchSize

	for _, record := range records {
		if cfg.FileProcessing.HasHeader {
			// Skip the first line if it is a header
			cfg.FileProcessing.HasHeader = false
			continue
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

// readCSVToMemory reads an entire CSV file into memory and returns the records as a slice of slices of strings.
func readCSVToMemory(filePath string, logger *logger.Logger) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Error closing file", err)
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo CSV: %w", err)
	}

	return records, nil
}
