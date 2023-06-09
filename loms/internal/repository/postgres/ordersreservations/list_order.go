package ordersreservations

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) ListOrder(ctx context.Context, orderId int64) ([]domain.OrderItem, error) {
	query := sq.
		Select("order_id", "warehouse_id", "sku", "count").
		From(TableNameOrdersReservations).
		Where(sq.Eq{"order_id": orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query orders_reservations.ListOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.OrdersReservationsItem
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query orders_reservations.ListOrder: %s", err)
	}

	return ConvertOrderItems(result), nil
}

func ConvertOrderItems(itemsSchema []schema.OrdersReservationsItem) []domain.OrderItem {
	itemsDomain := make([]domain.OrderItem, 0, len(itemsSchema))

	for _, item := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.OrderItem{
			Sku:   item.SKU,
			Count: int32(item.Count),
		})
	}

	return itemsDomain
}
