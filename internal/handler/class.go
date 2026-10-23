package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/KoIvRo/backend/internal/domain"
	"github.com/KoIvRo/backend/internal/repository"
)

type Handler struct {
	classes repository.ClassRepository
}

func New(classes repository.ClassRepository) *Handler {
	return &Handler{classes: classes}
}

type createClassRequest struct {
	Name     string    `json:"name"`
	Trainer  string    `json:"trainer"`
	StartsAt time.Time `json:"starts_at"`
	Capacity int       `json:"capacity"`
}

func (h *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var request createClassRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON", "")
		return
	}

	class := domain.Class{
		Name:     strings.TrimSpace(request.Name),
		Trainer:  strings.TrimSpace(request.Trainer),
		StartsAt: request.StartsAt,
		Capacity: request.Capacity,
	}

	if validationError := class.Validate(); validationError != nil {
		writeError(w, http.StatusBadRequest, validationError.Message, validationError.Field)
		return
	}

	created, err := h.classes.Create(r.Context(), class)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error", "")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.classes.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error", "")
		return
	}

	writeJSON(w, http.StatusOK, classes)
}
