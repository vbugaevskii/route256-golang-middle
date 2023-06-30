package main

import (
	"context"
	"fmt"
	"log"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalf("config init: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	group, err := kafka.NewConsumerGroup(
		config.AppConfig.Kafka.Brokers,
		config.AppConfig.Kafka.Group,
		config.AppConfig.Kafka.Topic,
	)
	defer func() {
		if err = group.Close(); err != nil {
			log.Fatalf("kafka consumer close: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("kafka consumer init: %v", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := group.Consume(context.Background()); err != nil {
				log.Fatalf("kafka consumer read: %v", err)
			}
		}
	}()

	<-group.Ready()

	for order := range group.Subscribe() {
		msg := tgbotapi.NewMessage(
			config.AppConfig.Telegram.ChatId,
			fmt.Sprintf("OrderId = %d; Status = %s", order.OrderId, order.Status),
		)
		bot.Send(msg)
	}

	wg.Wait()
}
