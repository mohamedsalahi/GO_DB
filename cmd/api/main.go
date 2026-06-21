package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohamed/go-clean-architecture/config"
	"github.com/mohamed/go-clean-architecture/internal/domain"
	"github.com/mohamed/go-clean-architecture/internal/delivery/handler"
	"github.com/mohamed/go-clean-architecture/internal/infra/auth"
	"github.com/mohamed/go-clean-architecture/internal/infra/db"
	"github.com/mohamed/go-clean-architecture/internal/infra/logger"
	"github.com/mohamed/go-clean-architecture/internal/repository"
	"github.com/mohamed/go-clean-architecture/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger.SetupLogger(cfg.App.Env, cfg.Log.Level)
	slog.Info("starting application",
		slog.String("name", cfg.App.Name),
		slog.String("version", cfg.App.Version),
		slog.String("env", cfg.App.Env),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := db.NewPostgresPool(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	rdb, err := db.NewRedisClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer rdb.Close()

	tokenMgr := auth.NewTokenManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiration,
		cfg.JWT.RefreshExpiration,
	)

	userRepo := repository.NewUserRepo(dbPool)
	taskRepo := repository.NewTaskRepo(dbPool)

	userService := usecase.NewUserService(userRepo, tokenMgr)
	taskService := usecase.NewTaskService(taskRepo, userRepo)
	adminService := usecase.NewAdminService(userRepo, taskRepo)

	// Seed the god user if not exists
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer seedCancel()
	if err := seedGodUser(seedCtx, userRepo, cfg); err != nil {
		slog.Warn("failed to seed god user", slog.String("error", err.Error()))
	}

	// Check for frontend dist directory (for production builds)
	frontendDir := "frontend/dist"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = ""
	}

	router := handler.NewRouter(&handler.Dependencies{
		UserService:    userService,
		TaskService:    taskService,
		AdminService:   adminService,
		TokenMgr:       tokenMgr,
		RedisClient:    rdb,
		FrontendDir:    frontendDir,
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}

func seedGodUser(ctx context.Context, userRepo domain.UserRepository, cfg *config.Config) error {
	existing, err := userRepo.GetByEmail(ctx, cfg.GodUser.Email)
	if err == nil && existing != nil {
		if existing.Role != domain.RoleAdmin {
			existing.Role = domain.RoleAdmin
			if err := userRepo.Update(ctx, existing); err != nil {
				return fmt.Errorf("update god user role: %w", err)
			}
			slog.Info("god user role updated to admin", slog.String("email", cfg.GodUser.Email))
		}
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.GodUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash god user password: %w", err)
	}

	user := &domain.User{
		Name:         cfg.GodUser.Name,
		Email:        cfg.GodUser.Email,
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleAdmin,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("create god user: %w", err)
	}

	slog.Info("god user seeded", slog.String("email", cfg.GodUser.Email))
	return nil
}
