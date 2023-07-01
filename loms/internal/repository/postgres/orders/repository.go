package orders

import (
	"context"
	"fmt"
	"log"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
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
	query := sq.
		Select(ColumnOrderId, ColumnUserId, ColumnStatus).
		From(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query ListOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	row := r.GetQuerier(ctx).QueryRow(ctx, queryRaw, queryArgs...)

	var order schema.Order
	if err := row.Scan(&order.OrderId, &order.UserId, &order.Status); err != nil {
		return domain.Order{}, fmt.Errorf("exec query ListOrder: %s", err)
	}

	return converter.ConvOrderSchemaDomain(order), nil
}

func (r *Repository) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	query := sq.
		Insert(TableName).
		Columns(ColumnUserId, ColumnStatus, ColumnCreatedAt).
		Values(userId, schema.StatusNew, time.Now()).
		Suffix(fmt.Sprintf("RETURNING %s", ColumnOrderId))

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query CreateOrder: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	row := r.GetQuerier(ctx).QueryRow(ctx, queryRaw, queryArgs...)

	var orderId int64
	if err := row.Scan(&orderId); err != nil {
		return 0, fmt.Errorf("exec query CreateOrder: %s", err)
	}

	return orderId, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderId int64, status domain.StatusType) error {
	query := sq.
		Update(TableName).
		Set(ColumnStatus, converter.ConvStatusDomainSchema(status)).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query UpdateOrderStatus: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query orders.SetOrderStatus: %s", err)
	}

	return nil
}

func (r *Repository) ListOrderOutdated(ctx context.Context) ([]domain.Order, error) {
	query := sq.
		Select(ColumnOrderId, ColumnUserId, ColumnStatus, ColumnCreatedAt).
		From(TableName).
		Where(sq.Eq{ColumnStatus: schema.StatusAwaitingPayment}).
		Where("created_at < now() - interval '10 minutes'")

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListOrder: %s", err)
	}

	// log.Printf("SQL: %s\n", queryRaw)
	// log.Printf("SQL: %+v\n", queryArgs)

	var ordersSchema []schema.Order
	if err := pgxscan.Select(ctx, r.GetQuerier(ctx), &ordersSchema, queryRaw, queryArgs...); err != nil {
		return nil, err
	}

	ordersDomain := make([]domain.Order, 0, len(ordersSchema))
	for _, order := range ordersSchema {
		ordersDomain = append(ordersDomain, converter.ConvOrderSchemaDomain(order))
	}
	return ordersDomain, nil
}
