package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superboomer/map-tile-provider/app/provider"
	"github.com/superboomer/map-tile-provider/app/server/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// Mocked dependencies
type (
	MockedProviders struct {
		mock.Mock
	}
)

func (m *MockedProviders) GetAllKey() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockedProviders) Get(key string) (provider.Provider, error) {
	args := m.Called(key)
	return args.Get(0).(provider.Provider), args.Error(1)
}

// TestProviderHandler tests the Provider handler
func TestProviderHandler(t *testing.T) {
	// Setup mocked logger and observer for logging
	observedLogs, _ := observer.New(zap.InfoLevel)
	logger := zap.New(observedLogs)

	// Create a mock for Providers
	providersMock := new(MockedProviders)
	providersMock.On("GetAllKey").Return([]string{"google", "arcgis"})
	providersMock.On("Get", "google").Return(provider.Google(), nil)
	providersMock.On("Get", "arcgis").Return(provider.ArcGIS(), nil)

	// Initialize the API struct with mocked dependencies
	apiPkg := &api.API{
		Logger:    logger,
		Providers: providersMock,
	}

	// Setup request and response recorder
	req, err := http.NewRequest("GET", "/provider", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(apiPkg.Provider)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code, "Handler did not return expected status code")

	// Check the response body
	expectedBody := `[{"name":"Google Maps (Satellite)","key":"google","max_zoom":21},{"name":"ArcGIS (Satellite)","key":"arcgis","max_zoom":19}]`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}
