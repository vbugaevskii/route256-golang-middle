package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func NewConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false

	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
