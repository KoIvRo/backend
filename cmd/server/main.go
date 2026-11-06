package main

import (
	"context"
	"net/http"
	"os"
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
	api := handler.New(classRepository, logger)

	mux := http.NewServeMux()
	apidocs.Register(mux)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/classes", api.CreateClass)
	mux.HandleFunc("GET /api/classes", api.ListClasses)
	mux.HandleFunc("GET /api/classes/{id}", api.GetClass)
	mux.HandleFunc("PUT /api/classes/{id}", api.ReplaceClass)
	mux.HandleFunc("PATCH /api/classes/{id}", api.UpdateClass)
	mux.HandleFunc("DELETE /api/classes/{id}", api.DeleteClass)

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
