package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Provider godoc
// @Summary handler return all registered providers name
// @Description reutrn JSON array with avalible provders
// @Accept  text/plain
// @Produce  application/json
// @Success 200 {array} string
// @Header 200 {string} X-Request-Id "request_id"
// @Router /provider [get]
func (a *API) Provider(w http.ResponseWriter, req *http.Request) {

	names := a.Providers.GetAllNames()

	results, err := json.Marshal(names)
	if err != nil {
		a.Logger.Error("error on provder handler", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
		return
	}

	_, err = w.Write(results)
	if err != nil {
		a.Logger.Error("error on provider handler", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
	}
}
