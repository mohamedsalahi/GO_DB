package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type TaskHandler struct {
	taskService domain.TaskService
	validate    *validator.Validate
}

func NewTaskHandler(taskService domain.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		validate:    validator.New(),
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.taskService.Create(r.Context(), *userID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.taskService.GetByID(r.Context(), *userID, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			respondError(w, http.StatusForbidden, "access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tasks, err := h.taskService.List(r.Context(), *userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req domain.UpdateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.taskService.Update(r.Context(), *userID, taskID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			respondError(w, http.StatusForbidden, "access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.taskService.Delete(r.Context(), *userID, taskID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			respondError(w, http.StatusForbidden, "access denied")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
