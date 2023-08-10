package metrics

import (
	. "git-tag-exporter/internal/config"
	. "git-tag-exporter/internal/gitlab/v4"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

type CustomCollector struct {
	LatestTag   *prometheus.Desc
	GenDuration *prometheus.Desc
	gl          *Gitlab
	cfg         *Config
	log         *slog.Logger
}
