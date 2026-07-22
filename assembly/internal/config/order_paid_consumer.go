package config

import "github.com/IBM/sarama"

type orderPaidConsumerConfig struct {
	Topic   string `yaml:"topic" env:"ORDER_PAID_TOPIC" env-default:"order.paid"`
	GroupID string `yaml:"group_id" env:"ORDER_PAID_CONSUMER_GROUP_ID" env-default:"assembly-service"`
}

func (c *orderPaidConsumerConfig) TopicName() string {
	return c.Topic
}

func (c *orderPaidConsumerConfig) ConsumerGroupID() string {
	return c.GroupID
}

func (c *orderPaidConsumerConfig) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	return cfg
}
