package router_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcb1310/bookstore/internal/interfaces"
	"github.com/alcb1310/bookstore/internal/mocks"
	"github.com/alcb1310/bookstore/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestHealthRoute(t *testing.T) {
	db := mocks.NewService(t)
	s := router.New(8080, db)
	s.Router()
	assert.NotNil(t, s)

	testURL := "/health"

	testCases := []struct {
		name     string
		status   int
		response map[string]any
		check    *mocks.Service_HealthCheck_Call
	}{
		{
			name:   "should return ok",
			status: http.StatusOK,
			response: map[string]any{
				"status": "ok",
			},
			check: db.EXPECT().HealthCheck().Return(nil),
		},
		{
			name:   "database is not available",
			status: http.StatusGatewayTimeout,
			response: map[string]any{
				"error": "Database is not available",
			},
			check: db.EXPECT().HealthCheck().Return(&interfaces.APIError{
				Code:          http.StatusGatewayTimeout,
				Msg:           "Database is not available",
				OriginalError: fmt.Errorf("database is not available"),
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.check != nil {
				tc.check.Times(1)
			}

			req, err := http.NewRequest(http.MethodGet, testURL, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			responseBody := map[string]any{}
			err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.status, rec.Code)
			assert.Equal(t, tc.response, responseBody)
		})
	}
}
