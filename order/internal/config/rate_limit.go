package config

// RateLimitConfig — публичный контракт секции rate_limit.
type RateLimitConfig interface {
	Address() string
	Rate() int
	Burst() int
}

type rateLimitConfig struct {
	RedisAddress string `yaml:"redis_address" env:"RATE_LIMIT_REDIS_ADDRESS" env-default:"localhost:6379"`
	RateValue    int    `yaml:"rate" env:"RATE_LIMIT_RATE" env-default:"100"`
	BurstValue   int    `yaml:"burst" env:"RATE_LIMIT_BURST" env-default:"200"`
}

func (c rateLimitConfig) Address() string {
	return c.RedisAddress
}

func (c rateLimitConfig) Rate() int {
	if c.RateValue <= 0 {
		return 100
	}
	return c.RateValue
}

func (c rateLimitConfig) Burst() int {
	if c.BurstValue <= 0 {
		return 200
	}
	return c.BurstValue
}

// RateLimit возвращает конфигурацию distributed rate limiter.
func (c *Config) RateLimit() RateLimitConfig {
	return c.RateLimitSection
}
