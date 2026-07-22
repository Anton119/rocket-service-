package config

import "github.com/IBM/sarama"

type shipAssembledProducerConfig struct {
	Topic string `yaml:"topic" env:"SHIP_ASSEMBLED_TOPIC" env-default:"assembly.ship-assembled"`
}

func (c *shipAssembledProducerConfig) TopicName() string {
	return c.Topic
}

func (c *shipAssembledProducerConfig) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Producer.Return.Successes = true
	return cfg
}
