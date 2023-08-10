package config

import (
	"testing"
)

func TestMustConfig(t *testing.T) {
	t.Run("Must load config from file and env variables", func(t *testing.T) {
		t.Setenv("GITLAB_API_TOKEN", "def4561")
		t.Setenv("EXPORTER_CONFIG_PATH", "../../config/local.yaml")

		cfg := MustLoad()

		expected := &Config{
			Exporter: &Exporter{
				Address:        cfg.Exporter.Address,
				GitlabApiToken: "def4561",
			},
		}

		if cfg.Exporter.Address != expected.Exporter.Address || cfg.Exporter.GitlabApiToken != expected.Exporter.GitlabApiToken {
			t.Errorf("Expected config to be %+v, but got %+v", expected, cfg)
		}
	})
}
