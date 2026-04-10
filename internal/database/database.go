package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alcb1310/bookstore/internal/interfaces"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface {
	HealthCheck() error
}

type service struct {
	DB *sql.DB
}

func New() (Service, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return &service{
		DB: conn,
	}, nil
}

func (s *service) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.DB.PingContext(ctx); err != nil {
		if e, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch e.Code {
			case "3D000":
				return &interfaces.APIError{
					Code:          http.StatusGatewayTimeout,
					Msg:           "Database is not available",
					OriginalError: e,
				}
			case "28000":
				return &interfaces.APIError{
					Code:          http.StatusBadGateway,
					Msg:           "Database is not available",
					OriginalError: e,
				}
			default:
				return &interfaces.APIError{
					Code:          http.StatusInternalServerError,
					Msg:           "Database is not available",
					OriginalError: e,
				}
			}
		}
		if _, ok := errors.AsType[*pgconn.ConnectError](err); ok {
			return &interfaces.APIError{
				Code:          http.StatusServiceUnavailable,
				Msg:           "Database is not available",
				OriginalError: err,
			}
		}

		return err
	}

	return nil
}
