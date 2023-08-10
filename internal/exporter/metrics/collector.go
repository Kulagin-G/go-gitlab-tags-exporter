package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	. "go-gitlab-tags-exporter/internal/config"
	. "go-gitlab-tags-exporter/internal/gitlab/v4"
	"golang.org/x/exp/slog"
	"time"
)

func NewCustomCollector(cfg *Config, log *slog.Logger, gl *Gitlab) *CustomCollector {
	return &CustomCollector{
		LatestTag: prometheus.NewDesc(
			"gitlab_tag_latest_info",
			"Returns the latest tag for a repo based on tag type",
			[]string{"project_name", "repository", "tag_type", "tag_name"}, nil,
		),
		GenDuration: prometheus.NewDesc(
			"gitlab_tag_parsing_duration_seconds",
			"Returns the time it took to parse all tags",
			[]string{}, nil,
		),
		cfg: cfg,
		log: log,
		gl:  gl,
	}
}

func (collector *CustomCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.LatestTag
	ch <- collector.GenDuration
}

func (collector *CustomCollector) Collect(ch chan<- prometheus.Metric) {
	generator := NewGenerator(collector.cfg, collector.log)
	start := time.Now()
	data := generator.GenerateData(collector.gl)

	if data == nil {
		collector.log.Warn("Skipping metric generation due to no data")
		return
	}

	elapsed := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(
		collector.GenDuration,
		prometheus.GaugeValue,
		elapsed,
	)

	for _, project := range data {
		ch <- prometheus.MustNewConstMetric(
			collector.LatestTag,
			prometheus.GaugeValue,
			1,
			project.Name,
			project.Repository,
			latestReleaseTagType,
			project.LatestReleaseTag.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.LatestTag,
			prometheus.GaugeValue,
			1,
			project.Name,
			project.Repository,
			latestReleaseCandidateTagType,
			project.LatestReleaseCandidateTag.Name,
		)

		collector.log.Info("Metrics have been generated", slog.String("project", project.Name))
	}
}
