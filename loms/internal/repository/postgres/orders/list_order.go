package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
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

	var result []schema.Order
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return domain.Order{}, fmt.Errorf("exec query for filter: %s", err)
	}
	if len(result) == 0 {
		return domain.Order{}, fmt.Errorf("order not found")
	}

	return ConvertOrder(result[0]), nil
}

func ConvertOrder(orderSchema schema.Order) domain.Order {
	return domain.Order{
		Status: ConvStatusSchemaDomain(orderSchema.Status),
		User:   orderSchema.UserId,
	}
}
