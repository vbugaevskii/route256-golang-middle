package stocks

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) Stocks(ctx context.Context, sku uint32) ([]domain.StocksItem, error) {
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
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query stocks.Stocks: %s", err)
	}

	return ConvertStocksItems(result), nil
}

func ConvertStocksItems(itemsSchema []schema.StocksItem) []domain.StocksItem {
	itemsDomain := make([]domain.StocksItem, 0, len(itemsSchema))

	for _, item := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.StocksItem{
			WarehouseId: item.WarehouseId,
			Count:       item.Count,
		})
	}

	return itemsDomain
}
