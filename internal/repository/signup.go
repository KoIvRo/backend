package repository

import (
	"context"
	"errors"

	"github.com/KoIvRo/backend/internal/domain"
)

var (
	ErrClassFull       = errors.New("class has no available places")
	ErrDuplicateSignup = errors.New("client is already signed up for this class")
)

type SignupRepository interface {
	Create(ctx context.Context, classID int, signup domain.Signup) (domain.Signup, error)
	ListByClassID(ctx context.Context, classID int) ([]domain.Signup, error)
}
