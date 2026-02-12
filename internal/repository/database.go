package repository

import (
	"context"

	"github.com/alexperezortuno/go-batch/internal/config"
	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
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

func IsUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *Database) insertWithConflict(
	ctx context.Context,
	users []domain.User,
) error {
	tx := r.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, u := range users {
		if err := tx.Exec(`
        INSERT INTO users (username, password, email, name, age)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT (email) DO NOTHING
    `, u.Username, u.Password, u.Email, u.Name, u.Age).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
