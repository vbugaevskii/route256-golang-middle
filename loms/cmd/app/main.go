package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"route256/loms/internal/api"
	"route256/loms/pkg/loms"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":8081"
	httpPort = ":8091"
)

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, api.NewService())

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
	loms.RegisterLomsHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	log.Printf("Serving gRPC-Gateway on %s\n", httpPort)
	log.Fatalln(httpServer.ListenAndServe())
}
