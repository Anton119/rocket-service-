package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var appConfig *Config

// Config — корневая конфигурация InventoryService.
type Config struct {
	Logger    loggerConfig    `yaml:"logger"`
	GRPC      grpcConfig      `yaml:"grpc"`
	PG        pgConfig        `yaml:"pg"`
	IAMClient iamClientConfig `yaml:"iam_client"`
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
