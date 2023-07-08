package main

import (
	"context"
	"net"
	"net/http"
	"route256/checkout/internal/api"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	pgcartitems "route256/checkout/internal/repository/postgres/cartitems"
	"route256/checkout/pkg/checkout"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"
	"strconv"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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
	metrics.Init(config.AppConfig.Name)

	connLoms, err := grpc.Dial(
		config.AppConfig.Services.Loms.Netloc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed to connect to server", zap.Error(err))
	}
	defer connLoms.Close()

	connProduct, err := grpc.Dial(
		config.AppConfig.Services.ProductService.Netloc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed to connect to server", zap.Error(err))
	}
	defer connProduct.Close()

	pool, err := pgxpool.Connect(
		context.Background(),
		config.AppConfig.Postgres.URL(),
	)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer pool.Close()

	model := domain.New(
		cliloms.NewLomsClient(connLoms),
		cliproduct.NewProductClient(
			connProduct,
			config.AppConfig.Services.ProductService.Token,
			config.AppConfig.Services.ProductService.RPS,
		),
		pgcartitems.NewCartItemsRepository(pool),
	)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port.GRPC))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				metrics.ServerMetricsInterceptor,
				logger.LoggingInterceptor,
			),
		),
	)
	reflection.Register(grpcServer)

	prometheus.Register(grpcServer)
	go func() {
		err := metrics.ListenAndServeMetrics(config.AppConfig.Metrics.Port)
		if err != nil {
			logger.Fatal("failed to serve metrics", zap.Error(err))
		}
	}()

	checkout.RegisterCheckoutServer(grpcServer, api.NewService(model))
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
	err = checkout.RegisterCheckoutHandler(context.Background(), mux, conn)
	if err != nil {
		logger.Fatal("failed to register gateway", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Port.HTTP),
		Handler: mux,
	}

	logger.Infof("Serving gRPC-Gateway on: %d", config.AppConfig.Port.HTTP)
	logger.Fatal("failed to serve http", zap.Error(httpServer.ListenAndServe()))
}
