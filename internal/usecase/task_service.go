package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type taskService struct {
	taskRepo domain.TaskRepository
	userRepo domain.UserRepository
}

func NewTaskService(taskRepo domain.TaskRepository, userRepo domain.UserRepository) domain.TaskService {
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *taskService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateTaskRequest) (*domain.Task, error) {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, domain.ErrUnauthorized
	}

	task := &domain.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatusPending,
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, userID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	if task.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return task, nil
}

func (s *taskService) List(ctx context.Context, userID uuid.UUID) ([]domain.Task, error) {
	filter := domain.TaskFilter{UserID: userID}
	tasks, err := s.taskRepo.ListByUserID(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) Update(ctx context.Context, userID, taskID uuid.UUID, req domain.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	if task.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, userID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("get task by id: %w", err)
	}

	if task.UserID != userID {
		return domain.ErrForbidden
	}

	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
