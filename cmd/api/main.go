package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/alcb1310/bookstore/internal/database"
	"github.com/alcb1310/bookstore/internal/router"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port64, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		slog.Error("Error parsing port", "error", err)
		panic(err)
	}
	port := uint16(port64)

	db, err := database.New()
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		panic(err)
	}

	s := router.New(port, db)
	if err := s.Router(); err != nil {
		panic(err)
	}
}
