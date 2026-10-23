package domain

import "time"

type Class struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Trainer  string    `json:"trainer"`
	StartsAt time.Time `json:"starts_at"`
	Capasity int       `json:"capasity"`
}
