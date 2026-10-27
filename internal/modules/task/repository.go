package task

import (
	"context"

	"github.com/Shazzadul-Shakib/go-tasks-api/internal/models"
)

// Repository is the storage the task module needs. storage.MemoryStore satisfies it.
type Repository interface {
	List(ctx context.Context, f models.ListFilter) ([]models.Task, error)
}
