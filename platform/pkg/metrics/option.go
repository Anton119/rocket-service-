package metrics

import (
	"time"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const defaultCollectorEndpoint = "localhost:4317"

// options хранит настройки MeterProvider.
type options struct {
	// interval — интервал экспорта метрик в коллектор.
	// По умолчанию 10 секунд (вместо дефолтных 60 у OTel SDK) —
	// чтобы метрики быстрее появлялись в Prometheus/Grafana при локальной разработке.
	interval time.Duration

	// collectorEndpoint — адрес OTLP коллектора (host:port).
	collectorEndpoint string

	// views — пользовательские View для переопределения агрегации метрик.
	// Например, кастомные бакеты гистограмм для конкретных инструментов.
	views []sdkmetric.View
}

// Option — функциональная опция для настройки MeterProvider.
type Option func(*options)

// WithCollectorEndpoint задаёт адрес OTLP коллектора (например, "localhost:4317").
func WithCollectorEndpoint(endpoint string) Option {
	return func(o *options) {
		if endpoint != "" {
			o.collectorEndpoint = endpoint
		}
	}
}

// WithInterval задаёт интервал экспорта метрик.
// По умолчанию 10 секунд — подходит для локальной разработки.
// В production рекомендуется 15–60 секунд.
func WithInterval(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.interval = d
		}
	}
}

// WithView добавляет View для переопределения агрегации конкретной метрики.
//
// Пример — кастомные бакеты для гистограммы rpc.server.call.duration:
//
//	metrics.WithView(sdkmetric.NewView(
//	    sdkmetric.Instrument{Name: "rpc.server.call.duration"},
//	    sdkmetric.Stream{
//	        Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
//	            Boundaries: []float64{0.0001, 0.001, 0.01, 0.1, 0.5, 1, 5},
//	        },
//	    },
//	))
func WithView(v sdkmetric.View) Option {
	return func(o *options) {
		o.views = append(o.views, v)
	}
}
