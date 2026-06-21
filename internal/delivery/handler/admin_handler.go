package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type adminHandler struct {
	adminSvc domain.AdminService
}

func NewAdminHandler(adminSvc domain.AdminService) *adminHandler {
	return &adminHandler{adminSvc: adminSvc}
}

func (h *adminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.adminSvc.ListUsers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *adminHandler) ListAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.adminSvc.ListAllTasks(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

type listUserTasksResponse struct {
	UserID uuid.UUID    `json:"user_id"`
	Tasks  []domain.Task `json:"tasks"`
}

func (h *adminHandler) ListUserTasks(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	tasks, err := h.adminSvc.ListUserTasks(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list user tasks")
		return
	}

	respondJSON(w, http.StatusOK, listUserTasksResponse{
		UserID: userID,
		Tasks:  tasks,
	})
}

func (h *adminHandler) PromoteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.adminSvc.PromoteToAdmin(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to promote user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
