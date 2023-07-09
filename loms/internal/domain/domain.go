package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	tx "route256/libs/txmanager/postgres"
	"time"

	"go.uber.org/zap"
)

type StocksRepository interface {
	ListStocks(ctx context.Context, sku uint32) ([]StocksItem, error)
	RemoveStocks(ctx context.Context, sku uint32, item StocksItem) error
}

type OrdersRepository interface {
	ListOrder(ctx context.Context, orderId int64) (Order, error)
	CreateOrder(ctx context.Context, userId int64) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, status StatusType) error
	ListOrderOutdated(ctx context.Context) ([]Order, error)
}

type OrdersReservationsRepository interface {
	ListOrderReservations(ctx context.Context, orderId int64) ([]OrdersReservationsItem, error)
	InsertOrderReservations(ctx context.Context, orderId int64, items []OrdersReservationsItem) error
	ListSkuReservations(ctx context.Context, sku uint32) ([]OrdersReservationsItem, error)
	DeleteOrderReservations(ctx context.Context, orderId int64) error
}

type KafkaProducer interface {
	SendOrderStatus(message Notification) error
}

type Notification struct {
	RecordId  int64      `json:"record_id"`
	OrderId   int64      `json:"order_id"`
	Status    StatusType `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type NotificationsOutboxRepository interface {
	CreateNotification(ctx context.Context, orderId int64, status StatusType) (int64, error)
	SetNotificationDelivered(ctx context.Context, recordId int64) error
	ListNotificationsWaiting(ctx context.Context) ([]Notification, error)
	DeleteNotificationsDelivered(ctx context.Context) error
}

type Model struct {
	txManager *tx.Manager

	producer      KafkaProducer
	notifications NotificationsOutboxRepository

	stocks       StocksRepository
	orders       OrdersRepository
	reservations OrdersReservationsRepository
}

func NewModel(
	txManager *tx.Manager,
	producer KafkaProducer,
	notifications NotificationsOutboxRepository,
	stocks StocksRepository,
	orders OrdersRepository,
	reservations OrdersReservationsRepository,
) *Model {
	return &Model{
		txManager: txManager,

		producer:      producer,
		notifications: notifications,

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
	OrderId int64
	Status  StatusType
	User    int64
	Items   []OrderItem
}

func (m *Model) ListOrder(ctx context.Context, orderId int64) (Order, error) {
	var order Order

	err := m.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var (
			itemsReserved []OrdersReservationsItem
			err           error
		)

		order, err = m.orders.ListOrder(ctxTx, orderId)
		logger.Infof("Orders.ListOrder: %+v", order)
		if err != nil {
			return err
		}

		itemsReserved, err = m.reservations.ListOrderReservations(ctxTx, orderId)
		logger.Infof("OrdersReservations.ListOrderReservations: %+v", itemsReserved)
		if err != nil {
			return err
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

		return nil
	})

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
	var stocks []StocksItem

	err := m.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var (
			stocksResevered []OrdersReservationsItem
			err             error
		)

		stocksResevered, err = m.reservations.ListSkuReservations(ctxTx, sku)
		logger.Infof("OrdersReservations.ListSkuReservations: %+v", stocksResevered)
		if err != nil {
			return err
		}

		stocksReseveredMap := make(map[int64]uint16)
		for _, item := range stocksResevered {
			stocksReseveredMap[item.WarehouseId] += item.Count
		}

		stocks, err = m.stocks.ListStocks(ctxTx, sku)
		logger.Infof("Stocks.ListStocks: %+v", stocks)
		if err != nil {
			return err
		}

		for i, item := range stocks {
			if cnt, exists := stocksReseveredMap[item.WarehouseId]; exists {
				if item.Count < cnt {
					return fmt.Errorf("incosistent stocks for sku=%d and warehouse_id=%d", sku, item.WarehouseId)
				}
				stocks[i].Count -= cnt
			}
		}
		logger.Infof("Stocks.ListStocks: %+v\n modified", stocks)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return stocks, nil
}

type OrdersReservationsItem struct {
	WarehouseId int64
	Sku         uint32
	Count       uint16
}

func (m *Model) CreateOrder(ctx context.Context, userId int64, items []OrderItem) (int64, error) {
	var orderId int64

	err := m.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var err error

		orderId, err = m.orders.CreateOrder(ctxTx, userId)
		logger.Infof("Orders.CreateOrder: %+v", orderId)
		if err != nil {
			return err
		}

		if _, err = m.notifications.CreateNotification(ctxTx, orderId, StatusNew); err != nil {
			logger.Error("failed Notifications.CreateNotification", zap.Error(err))
		}

		defer func() {
			if err == nil {
				return
			}

			if err = m.orders.UpdateOrderStatus(ctxTx, orderId, StatusFailed); err != nil {
				logger.Error("failed Orders.UpdateOrderStatus", zap.Error(err))
				return
			}

			if _, err = m.notifications.CreateNotification(ctxTx, orderId, StatusFailed); err != nil {
				logger.Error("failed Notifications.CreateNotification", zap.Error(err))
				return
			}
		}()

		itemsReservered := make([]OrdersReservationsItem, 0, len(items))
		for _, item := range items {
			var stocks []StocksItem // to make defer work

			stocks, err = m.stocks.ListStocks(ctxTx, item.Sku)
			logger.Infof("Stocks.ListStocks: %+v", stocks)
			if err != nil {
				return err
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
				return err
			}
		}

		err = m.reservations.InsertOrderReservations(ctxTx, orderId, itemsReservered)
		if err != nil {
			return err
		}

		err = m.orders.UpdateOrderStatus(ctxTx, orderId, StatusAwaitingPayment)
		if err != nil {
			return err
		}

		if _, err = m.notifications.CreateNotification(ctxTx, orderId, StatusAwaitingPayment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return orderId, err
	}

	return orderId, nil
}

func (m *Model) CancelOrder(ctx context.Context, orderId int64) error {
	err := m.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var (
			order Order
			err   error
		)

		order, err = m.orders.ListOrder(ctxTx, orderId)
		logger.Infof("Orders.ListOrder: %+v", order)
		if err != nil {
			return err
		}

		if order.Status != StatusAwaitingPayment {
			return fmt.Errorf("order can be canceled cause wrong status; status = %s", order.Status)
		}

		err = m.reservations.DeleteOrderReservations(ctxTx, orderId)
		if err != nil {
			return err
		}

		err = m.orders.UpdateOrderStatus(ctxTx, orderId, StatusCancelled)
		if err != nil {
			return err
		}

		if _, err = m.notifications.CreateNotification(ctxTx, orderId, StatusCancelled); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) OrderPayed(ctx context.Context, orderId int64) error {
	var itemsReserved []OrdersReservationsItem

	err := m.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var (
			order Order
			err   error
		)

		order, err = m.orders.ListOrder(ctxTx, orderId)
		if err != nil {
			return err
		}

		if order.Status != StatusAwaitingPayment {
			return fmt.Errorf("order can be payed cause wrong status; status = %s", order.Status)
		}

		itemsReserved, err = m.reservations.ListOrderReservations(ctxTx, orderId)
		logger.Infof("OrdersReservations.ListOrderReservations: %+v", itemsReserved)
		if err != nil {
			return err
		}

		err = m.reservations.DeleteOrderReservations(ctxTx, orderId)
		if err != nil {
			return err
		}

		for _, item := range itemsReserved {
			err = m.stocks.RemoveStocks(ctxTx, item.Sku, StocksItem{
				WarehouseId: item.WarehouseId,
				Count:       item.Count,
			})
			if err != nil {
				return err
			}
		}

		err = m.orders.UpdateOrderStatus(ctxTx, orderId, StatusPayed)
		if err != nil {
			return err
		}

		if _, err = m.notifications.CreateNotification(ctxTx, orderId, StatusPayed); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) RunCancelOrderByTimeout(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			orders, err := m.orders.ListOrderOutdated(ctx)
			if err != nil {
				return err
			}

			if orders == nil {
				continue
			}

			for _, order := range orders {
				err = m.CancelOrder(ctx, order.OrderId)
				if err != nil {
					return err
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Model) RunNotificationsSender(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			orders, err := m.notifications.ListNotificationsWaiting(ctx)
			logger.Infof("Notifications.ListNotificationsWaiting: %+v", orders)

			if err != nil {
				return err
			}

			if orders == nil {
				continue
			}

			for _, order := range orders {
				if err = m.producer.SendOrderStatus(order); err == nil {
					err = m.notifications.SetNotificationDelivered(ctx, order.RecordId)
					if err != nil {
						return err
					}
				} else {
					logger.Error("failed to write message to kafka", zap.Error(err))
				}
			}

			err = m.notifications.DeleteNotificationsDelivered(ctx)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
