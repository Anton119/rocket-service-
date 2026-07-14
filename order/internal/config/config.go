package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigPath = "order/config.local.yaml"

// Config — конфигурация OrderService.
type Config struct {
	HTTP            httpConfig            `yaml:"http"`
	Logger          loggerConfig          `yaml:"logger"`
	PG              pgConfig              `yaml:"pg"`
	InventoryClient inventoryClientConfig `yaml:"inventory_client"`
	PaymentClient   paymentClientConfig   `yaml:"payment_client"`
}

// ResolveConfigPath определяет путь к конфиг-файлу: флаг -config > env CONFIG_PATH > default.
func ResolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу (например, order/config.production.yaml)")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return defaultConfigPath
}

// Load загружает конфигурацию: YAML-файл + env-переменные поверх.
func Load(path string) (*Config, error) {
	var cfg Config

	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("загрузить конфиг из %q: %w", path, err)
		}

		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("загрузить конфиг из env: %w", err)
	}

	return &cfg, nil
}
