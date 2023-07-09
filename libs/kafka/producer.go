package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll

	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) BuildMessage(key string, value any) (*sarama.ProducerMessage, error) {
	val, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Failed to serialize value", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     p.topic,
		Partition: -1,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(val),
	}, nil
}

func (p *Producer) SendMessage(msg *sarama.ProducerMessage) error {
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
