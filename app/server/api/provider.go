package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// providerModel contains data about provider
type providerModel struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	MaxZoom int    `json:"max_zoom"`
}

// Provider godoc
// @Summary handler return all registered providers
// @Description reutrn JSON array with avalible provders
// @Accept  text/plain
// @Produce  application/json
// @Success		200	{array}	providerModel
// @Header 200 {string} X-Request-Id "request_id"
// @Router /provider [get]
func (a *API) Provider(w http.ResponseWriter, req *http.Request) {

	var allProviders = make([]providerModel, 0)

	for _, key := range a.Providers.GetAllID() {
		p, err := a.Providers.Get(key)
		if err != nil {
			continue
		}

		allProviders = append(allProviders, providerModel{Name: p.Name(), Key: p.ID(), MaxZoom: p.MaxZoom()})
	}

	results, err := json.Marshal(allProviders)
	if err != nil {
		a.Logger.Error("error on provder handler", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(results)
}
