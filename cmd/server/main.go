package main

import (
	"log"
	"net/http"

	"github.com/KoIvRo/backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)

	log.Println("server started on http://localhost:8000")

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
