package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/maptile/app/cache"
	"github.com/superboomer/maptile/app/downloader"
	"github.com/superboomer/maptile/app/provider"
	"github.com/superboomer/maptile/app/tile"

	"go.uber.org/zap"
)

// Initialize the API with the mock logger
var apiPkg = &API{
	Logger: zap.NewNop(),
	Providers: &provider.ListMock{
		GetFunc: func(key string) (provider.Provider, error) {
			if key == "example2" {
				return nil, fmt.Errorf("not found")
			}
			return &provider.ProviderMock{
				MaxZoomFunc:    func() int { return 2 },
				NameFunc:       func() string { return "example" },
				IDFunc:         func() string { return "ex" },
				GetTileFunc:    func(lat, long, scale float64) tile.Tile { return tile.Tile{X: 0, Y: 0, Z: 0} },
				MaxJobsFunc:    func() int { return 1 },
				GetRequestFunc: func(t *tile.Tile) *http.Request { return &http.Request{} },
			}, nil
		},
	},
	MaxSide: 10,
	Downloader: &downloader.DownloaderMock{
		DownloadFunc: func(c cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error) {
			return []tile.Tile{}, nil
		},
		MergeFunc: func(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error) { return []byte{}, nil },
	},
}

func TestMapHandler_ValidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=40.7128&long=74.0060&zoom=1&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"))
}

func TestMapHandler_MissingRequiredParameterProvider(t *testing.T) {
	// Omitting one or more required parameters
	req, err := http.NewRequest("GET", "/map?lat=40.7128&long=-74.0060&zoom=1&side=3", http.NoBody) // Missing vendor parameter
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedBody := `{"status":400,"body":"provider parameter error: not specified"}`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}

func TestMapHandler_MissingRequiredParameterZoom(t *testing.T) {
	// Omitting one or more required parameters
	req, err := http.NewRequest("GET", "/map?provider=example&lat=40.7128&long=-74.0060&side=3", http.NoBody) // Missing vendor parameter
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "not specified")
}

func TestMapHandler_InvalidParameterMaxZoom(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=-74.0060&zoom=15&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedBody := `{"status":400,"body":"zoom parameter error: max zoom for provider example - 2"}`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}

func TestMapHandler_InvalidParameterZoom(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=-74.0060&zoom=invalid&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "zoom parameter error")
}

func TestMapHandler_InvalidParameterLatitude(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=invalid&long=-74.0060&zoom=15&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "lat parameter error")
}

func TestMapHandler_InvalidParameterLongtitude(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=invalid&zoom=15&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "long parameter error")
}

func TestMapHandler_InvalidParameterProvider(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example2&lat=12.0&long=-74.0060&zoom=1&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedBody := `{"status":400,"body":"provider parameter error: example2 not found"}`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}

func TestMapHandler_InvalidParameterSide(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=-74.0060&zoom=1&side=invalid", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "side parameter error")
}

func TestMapHandler_InvalidParameterSideMax(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=-74.0060&zoom=1&side=100", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "side parameter error: must be greater or equal")
}

func TestMapHandler_InvalidParameterSideMin(t *testing.T) {
	req, err := http.NewRequest("GET", "/map?provider=example&lat=12.0&long=-74.0060&zoom=1&side=0", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	assert.Contains(t, rr.Body.String(), "side parameter error: must be greater or equal")
}

func TestMapHandler_ErrorDuringDownload(t *testing.T) {

	// Initialize the API with the mock logger and error DownloadFunc
	var apiPkg = &API{
		Logger: zap.NewNop(),
		Providers: &provider.ListMock{
			GetFunc: func(key string) (provider.Provider, error) {
				return &provider.ProviderMock{
					MaxZoomFunc:    func() int { return 2 },
					NameFunc:       func() string { return "example" },
					IDFunc:         func() string { return "ex" },
					GetTileFunc:    func(lat, long, scale float64) tile.Tile { return tile.Tile{X: 0, Y: 0, Z: 0} },
					MaxJobsFunc:    func() int { return 1 },
					GetRequestFunc: func(t *tile.Tile) *http.Request { return &http.Request{} },
				}, nil
			},
		},
		MaxSide: 10,
		Downloader: &downloader.DownloaderMock{
			DownloadFunc: func(c cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error) {
				return nil, fmt.Errorf("mock error")
			},
			MergeFunc: func(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error) { return []byte{}, nil },
		},
	}

	req, err := http.NewRequest("GET", "/map?provider=example&lat=40.7128&long=-74.0060&zoom=1&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedBody := `{"status":500,"body":"error occurred when dowloading tiles: mock error"}`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}

func TestMapHandler_ErrorDuringMerging(t *testing.T) {

	// Initialize the API with the mock logger and error DownloadFunc
	var apiPkg = &API{
		Logger: zap.NewNop(),
		Providers: &provider.ListMock{
			GetFunc: func(key string) (provider.Provider, error) {
				return &provider.ProviderMock{
					MaxZoomFunc:    func() int { return 2 },
					NameFunc:       func() string { return "example" },
					IDFunc:         func() string { return "ex" },
					GetTileFunc:    func(lat, long, scale float64) tile.Tile { return tile.Tile{X: 0, Y: 0, Z: 0} },
					MaxJobsFunc:    func() int { return 1 },
					GetRequestFunc: func(t *tile.Tile) *http.Request { return &http.Request{} },
				}, nil
			},
		},
		MaxSide: 10,
		Downloader: &downloader.DownloaderMock{
			DownloadFunc: func(c cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error) {
				return []tile.Tile{}, nil
			},
			MergeFunc: func(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error) {
				return []byte{}, fmt.Errorf("mock error")
			},
		},
	}

	req, err := http.NewRequest("GET", "/map?provider=example&lat=40.7128&long=-74.0060&zoom=1&side=3", http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	apiPkg.Map(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedBody := `{"status":500,"body":"error occurred when merging tiles: mock error"}`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Response body did not match expected JSON")
}
