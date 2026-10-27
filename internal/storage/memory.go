package storage

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/Shazzadul-Shakib/go-tasks-api/internal/models"
)

// MemoryStore keeps tasks in a map. It is safe for concurrent use.
type MemoryStore struct {
	mu     sync.RWMutex
	tasks  map[int64]models.Task
	nextID int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: make(map[int64]models.Task)}
}

func (s *MemoryStore) List(_ context.Context, f models.ListFilter) ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Task, 0, len(s.tasks)) // non-nil: encodes as [] not null
	for _, t := range s.tasks {
		if f.Done != nil && t.Done != *f.Done {
			continue
		}
		out = append(out, t)
	}
	// Map order is random; sort so clients get a stable order.
	slices.SortFunc(out, func(a, b models.Task) int { return cmp.Compare(a.ID, b.ID) })
	return out, nil
}

func (s *MemoryStore) Get(_ context.Context, id int64) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("get task %d: %w", id, models.ErrNotFound)
	}
	return t, nil // a copy: callers can't modify our map entry
}

func (s *MemoryStore) Create(_ context.Context, t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	now := time.Now().UTC()
	t.ID = s.nextID
	t.CreatedAt, t.UpdatedAt = now, now
	s.tasks[t.ID] = t
	return t, nil
}

func (s *MemoryStore) Update(_ context.Context, t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.tasks[t.ID]
	if !ok {
		return models.Task{}, fmt.Errorf("update task %d: %w", t.ID, models.ErrNotFound)
	}
	t.CreatedAt = existing.CreatedAt // clients can't change these
	t.UpdatedAt = time.Now().UTC()
	s.tasks[t.ID] = t
	return t, nil
}

func (s *MemoryStore) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("delete task %d: %w", id, models.ErrNotFound)
	}
	delete(s.tasks, id)
	return nil
}
