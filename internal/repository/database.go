package repository

import (
	"context"

	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db   *gorm.DB
	Pool *pgxpool.Pool
}

func NewPgxPool(cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 10

	return pgxpool.NewWithConfig(context.Background(), config)
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	if cfg.DB.Migrate {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			return nil, err
		}
	}

	return &Database{Db: db}, nil
}
