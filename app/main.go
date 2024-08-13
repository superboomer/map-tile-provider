package main

import (
	"context"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/superboomer/maptile/app/options"
	"github.com/superboomer/maptile/app/server"
	"github.com/umputun/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// Version contains build version
var Version = "dev"

// @title Map Satellite provider
// @version 1.0.0
// @description This is a easy HTTP API which provide map tiles
func main() {
	var Opts = &options.Opts{}
	p := flags.NewParser(Opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Fatalf("flags error: %v", err)
		}
		os.Exit(1)
	}

	logger := createLogger(&Opts.Log)

	logger.Info("build version", zap.String("build", Version))

	if err := run(Opts, logger); err != nil {
		logger.Fatal("fatal error", zap.Error(err))
	}
}

func run(opts *options.Opts, logger *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	return server.Run(ctx, logger, opts)
}

func createLogger(opts *options.Log) *zap.Logger {
	// Setting up logging to file with rotation.
	//
	// Log to file, so we don't interfere with prompts and messages to user.
	swSugar := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	)

	if opts.Save {
		logWriter := zapcore.AddSync(&lj.Logger{
			Filename:   opts.Path,
			MaxBackups: opts.MaxBackups,
			MaxSize:    opts.MaxSize, // megabytes
			MaxAge:     opts.MaxAge,  // days
		})

		swSugar = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			logWriter,
		)
	}

	encoder := zap.NewProductionEncoderConfig()

	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		swSugar,
		zap.InfoLevel,
	)
	logger := zap.New(logCore)
	defer func() { _ = logger.Sync() }()

	if opts.Save {
		logger.Info("log writer enabled", zap.String("path", opts.Path))
	}

	return logger
}
