package config

import (
	"net"
	"time"
)

type httpConfig struct {
	Host                 string `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
	Port                 string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadHeaderTimeoutSec int    `yaml:"read_header_timeout_sec" env:"HTTP_READ_HEADER_TIMEOUT_SEC" env-default:"5"`
	ReadTimeoutSec       int    `yaml:"read_timeout_sec" env:"HTTP_READ_TIMEOUT_SEC" env-default:"15"`
	WriteTimeoutSec      int    `yaml:"write_timeout_sec" env:"HTTP_WRITE_TIMEOUT_SEC" env-default:"15"`
	IdleTimeoutSec       int    `yaml:"idle_timeout_sec" env:"HTTP_IDLE_TIMEOUT_SEC" env-default:"60"`
	ShutdownTimeoutSec   int    `yaml:"shutdown_timeout_sec" env:"HTTP_SHUTDOWN_TIMEOUT_SEC" env-default:"10"`
}

func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *httpConfig) ReadHeaderTimeout() time.Duration {
	return time.Duration(c.ReadHeaderTimeoutSec) * time.Second
}

func (c *httpConfig) ReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeoutSec) * time.Second
}

func (c *httpConfig) WriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeoutSec) * time.Second
}

func (c *httpConfig) IdleTimeout() time.Duration {
	return time.Duration(c.IdleTimeoutSec) * time.Second
}

func (c *httpConfig) ShutdownTimeout() time.Duration {
	return time.Duration(c.ShutdownTimeoutSec) * time.Second
}
