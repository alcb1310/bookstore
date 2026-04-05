package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface{}

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
