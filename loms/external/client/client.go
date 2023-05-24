package client

import (
	"context"
	"route256/libs/cliwrapper"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/stocks"
)

type Client struct {
	StocksHandler      *cliwrapper.Wrapper[*stocks.Request, stocks.Response]
	CreateOrderHandler *cliwrapper.Wrapper[*createorder.Request, createorder.Response]
}

func New(netloc string) *Client {
	return &Client{
		StocksHandler:      cliwrapper.New[*stocks.Request, stocks.Response](netloc),
		CreateOrderHandler: cliwrapper.New[*createorder.Request, createorder.Response](netloc),
	}
}

func (c *Client) Stocks(ctx context.Context, sku uint32) (stocks.Response, error) {
	req := stocks.Request{
		SKU: sku,
	}
	return c.StocksHandler.Retrieve(ctx, &req)
}

type CreateOrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint64 `json:"count"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64, items []CreateOrderItem) (createorder.Response, error) {
	itemsCast := make([]createorder.OrderItem, 0, len(items))
	for _, v := range items {
		itemsCast = append(itemsCast, createorder.OrderItem{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}

	req := createorder.Request{
		User:  user,
		Items: itemsCast,
	}
	return c.CreateOrderHandler.Retrieve(ctx, &req)
}
