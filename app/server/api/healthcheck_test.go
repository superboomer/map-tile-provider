package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/server/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHealthCheckHandler(t *testing.T) {
	// Setup mocked logger and observer for logging
	observedLogs, _ := observer.New(zap.InfoLevel)
	logger := zap.New(observedLogs)

	// Initialize the API with the mock logger
	apiPkg := &api.API{
		Logger: logger,
	}

	req, err := http.NewRequest("GET", "/healthcheck", http.NoBody)
	if err != nil {
		t.Fatalf("An error occurred while creating the request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(apiPkg.HealthCheck)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp api.HealthCheckModel
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	assert.Equal(t, 200, resp.Status)
	assert.Equal(t, "OK", resp.Body)
}
