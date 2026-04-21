package router

import (
	"time"

	"github.com/alcb1310/bookstore/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

type Router struct {
	port uint16
	db   database.Service
}

func New(port uint16, db database.Service) *Router {
	return &Router{
		port: port,
		db:   db,
	}
}

func (s *Router) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Get("/", HandleErrors(HomeRoute))
	r.Get("/health", HandleErrors(s.HealthRoute))

	return r
}
