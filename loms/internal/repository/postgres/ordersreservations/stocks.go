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

func (r *Repository) Stocks(ctx context.Context, sku uint32) ([]domain.OrdersReservationsItem, error) {
	query := sq.
		Select(ColumnOrderId, ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnSKU: sku}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for orders_reservations.Stocks: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.OrdersReservationsItem
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query for orders_reservations.Stocks: %s", err)
	}

	return converter.ConvOrdersReservationsSchemaDomain(result), nil
}
