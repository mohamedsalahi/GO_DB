package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mohamed/go-clean-architecture/internal/domain"
)

type UserHandler struct {
	userService domain.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.userService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyInUse) {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.userService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCreds) {
			respondError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userService.GetProfile(r.Context(), *userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
