package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {
	handAddToCart := addtocart.Handler{}
	http.Handle("/addToCart", srvwrapper.New(handAddToCart.Handle))

	handDeleteFromCart := deletefromcart.Handler{}
	http.Handle("/deleteFromCart", srvwrapper.New(handDeleteFromCart.Handle))

	handListCart := listcart.Handler{}
	http.Handle("/listCart", srvwrapper.New(handListCart.Handle))

	handPurchase := purchase.Handler{}
	http.Handle("/purchase", srvwrapper.New(handPurchase.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
