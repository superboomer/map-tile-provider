package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log use for logging all http requests
func (m *MD) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var start, end int64
		start = time.Now().UnixNano()
		path := req.URL.Path

		next.ServeHTTP(w, req)

		end = time.Now().UnixNano()
		duration := end - start
		t := ""
		if duration >= 1000000000 {
			t = fmt.Sprintf("%.2fs", float64(duration)/1000000000)
		} else if duration >= 1000000 {
			t = fmt.Sprintf("%.2fms", float64(duration)/1000000)
		} else if duration >= 1000 {
			t = fmt.Sprintf("%.2fµs", float64(duration)/1000)
		} else {
			t = fmt.Sprintf("%ddns", duration)
		}

		m.Logger.Info("request completed", zap.String("method", req.Method), zap.String("path", path), zap.String("duration", t), zap.String("req_id", req.Header.Get("X-Request-ID")))
	})
}
