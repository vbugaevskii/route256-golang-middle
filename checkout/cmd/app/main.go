package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {
	model := domain.New("http://localhost:8081")

	handAddToCart := addtocart.Handler{
		Model: model,
	}
	http.Handle(addtocart.Endpoint, srvwrapper.New(handAddToCart.Handle))

	handDeleteFromCart := deletefromcart.Handler{}
	http.Handle(deletefromcart.Endpoint, srvwrapper.New(handDeleteFromCart.Handle))

	handListCart := listcart.Handler{}
	http.Handle(listcart.Endpoint, srvwrapper.New(handListCart.Handle))

	handPurchase := purchase.Handler{
		Model: model,
	}
	http.Handle(purchase.Endpoint, srvwrapper.New(handPurchase.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
