package api

import (
	"encoding/json"
	"net/http"
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
func (a *API) HealthCheck(w http.ResponseWriter, _ *http.Request) {

	results, _ := json.Marshal(healthCheckModel{
		Status: 200,
		Body:   "OK",
	})

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(results)
}
