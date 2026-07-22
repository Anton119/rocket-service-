package config

type paymentClientConfig struct {
	Address string `yaml:"address" env:"PAYMENT_GRPC_ADDRESS" env-default:"localhost:50052"`
}

func (c *paymentClientConfig) GRPCAddress() string {
	return c.Address
}
