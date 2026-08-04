package config

type iamClientConfig struct {
	Address string `yaml:"address" env:"IAM_GRPC_ADDRESS" env-default:"localhost:50053"`
}

func (c *iamClientConfig) GRPCAddress() string {
	return c.Address
}
