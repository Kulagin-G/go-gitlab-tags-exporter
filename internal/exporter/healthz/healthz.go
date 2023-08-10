package healthz

import (
	"fmt"
	. "git-tag-exporter/internal/config"
	"git-tag-exporter/internal/gitlab/v4"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"time"
)

type HealthCheck struct {
	cfg *Config
	log *slog.Logger
}

func NewHealth(cfg *Config, log *slog.Logger) *HealthCheck {
	h := &HealthCheck{
		cfg: cfg,
		log: log,
	}

	return h
}

func (h *HealthCheck) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "{'live': true}")
}

func (h *HealthCheck) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := gitlab.NewClient(h.cfg, h.log)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("{'ready': false, 'err': %v}", err))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "{'ready': 'true'}")
}

func (h *HealthCheck) StartHandlers() {
	addr := fmt.Sprintf("%s:%s", h.cfg.Exporter.Address, h.cfg.Exporter.Health.Port)
	server := http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	h.log.Info("Starting health handlers", slog.String("addr", addr))

	http.HandleFunc("/healthz/live", h.LivenessHandler)
	http.HandleFunc("/healthz/ready", h.ReadinessHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
}
