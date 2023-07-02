package kafka

import (
	"context"
	"route256/libs/kafka"
)

type ConsumerGroup struct {
	group   *kafka.ConsumerGroup
	handler ConsumerGroupHandler
}

func NewConsumerGroup(brokers []string, groupName string, topic string) (*ConsumerGroup, error) {
	group, err := kafka.NewConsumerGroup(brokers, groupName, topic)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group:   group,
		handler: NewConsumerGroupHandler(),
	}, nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context) error {
	return cg.group.Consume(ctx, &cg.handler)
}

func (cg *ConsumerGroup) Ready() <-chan bool {
	return cg.handler.Ready()
}

func (cg *ConsumerGroup) Subscribe() <-chan Order {
	return cg.handler.Subscribe()
}

func (cg *ConsumerGroup) Close() error {
	cg.handler.Close()
	return cg.group.Close()
}
