package kafka

import (
	"context"
	"fmt"
	"route256/libs/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type KafkaListener struct {
	group *ConsumerGroup
	bot   *tgbotapi.BotAPI
}

func NewKafkaListener(group *ConsumerGroup, bot *tgbotapi.BotAPI) *KafkaListener {
	return &KafkaListener{
		group: group,
		bot:   bot,
	}
}

func (kl *KafkaListener) RunListener(ctx context.Context) {
	for {
		if err := kl.group.Consume(context.Background()); err != nil {
			logger.Fatal("failed kafka consumer read", zap.Error(err))
		}
	}
}

func (kl *KafkaListener) RunServicer(ctx context.Context, chatId int64) {
	<-kl.group.Ready()
	logger.Info("service ready to listen to kafka")

	for order := range kl.group.Subscribe() {
		msg := tgbotapi.NewMessage(
			chatId,
			fmt.Sprintf("[%v] OrderId = %d; Status = %s", order.CreatedAt, order.OrderId, order.Status),
		)
		if _, err := kl.bot.Send(msg); err != nil {
			logger.Infof("failed to send message %v", err)
		}
	}
}
