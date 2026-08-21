package config

import platformLogger "github.com/Anton119/rocket-service-/platform/pkg/logger"

const defaultEnvironment = "development"

type loggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
}

// LoggerPlatformConfig возвращает конфигурацию платформенного логгера (stdout + OTLP).
func (c *Config) LoggerPlatformConfig() platformLogger.Config {
	return platformLogger.Config{
		Level:             c.Logger.Level,
		ServiceName:       c.Otel.GetServiceName(),
		Environment:       defaultEnvironment,
		EnableOTLP:        true,
		CollectorEndpoint: c.Otel.CollectorEndpoint(),
	}
}
