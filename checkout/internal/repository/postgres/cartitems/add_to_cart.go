package cartitems

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	query := sq.
		Insert(TableName).
		Columns(ColumnUserId, ColumnSKU, ColumnCount).
		Values(user, sku, count).
		Suffix(fmt.Sprintf(
			"ON CONFLICT (%s, %s) DO UPDATE SET %s = EXCLUDED.%s",
			ColumnUserId, ColumnSKU, ColumnCount, ColumnCount,
		))

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query cart_items.AddToCart: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
