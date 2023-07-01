package loms

import (
	"context"
	pbloms "route256/checkout/pkg/loms"

	"google.golang.org/grpc"
)

type RequestStocks struct {
	SKU uint32
}

type ResponseStockItem struct {
	WarehouseID int64
	Count       uint64
}

type ResponseStocks struct {
	Stocks []ResponseStockItem
}

type RequestCreateOrderItem struct {
	SKU   uint32
	Count uint64
}

type RequestCreateOrder struct {
	User  int64
	Items []RequestCreateOrderItem
}

type ResponseCreateOrder struct {
	OrderId int64
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
		Items: make([]*pbloms.RequestCreateOrder_OrderItem, 0, len(items)),
	}

	for _, item := range items {
		reqProto.Items = append(reqProto.Items, &pbloms.RequestCreateOrder_OrderItem{
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

func (cli *LomsService) CancelOrder(ctx context.Context, orderId int64) error {
	reqProto := pbloms.RequestCancelOrder{
		OrderID: orderId,
	}

	_, err := cli.client.CancelOrder(ctx, &reqProto)
	if err != nil {
		return err
	}

	return nil
}
