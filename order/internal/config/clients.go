package config

type inventoryClientConfig struct {
	Address string `yaml:"address" env:"INVENTORY_CLIENT_ADDRESS" env-default:"localhost:50051"`
}

type paymentClientConfig struct {
	Address string `yaml:"address" env:"PAYMENT_CLIENT_ADDRESS" env-default:"localhost:50052"`
}
