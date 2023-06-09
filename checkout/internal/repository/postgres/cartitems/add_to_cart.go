package cartitems

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	query := sq.
		Insert(TableNameCartItems).
		Columns("user_id", "sku", "count").
		Values(user, sku, count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = EXCLUDED.count")

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query for filter: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}
