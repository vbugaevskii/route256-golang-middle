package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
	"strconv"

	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalln("config init", err)
	}

	model := domain.New(
		cliloms.NewLomsClient(config.AppConfig.Services.Loms),
		cliproduct.NewProduct(config.AppConfig.Services.ProductService),
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
