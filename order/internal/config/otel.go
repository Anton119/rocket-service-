package config

// OTelConfig — общая конфигурация observability (логи, метрики, трейсы).
type OTelConfig interface {
	CollectorEndpoint() string
	GetServiceName() string
}

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
