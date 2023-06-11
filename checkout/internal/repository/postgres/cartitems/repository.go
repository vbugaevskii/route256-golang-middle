package cartitems

import (
	"context"
	"fmt"
	"log"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewCartItemsRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const (
	TableName = "cart_items"

	ColumnUserId = "user_id"
	ColumnSKU    = "sku"
	ColumnCount  = "count"
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
		return fmt.Errorf("build query AddToCart: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query AddToCart: %s", err)
	}

	return nil
}

func (r *Repository) ListCart(ctx context.Context, user int64) ([]domain.CartItem, error) {
	query := sq.
		Select(ColumnUserId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnUserId: user}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListCart: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.CartItem
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query ListCart: %s", err)
	}
	return converter.ConvCartItemsSchemaDomain(result), nil
}

func (r *Repository) DeleteCart(ctx context.Context, user int64) error {
	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnUserId: user})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query DeleteCart: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return fmt.Errorf("exec query for DeleteCart: %s", err)
	}

	return nil
}
