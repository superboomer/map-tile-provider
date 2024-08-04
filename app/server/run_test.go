package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/superboomer/map-tile-provider/app/options"
	"github.com/superboomer/map-tile-provider/app/server"
	"go.uber.org/zap"
)

// Test RunHTTP
func TestRunHTTP(t *testing.T) {
	s := server.NewServer(&options.Opts{APIPort: "8080"}, zap.NewNop())

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

// Test Run
func TestRun(t *testing.T) {
	logger := zap.NewNop()
	opts := &options.Opts{APIPort: "8080"}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond) // Allow server to start
		cancel()                           // Cancel the context to trigger shutdown
	}()

	err := server.Run(ctx, logger, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
