package server

import (
	. "git-tag-exporter/internal/config"
	. "git-tag-exporter/internal/exporter/metrics"
	. "git-tag-exporter/internal/gitlab/v4"
	"git-tag-exporter/internal/lib/logger/sl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"time"
)

const (
	ReadHeaderTimeout = 60 * time.Second
)

type Server struct {
	cfg *Config
	log *slog.Logger
	reg prometheus.Registerer
	gl  *Gitlab
}

func NewServer(cfg *Config, log *slog.Logger, gl *Gitlab) *Server {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
		NewCustomCollector(cfg, log, gl),
	)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	http.Handle(cfg.Exporter.MetricsEndpoint, promHandler)

	return &Server{
		cfg: cfg,
		log: log,
		reg: reg,
		gl:  gl,
	}
}

func (s *Server) StartExporter() {
	go func() {
		s.log.Info("Starting exporter",
			slog.String("address", s.cfg.Exporter.Address),
			slog.String("port", s.cfg.Exporter.Port),
			slog.String("endpoint", s.cfg.Exporter.MetricsEndpoint))

		server := &http.Server{
			Addr:              s.cfg.Exporter.Address + ":" + s.cfg.Exporter.Port,
			ReadHeaderTimeout: ReadHeaderTimeout,
		}

		if err := server.ListenAndServe(); err != nil {
			s.log.Error("Failed to start exporter server", sl.Err(err))
			os.Exit(1)
		}
	}()
}
