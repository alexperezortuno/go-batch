package repository

import (
	"strconv"

	"github.com/alexperezortuno/go-batch/internal/domain"
)

type userCopySource struct {
	users []domain.User
	idx   int
	err   error
}

func newUserCopySource(users []domain.User) *userCopySource {
	return &userCopySource{
		users: users,
	}
}

func (s *userCopySource) Next() bool {
	if s.err != nil {
		return false
	}

	return s.idx < len(s.users)
}

func (s *userCopySource) Values() ([]any, error) {
	u := s.users[s.idx]
	s.idx++

	return []any{
		u.Username,
		u.Password,
		u.Email,
		u.Name,
		strconv.Itoa(u.Age),
	}, nil
}

func (s *userCopySource) Err() error {
	return s.err
}
