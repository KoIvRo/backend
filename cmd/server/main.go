package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KoIvRo/backend/internal/apidocs"
	"github.com/KoIvRo/backend/internal/database"
	"github.com/KoIvRo/backend/internal/handler"
	"github.com/KoIvRo/backend/internal/logging"
	"github.com/KoIvRo/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := logging.New(os.Stdout)

	databaseConfig, err := database.PoolConfigFromEnv()
	if err != nil {
		logger.Error(
			"database_configuration_failed",
			"event", "database_configuration_failed",
			"error", err,
		)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, databaseConfig)
	if err != nil {
		logger.Error(
			"database_pool_creation_failed",
			"event", "database_pool_creation_failed",
			"error", err,
		)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error(
			"database_connection_failed",
			"event", "database_connection_failed",
			"error", err,
		)
		os.Exit(1)
	}

	logger.Info(
		"database_connected",
		"event", "database_connected",
	)

	classRepository := repository.NewPostgresClassRepository(pool)
	signupRepository := repository.NewPostgresSignupRepository(pool)
	api := handler.New(classRepository, classRepository, signupRepository, logger)
	users := &repository.UserRepository{Pool: pool}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Error("JWT_SECRET is required")
		os.Exit(1)
	}
	ttlSeconds, err := strconv.Atoi(os.Getenv("JWT_TTL_SECONDS"))
	if err != nil || ttlSeconds <= 0 {
		logger.Error("JWT_TTL_SECONDS must be positive")
		os.Exit(1)
	}
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" || adminPass == "" {
		logger.Error("admin credentials required")
		os.Exit(1)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("invalid admin password")
		os.Exit(1)
	}
	if err := users.EnsureAdmin(ctx, adminUser, string(hash)); err != nil {
		logger.Error("admin creation failed", "error", err)
		os.Exit(1)
	}
	if _, err := pool.Exec(ctx, "UPDATE classes SET owner_id=(SELECT id FROM users WHERE username=$1) WHERE owner_id IS NULL", adminUser); err != nil {
		logger.Error("class ownership migration failed", "error", err)
		os.Exit(1)
	}
	auth := handler.NewAuth(users, secret, time.Duration(ttlSeconds)*time.Second)

	mux := http.NewServeMux()
	apidocs.Register(mux)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /auth/register", auth.Register)
	mux.HandleFunc("POST /auth/login", auth.Login)
	mux.HandleFunc("GET /auth/me", auth.Required(auth.Me))
	mux.HandleFunc("POST /api/classes", auth.Required(api.CreateClass))
	mux.HandleFunc("GET /api/classes", api.ListClasses)
	mux.HandleFunc("GET /api/classes/stats", api.ClassStats)
	mux.HandleFunc("GET /api/classes/{id}", api.GetClass)
	mux.HandleFunc("PUT /api/classes/{id}", auth.Owned(classRepository, api.ReplaceClass))
	mux.HandleFunc("PATCH /api/classes/{id}", auth.Owned(classRepository, api.UpdateClass))
	mux.HandleFunc("DELETE /api/classes/{id}", auth.Owned(classRepository, api.DeleteClass))
	mux.HandleFunc("POST /api/classes/{id}/signups", auth.Required(api.CreateSignup))
	mux.HandleFunc("GET /api/classes/{id}/signups", api.ListSignups)

	logger.Info(
		"server_started",
		"event", "server_started",
		"address", "http://localhost:8000",
	)

	if err := http.ListenAndServe(":8000", handler.LoggingMiddleware(logger, mux)); err != nil {
		logger.Error(
			"server_stopped",
			"event", "server_stopped",
			"error", err,
		)
		os.Exit(1)
	}
}
