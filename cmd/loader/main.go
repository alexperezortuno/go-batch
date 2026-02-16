package main

import (
	"context"
	"log"

	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/handler"
	"github.com/alexperezortuno/go-batch/internal/metrics"
	"github.com/alexperezortuno/go-batch/internal/repository"
	"github.com/alexperezortuno/go-batch/internal/utils/logger"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize metrics
	if cfg != nil {
		if cfg.Metrics.Enabled {
			metrics.InitMetrics()
		}

		// Initialize logger
		appLogger := logger.NewLogger(*cfg)
		defer func(appLogger *logger.Logger) {
			err := appLogger.Close()
			if err != nil {
				log.Fatalf("Failed to close logger: %v", err)
			}
		}(appLogger)

		// Context with cancel
		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Initialize service
		appLogger.Info("Starting application")

		pool, err := repository.NewPgxPool(cfg)
		if err != nil {
			log.Fatal(err)
		}

		db, err := repository.NewDatabase(cfg)
		if err != nil {
			log.Fatal(err)
		}

		dbInstance := &repository.Database{
			Db:   db.Db,
			Pool: pool,
		}

		err = handler.ProcessUserCSV(*cfg, dbInstance, appLogger)

		if err != nil {
			appLogger.Error("Error processing CSV file", err)
		}

		appLogger.Info("Application shutdown completed")
	}

	if cfg == nil {
		log.Fatalf("Failed to load config")
	}
}
