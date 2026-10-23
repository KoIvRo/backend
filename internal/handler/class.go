package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.classes.List(r.Context())
	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, classes)
}

func getIdFromRequest(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 0 {
		return 0, errors.New("Invalid Class ID")
	}

	return id, nil

}

func (h *Handler) GetClass(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID", "ID")
		return
	}

	class, err := h.classes.GetByID(r.Context(), id)

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusBadRequest, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, class)
}
