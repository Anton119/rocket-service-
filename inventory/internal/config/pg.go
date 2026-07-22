package config

import sharedconfig "github.com/Anton119/rocket-service-/shared/pkg/config"

type pgConfig struct {
	URI string `yaml:"uri" env:"DB_URI" env-default:"postgres://inventory-service-user:inventory-service-password@localhost:5433/inventory-service?sslmode=disable"`
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *pgConfig) DSN() (string, error) {
	if c.URI != "" {
		return c.URI, nil
	}
	return sharedconfig.InventoryDBURI()
}
