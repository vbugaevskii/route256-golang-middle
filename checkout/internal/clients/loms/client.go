package loms

import (
	"context"
	pbloms "route256/checkout/pkg/loms"

	"google.golang.org/grpc"
)

type RequestStocks struct {
	SKU uint32 `json:"sku"`
}

type ResponseStockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type ResponseStocks struct {
	Stocks []ResponseStockItem `json:"stocks"`
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
	client pbloms.LomsClient
}

func NewLomsClient(con grpc.ClientConnInterface) *LomsService {
	return &LomsService{
		client: pbloms.NewLomsClient(con),
	}
}

func (cli *LomsService) Stocks(ctx context.Context, sku uint32) (ResponseStocks, error) {
	reqProto := pbloms.RequestStocks{
		Sku: sku,
	}
	resProto, err := cli.client.Stocks(ctx, &reqProto)

	if err != nil {
		return ResponseStocks{}, err
	}

	res := ResponseStocks{
		Stocks: make([]ResponseStockItem, 0, len(resProto.Stocks)),
	}

	for _, item := range resProto.Stocks {
		res.Stocks = append(res.Stocks, ResponseStockItem{
			WarehouseID: item.WarehouseID,
			Count:       item.Count,
		})
	}

	return res, nil
}

func (cli *LomsService) CreateOrder(
	ctx context.Context,
	user int64,
	items []RequestCreateOrderItem,
) (ResponseCreateOrder, error) {
	reqProto := pbloms.RequestCreateOrder{
		User:  user,
		Items: make([]*pbloms.OrderItem, 0, len(items)),
	}

	for _, item := range items {
		reqProto.Items = append(reqProto.Items, &pbloms.OrderItem{
			Sku:   item.SKU,
			Count: item.Count,
		})
	}

	resProto, err := cli.client.CreateOrder(ctx, &reqProto)
	if err != nil {
		return ResponseCreateOrder{}, err
	}

	res := ResponseCreateOrder{
		OrderId: resProto.OrderID,
	}
	return res, nil
}
