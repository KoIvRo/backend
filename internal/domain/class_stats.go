package domain

type ClassStats struct {
	ClassID     int    `json:"class_id"`
	Name        string `json:"name"`
	SignupCount int    `json:"signup_count"`
}
