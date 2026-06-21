package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mohamed/go-clean-architecture/internal/delivery/handler/middleware"
)

func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func getUserID(r *http.Request) *uuid.UUID {
	id := middleware.GetUserID(r.Context())
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, param))
}
