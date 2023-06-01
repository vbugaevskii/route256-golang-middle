package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"route256/checkout/internal/api"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalln("config init", err)
	}

	connLoms, err := grpc.Dial(
		config.AppConfig.Services.Loms.Netloc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer connLoms.Close()

	connProduct, err := grpc.Dial(
		config.AppConfig.Services.ProductService.Netloc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer connProduct.Close()

	model := domain.New(
		cliloms.NewLomsClient(connLoms),
		cliproduct.NewProductClient(
			connProduct,
			config.AppConfig.Services.ProductService.Token,
		),
	)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port.GRPC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	checkout.RegisterCheckoutServer(grpcServer, api.NewService(model))

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
	checkout.RegisterCheckoutHandler(context.Background(), mux, conn)
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
