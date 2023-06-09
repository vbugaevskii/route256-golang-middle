package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) ListOrder(ctx context.Context, orderId int64) (domain.Order, error) {
	query := sq.
		Select(ColumnOrderId, ColumnUserId, ColumnStatus).
		From(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query orders.ListOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	row := r.pool.QueryRow(ctx, queryRaw, queryArgs...)

	var order schema.Order
	if err := row.Scan(&order.OrderId, &order.UserId, &order.Status); err != nil {
		return domain.Order{}, fmt.Errorf("exec query orders.ListOrder: %s", err)
	}

	return converter.ConvOrderSchemaDomain(order), nil
}
