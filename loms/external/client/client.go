package client

import (
	"context"
	"route256/libs/cliwarpper"
	"route256/loms/internal/handlers/stocks"
)

type Client struct {
	StocksHandler *cliwarpper.Wrapper[*stocks.Request, stocks.Response]
}

func New(netloc string) *Client {
	return &Client{
		StocksHandler: cliwarpper.New[*stocks.Request, stocks.Response](netloc),
	}
}

func (c *Client) Stocks(ctx context.Context, sku uint32) (stocks.Response, error) {
	req := stocks.Request{
		SKU: sku,
	}
	return c.StocksHandler.Retrieve(ctx, &req)
}
