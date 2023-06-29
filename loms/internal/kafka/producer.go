package kafka

import (
	"fmt"
	"route256/libs/kafka"
)

type Order struct {
	OrderId int64  `json:"order_id"`
	Status  string `json:"status"`
}

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

func (p *Producer) SendOrderStatus(orderId int64, status string) error {
	msg, err := p.BuildMessage(fmt.Sprint(orderId), Order{
		OrderId: orderId,
		Status:  status,
	})
	if err != nil {
		return err
	}

	if err = p.SendMessage(msg); err != nil {
		return err
	}
	return nil
}
