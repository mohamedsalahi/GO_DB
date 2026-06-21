package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTaskService() (*MockUserRepo, *MockTaskRepo, domain.TaskService) {
	mockUserRepo := new(MockUserRepo)
	mockTaskRepo := new(MockTaskRepo)
	svc := NewTaskService(mockTaskRepo, mockUserRepo)
	return mockUserRepo, mockTaskRepo, svc
}

func TestTaskService_Create_Success(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)
	mockTaskRepo.On("Create", mock.Anything, mock.MatchedBy(func(t *domain.Task) bool {
		return t.Title == "My Task" && t.UserID == userID && t.Status == domain.TaskStatusPending
	})).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(1).(*domain.Task)
		task.ID = uuid.New()
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
	})

	dueDate := time.Now().Add(24 * time.Hour)
	task, err := svc.Create(context.Background(), userID, domain.CreateTaskRequest{
		Title:   "My Task",
		DueDate: &dueDate,
	})

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "My Task", task.Title)
	assert.Equal(t, domain.TaskStatusPending, task.Status)
	assert.NotEqual(t, uuid.Nil, task.ID)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Create_Unauthorized(t *testing.T) {
	mockUserRepo, _, svc := setupTaskService()
	userID := uuid.New()

	mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)

	task, err := svc.Create(context.Background(), userID, domain.CreateTaskRequest{
		Title: "My Task",
	})

	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	assert.Nil(t, task)
	mockUserRepo.AssertExpectations(t)
}

func TestTaskService_GetByID_Success(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: userID,
		Title:  "My Task",
		Status: domain.TaskStatusPending,
	}, nil)

	task, err := svc.GetByID(context.Background(), userID, taskID)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, taskID, task.ID)
	assert.Equal(t, userID, task.UserID)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_GetByID_Forbidden(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()
	otherUserID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: otherUserID,
		Title:  "My Task",
	}, nil)

	task, err := svc.GetByID(context.Background(), userID, taskID)

	assert.ErrorIs(t, err, domain.ErrForbidden)
	assert.Nil(t, task)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_List_Success(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()

	expectedTasks := []domain.Task{
		{ID: uuid.New(), UserID: userID, Title: "Task 1", Status: domain.TaskStatusPending},
		{ID: uuid.New(), UserID: userID, Title: "Task 2", Status: domain.TaskStatusCompleted},
	}

	mockTaskRepo.On("ListByUserID", mock.Anything, mock.MatchedBy(func(f domain.TaskFilter) bool {
		return f.UserID == userID
	})).Return(expectedTasks, nil)

	tasks, err := svc.List(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Title)
	assert.Equal(t, "Task 2", tasks[1].Title)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Update_Success(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()
	newStatus := domain.TaskStatusCompleted

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: userID,
		Title:  "Old Title",
		Status: domain.TaskStatusPending,
	}, nil)
	mockTaskRepo.On("Update", mock.Anything, mock.MatchedBy(func(t *domain.Task) bool {
		return t.ID == taskID && t.Title == "New Title" && t.Status == domain.TaskStatusCompleted
	})).Return(nil)

	task, err := svc.Update(context.Background(), userID, taskID, domain.UpdateTaskRequest{
		Title:  strPtr("New Title"),
		Status: &newStatus,
	})

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "New Title", task.Title)
	assert.Equal(t, domain.TaskStatusCompleted, task.Status)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Update_Forbidden(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: uuid.New(),
		Title:  "Task",
	}, nil)

	task, err := svc.Update(context.Background(), userID, taskID, domain.UpdateTaskRequest{
		Title: strPtr("New Title"),
	})

	assert.ErrorIs(t, err, domain.ErrForbidden)
	assert.Nil(t, task)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Delete_Success(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: userID,
		Title:  "Task",
	}, nil)
	mockTaskRepo.On("Delete", mock.Anything, taskID).Return(nil)

	err := svc.Delete(context.Background(), userID, taskID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Delete_Forbidden(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(&domain.Task{
		ID:     taskID,
		UserID: uuid.New(),
		Title:  "Task",
	}, nil)

	err := svc.Delete(context.Background(), userID, taskID)

	assert.ErrorIs(t, err, domain.ErrForbidden)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_Delete_NotFound(t *testing.T) {
	mockUserRepo, mockTaskRepo, svc := setupTaskService()
	userID := uuid.New()
	taskID := uuid.New()

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, domain.ErrNotFound)

	err := svc.Delete(context.Background(), userID, taskID)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockUserRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}
