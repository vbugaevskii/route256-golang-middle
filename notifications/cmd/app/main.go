package main

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Fatalf("config init: %v", err)
	}

	logger.Init(config.AppConfig.LogLevel)

	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.Token)
	if err != nil {
		logger.Fatalf("telegram bot init: %v", err)
	}
	bot.Debug = true

	group, err := kafka.NewConsumerGroup(
		config.AppConfig.Kafka.Brokers,
		config.AppConfig.Kafka.Group,
		config.AppConfig.Kafka.Topic,
	)
	defer func() {
		if err = group.Close(); err != nil {
			logger.Fatalf("kafka consumer close: %v", err)
		}
	}()

	if err != nil {
		logger.Fatalf("kafka consumer init: %v", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := group.Consume(context.Background()); err != nil {
				logger.Fatalf("kafka consumer read: %v", err)
			}
		}
	}()

	<-group.Ready()
	logger.Info("service ready to listen to kafka")

	for order := range group.Subscribe() {
		msg := tgbotapi.NewMessage(
			config.AppConfig.Telegram.ChatId,
			fmt.Sprintf("[%v] OrderId = %d; Status = %s", order.CreatedAt, order.OrderId, order.Status),
		)
		if _, err := bot.Send(msg); err != nil {
			logger.Infof("failed to send message %v", err)
		}
	}

	wg.Wait()
}
