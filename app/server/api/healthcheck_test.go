package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthCheckHandler(t *testing.T) {
	// Initialize the API with the mock logger
	apiPkg := &API{
		Logger: zap.NewNop(),
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

	var resp healthCheckModel
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	assert.Equal(t, 200, resp.Status)
	assert.Equal(t, "OK", resp.Body)
}
