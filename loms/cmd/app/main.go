package main

import (
	"context"
	"net"
	"net/http"
	"route256/libs/logger"
	"route256/libs/tracing"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/api"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/kafka"
	"route256/loms/internal/repository/postgres/notificationsoutbox"
	"route256/loms/internal/repository/postgres/orders"
	"route256/loms/internal/repository/postgres/ordersreservations"
	"route256/loms/internal/repository/postgres/stocks"
	"route256/loms/pkg/loms"
	"strconv"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	logger.Init(config.AppConfig.LogLevel)
	tracing.Init(config.AppConfig.Name)

	pool, err := pgxpool.Connect(
		context.Background(),
		config.AppConfig.Postgres.URL(),
	)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer pool.Close()

	producer, err := kafka.NewProducer(
		config.AppConfig.Kafka.Brokers,
		config.AppConfig.Kafka.Topic,
	)
	if err != nil {
		logger.Fatal("failed to create kafka producer", zap.Error(err))
	}
	defer producer.Close()

	model := domain.NewModel(
		tx.NewTxManager(pool),
		producer,
		notificationsoutbox.NewNotificationsOutboxRepository(pool),
		stocks.NewStocksRepository(pool),
		orders.NewOrdersRepository(pool),
		ordersreservations.NewOrdersReservationsRepository(pool),
	)
	go func() {
		err := model.RunCancelOrderByTimeout(context.Background())
		if err != nil {
			logger.Fatal("failed to cancel order by timeout", zap.Error(err))
		}
	}()
	go func() {
		err := model.RunNotificationsSender(context.Background())
		if err != nil {
			logger.Fatal("failed to send notifications", zap.Error(err))
		}
	}()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port.GRPC))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
			),
		),
	)
	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, api.NewService(model))

	logger.Infof("server listening at %v", lis.Addr())

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// NOTE: https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		lis.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed to dial server", zap.Error(err))
	}

	mux := runtime.NewServeMux()
	err = loms.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		logger.Fatalf("failed to register gateway", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Port.HTTP),
		Handler: mux,
	}

	logger.Infof("Serving gRPC-Gateway on: %d", config.AppConfig.Port.HTTP)
	logger.Fatal("failed to serve http", zap.Error(httpServer.ListenAndServe()))
}
