package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"route256/loms/internal/api"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/postgres/orders"
	"route256/loms/internal/repository/postgres/ordersreservations"
	"route256/loms/internal/repository/postgres/stocks"
	"route256/loms/internal/repository/postgres/tx"
	"route256/loms/pkg/loms"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalln("config init", err)
	}

	pool, err := pgxpool.Connect(
		context.Background(),
		config.AppConfig.Postgres.URL(),
	)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	model := domain.NewModel(
		tx.NewTxManager(pool),
		stocks.NewStocksRepository(pool),
		orders.NewOrdersRepository(pool),
		ordersreservations.NewOrdersReservationsRepository(pool),
	)
	go model.RunCancelOrderByTimeout(context.Background())

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port.GRPC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, api.NewService(model))

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
		log.Fatalln("Failed to dial server:", err)
	}

	mux := runtime.NewServeMux()
	err = loms.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.AppConfig.Port.HTTP),
		Handler: mux,
	}

	log.Printf("Serving gRPC-Gateway on :%d\n", config.AppConfig.Port.HTTP)
	log.Fatalln(httpServer.ListenAndServe())
}
