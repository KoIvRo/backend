package repository

import (
	"context"
	"errors"

	"github.com/KoIvRo/backend/internal/domain"
)

var ErrNotFound = errors.New("class not found")

type ClassRepository interface {
	Create(ctx context.Context, class domain.Class) (domain.Class, error)
	List(ctx context.Context, options ListOptions) ([]domain.Class, error)
	GetByID(ctx context.Context, id int) (domain.Class, error)
	Update(ctx context.Context, id int, class domain.Class) (domain.Class, error)
	Delete(ctx context.Context, id int) error
}

type ListOptions struct {
	Name    string
	Trainer string
	Limit   int
	Offset  int
}
