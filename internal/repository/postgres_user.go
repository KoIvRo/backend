package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}
type UserRepository struct{ Pool *pgxpool.Pool }

var ErrUserExists = errors.New("username already exists")

func (r *UserRepository) Create(ctx context.Context, username, hash string) (User, error) {
	var u User
	err := r.Pool.QueryRow(ctx, "INSERT INTO users(username,password_hash,role) VALUES($1,$2,'user') RETURNING id,username,role", username, hash).Scan(&u.ID, &u.Username, &u.Role)
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == "23505" {
		return User{}, ErrUserExists
	}
	return u, err
}
func (r *UserRepository) ByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := r.Pool.QueryRow(ctx, "SELECT id,username,role,password_hash FROM users WHERE username=$1", username).Scan(&u.ID, &u.Username, &u.Role, &u.PasswordHash)
	return u, err
}
func (r *UserRepository) ByID(ctx context.Context, id int) (User, error) {
	var u User
	err := r.Pool.QueryRow(ctx, "SELECT id,username,role,password_hash FROM users WHERE id=$1", id).Scan(&u.ID, &u.Username, &u.Role, &u.PasswordHash)
	return u, err
}
func (r *UserRepository) EnsureAdmin(ctx context.Context, username, hash string) error {
	if username == "" || hash == "" {
		return errors.New("admin credentials required")
	}
	_, err := r.Pool.Exec(ctx, "INSERT INTO users(username,password_hash,role) VALUES($1,$2,'admin') ON CONFLICT(username) DO NOTHING", username, hash)
	return err
}

var _ = pgx.ErrNoRows
