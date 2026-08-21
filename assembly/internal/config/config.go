package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var appConfig *Config

// Config — корневая конфигурация AssemblyService.
type Config struct {
	Logger                loggerConfig                `yaml:"logger"`
	Otel                  otelConfig                  `yaml:"otel"`
	Kafka                 kafkaConfig                 `yaml:"kafka"`
	OrderPaidConsumer     orderPaidConsumerConfig     `yaml:"order_paid_consumer"`
	ShipAssembledProducer shipAssembledProducerConfig `yaml:"ship_assembled_producer"`
	Assembler             assemblerConfig             `yaml:"assembler"`
}

const defaultConfigPath = "config.local.yaml"

// ResolveConfigPath определяет путь к YAML: -config > CONFIG_PATH > config.local.yaml.
func ResolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}
	return defaultConfigPath
}

// MustLoad загружает конфиг (YAML + env) и сохраняет в AppConfig.
func MustLoad(path string) {
	var cfg Config
	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic(fmt.Sprintf("не удалось загрузить конфиг из %q: %v", path, err))
		}
	} else if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("не удалось загрузить конфиг из env: %v", err))
	}
	appConfig = &cfg
}

// AppConfig возвращает загруженный конфиг.
func AppConfig() *Config {
	return appConfig
}

// OTel возвращает конфигурацию observability.
func (c *Config) OTel() OTelConfig {
	return c.Otel
}
