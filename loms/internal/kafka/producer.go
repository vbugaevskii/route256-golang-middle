package kafka

import (
	"fmt"
	"route256/libs/kafka"
	"route256/loms/internal/domain"
)

type Producer struct {
	*kafka.Producer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		return nil, err
	}
	return &Producer{producer}, nil
}

func (p *Producer) SendOrderStatus(message domain.Notification) error {
	msg, err := p.BuildMessage(fmt.Sprint(message.OrderId), message)
	if err != nil {
		return err
	}

	if err = p.SendMessage(msg); err != nil {
		return err
	}
	return nil
}
