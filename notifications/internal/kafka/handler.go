package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
)

type Order struct {
	OrderId int64  `json:"order_id"`
	Status  string `json:"status"`
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
	close(h.ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	close(h.output)
	return nil
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

			log.Printf("Message claimed: value = %v, timestamp = %v, topic = %s",
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
