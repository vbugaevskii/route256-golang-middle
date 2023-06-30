package main

import (
	"context"
	"log"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"sync"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalf("config init: %v", err)
	}

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
				log.Panicf("Error from consumer: %v", err)
			}
		}
	}()

	<-group.Ready()
	wg.Wait()
}
