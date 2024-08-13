package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/superboomer/maptile/app/provider"
)

// mapErrorModel contains data about error query
type mapErrorModel struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// Map godoc
// @Summary handler for generating satellite map for specified lat long and from specified vendor
// @Description return merged satellite tiles in one image
// @Accept  text/plain
// @Produce image/jpeg
// @Param provider query string true "tile provider"
// @Param lat query		 number	 true "latitude"
// @Param long query		 number	 true "longitude"
// @Param zoom query		 int true "zoom of image"
// @Param side query		 int false "count of tile of result image square" default(3) minimum(1)		maximum(10)
// @Success 200 {file} image/jpeg
// @Failure 400 {object} mapErrorModel
// @Header 200 {string} X-Request-Id "request_id"
// @Router /map [get]
func (a *API) Map(w http.ResponseWriter, req *http.Request) {

	params, vendor, err := a.parseRequest(req)
	if err != nil {
		results, _ := json.Marshal(mapErrorModel{
			Status: http.StatusBadRequest,
			Body:   err.Error(),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(results)
		return
	}

	var centerTile = vendor.GetTile(params.Latitude, params.Longitude, params.Zoom)

	tiles, err := a.Downloader.Download(a.Cache, vendor, centerTile.GetNearby(params.Side)...)
	if err != nil {
		a.Logger.Error("error occurred when downloading tiles", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))

		results, _ := json.Marshal(mapErrorModel{
			Status: http.StatusInternalServerError,
			Body:   fmt.Sprintf("error occurred when dowloading tiles: %s", err.Error()),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(results)
		return
	}

	merged, err := a.Downloader.Merge(params.Side, centerTile, tiles...)
	if err != nil {
		a.Logger.Error("error occurred when merging tiles", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))

		results, _ := json.Marshal(mapErrorModel{
			Status: http.StatusInternalServerError,
			Body:   fmt.Sprintf("error occurred when merging tiles: %s", err.Error()),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(results)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	_, _ = w.Write(merged)

	a.Logger.Info("new map download request", zap.Float64("lat", params.Latitude), zap.Float64("long", params.Longitude), zap.Int("side", params.Side), zap.String("vendor", vendor.Name()), zap.String("req_id", req.Header.Get("X-Request-ID")))
}

type mapParams struct {
	Latitude  float64
	Longitude float64
	Zoom      float64
	Side      int
}

func (a *API) validateZoom(zoom float64, maxZoom int, vendorName string) error {
	if zoom < 1 || zoom > float64(maxZoom) {
		return fmt.Errorf("max zoom for provider %s - %d", vendorName, maxZoom)
	}
	return nil
}

func (a *API) parseRequest(req *http.Request) (*mapParams, provider.Provider, error) {

	var params = mapParams{
		Side: 3,
	}

	pVendor := req.URL.Query().Get("provider")
	if pVendor == "" {
		return nil, nil, fmt.Errorf("provider parameter error: not specified")
	}

	vendor, err := a.Providers.Get(pVendor)
	if err != nil {
		return nil, nil, fmt.Errorf("provider parameter error: %s not found", pVendor)
	}

	pLat := req.URL.Query().Get("lat")
	params.Latitude, err = parseFloatParam(pLat)
	if err != nil {
		return nil, nil, fmt.Errorf("lat parameter error: %w", err)
	}

	pLong := req.URL.Query().Get("long")
	params.Longitude, err = parseFloatParam(pLong)
	if err != nil {
		return nil, nil, fmt.Errorf("long parameter error: %w", err)
	}

	pZoom := req.URL.Query().Get("zoom")
	params.Zoom, err = parseFloatParam(pZoom)
	if err != nil {
		return nil, nil, fmt.Errorf("zoom parameter error: %w", err)
	}

	if zoomErr := a.validateZoom(params.Zoom, vendor.MaxZoom(), vendor.Name()); zoomErr != nil {
		return nil, nil, fmt.Errorf("zoom parameter error: %w", zoomErr)
	}

	pSide := req.URL.Query().Get("side")
	if pSide != "" {
		sideInt, err := strconv.Atoi(pSide)
		if err != nil {
			return nil, nil, fmt.Errorf("side parameter error: %w", err)
		}

		if sideInt < 1 || sideInt > a.MaxSide {
			return nil, nil, fmt.Errorf("side parameter error: must be greater or equal to 1 and less than %d", a.MaxSide)
		}
		params.Side = sideInt
	}

	return &params, vendor, nil
}

func parseFloatParam(param string) (float64, error) {
	valueStr := strings.TrimSpace(param)

	if valueStr == "" {
		return 0, fmt.Errorf("not specified")
	}

	valueFloat, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}

	return valueFloat, nil
}
