package task

import (
	"context"

	"github.com/Shazzadul-Shakib/go-tasks-api/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, f models.ListFilter) ([]models.Task, error) {
	return s.repo.List(ctx, f)
}
