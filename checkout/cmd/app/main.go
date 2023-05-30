package main

import (
	"log"
	"net"
	"route256/checkout/internal/api"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout"
	"strconv"

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

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	checkout.RegisterCheckoutServer(s, api.NewService(model))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
