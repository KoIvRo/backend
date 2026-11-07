package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KoIvRo/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSignupRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSignupRepository(pool *pgxpool.Pool) *PostgresSignupRepository {
	return &PostgresSignupRepository{pool: pool}
}

func (r *PostgresSignupRepository) Create(
	ctx context.Context,
	classID int,
	signup domain.Signup,
) (domain.Signup, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Signup{}, fmt.Errorf("begin signup transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const classQuery = `
		SELECT capacity
		FROM classes
		WHERE id = $1
		FOR UPDATE
	`

	var capacity int

	if err := tx.QueryRow(ctx, classQuery, classID).Scan(&capacity); errors.Is(err, pgx.ErrNoRows) {
		return domain.Signup{}, ErrNotFound
	} else if err != nil {
		return domain.Signup{}, fmt.Errorf("lock class for signup: %w", err)
	}

	const duplicateQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM signups
			WHERE class_id = $1
			  AND LOWER(client) = LOWER($2)
		)
	`

	var duplicate bool

	if err := tx.QueryRow(ctx, duplicateQuery, classID, signup.Client).Scan(&duplicate); err != nil {
		return domain.Signup{}, fmt.Errorf("check duplicate signup: %w", err)
	}
	if duplicate {
		return domain.Signup{}, ErrDuplicateSignup
	}

	const countQuery = `SELECT COUNT(*) FROM signups WHERE class_id = $1`

	var signupCount int

	if err := tx.QueryRow(ctx, countQuery, classID).Scan(&signupCount); err != nil {
		return domain.Signup{}, fmt.Errorf("count class signups: %w", err)
	}
	if signupCount >= capacity {
		return domain.Signup{}, ErrClassFull
	}

	const insertQuery = `
		INSERT INTO signups (class_id, client, phone)
		VALUES ($1, $2, $3)
		RETURNING id, class_id, client, phone
	`

	created, err := scanSignup(tx.QueryRow(
		ctx,
		insertQuery,
		classID,
		signup.Client,
		signup.Phone,
	))
	if err != nil {
		var databaseError *pgconn.PgError
		if errors.As(err, &databaseError) && databaseError.Code == "23505" {
			return domain.Signup{}, ErrDuplicateSignup
		}

		return domain.Signup{}, fmt.Errorf("create signup: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Signup{}, fmt.Errorf("commit signup transaction: %w", err)
	}

	return created, nil
}

func (r *PostgresSignupRepository) ListByClassID(
	ctx context.Context,
	classID int,
) ([]domain.Signup, error) {
	const query = `
		SELECT id, class_id, client, phone
		FROM signups
		WHERE class_id = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, classID)
	if err != nil {
		return nil, fmt.Errorf("list signups: %w", err)
	}
	defer rows.Close()

	signups := make([]domain.Signup, 0)

	for rows.Next() {
		signup, err := scanSignup(rows)
		if err != nil {
			return nil, fmt.Errorf("scan signup: %w", err)
		}

		signups = append(signups, signup)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate signups: %w", err)
	}

	return signups, nil
}

func scanSignup(row rowScanner) (domain.Signup, error) {
	var signup domain.Signup

	err := row.Scan(
		&signup.ID,
		&signup.ClassID,
		&signup.Client,
		&signup.Phone,
	)

	return signup, err
}

var _ SignupRepository = (*PostgresSignupRepository)(nil)
