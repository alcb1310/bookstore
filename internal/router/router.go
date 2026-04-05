package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type service struct {
	port uint16
}

func New(port uint16) *service {
	return &service{
		port: port,
	}
}

func (s *service) Router() error {
	r := chi.NewRouter()

	r.Get("/", HandleErrors(HomeRoute))

	port := fmt.Sprintf(":%d", s.port)
	slog.Info("Starting server", "port", port)
	return http.ListenAndServe(port, r)
}
