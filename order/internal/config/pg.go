package config

import sharedconfig "github.com/Anton119/rocket-service-/shared/pkg/config"

type pgConfig struct {
	URI string `yaml:"uri" env:"DB_URI" env-default:"postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable"`
}

// DSN возвращает строку подключения к PostgreSQL.
// Приоритет: DB_URI из env/yaml; при пустом значении — shared/pkg/config.DBURI().
func (c *pgConfig) DSN() (string, error) {
	if c.URI != "" {
		return c.URI, nil
	}
	return sharedconfig.DBURI()
}
