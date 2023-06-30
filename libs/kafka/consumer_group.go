package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
)

type ConsumerGroup struct {
	group sarama.ConsumerGroup
	topic string
}

func NewConsumerGroup(brokers []string, groupName string, topic string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()

	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	// config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.BalanceStrategyRoundRobin,
	}

	group, err := sarama.NewConsumerGroup(brokers, groupName, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group: group,
		topic: topic,
	}, nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
	return cg.group.Consume(ctx, []string{cg.topic}, handler)
}

func (cg *ConsumerGroup) Close() error {
	return cg.group.Close()
}
