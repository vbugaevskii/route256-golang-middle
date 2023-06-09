package ordersreservations

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CancelOrder(ctx context.Context, orderId int64) error {
	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query orders_reservations.CancelOrder: %s", err)
	}

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
