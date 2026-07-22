package config

import "github.com/IBM/sarama"

type orderPaidProducerConfig struct {
	Topic string `yaml:"topic" env:"ORDER_PAID_TOPIC" env-default:"order.paid"`
}

func (c *orderPaidProducerConfig) TopicName() string {
	return c.Topic
}

func (c *orderPaidProducerConfig) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Producer.Return.Successes = true
	return cfg
}
