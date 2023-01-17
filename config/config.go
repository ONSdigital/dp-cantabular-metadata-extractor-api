package config

import (
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-cantabular-metadata-extractor-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	CantabularMetadataURL      string        `envconfig:"CANTABULAR_METADATA_URL"`
	AuthorisationConfig        *authorisation.Config
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	auth := authorisation.NewDefaultConfig()
	auth.Enabled = false

	cfg = &Config{
		BindAddr:                   "localhost:28300",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		CantabularMetadataURL:      "http://localhost:8493",
		AuthorisationConfig:        auth,
	}

	return cfg, envconfig.Process("", cfg)
}
