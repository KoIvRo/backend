package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KoIvRo/backend/internal/apidocs"
	"github.com/KoIvRo/backend/internal/handler"
	"github.com/KoIvRo/backend/internal/logging"
	"github.com/KoIvRo/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func databaseConfigFromEnv() (*pgxpool.Config, error) {
	host := os.Getenv("POSTGRES_HOST")
	portValue := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	if host == "" || portValue == "" || database == "" || user == "" || password == "" {
		return nil, errors.New("database environment variables are not fully set")
	}

	port, err := strconv.ParseUint(portValue, 10, 16)
	if err != nil || port == 0 {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %q", portValue)
	}

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("create database config: %w", err)
	}

	config.ConnConfig.Host = host
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.Database = database
	config.ConnConfig.User = user
	config.ConnConfig.Password = password

	return config, nil
}

func main() {
	logger := logging.New(os.Stdout)

	databaseConfig, err := databaseConfigFromEnv()
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

	classRepository := repository.NewMemoryClassRepository()
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
