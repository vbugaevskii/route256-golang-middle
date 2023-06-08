package domain

import "context"

type StocksRepository interface {
	Stocks(ctx context.Context, sku uint32) ([]StocksItem, error)
}

type OrdersRepository interface {
	ListOrder(ctx context.Context, orderId int64) (Order, error)
}

type OrdersReservationsRepository interface {
	ListOrder(ctx context.Context, orderId int64) ([]OrderItem, error)
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
	SKU         uint32
	Count       uint16
}

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]StocksItem, error) {
	return m.stocks.Stocks(ctx, sku)
}
