package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/mohamed/go-clean-architecture/internal/repository/dbgen"
)

type taskRepo struct {
	q *dbgen.Queries
}

func NewTaskRepo(pool *pgxpool.Pool) domain.TaskRepository {
	return &taskRepo{q: dbgen.New(pool)}
}

func (r *taskRepo) Create(ctx context.Context, task *domain.Task) error {
	params := dbgen.CreateTaskParams{
		UserID:      uuidToPgtype(task.UserID),
		Title:       task.Title,
		Description: stringToPgtypeText(derefOrEmpty(task.Description)),
		Status:      string(task.Status),
		DueDate:     dueDateToPgtype(task.DueDate),
	}

	dbTask, err := r.q.CreateTask(ctx, params)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	*task = *r.toDomain(&dbTask)
	return nil
}

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	dbTask, err := r.q.GetTaskByID(ctx, uuidToPgtype(id))
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	return r.toDomain(&dbTask), nil
}

func (r *taskRepo) ListByUserID(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, error) {
	dbTasks, err := r.q.ListTasksByUserID(ctx, uuidToPgtype(filter.UserID))
	if err != nil {
		return nil, fmt.Errorf("list tasks by user id: %w", err)
	}

	tasks := make([]domain.Task, 0, len(dbTasks))
	for i := range dbTasks {
		tasks = append(tasks, *r.toDomain(&dbTasks[i]))
	}

	return tasks, nil
}

func (r *taskRepo) Update(ctx context.Context, task *domain.Task) error {
	params := dbgen.UpdateTaskParams{
		ID:          uuidToPgtype(task.ID),
		Title:       task.Title,
		Description: stringToPgtypeText(derefOrEmpty(task.Description)),
		Status:      string(task.Status),
		DueDate:     dueDateToPgtype(task.DueDate),
	}

	dbTask, err := r.q.UpdateTask(ctx, params)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	*task = *r.toDomain(&dbTask)
	return nil
}

func (r *taskRepo) ListAll(ctx context.Context) ([]domain.Task, error) {
	dbTasks, err := r.q.ListAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all tasks: %w", err)
	}

	tasks := make([]domain.Task, 0, len(dbTasks))
	for i := range dbTasks {
		tasks = append(tasks, *r.toDomain(&dbTasks[i]))
	}

	return tasks, nil
}

func (r *taskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteTask(ctx, uuidToPgtype(id))
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (r *taskRepo) toDomain(t *dbgen.Task) *domain.Task {
	return &domain.Task{
		ID:          pgtypeUUIDToDomain(t.ID),
		UserID:      pgtypeUUIDToDomain(t.UserID),
		Title:       t.Title,
		Description: pgtypeTextToStringPtr(t.Description),
		Status:      domain.TaskStatus(t.Status),
		DueDate:     pgtypeTimestamptzToTimePtr(t.DueDate),
		CreatedAt:   timestamptzToDomain(t.CreatedAt),
		UpdatedAt:   timestamptzToDomain(t.UpdatedAt),
	}
}

func derefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func dueDateToPgtype(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func pgtypeTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
