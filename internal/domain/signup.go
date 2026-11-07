package domain

import "strings"

type Signup struct {
	ID      int    `json:"id"`
	ClassID int    `json:"class_id"`
	Client  string `json:"client"`
	Phone   string `json:"phone"`
}

func (s Signup) Validate() *ValidationError {
	if strings.TrimSpace(s.Client) == "" {
		return &ValidationError{
			Field:   "client",
			Message: "client is required",
		}
	}

	if strings.TrimSpace(s.Phone) == "" {
		return &ValidationError{
			Field:   "phone",
			Message: "phone is required",
		}
	}

	return nil
}
