package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) SetOrderStatus(ctx context.Context, orderId int64, status domain.StatusType) error {
	query := sq.
		Update(TableName).
		Set(ColumnStatus, converter.ConvStatusDomainSchema(status)).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query orders.SetOrderStatus: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query orders.SetOrderStatus: %s", err)
	}

	return nil
}
