package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTaskRequest is the DTO for creating a new task
type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// UpdateTaskRequest is the DTO for updating an existing task
type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      *TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
	DueDate     *time.Time  `json:"due_date,omitempty"`
}

// TaskFilter defines filtering/pagination for task listing
type TaskFilter struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type (
	TaskRepository interface {
		Create(ctx context.Context, task *Task) error
		GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
		ListByUserID(ctx context.Context, filter TaskFilter) ([]Task, error)
		Update(ctx context.Context, task *Task) error
		Delete(ctx context.Context, id uuid.UUID) error
		ListAll(ctx context.Context) ([]Task, error)
	}

	TaskService interface {
		Create(ctx context.Context, userID uuid.UUID, req CreateTaskRequest) (*Task, error)
		GetByID(ctx context.Context, userID, taskID uuid.UUID) (*Task, error)
		List(ctx context.Context, userID uuid.UUID) ([]Task, error)
		Update(ctx context.Context, userID, taskID uuid.UUID, req UpdateTaskRequest) (*Task, error)
		Delete(ctx context.Context, userID, taskID uuid.UUID) error
	}
)
