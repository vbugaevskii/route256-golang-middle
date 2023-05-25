package loms

import (
	"context"
	"net/http"
	"route256/checkout/internal/config"
	"route256/libs/cliwrapper"
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

type LomsService struct {
	StocksHandler      *cliwrapper.Wrapper[*RequestStocks, ResponseStocks]
	CreateOrderHandler *cliwrapper.Wrapper[*RequestCreateOrder, ResponseCreateOrder]
}

func NewLomsClient(cfg config.ConfigService) *LomsService {
	return &LomsService{
		StocksHandler: cliwrapper.New[*RequestStocks, ResponseStocks](
			cfg.Netloc,
			"/stocks",
			http.MethodPost,
		),
		CreateOrderHandler: cliwrapper.New[*RequestCreateOrder, ResponseCreateOrder](
			cfg.Netloc,
			"/createOrder",
			http.MethodPost,
		),
	}
}

func (cli *LomsService) Stocks(ctx context.Context, sku uint32) (ResponseStocks, error) {
	req := RequestStocks{
		SKU: sku,
	}
	return cli.StocksHandler.Retrieve(ctx, &req)
}

func (cli *LomsService) CreateOrder(
	ctx context.Context,
	user int64,
	items []RequestCreateOrderItem,
) (ResponseCreateOrder, error) {
	req := RequestCreateOrder{
		User:  user,
		Items: items,
	}
	return cli.CreateOrderHandler.Retrieve(ctx, &req)
}
