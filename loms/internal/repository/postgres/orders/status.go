package orders

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
)

func ConvStatusDomainSchema(statusDomain domain.StatusType) schema.StatusType {
	var statusSchema schema.StatusType

	switch statusDomain {
	case domain.New:
		statusSchema = schema.New
	case domain.AwaitingPayment:
		statusSchema = schema.AwaitingPayment
	case domain.Failed:
		statusSchema = schema.Failed
	case domain.Payed:
		statusSchema = schema.Payed
	case domain.Cancelled:
		statusSchema = schema.Cancelled
	}

	return statusSchema
}

func ConvStatusSchemaDomain(statusSchema schema.StatusType) domain.StatusType {
	var statusDomain domain.StatusType

	switch statusSchema {
	case schema.New:
		statusDomain = domain.New
	case schema.AwaitingPayment:
		statusDomain = domain.AwaitingPayment
	case schema.Failed:
		statusDomain = domain.Failed
	case schema.Payed:
		statusDomain = domain.Payed
	case schema.Cancelled:
		statusDomain = domain.Cancelled
	}

	return statusDomain
}

func (r *Repository) SetOrderStatus(ctx context.Context, orderId int64, status domain.StatusType) error {
	query := sq.
		Update(TableNameOrders).
		Set("status", ConvStatusDomainSchema(status)).
		Where(sq.Eq{"order_id": orderId})

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
