package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/mohamed/go-clean-architecture/internal/infra/auth"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepository
	tokenMgr *auth.TokenManager
}

func NewUserService(userRepo domain.UserRepository, tokenMgr *auth.TokenManager) domain.UserService {
	return &userService{
		userRepo: userRepo,
		tokenMgr: tokenMgr,
	}
}

func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenMgr.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCreds
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCreds
	}

	accessToken, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenMgr.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
