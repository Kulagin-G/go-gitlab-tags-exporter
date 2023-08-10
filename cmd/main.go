package main

import (
	"go-gitlab-tags-exporter/internal/config"
	. "go-gitlab-tags-exporter/internal/exporter/healthz"
	. "go-gitlab-tags-exporter/internal/exporter/server"
	. "go-gitlab-tags-exporter/internal/gitlab/v4"
	"go-gitlab-tags-exporter/internal/lib/logger"
	"go-gitlab-tags-exporter/internal/lib/logger/sl"
	_ "go.uber.org/automaxprocs"
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
