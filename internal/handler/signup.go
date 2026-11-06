package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/KoIvRo/backend/internal/domain"
	"github.com/KoIvRo/backend/internal/repository"
)

type createSignupRequest struct {
	Client string `json:"client"`
	Phone  string `json:"phone"`
}

func (r createSignupRequest) toSignup(classID int) domain.Signup {
	return domain.Signup{
		ClassID: classID,
		Client:  strings.TrimSpace(r.Client),
		Phone:   strings.TrimSpace(r.Phone),
	}
}

func (h *Handler) CreateSignup(w http.ResponseWriter, r *http.Request) {
	classID, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	var request createSignupRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON", "")
		return
	}

	signup := request.toSignup(classID)

	if validationError := signup.Validate(); validationError != nil {
		writeError(w, http.StatusBadRequest, validationError.Message, validationError.Field)
		return
	}

	created, err := h.signups.Create(r.Context(), classID, signup)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	}
	if errors.Is(err, repository.ErrDuplicateSignup) || errors.Is(err, repository.ErrClassFull) {
		writeError(w, http.StatusConflict, err.Error(), "")
		return
	}
	if err != nil {
		writeInternalServerError(w)
		return
	}

	h.logger.Info(
		"signup_created",
		"event", "signup_created",
		"signup_id", created.ID,
		"class_id", created.ClassID,
		"client", created.Client,
	)

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListSignups(w http.ResponseWriter, r *http.Request) {
	classID, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid class ID", "id")
		return
	}

	if _, err := h.classes.GetByID(r.Context(), classID); errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "class not found", "")
		return
	} else if err != nil {
		writeInternalServerError(w)
		return
	}

	signups, err := h.signups.ListByClassID(r.Context(), classID)
	if err != nil {
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, signups)
}
