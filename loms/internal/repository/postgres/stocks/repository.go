package stocks

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tracing"
	tx "route256/libs/txmanager/postgres"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
)

type Repository struct {
	tx.Manager
}

func NewStocksRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{tx.Manager{Pool: pool}}
}

const (
	TableName = "stocks"

	ColumnWarehouseId = "warehouse_id"
	ColumnSKU         = "sku"
	ColumnCount       = "count"
)

func (r *Repository) ListStocks(ctx context.Context, sku uint32) ([]domain.StocksItem, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/stocks/list_stocks")
	defer span.Finish()

	query := sq.
		Select(ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnSKU: sku}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query ListStocks: %s", err))
		return nil, err
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.StocksItem
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query ListStocks: %s", err))
		return nil, err
	}

	return converter.ConvStocksItemsSchemaDomain(result), nil
}

func (r *Repository) RemoveStocks(ctx context.Context, sku uint32, item domain.StocksItem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "loms/stocks/remove_stocks")
	defer span.Finish()

	query := sq.
		Update(TableName).
		Set(ColumnCount, sq.ConcatExpr(ColumnCount, sq.Expr(" - ?", item.Count))).
		Where(sq.Eq{ColumnSKU: sku, ColumnWarehouseId: item.WarehouseId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("build query CancelOrder: %s", err))
	}

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query CancelOrder: %s", err))
	}

	return nil
}
