package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	. "go-gitlab-tags-exporter/internal/config"
	. "go-gitlab-tags-exporter/internal/gitlab/v4"
	"golang.org/x/exp/slog"
)

type CustomCollector struct {
	LatestTag   *prometheus.Desc
	GenDuration *prometheus.Desc
	gl          *Gitlab
	cfg         *Config
	log         *slog.Logger
}
