package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KoIvRo/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClassRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClassRepository(pool *pgxpool.Pool) *PostgresClassRepository {
	return &PostgresClassRepository{pool: pool}
}

func (r *PostgresClassRepository) Create(ctx context.Context, class domain.Class) (domain.Class, error) {
	const query = `
		INSERT INTO classes (name, trainer, starts_at, capacity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, trainer, starts_at, capacity
	`

	created, err := scanClass(r.pool.QueryRow(
		ctx,
		query,
		class.Name,
		class.Trainer,
		class.StartsAt,
		class.Capacity,
	))
	if err != nil {
		return domain.Class{}, fmt.Errorf("create class: %w", err)
	}

	return created, nil
}

func (r *PostgresClassRepository) List(ctx context.Context, options ListOptions) ([]domain.Class, error) {
	const query = `
		SELECT id, name, trainer, starts_at, capacity
		FROM classes
		WHERE ($1::text = '' OR LOWER(name) = LOWER($1))
		  AND ($2::text = '' OR LOWER(trainer) = LOWER($2))
		ORDER BY id
		LIMIT NULLIF($3::integer, 0)
		OFFSET $4
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		options.Name,
		options.Trainer,
		options.Limit,
		options.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list classes: %w", err)
	}
	defer rows.Close()

	classes := make([]domain.Class, 0)

	for rows.Next() {
		class, err := scanClass(rows)
		if err != nil {
			return nil, fmt.Errorf("scan class: %w", err)
		}

		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate classes: %w", err)
	}

	return classes, nil
}

func (r *PostgresClassRepository) GetByID(ctx context.Context, id int) (domain.Class, error) {
	const query = `
		SELECT id, name, trainer, starts_at, capacity
		FROM classes
		WHERE id = $1
	`

	class, err := scanClass(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Class{}, ErrNotFound
	}
	if err != nil {
		return domain.Class{}, fmt.Errorf("get class by ID: %w", err)
	}

	return class, nil
}

func (r *PostgresClassRepository) Update(
	ctx context.Context,
	id int,
	class domain.Class,
) (domain.Class, error) {
	const query = `
		UPDATE classes
		SET name = $2,
		    trainer = $3,
		    starts_at = $4,
		    capacity = $5
		WHERE id = $1
		RETURNING id, name, trainer, starts_at, capacity
	`

	updated, err := scanClass(r.pool.QueryRow(
		ctx,
		query,
		id,
		class.Name,
		class.Trainer,
		class.StartsAt,
		class.Capacity,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Class{}, ErrNotFound
	}
	if err != nil {
		return domain.Class{}, fmt.Errorf("update class: %w", err)
	}

	return updated, nil
}

func (r *PostgresClassRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM classes WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete class: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanClass(row rowScanner) (domain.Class, error) {
	var class domain.Class

	err := row.Scan(
		&class.ID,
		&class.Name,
		&class.Trainer,
		&class.StartsAt,
		&class.Capacity,
	)

	return class, err
}

var _ ClassRepository = (*PostgresClassRepository)(nil)
