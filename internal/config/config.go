package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

func MustLoad() *Config {
	configPath := os.Getenv("EXPORTER_CONFIG_PATH")
	gitlabApiToken := os.Getenv("GITLAB_API_TOKEN")

	if configPath == "" {
		log.Fatalln("EXPORTER_CONFIG_PATH is not set!")
	}

	if gitlabApiToken == "" {
		log.Fatalln("GITLAB_API_TOKEN is not set!")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s is not found: %v", configPath, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Config file %s is not readable: %v", configPath, err)
	}

	cfg.Exporter.GitlabApiToken = gitlabApiToken

	return &cfg
}
