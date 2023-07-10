package main

import (
	"context"
	"log"
	"net"
	"route256/libs/logger"
	"route256/notifications/internal/api"
	"route256/notifications/internal/config"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/listener"
	pgnotify "route256/notifications/internal/repository/postgres/notifications"
	pbnotify "route256/notifications/pkg/notifications"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	lrucache "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	pool, err := pgxpool.Connect(
		context.Background(),
		config.AppConfig.Postgres.URL(),
	)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer pool.Close()

	cache, err := lrucache.New[int64, []domain.Notification](128)
	if err != nil {
		logger.Fatal("failed to init cache", zap.Error(err))
	}

	wg := &sync.WaitGroup{}
	klistener := listener.NewKafkaListener(
		group,
		bot,
		pgnotify.NewNotificationsRepository(pool),
		cache,
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		klistener.RunListener(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		klistener.RunServicer(context.Background(), config.AppConfig.Telegram.ChatId)
	}()

	model := domain.NewModel(
		pgnotify.NewNotificationsRepository(pool),
		cache,
	)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port.GRPC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pbnotify.RegisterNotificationsServer(grpcServer, api.NewService(model))

	logger.Infof("server listening at %v", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}

	wg.Wait()
}
