package handler

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/alexperezortuno/go-batch/internal/metrics"
	"github.com/alexperezortuno/go-batch/internal/repository"
	"github.com/alexperezortuno/go-batch/internal/service"
	"github.com/alexperezortuno/go-batch/internal/utils/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProcessUserCSV(cfg config.Config, db *pgxpool.Pool, logger *logger.Logger) error {
	// Check if the file exists
	if _, err := os.Stat(cfg.FileProcessing.Paths.UserFile); os.IsNotExist(err) {
		logger.Error("File does not exist", err)
		return fmt.Errorf("file does not exist: %w", err)
	}

	// read the file an load it into memory
	records, err := readCSVToMemory(cfg, logger)
	if err != nil {
		logger.Error("Error reading file", err)
		return err
	}

	repo := &repository.Database{Pool: db}
	svc := &service.LoaderService{Repo: repo}

	var batch []domain.User
	var validate = validator.New()
	batchSize := cfg.FileProcessing.BatchSize

	for _, record := range records {
		if cfg.FileProcessing.Header {
			// Skip the first line if it is a header
			cfg.FileProcessing.Header = false
			continue
		}

		age, _ := strconv.Atoi(record[4])
		user := domain.User{
			Username: record[0],
			Password: record[1],
			Email:    record[2],
			Name:     record[3],
			Age:      age,
		}

		// Validate the user
		err := validate.Struct(user)
		if err != nil {
			logger.Error("Validation error with value %v", err, user)
			continue
		}

		if cfg.Metrics.Enabled {
			metrics.RecordsInserted.Inc()
		}

		batch = append(batch, user)
	}

	//return svc.InsertUsers(batch, batchSize)
	return svc.ProcessUsers(batch, batchSize)
}

// readCSVToMemory reads an entire CSV file into memory and returns the records as a slice of slices of strings.
func readCSVToMemory(cfg config.Config, logger *logger.Logger) ([][]string, error) {
	filePath := cfg.FileProcessing.Paths.UserFile
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

	if cfg.Metrics.Enabled {
		metrics.RecordsProcessed.Inc()
	}

	return records, nil
}
