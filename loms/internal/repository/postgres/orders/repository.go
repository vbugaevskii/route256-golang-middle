package orders

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tracing"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type Repository struct {
	tx.Manager
}

func NewOrdersRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{tx.Manager{Pool: pool}}
}

const (
	TableName = "orders"

	ColumnOrderId = "order_id"
	ColumnUserId  = "user_id"
	ColumnStatus  = "status"

	ColumnCreatedAt = "created_at"
)

func (r *Repository) ListOrder(ctx context.Context, orderId int64) (domain.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/orders/list_order")
	defer span.Finish()

	query := sq.
		Select(ColumnOrderId, ColumnUserId, ColumnStatus).
		From(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query ListOrder: %s", err))
		return domain.Order{}, err
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	row := r.GetQuerier(ctx).QueryRow(ctx, queryRaw, queryArgs...)

	var order schema.Order
	if err := row.Scan(&order.OrderId, &order.UserId, &order.Status); err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query ListOrder: %s", err))
		return domain.Order{}, err
	}

	return converter.ConvOrderSchemaDomain(order), nil
}

func (r *Repository) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/orders/create_order")
	defer span.Finish()

	query := sq.
		Insert(TableName).
		Columns(ColumnUserId, ColumnStatus, ColumnCreatedAt).
		Values(userId, schema.StatusNew, time.Now()).
		Suffix(fmt.Sprintf("RETURNING %s", ColumnOrderId))

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query CreateOrder: %s", err))
		return 0, err
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	row := r.GetQuerier(ctx).QueryRow(ctx, queryRaw, queryArgs...)

	var orderId int64
	if err := row.Scan(&orderId); err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query CreateOrder: %s", err))
		return 0, err
	}

	return orderId, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderId int64, status domain.StatusType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/orders/update_order_status")
	defer span.Finish()

	query := sq.
		Update(TableName).
		Set(ColumnStatus, converter.ConvStatusDomainSchema(status)).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("build query UpdateOrderStatus: %s", err))
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query UpdateOrderStatus: %s", err))
	}

	return nil
}

func (r *Repository) ListOrderOutdated(ctx context.Context) ([]domain.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/orders/list_order_outdated")
	defer span.Finish()

	query := sq.
		Select(ColumnOrderId, ColumnUserId, ColumnStatus, ColumnCreatedAt).
		From(TableName).
		Where(sq.Eq{ColumnStatus: schema.StatusAwaitingPayment}).
		Where("created_at < now() - interval '10 minutes'")

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query ListOrder: %s", err))
		return nil, err
	}

	// logger.Debugf("SQL: %s", queryRaw)
	// logger.Debugf("SQL: %+v", queryArgs)

	var ordersSchema []schema.Order
	if err := pgxscan.Select(ctx, r.GetQuerier(ctx), &ordersSchema, queryRaw, queryArgs...); err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query ListOrder: %s", err))
		return nil, err
	}

	ordersDomain := make([]domain.Order, 0, len(ordersSchema))
	for _, order := range ordersSchema {
		ordersDomain = append(ordersDomain, converter.ConvOrderSchemaDomain(order))
	}
	return ordersDomain, nil
}
