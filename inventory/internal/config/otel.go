package config

import platformTracing "github.com/Anton119/rocket-service-/platform/pkg/tracing"

// OTelConfig — общая конфигурация observability (логи, метрики, трейсы).
type OTelConfig interface {
	CollectorEndpoint() string
	GetServiceName() string
}

const tracingServiceVersion = "1.0.0"

// otelConfig — секция otel в YAML-конфиге сервиса.
type otelConfig struct {
	Endpoint    string `yaml:"endpoint" env:"OTEL_EXPORTER_OTLP_ENDPOINT" env-default:"localhost:4317"`
	ServiceName string `yaml:"service_name" env:"OTEL_SERVICE_NAME"`
}

func (c otelConfig) CollectorEndpoint() string {
	return c.Endpoint
}

func (c otelConfig) GetServiceName() string {
	return c.ServiceName
}

// TracingPlatformConfig возвращает конфигурацию платформенного трейсера.
func (c *Config) TracingPlatformConfig() platformTracing.Config {
	return platformTracing.Config{
		CollectorEndpoint: c.Otel.CollectorEndpoint(),
		ServiceName:       c.Otel.GetServiceName(),
		Environment:       defaultEnvironment,
		ServiceVersion:    tracingServiceVersion,
		SamplingRatio:     1.0,
	}
}
