package config

import "time"

type Config struct {
	Exporter       *Exporter         `yaml:"exporter"`
	MetricsOptions *MetricsOptions   `yaml:"metricsOptions"`
	Projects       []*ProjectConfigs `yaml:"projects"`
}

type MetricsOptions struct {
	ReleaseTagRegex          string `yaml:"releaseTagPattern"`
	ReleaseCandidateTagRegex string `yaml:"releaseCandidateTagPattern"`
}

type Exporter struct {
	Address           string        `yaml:"address" env:"EXPORTER_ADDRESS_BIND" env-default:"0.0.0.0"`
	Port              string        `yaml:"port" env:"EXPORTER_PORT_BIND" env-default:"8090"`
	LogLevel          string        `yaml:"logLevel" env:"EXPORTER_LOG_LEVEL" env-default:"info"`
	GitlabUrl         string        `yaml:"gitlabUrl" env:"GITLAB_URL" env-default:"https://git.ringcentral.com"`
	GitlabApiToken    string        `yaml:"gitlabApiToken,omitempty"`
	GitlabApiRetryMax int           `yaml:"gitlabApiRetryMax" env:"EXPORTER_API_RETRY_MAX" env-default:"5"`
	GoroutinesMax     int           `yaml:"goroutinesMax" env:"EXPORTER_GOROUTINES_MAX" env-default:"10"`
	GoroutinesTimeout time.Duration `yaml:"goroutinesTimeout" env:"EXPORTER_GOROUTINES_TIMEOUT" env-default:"60s"`
	MetricsEndpoint   string        `yaml:"metricsEndpoint" env:"EXPORTER_METRICS_ENDPOINT" env-default:"/metrics"`
	Health            *Health       `yaml:"health"`
}

type ProjectConfigs struct {
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	Repository string `yaml:"repository,omitempty"`
}

type Health struct {
	Port string `yaml:"port" env:"EXPORTER_HEALTH_PORT" env-default:"8091"`
}
