package repository

import (
	"context"
	"sync"

	"github.com/KoIvRo/backend/internal/domain"
)

type MemoryClassRepository struct {
	mu      sync.RWMutex
	classes map[int]domain.Class
	nextID  int
}

func NewMemoryClassRepository() *MemoryClassRepository {
	return &MemoryClassRepository{classes: make(map[int]domain.Class), nextID: 1}
}

func (r *MemoryClassRepository) Create(ctx context.Context, class domain.Class) (domain.Class, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	class.ID = r.nextID
	r.nextID++

	r.classes[class.ID] = class

	return class, nil
}

func (r *MemoryClassRepository) List(ctx context.Context) ([]domain.Class, error) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	list := make([]domain.Class, 0, len(r.classes))

	for _, class := range r.classes {
		list = append(list, class)
	}

	return list, nil

}

func (r *MemoryClassRepository) GetByID(ctx context.Context, id int) (domain.Class, error) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	class, ok := r.classes[id]

	if !ok {
		return domain.Class{}, ErrNotFound
	}

	return class, nil
}

func (r *MemoryClassRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, ok := r.classes[id]; ok {
		delete(r.classes, id)
		return nil
	} else {
		return ErrNotFound
	}
}

func (r *MemoryClassRepository) Update(ctx context.Context, id int, class domain.Class) (domain.Class, error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, ok := r.classes[id]; !ok {
		return domain.Class{}, ErrNotFound
	}

	class.ID = id
	r.classes[id] = class

	return r.classes[id], nil
}
