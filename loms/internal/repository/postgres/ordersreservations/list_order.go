package ordersreservations

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) ListOrder(ctx context.Context, orderId int64) ([]domain.OrdersReservationsItem, error) {
	query := sq.
		Select(ColumnOrderId, ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

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

	return converter.ConvOrdersReservationsSchemaDomain(result), nil
}
