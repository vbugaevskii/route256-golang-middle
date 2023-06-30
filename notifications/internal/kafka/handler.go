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
	ready chan bool
}

func NewConsumerGroupHandler() ConsumerGroupHandler {
	return ConsumerGroupHandler{
		ready: make(chan bool),
	}
}

func (h *ConsumerGroupHandler) Ready() <-chan bool {
	return h.ready
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
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

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
