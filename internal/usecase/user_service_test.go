package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/mohamed/go-clean-architecture/internal/infra/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupUserService() (*MockUserRepo, *auth.TokenManager, domain.UserService) {
	mockRepo := new(MockUserRepo)
	tokenMgr := auth.NewTokenManager("test-secret-key", 15*time.Minute, 7*24*time.Hour)
	svc := NewUserService(mockRepo, tokenMgr)
	return mockRepo, tokenMgr, svc
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo, _, svc := setupUserService()

	mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Name == "John Doe" && u.Email == "john@example.com"
	})).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*domain.User)
		u.ID = uuid.New()
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
	})

	resp, err := svc.Register(context.Background(), domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "John Doe", resp.User.Name)
	assert.Equal(t, "john@example.com", resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo, _, svc := setupUserService()

	existingUser := &domain.User{ID: uuid.New(), Email: "john@example.com"}
	mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(existingUser, nil)

	resp, err := svc.Register(context.Background(), domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword",
	})

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyInUse)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo, _, svc := setupUserService()
	userID := uuid.New()
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(&domain.User{
		ID:           userID,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: string(hashedPW),
	}, nil)

	resp, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "john@example.com",
		Password: "correctpassword",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, userID, resp.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	mockRepo, _, svc := setupUserService()
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(&domain.User{
		ID:           uuid.New(),
		Email:        "john@example.com",
		PasswordHash: string(hashedPW),
	}, nil)

	resp, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCreds)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo, _, svc := setupUserService()

	mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, domain.ErrNotFound)

	resp, err := svc.Login(context.Background(), domain.LoginRequest{
		Email:    "notfound@example.com",
		Password: "anypassword",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCreds)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_Success(t *testing.T) {
	mockRepo, _, svc := setupUserService()
	userID := uuid.New()

	mockRepo.On("GetByID", mock.Anything, userID).Return(&domain.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil)

	user, err := svc.GetProfile(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_NotFound(t *testing.T) {
	mockRepo, _, svc := setupUserService()
	userID := uuid.New()

	mockRepo.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrNotFound)

	user, err := svc.GetProfile(context.Background(), userID)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
