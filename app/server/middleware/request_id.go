package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

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
