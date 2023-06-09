package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	query := sq.
		Insert(TableNameOrders).
		Columns("user_id", "status").
		Values(userId, schema.New).
		Suffix(fmt.Sprintf("RETURNING %s", "order_id"))

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query orders.CreateOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	row := r.pool.QueryRow(ctx, queryRaw, queryArgs...)

	var orderId int64
	if err := row.Scan(&orderId); err != nil {
		return 0, fmt.Errorf("exec query orders.CreateOrder: %s", err)
	}

	return orderId, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, orderId int64, status domain.StatusType) error {
	query := sq.
		Update(TableNameOrders).
		Set("status", ConvStatusDomainSchema(status)).
		Where(sq.Eq{"order_id": orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query orders.UpdateOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query orders.UpdateOrder: %s", err)
	}

	return nil
}
