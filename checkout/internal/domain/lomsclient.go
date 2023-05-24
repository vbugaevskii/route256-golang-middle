package domain

import (
	"context"
	"net/http"
	"route256/libs/cliwrapper"
	"route256/loms/external/client"
)

type LomsClient struct {
	StocksHandler      *cliwrapper.Wrapper[*client.RequestStocks, client.ResponseStocks]
	CreateOrderHandler *cliwrapper.Wrapper[*client.RequestCreateOrder, client.ResponseCreateOrder]
}

func NewLomsClient(netloc string) *LomsClient {
	return &LomsClient{
		StocksHandler: cliwrapper.New[*client.RequestStocks, client.ResponseStocks](
			netloc,
			"/stocks",
			http.MethodPost,
		),
		CreateOrderHandler: cliwrapper.New[*client.RequestCreateOrder, client.ResponseCreateOrder](
			netloc,
			"/createOrder",
			http.MethodPost,
		),
	}
}

func (cli *LomsClient) Stocks(ctx context.Context, sku uint32) (client.ResponseStocks, error) {
	req := client.RequestStocks{
		SKU: sku,
	}
	return cli.StocksHandler.Retrieve(ctx, &req)
}

func (cli *LomsClient) CreateOrder(
	ctx context.Context,
	user int64,
	items []client.RequestCreateOrderItem,
) (client.ResponseCreateOrder, error) {
	req := client.RequestCreateOrder{
		User:  user,
		Items: items,
	}
	return cli.CreateOrderHandler.Retrieve(ctx, &req)
}