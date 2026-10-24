package main

import (
	"net/http"
	"os"

	"github.com/KoIvRo/backend/internal/apidocs"
	"github.com/KoIvRo/backend/internal/handler"
	"github.com/KoIvRo/backend/internal/logging"
	"github.com/KoIvRo/backend/internal/repository"
)

func main() {
	logger := logging.New(os.Stdout)
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
