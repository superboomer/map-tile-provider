package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Log(t *testing.T) {
	md := &MD{
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

func TestMiddleware_LogDurationsLong(t *testing.T) {
	md := &MD{
		Logger: zap.NewNop(),
	}

	handler := md.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		_, _ = w.Write([]byte("request completed"))
	}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "request completed")
}

func TestMiddleware_LogDurationsMedium(t *testing.T) {
	md := &MD{
		Logger: zap.NewNop(),
	}

	handler := md.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		_, _ = w.Write([]byte("request completed"))
	}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "request completed")
}

func TestMiddleware_LogDurationsMinimum(t *testing.T) {
	md := &MD{
		Logger: zap.NewNop(),
	}

	handler := md.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMiddleware_RequestID(t *testing.T) {
	md := &MD{
		Logger: zap.NewNop(),
	}

	handler := md.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req, _ := http.NewRequest("GET", "/", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, req.Header.Get("X-Request-ID"))
}
