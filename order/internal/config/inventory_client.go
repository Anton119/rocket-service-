package config

type inventoryClientConfig struct {
	Address string `yaml:"address" env:"INVENTORY_GRPC_ADDRESS" env-default:"localhost:50051"`
}

func (c *inventoryClientConfig) GRPCAddress() string {
	return c.Address
}
