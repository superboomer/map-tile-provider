package api

import (
	"go.uber.org/zap"
	"net/http"
)

// HealthCheck godoc
// @Summary handler for health check
// @Description just return 200 with string
// @Accept  json
// @Produce  json
// @Success 200 {string} Data "service ok"
// @Header 200 {string} X-Request-Id "request_id"
// @Router /healthcheck [get]
func (a *API) HealthCheck(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("service ok"))
	if err != nil {
		a.Logger.Error("error on health check", zap.Error(err), zap.String("req_id", req.Header.Get("X-Request-ID")))
	}
}
