package router

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alcb1310/bookstore/internal/interfaces"
)

type ErrorResponse func(w http.ResponseWriter, r *http.Request) error

func HandleErrors(h ErrorResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if e, ok := errors.AsType[*interfaces.APIError](err); ok {
				slog.Error("API error", "code", e.Code, "error", e.OriginalError)
				_ = JSONResponse(w, e.Code, map[string]any{"error": e.Msg})
				return

			}

			slog.Error("Internal server error", "error", err)
			_ = JSONResponse(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
	}
}

func HomeRoute(w http.ResponseWriter, r *http.Request) error {
	return JSONResponse(w, http.StatusOK, map[string]any{"data": "Hello world"})
}

func JSONResponse(w http.ResponseWriter, code int, data map[string]any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
