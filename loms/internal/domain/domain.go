package domain

import (
	"context"
	"fmt"
	"log"
)

type StocksRepository interface {
	ListStocks(ctx context.Context, sku uint32) ([]StocksItem, error)
	RemoveStocks(ctx context.Context, sku uint32, item StocksItem) error
}

type OrdersRepository interface {
	ListOrder(ctx context.Context, orderId int64) (Order, error)
	CreateOrder(ctx context.Context, userId int64) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, status StatusType) error
}

type OrdersReservationsRepository interface {
	ListOrderReservations(ctx context.Context, orderId int64) ([]OrdersReservationsItem, error)
	InsertOrderReservations(ctx context.Context, orderId int64, items []OrdersReservationsItem) error
	ListSkuReservations(ctx context.Context, sku uint32) ([]OrdersReservationsItem, error)
	DeleteOrderReservations(ctx context.Context, orderId int64) error
}

type Model struct {
	stocks       StocksRepository
	orders       OrdersRepository
	reservations OrdersReservationsRepository
}

func NewModel(
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

type StatusType string

const (
	StatusNew             StatusType = "new"
	StatusAwaitingPayment StatusType = "awaiting payment"
	StatusFailed          StatusType = "failed"
	StatusPayed           StatusType = "payed"
	StatusCancelled       StatusType = "cancelled"
)

type Order struct {
	Status StatusType
	User   int64
	Items  []OrderItem
}

func (m *Model) ListOrder(ctx context.Context, orderId int64) (Order, error) {
	order, err := m.orders.ListOrder(ctx, orderId)
	if err != nil {
		return Order{}, err
	}

	itemsReserved, err := m.reservations.ListOrderReservations(ctx, orderId)
	if err != nil {
		return Order{}, err
	}

	itemsReserveredMap := make(map[uint32]uint16)
	for _, item := range itemsReserved {
		itemsReserveredMap[item.Sku] += item.Count
	}

	order.Items = make([]OrderItem, 0, len(itemsReserveredMap))
	for sku, count := range itemsReserveredMap {
		order.Items = append(order.Items, OrderItem{
			Sku:   sku,
			Count: int32(count),
		})
	}

	return order, nil
}

type StocksItem struct {
	WarehouseId int64
	Count       uint16
}

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]StocksItem, error) {
	stocksResevered, err := m.reservations.ListSkuReservations(ctx, sku)
	log.Printf("OrdersReservations.ListSkuReservations: %+v\n", stocksResevered)
	if err != nil {
		return nil, err
	}

	stocksReseveredMap := make(map[int64]uint16)
	for _, item := range stocksResevered {
		stocksReseveredMap[item.WarehouseId] += item.Count
	}

	stocks, err := m.stocks.ListStocks(ctx, sku)
	log.Printf("Stocks.ListStocks: %+v\n", stocks)
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

type OrdersReservationsItem struct {
	WarehouseId int64
	Sku         uint32
	Count       uint16
}

func (m *Model) CreateOrder(ctx context.Context, userId int64, items []OrderItem) (int64, error) {
	orderId, err := m.orders.CreateOrder(ctx, userId)
	log.Printf("Orders.CreateOrder: %+v\n", orderId)
	if err != nil {
		return orderId, err
	}

	defer func() {
		if err != nil {
			m.orders.UpdateOrderStatus(ctx, orderId, StatusFailed)
		} else {
			m.orders.UpdateOrderStatus(ctx, orderId, StatusAwaitingPayment)
		}
	}()

	itemsReservered := make([]OrdersReservationsItem, 0, len(items))
	for _, item := range items {
		var stocks []StocksItem // to make defer work

		stocks, err = m.stocks.ListStocks(ctx, item.Sku)
		log.Printf("Stocks.ListStocks: %+v\n", stocks)
		if err != nil {
			return orderId, err
		}

		countLeft := uint16(item.Count)
		for _, stock := range stocks {
			if countLeft == 0 {
				break
			}

			var countAdded uint16

			if countLeft > stock.Count {
				countAdded = stock.Count
				countLeft -= stock.Count
			} else {
				countAdded = countLeft
				countLeft = 0
			}

			itemsReservered = append(itemsReservered, OrdersReservationsItem{
				WarehouseId: stock.WarehouseId,
				Sku:         item.Sku,
				Count:       countAdded,
			})
		}

		if countLeft > 0 {
			err = fmt.Errorf("insufficent stocks; sku = %d", item.Sku)
			return orderId, err
		}
	}

	err = m.reservations.InsertOrderReservations(ctx, orderId, itemsReservered)
	if err != nil {
		return orderId, err
	}

	return orderId, nil
}

func (m *Model) CancelOrder(ctx context.Context, orderId int64) error {
	var err error

	err = m.reservations.DeleteOrderReservations(ctx, orderId)
	if err != nil {
		return err
	}

	err = m.orders.UpdateOrderStatus(ctx, orderId, StatusCancelled)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) OrderPayed(ctx context.Context, orderId int64) error {
	itemsReserved, err := m.reservations.ListOrderReservations(ctx, orderId)
	log.Printf("OrdersReservations.ListOrderReservations: %+v\n", itemsReserved)
	if err != nil {
		return err
	}

	err = m.reservations.DeleteOrderReservations(ctx, orderId)
	if err != nil {
		return err
	}

	for _, item := range itemsReserved {
		err = m.stocks.RemoveStocks(ctx, item.Sku, StocksItem{
			WarehouseId: item.WarehouseId,
			Count:       item.Count,
		})
		if err != nil {
			return err
		}
	}

	err = m.orders.UpdateOrderStatus(ctx, orderId, StatusPayed)
	if err != nil {
		return err
	}

	return nil
}
