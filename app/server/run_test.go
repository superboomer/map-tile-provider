package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/superboomer/maptile/app/options"
	"go.uber.org/zap"
)

func TestRunHTTP_Success(t *testing.T) {
	s := NewServer(&options.Opts{APIPort: "8080"}, zap.NewNop())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond) // Allow server to start
		cancel()                           // Cancel the context to trigger shutdown
	}()

	err := s.RunHTTP(ctx)
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_Success(t *testing.T) {
	logger := zap.NewNop()
	opts := &options.Opts{APIPort: "8080", Swagger: true, Schema: "./../../example/providers.json"}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond) // Allow server to start
		cancel()                           // Cancel the context to trigger shutdown
	}()

	err := Run(ctx, logger, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_Failed(t *testing.T) {
	logger := zap.NewNop()
	opts := &options.Opts{APIPort: "8080", Swagger: true, Schema: ""}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond) // Allow server to start
		cancel()                           // Cancel the context to trigger shutdown
	}()

	err := Run(ctx, logger, opts)
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}
