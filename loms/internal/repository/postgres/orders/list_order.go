package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/api"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) ListOrder(ctx context.Context, orderId int64) (domain.Order, error) {
	query := sq.
		Select("order_id", "user_id", "status", "created_at", "updated_at").
		From(TableNameOrders).
		Where(sq.Eq{"order_id": orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query for filter: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.Order
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return domain.Order{}, fmt.Errorf("exec query for filter: %s", err)
	}
	if len(result) == 0 {
		return domain.Order{}, api.ErrOrderNotFound
	}

	return ConvertOrder(result[0]), nil
}

func ConvertOrder(orderSchema schema.Order) domain.Order {
	var status api.StatusType

	switch orderSchema.Status {
	case schema.New:
		status = api.New
	case schema.AwaitingPayment:
		status = api.AwaitingPayment
	case schema.Failed:
		status = api.Failed
	case schema.Payed:
		status = api.Payed
	case schema.Cancelled:
		status = api.Cancelled
	}

	return domain.Order{
		Status: string(status),
		User:   orderSchema.UserId,
	}
}
