package main

import (
	"log"
	"net/http"

	"github.com/KoIvRo/backend/internal/handler"
	"github.com/KoIvRo/backend/internal/repository"
)

func main() {
	classRepository := repository.NewMemoryClassRepository()
	api := handler.New(classRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/classes", api.CreateClass)
	mux.HandleFunc("GET /api/classes", api.ListClasses)
	mux.HandleFunc("GET /api/classes/{id}", api.GetClass)

	log.Println("server started on http://localhost:8000")

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
