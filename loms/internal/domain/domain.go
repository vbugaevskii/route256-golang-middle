package domain

import (
	"context"
	"fmt"
)

type StocksRepository interface {
	Stocks(ctx context.Context, sku uint32) ([]StocksItem, error)
}

type OrdersRepository interface {
	ListOrder(ctx context.Context, orderId int64) (Order, error)
}

type OrdersReservationsRepository interface {
	ListOrder(ctx context.Context, orderId int64) ([]OrderItem, error)
	Stocks(ctx context.Context, sku uint32) ([]StocksItem, error)
}

type Model struct {
	stocks       StocksRepository
	orders       OrdersRepository
	reservations OrdersReservationsRepository
}

func New(
	stocks StocksRepository,
	orders OrdersRepository,
	reservations OrdersReservationsRepository,
) *Model {
	return &Model{
		stocks:       stocks,
		orders:       orders,
		reservations: reservations,
	}
}

type OrderItem struct {
	Sku   uint32
	Count int32
}

type Order struct {
	Status string
	User   int64
	Items  []OrderItem
}

func (m *Model) ListOrder(ctx context.Context, orderId int64) (Order, error) {
	order, err := m.orders.ListOrder(ctx, orderId)
	if err != nil {
		return Order{}, err
	}

	order.Items, err = m.reservations.ListOrder(ctx, orderId)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

type StocksItem struct {
	WarehouseId int64
	Count       uint16
}

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]StocksItem, error) {
	// FIXME: Don't take into account reserved items

	stocksResevered, err := m.reservations.Stocks(ctx, sku)
	if err != nil {
		return nil, err
	}

	stocksReseveredMap := make(map[int64]uint16)
	for _, item := range stocksResevered {
		stocksReseveredMap[item.WarehouseId] += item.Count
	}

	stocks, err := m.stocks.Stocks(ctx, sku)
	if err != nil {
		return nil, err
	}

	for _, item := range stocks {
		if cnt, exists := stocksReseveredMap[item.WarehouseId]; exists {
			if item.Count < cnt {
				return nil, fmt.Errorf("incosistent stocks for sku=%d and warehouse_id=%d", sku, item.WarehouseId)
			}
			item.Count -= cnt
		}
	}

	return stocks, nil
}
