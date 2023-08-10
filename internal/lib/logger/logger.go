package logger

import (
	. "git-tag-exporter/internal/config"
	"golang.org/x/exp/slog"
	"os"
	"strings"
)

const (
	debug = "debug"
	info  = "info"
)

func SetupLogger(cfg *Config) *slog.Logger {
	var log *slog.Logger

	logLevel := strings.ToLower(cfg.Exporter.LogLevel)

	switch logLevel {
	case debug:
		// TODO: Add pretty print
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case info:
		// TODO: Add pretty print
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: false}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: false}))
	}

	return log
}
