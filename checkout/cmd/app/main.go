package main

import (
	"log"
	"net/http"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	handAddToCart := addtocart.Handler{
		Model: model,
	}
	http.Handle("/addToCart", srvwrapper.New(handAddToCart.Handle))

	handDeleteFromCart := deletefromcart.Handler{}
	http.Handle("/deleteFromCart", srvwrapper.New(handDeleteFromCart.Handle))

	handListCart := listcart.Handler{
		Model: model,
	}
	http.Handle("/listCart", srvwrapper.New(handListCart.Handle))

	handPurchase := purchase.Handler{
		Model: model,
	}
	http.Handle("/purchase", srvwrapper.New(handPurchase.Handle))

	err = http.ListenAndServe(":"+strconv.Itoa(config.AppConfig.Port), nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
