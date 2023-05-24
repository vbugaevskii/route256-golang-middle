package client

import (
	"context"
)

type RequestStocks struct {
	SKU uint32 `json:"sku"`
}

type ResponseStocks struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

type RequestCreateOrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint64 `json:"count"`
}

type RequestCreateOrder struct {
	User  int64                    `json:"user"`
	Items []RequestCreateOrderItem `json:"items"`
}

type ResponseCreateOrder struct {
	OrderId int64 `json:"orderID"`
}

type Client interface {
	Stocks(ctx context.Context, sku uint32) (ResponseStocks, error)
	CreateOrder(ctx context.Context, user int64, items []RequestCreateOrderItem) (ResponseCreateOrder, error)
}
