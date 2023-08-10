package main

import (
	"git-tag-exporter/internal/config"
	. "git-tag-exporter/internal/exporter/healthz"
	. "git-tag-exporter/internal/exporter/server"
	. "git-tag-exporter/internal/gitlab/v4"
	"git-tag-exporter/internal/lib/logger"
	"git-tag-exporter/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg)
	gl, err := NewClient(cfg, log)

	if err != nil {
		log.Error("Failed to create Gitlab client", sl.Err(err))
		os.Exit(1)
	}

	health := NewHealth(cfg, log)
	health.StartHandlers()

	srv := NewServer(cfg, log, gl)
	srv.StartExporter()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	exit := <-sigChan
	log.Info("Stopped by signal", slog.String("signal", exit.String()))
}
