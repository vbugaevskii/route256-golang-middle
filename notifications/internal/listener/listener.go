package listener

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/kafka"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	ListNotifications(ctx context.Context, userId int64) ([]domain.Notification, error)
	SaveNotification(ctx context.Context, recordId int64, userId int64, message string) error
}

type KafkaListener struct {
	group *kafka.ConsumerGroup
	bot   *tgbotapi.BotAPI
	repo  NotificationRepository
}

func NewKafkaListener(group *kafka.ConsumerGroup, bot *tgbotapi.BotAPI, repo NotificationRepository) *KafkaListener {
	return &KafkaListener{
		group: group,
		bot:   bot,
		repo:  repo,
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
		msgTxt := fmt.Sprintf(
			"[%v] UserId = %d; OrderId = %d; Status = %s",
			order.CreatedAt,
			order.UserId,
			order.OrderId,
			order.Status,
		)

		msg := tgbotapi.NewMessage(chatId, msgTxt)
		if _, err := kl.bot.Send(msg); err != nil {
			logger.Infof("failed to send message %v", err)
		}

		err := kl.repo.SaveNotification(ctx, order.RecordId, order.UserId, msgTxt)
		if err != nil {
			logger.Infof("failed to save messag %v", err)
		}
	}
}
