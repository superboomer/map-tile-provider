package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// healthCheckModel contains data about health check
type healthCheckModel struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// HealthCheck godoc
// @Summary handler for health check
// @Description just return HealthCheckModel with API status (always return 200)
// @Accept  json
// @Produce  application/json
// @Success		200	{object}	healthCheckModel
// @Header 200 {string} X-Request-Id "request_id"
// @Router /healthcheck [get]
func (a *API) HealthCheck(w http.ResponseWriter, req *http.Request) {

	var res = healthCheckModel{
		Status: 200,
		Body:   "OK",
	}

	results, err := json.Marshal(res)
	if err != nil {
		a.Logger.Error("error onhealth check handler", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(results)
	if err != nil {
		a.Logger.Error("error on health check handler", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
	}
}
