package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MD represent struct for middlewares
type MD struct {
	Logger *zap.Logger
}

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

// RequestID generate reqID and set it to header
func (m *MD) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rid := req.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
			req.Header.Add("X-Request-ID", rid)
			w.Header().Add("X-Request-ID", rid)
		}
		next.ServeHTTP(w, req)
	})
}
