package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/maptile/app/provider"
	"github.com/superboomer/maptile/app/server/api"
	"go.uber.org/zap"
)

func TestProviderHandler_Success(t *testing.T) {

	list := &provider.ListMock{
		GetAllIDFunc: func() []string {
			return []string{"providerA", "providerB"}
		},
		GetFunc: func(key string) (provider.Provider, error) {
			switch key {
			case "providerA":
				return &provider.ProviderMock{
					NameFunc: func() string {
						return "providerA"
					},
					IDFunc: func() string {
						return "a"
					},
					MaxZoomFunc: func() int {
						return 2
					}}, nil
			case "providerB":
				return &provider.ProviderMock{
					NameFunc: func() string {
						return "providerB"
					},
					IDFunc: func() string {
						return "b"
					},
					MaxZoomFunc: func() int {
						return 3
					}}, nil
			default:
				return nil, fmt.Errorf("not found")
			}
		},
	}

	// Initialize the API struct with mocked dependencies
	apiPkg := &api.API{
		Logger:    zap.NewNop(),
		Providers: list,
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
	expectedBody := `[{"name":"providerA","key":"a","max_zoom":2},{"name":"providerB","key":"b","max_zoom":3}]`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}

func TestProviderHandler_FailedUnregistered(t *testing.T) {

	list := &provider.ListMock{
		GetAllIDFunc: func() []string {
			return []string{"providerA", "providerB"}
		},
		GetFunc: func(key string) (provider.Provider, error) {
			switch key {
			case "providerA":
				return &provider.ProviderMock{
					NameFunc: func() string {
						return "providerA"
					},
					IDFunc: func() string {
						return "a"
					},
					MaxZoomFunc: func() int {
						return 2
					}}, nil
			default:
				return nil, fmt.Errorf("not found")
			}
		},
	}

	// Initialize the API struct with mocked dependencies
	apiPkg := &api.API{
		Logger:    zap.NewNop(),
		Providers: list,
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
	expectedBody := `[{"name":"providerA","key":"a","max_zoom":2}]`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}
