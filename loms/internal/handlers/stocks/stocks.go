package stocks

import (
	"context"
	"log"
)

const Endpoint = "/stocks"

type Handler struct {
}

type Request struct {
	SKU uint32 `json:"sku"`
}

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Response struct {
	Stocks []StockItem `json:"stocks"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)

	return Response{
		Stocks: []StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}, nil
}
