package kafka

import (
	"encoding/json"
	"route256/libs/logger"
	"time"

	"github.com/Shopify/sarama"
)

type Order struct {
	UserId    int64     `json:"user_id"`
	OrderId   int64     `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ConsumerGroupHandler struct {
	ready  chan bool
	output chan Order
}

func NewConsumerGroupHandler() ConsumerGroupHandler {
	return ConsumerGroupHandler{
		ready:  make(chan bool),
		output: make(chan Order),
	}
}

func (h *ConsumerGroupHandler) Ready() <-chan bool {
	return h.ready
}

func (h *ConsumerGroupHandler) Subscribe() <-chan Order {
	return h.output
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	logger.Info("setup handler")
	close(h.ready)
	h.ready = nil
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	logger.Info("cleanup handler")
	h.ready = make(chan bool)
	return nil
}

func (h *ConsumerGroupHandler) Close() {
	close(h.output)
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			order := Order{}
			err := json.Unmarshal(message.Value, &order)
			if err != nil {
				return err
			}

			logger.Infof("Message claimed: value = %v, timestamp = %v, topic = %s",
				order,
				message.Timestamp,
				message.Topic,
			)

			h.output <- order

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
