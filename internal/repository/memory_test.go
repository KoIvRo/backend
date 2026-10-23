package repository

import (
	"context"
	"testing"
	"time"

	"github.com/KoIvRo/backend/internal/domain"
)

func TestMemoryClassRepository(t *testing.T) {
	repo := NewMemoryClassRepository()
	ctx := context.Background()

	startsAt := time.Date(
		2026, time.November, 2,
		18, 30, 0, 0,
		time.UTC,
	)

	input := domain.Class{
		Name:     "Йога",
		Trainer:  "Анна",
		StartsAt: startsAt,
		Capacity: 15,
	}

	created, err := repo.Create(ctx, input)

	if err != nil {
		t.Fatalf("Error in creating func: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("Unknown ID: %v", err)
	}

	found, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Error in get func: %v", err)
	}

	if found != created {
		t.Errorf("Found != Created: %v", err)
	}

	list, err := repo.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("Error in list func: %v", err)
	} else if len(list) != 1 {
		t.Fatalf("List() returned %d classes, want 1", len(list))
	} else if list[0] != found {
		t.Errorf("Error in list func: %v", list[0])
	}

	newTrainer := "Иван"

	found.Trainer = newTrainer

	updatedClass, err := repo.Update(ctx, found.ID, found)
	if err != nil {
		t.Fatalf("Error in update func: %v", err)
	}

	err = repo.Delete(ctx, updatedClass.ID)
	if err != nil {
		t.Fatalf("Error in delete func: %v", err)
	}

}
