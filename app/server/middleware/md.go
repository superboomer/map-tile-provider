package middleware

import "go.uber.org/zap"

// MD represent struct for middlewares
type MD struct {
	Logger *zap.Logger
}
