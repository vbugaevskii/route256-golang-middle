package ordersreservations

import (
	"context"
	"fmt"
	"route256/libs/logger"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	tx.Manager
}

func NewOrdersReservationsRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{tx.Manager{Pool: pool}}
}

const (
	TableName = "orders_reservations"

	ColumnOrderId     = "order_id"
	ColumnWarehouseId = "warehouse_id"
	ColumnSKU         = "sku"
	ColumnCount       = "count"
)

func (r *Repository) ListOrderReservations(ctx context.Context, orderId int64) ([]domain.OrdersReservationsItem, error) {
	query := sq.
		Select(ColumnOrderId, ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListOrderReservations: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.OrdersReservationsItem
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query ListOrderReservations: %s", err)
	}

	return converter.ConvOrdersReservationsSchemaDomain(result), nil
}

func (r *Repository) InsertOrderReservations(ctx context.Context, orderId int64, items []domain.OrdersReservationsItem) error {
	query := sq.
		Insert(TableName).
		Columns(ColumnOrderId, ColumnWarehouseId, ColumnSKU, ColumnCount)

	for _, item := range items {
		query = query.Values(
			orderId,
			item.WarehouseId,
			item.Sku,
			item.Count,
		)
	}

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query InsertOrderReservations: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query for InsertOrderReservations: %s", err)
	}

	return nil
}

func (r *Repository) ListSkuReservations(ctx context.Context, sku uint32) ([]domain.OrdersReservationsItem, error) {
	query := sq.
		Select(ColumnOrderId, ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnSKU: sku}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for ListSkuReservations: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.OrdersReservationsItem
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query for ListSkuReservations: %s", err)
	}

	return converter.ConvOrdersReservationsSchemaDomain(result), nil
}

func (r *Repository) DeleteOrderReservations(ctx context.Context, orderId int64) error {
	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnOrderId: orderId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query DeleteOrderReservations: %s", err)
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query for DeleteOrderReservations: %s", err)
	}

	return nil
}
