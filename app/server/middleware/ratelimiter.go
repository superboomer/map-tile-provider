package middleware

import (
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// GlobalRateLimiter is a rate limiter for /map query/ its affect all requests
func (m *MD) GlobalRateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 3) // 2 per second with batch 3 events
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte("too many requests. try again later"))
			if err != nil {
				m.Logger.Error("error on GlobalRateLimiter", zap.Error(err), zap.String("req_id", r.Header.Get("X-Request-ID")))
			} else {
				m.Logger.Info("request was limited by GlobalRateLimiter", zap.Error(err), zap.String("req_id", r.Header.Get("X-Request-ID")))
			}
			return
		}

		next(w, r)
	})
}
