package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/mohamed/go-clean-architecture/internal/repository/dbgen"
)

type userRepo struct {
	q *dbgen.Queries
}

func NewUserRepo(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepo{q: dbgen.New(pool)}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	params := dbgen.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	dbUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	user.ID = pgtypeUUIDToDomain(dbUser.ID)
	user.Name = dbUser.Name
	user.Email = dbUser.Email
	user.PasswordHash = dbUser.PasswordHash
	user.Role = dbUser.Role
	user.CreatedAt = timestamptzToDomain(dbUser.CreatedAt)
	user.UpdatedAt = timestamptzToDomain(dbUser.UpdatedAt)
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dbUser, err := r.q.GetUserByID(ctx, uuidToPgtype(id))
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &domain.User{
		ID:           pgtypeUUIDToDomain(dbUser.ID),
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    timestamptzToDomain(dbUser.CreatedAt),
		UpdatedAt:    timestamptzToDomain(dbUser.UpdatedAt),
	}, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &domain.User{
		ID:           pgtypeUUIDToDomain(dbUser.ID),
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    timestamptzToDomain(dbUser.CreatedAt),
		UpdatedAt:    timestamptzToDomain(dbUser.UpdatedAt),
	}, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	params := dbgen.UpdateUserParams{
		ID:           uuidToPgtype(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	dbUser, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	user.ID = pgtypeUUIDToDomain(dbUser.ID)
	user.Name = dbUser.Name
	user.Email = dbUser.Email
	user.PasswordHash = dbUser.PasswordHash
	user.Role = dbUser.Role
	user.CreatedAt = timestamptzToDomain(dbUser.CreatedAt)
	user.UpdatedAt = timestamptzToDomain(dbUser.UpdatedAt)
	return nil
}

func (r *userRepo) List(ctx context.Context) ([]domain.User, error) {
	dbUsers, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]domain.User, 0, len(dbUsers))
	for i := range dbUsers {
		users = append(users, domain.User{
			ID:           pgtypeUUIDToDomain(dbUsers[i].ID),
			Name:         dbUsers[i].Name,
			Email:        dbUsers[i].Email,
			PasswordHash: dbUsers[i].PasswordHash,
			Role:         dbUsers[i].Role,
			CreatedAt:    timestamptzToDomain(dbUsers[i].CreatedAt),
			UpdatedAt:    timestamptzToDomain(dbUsers[i].UpdatedAt),
		})
	}

	return users, nil
}
