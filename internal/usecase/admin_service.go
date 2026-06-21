package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type adminService struct {
	userRepo domain.UserRepository
	taskRepo domain.TaskRepository
}

func NewAdminService(userRepo domain.UserRepository, taskRepo domain.TaskRepository) domain.AdminService {
	return &adminService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

func (s *adminService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.List(ctx)
}

func (s *adminService) ListAllTasks(ctx context.Context) ([]domain.Task, error) {
	return s.taskRepo.ListAll(ctx)
}

func (s *adminService) ListUserTasks(ctx context.Context, userID uuid.UUID) ([]domain.Task, error) {
	return s.taskRepo.ListByUserID(ctx, domain.TaskFilter{UserID: userID})
}

func (s *adminService) PromoteToAdmin(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Role = domain.RoleAdmin
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
