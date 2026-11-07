package repository

import (
	"context"

	"github.com/KoIvRo/backend/internal/domain"
)

type ClassStatsRepository interface {
	MostAttended(ctx context.Context) ([]domain.ClassStats, error)
}
