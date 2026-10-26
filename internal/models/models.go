package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("Task not found")

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required, max=100"`
	Description string `json:"description" binding:"required, max=2000"`
	Priority    int    `json:"priority" binding:"omitempty, min=1, max=5"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty, max=100"`
	Description string `json:"description" binding:"omitempty, max=2000"`
	Done        *bool  `json:"done" binding:"required"`
	Priority    int    `json:"priority" binding:"omitempty, min=1, max=5"`
}

type ListFilter struct {
	Done *bool
}
