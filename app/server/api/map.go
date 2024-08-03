package api

import (
	"fmt"
	"net/http"

	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/superboomer/map-tile-provider/app/downloader"
)

// Map godoc
// @Summary handler for generating satellite map for specified lat long and from specified vendor
// @Description return merged satellite tiles in one image
// @Accept  text/plain
// @Produce image/jpeg
// @Param vendor query string true "Map vendor"
// @Param lat query		 number	 true "latitude"
// @Param long query		 number	 true "longitude"
// @Param zoom query		 int true "zoom of image"
// @Param side query		 int false "count of tile of result image square" default(3) minimum(1)		maximum(10)
// @Success 200 {file} image/jpeg
// @Header 200 {string} X-Request-Id "request_id"
// @Router /map [get]
func (a *API) Map(w http.ResponseWriter, req *http.Request) {

	pVendor := req.URL.Query().Get("vendor")
	if pVendor == "" {
		http.Error(w, "vendor provider not specified", http.StatusBadRequest)
		return
	}

	pLat := req.URL.Query().Get("lat")
	pLatFloat, err := strconv.ParseFloat(strings.TrimSpace(pLat), 64)
	if err != nil {
		http.Error(w, "lat parameter not specified", http.StatusBadRequest)
		return
	}

	pLong := req.URL.Query().Get("long")
	pLongFloat, err := strconv.ParseFloat(strings.TrimSpace(pLong), 64)
	if err != nil {
		http.Error(w, "long parameter not specified", http.StatusBadRequest)
		return
	}

	pZoom := req.URL.Query().Get("zoom")
	pZoomFloat, err := strconv.ParseFloat(strings.TrimSpace(pZoom), 64)
	if err != nil {
		http.Error(w, "zoom parameter not specified", http.StatusBadRequest)
		return
	}

	pSide := req.URL.Query().Get("side")
	if pSide == "" {
		pSide = "3"
	}
	pSideInt, err := strconv.Atoi(pSide)
	if err != nil {
		http.Error(w, "side parameter not valid", http.StatusBadRequest)
		return
	}

	if pSideInt < 1 || pSideInt > 10 {
		http.Error(w, "side parameter must be side>1 and side<10", http.StatusBadRequest)
		return
	}

	vendor, err := a.Providers.Get(pVendor)
	if err != nil {
		http.Error(w, "vendor param incorrect", http.StatusBadRequest)
		return
	}

	if pZoomFloat < 1 || pZoomFloat > float64(vendor.MaxZoom()) {
		http.Error(w, fmt.Sprintf("zoom paramater not valid. max zoom for provider %s - %d", vendor.Name(), vendor.MaxZoom()), http.StatusBadRequest)
		return
	}

	tiles, err := downloader.StartMultiDownload(a.Cache, vendor, vendor.GetTile(pLatFloat, pLongFloat, pZoomFloat).GetNearby(pSideInt)...)
	if err != nil {
		a.Logger.Error("error occurred when downloading tiles", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
		http.Error(w, fmt.Sprintf("error occurred when downloading tiles error=%s", err), http.StatusInternalServerError)
		return
	}

	merged, err := downloader.Merge(pSideInt, vendor.GetTile(pLatFloat, pLongFloat, pZoomFloat), tiles...)
	if err != nil {
		a.Logger.Error("error occurred when merging tiles", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
		http.Error(w, "error occurred when merging tiles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = w.Write(merged)
	if err != nil {
		a.Logger.Error("error on merging image", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
	}
	a.Logger.Info("new map download request", zap.String("lat", pLat), zap.String("long", pLong), zap.String("side", pSide), zap.String("vendor", pVendor), zap.String("req_id", req.Header.Get("X-Request-ID")))

}
