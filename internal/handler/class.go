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

func (r createClassRequest) toClass() domain.Class {
	return domain.Class{
		Name:     strings.TrimSpace(r.Name),
		Trainer:  strings.TrimSpace(r.Trainer),
		StartsAt: r.StartsAt,
		Capacity: r.Capacity,
	}
}

type updateClassRequest struct {
	Name     *string    `json:"name"`
	Trainer  *string    `json:"trainer"`
	StartsAt *time.Time `json:"starts_at"`
	Capacity *int       `json:"capacity"`
}

func (h *Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var request createClassRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON", "")
		return
	}

	class := request.toClass()

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

func optionsFromRequest(r *http.Request) (repository.ListOptions, *domain.ValidationError) {
	options := repository.ListOptions{Trainer: strings.TrimSpace(r.URL.Query().Get("trainer"))}

	if value := r.URL.Query().Get("limit"); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil || limit <= 0 || limit > 100 {
			return repository.ListOptions{}, &domain.ValidationError{
				Field:   "limit",
				Message: "limit must be between 1 and 100",
			}
		}

		options.Limit = limit
	}

	if value := r.URL.Query().Get("offset"); value != "" {
		offset, err := strconv.Atoi(value)
		if err != nil || offset < 0 {
			return repository.ListOptions{}, &domain.ValidationError{
				Field:   "offset",
				Message: "offset must be zero or greater",
			}
		}

		options.Offset = offset
	}

	return options, nil
}

func (h *Handler) ListClasses(w http.ResponseWriter, r *http.Request) {
	options, validationError := optionsFromRequest(r)
	if validationError != nil {
		writeError(w, http.StatusBadRequest, validationError.Message, validationError.Field)
		return
	}

	classes, err := h.classes.List(r.Context(), options)
	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, classes)
}

func getIDFromRequest(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id <= 0 {
		return 0, errors.New("invalid class ID")
	}

	return id, nil
}

func (h *Handler) GetClass(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	class, err := h.classes.GetByID(r.Context(), id)

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, class)
}

func (h *Handler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	err = h.classes.Delete(r.Context(), id)

	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ReplaceClass(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	var request createClassRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON", "")
		return
	}

	class := request.toClass()

	if validationError := class.Validate(); validationError != nil {
		writeError(w, http.StatusBadRequest, validationError.Message, validationError.Field)
		return
	}

	updated, err := h.classes.Update(r.Context(), id, class)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	class, err := h.classes.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	var request updateClassRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON", "")
		return
	}

	if request.Name != nil {
		class.Name = strings.TrimSpace(*request.Name)
	}

	if request.Trainer != nil {
		class.Trainer = strings.TrimSpace(*request.Trainer)
	}

	if request.StartsAt != nil {
		class.StartsAt = *request.StartsAt
	}

	if request.Capacity != nil {
		class.Capacity = *request.Capacity
	}

	if validationError := class.Validate(); validationError != nil {
		writeError(w, http.StatusBadRequest, validationError.Message, validationError.Field)
		return
	}

	updated, err := h.classes.Update(r.Context(), id, class)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}

	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
