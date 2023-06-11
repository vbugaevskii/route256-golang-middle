package stocks

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"context"
	"fmt"
	"log"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/postgres/tx"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
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
	query := sq.
		Select(ColumnWarehouseId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnSKU: sku}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query stocks.Stocks: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.StocksItem
	err = pgxscan.Select(ctx, r.GetQuerier(ctx), &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query stocks.Stocks: %s", err)
	}

	return converter.ConvStocksItemsSchemaDomain(result), nil
}

func (r *Repository) RemoveStocks(ctx context.Context, sku uint32, item domain.StocksItem) error {
	query := sq.
		Update(TableName).
		Set(ColumnCount, sq.ConcatExpr(ColumnCount, sq.Expr(" - ?", item.Count))).
		Where(sq.Eq{ColumnSKU: sku, ColumnWarehouseId: item.WarehouseId})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query orders_reservations.CancelOrder: %s", err)
	}

	_, err = r.GetQuerier(ctx).Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
