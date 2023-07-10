package main

import (
	"context"
	"route256/libs/logger"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	logger.Init(config.AppConfig.LogLevel)

	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.Token)
	if err != nil {
		logger.Fatal("failed telegram bot init", zap.Error(err))
	}
	bot.Debug = true

	group, err := kafka.NewConsumerGroup(
		config.AppConfig.Kafka.Brokers,
		config.AppConfig.Kafka.Group,
		config.AppConfig.Kafka.Topic,
	)
	if err != nil {
		logger.Fatal("failed kafka consumer init", zap.Error(err))
	}

	defer func() {
		if err = group.Close(); err != nil {
			logger.Fatal("failed kafka consumer close", zap.Error(err))
		}
	}()

	wg := &sync.WaitGroup{}
	listener := kafka.NewKafkaListener(group, bot)

	wg.Add(1)
	go func() {
		defer wg.Done()
		listener.RunListener(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		listener.RunServicer(context.Background(), config.AppConfig.Telegram.ChatId)
	}()

	wg.Wait()
}
