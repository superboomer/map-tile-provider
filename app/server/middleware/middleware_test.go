package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/server/middleware"
)

func TestMiddleware_Log(t *testing.T) {
	md := &middleware.MD{
		Logger: zap.NewNop(),
	}

	handler := md.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("request completed"))
	}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "request completed")
}

func TestMiddleware_RequestID(t *testing.T) {
	md := &middleware.MD{
		Logger: zap.NewNop(),
	}

	handler := md.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, req.Header.Get("X-Request-ID"))
}
