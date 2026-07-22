package config

import "github.com/IBM/sarama"

type shipAssembledConsumerConfig struct {
	Topic   string `yaml:"topic" env:"SHIP_ASSEMBLED_TOPIC" env-default:"assembly.ship-assembled"`
	GroupID string `yaml:"group_id" env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID" env-default:"order-service"`
}

func (c *shipAssembledConsumerConfig) TopicName() string {
	return c.Topic
}

func (c *shipAssembledConsumerConfig) ConsumerGroupID() string {
	return c.GroupID
}

func (c *shipAssembledConsumerConfig) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	return cfg
}
