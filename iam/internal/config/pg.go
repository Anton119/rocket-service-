package config

type pgConfig struct {
	URI string `yaml:"uri" env:"DB_URI" env-default:"postgres://iam-service-user:iam-service-password@localhost:5434/iam-service?sslmode=disable"`
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *pgConfig) DSN() (string, error) {
	return c.URI, nil
}
