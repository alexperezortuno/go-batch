package repository

import (
	"context"
	"errors"

	"github.com/alexperezortuno/go-batch/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Database) InsertBatchHeavy(
	ctx context.Context,
	users []domain.User,
) error {
	if len(users) == 0 {
		return nil
	}

	if r.Pool == nil {
		return errors.New("pgx pool not initialized")
	}

	src := newUserCopySource(users)

	_, err := r.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"username", "password", "email", "name", "age"},
		src,
	)

	if err == nil {
		return nil
	}

	if IsUniqueViolation(err) {
		return r.insertWithConflict(ctx, users)
	}

	return err
}
