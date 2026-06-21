package handler

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mohamed/go-clean-architecture/internal/delivery/handler/middleware"
	"github.com/mohamed/go-clean-architecture/internal/infra/auth"
	"github.com/redis/go-redis/v9"

	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type Dependencies struct {
	UserService    domain.UserService
	TaskService    domain.TaskService
	AdminService   domain.AdminService
	TokenMgr       *auth.TokenManager
	RedisClient    *redis.Client
	FrontendDir    string // path to built frontend dist, empty = no frontend
}

func NewRouter(deps *Dependencies) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(middleware.RequestLogger)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiter (100 req/min per IP)
	rateLimiter := middleware.NewRateLimiter(deps.RedisClient, 100, 1*time.Minute)
	r.Use(rateLimiter.Middleware)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes (public)
	r.Route("/api/v1/auth", func(r chi.Router) {
		userHandler := NewUserHandler(deps.UserService)
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(deps.TokenMgr))

		// User profile
		r.Route("/api/v1/users", func(r chi.Router) {
			userHandler := NewUserHandler(deps.UserService)
			r.Get("/me", userHandler.GetProfile)
		})

		// Task CRUD
		r.Route("/api/v1/tasks", func(r chi.Router) {
			taskHandler := NewTaskHandler(deps.TaskService)
			r.Get("/", taskHandler.List)
			r.Post("/", taskHandler.Create)
			r.Get("/{taskID}", taskHandler.GetByID)
			r.Put("/{taskID}", taskHandler.Update)
			r.Delete("/{taskID}", taskHandler.Delete)
		})

		// Admin routes (God's Eye)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminMiddleware)
			adminHandler := NewAdminHandler(deps.AdminService)
			r.Route("/api/v1/admin", func(r chi.Router) {
				r.Get("/users", adminHandler.ListUsers)
				r.Get("/tasks", adminHandler.ListAllTasks)
				r.Get("/users/{userID}/tasks", adminHandler.ListUserTasks)
				r.Put("/users/{userID}/promote", adminHandler.PromoteUser)
			})
		})
	})

	// Serve frontend SPA (if frontend directory is set)
	if deps.FrontendDir != "" {
		fileServer := http.FileServer(http.Dir(deps.FrontendDir))

		// Handle all non-API routes for SPA
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// Check if file exists
			path := deps.FrontendDir + r.URL.Path
			if _, err := os.Stat(path); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
			// Fallback to index.html for SPA routing
			http.ServeFile(w, r, deps.FrontendDir+"/index.html")
		})

		slog.Info("frontend SPA configured", slog.String("dir", deps.FrontendDir))
	}

	slog.Info("router configured with all routes and middleware")
	return r
}
