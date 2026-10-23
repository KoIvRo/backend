package domain

import (
	"strings"
	"time"
)

type Class struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Trainer  string    `json:"trainer"`
	StartsAt time.Time `json:"starts_at"`
	Capacity int       `json:"capacity"`
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (c Class) Validate() *ValidationError {
	if strings.TrimSpace(c.Name) == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	if strings.TrimSpace(c.Trainer) == "" {
		return &ValidationError{
			Field:   "trainer",
			Message: "trainer is required",
		}
	}

	if c.StartsAt.IsZero() {
		return &ValidationError{
			Field:   "starts_at",
			Message: "starts_at is required",
		}
	}

	if c.Capacity <= 0 {
		return &ValidationError{
			Field:   "capacity",
			Message: "capacity must be greater than zero",
		}
	}

	return nil
}
