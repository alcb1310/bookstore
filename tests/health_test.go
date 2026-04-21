package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHealthEndPoint(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("bookstore"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)

	assert.NotNil(t, pgContainer)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if pgContainer != nil {
			err = pgContainer.Terminate(ctx)
			assert.NoError(t, err)
		}
	})

	testURL := "/health"
	_ = testURL

	s, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	t.Run("Integration - should return ok", func(t *testing.T) {
		expected := map[string]any{"status": "ok"}
		req, err := http.NewRequest("GET", testURL, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		responseBody := map[string]any{}
		err = json.Unmarshal(res.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("Integration - database looses connection", func(t *testing.T) {
		err := pgContainer.Terminate(ctx)
		assert.NoError(t, err)
		pgContainer = nil
		expected := map[string]any{"error": "Database is not available"}
		req, err := http.NewRequest("GET", testURL, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusServiceUnavailable, res.Code)
		responseBody := map[string]any{}
		err = json.Unmarshal(res.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expected, responseBody)
	})
}
