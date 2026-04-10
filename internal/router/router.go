package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alcb1310/bookstore/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

type service struct {
	port uint16
	db   database.Service
}

func New(port uint16, db database.Service) *service {
	return &service{
		port: port,
		db:   db,
	}
}

func (s *service) Router() error {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, 1*time.Second))

	r.Get("/", HandleErrors(HomeRoute))
	r.Get("/health", HandleErrors(s.HealthRoute))

	port := fmt.Sprintf(":%d", s.port)
	slog.Info("Starting server", "port", port)
	return http.ListenAndServe(port, r)
}
