package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	handCreateOrder := createorder.Handler{}
	http.Handle(createorder.Endpoint, srvwrapper.New(handCreateOrder.Handle))

	handListOrder := &listorder.Handler{}
	http.Handle(listorder.Endpoint, srvwrapper.New(handListOrder.Handle))

	handOrderPayed := orderpayed.Handler{}
	http.Handle(orderpayed.Endpoint, srvwrapper.New(handOrderPayed.Handle))

	handCancelOrder := cancelorder.Handler{}
	http.Handle(cancelorder.Endpoint, srvwrapper.New(handCancelOrder.Handle))

	handStocks := stocks.Handler{}
	http.Handle(stocks.Endpoint, srvwrapper.New(handStocks.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
