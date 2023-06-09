package ordersreservations

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateOrder(ctx context.Context, orderId int64, items []domain.OrdersReservationsItem) error {
	query := sq.
		Insert(TableName).
		Columns("order_id", "warehouse_id", "sku", "count")

	for _, item := range items {
		query = query.Values(
			orderId,
			item.WarehouseId,
			item.Sku,
			item.Count,
		)
	}

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query orders_reservations.CreateOrder: %s", err)
	}

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
