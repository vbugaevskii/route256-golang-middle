package cartitems

import (
	"context"
	"fmt"
	"route256/checkout/internal/repository/schema"
	"route256/libs/logger"
	"route256/libs/tracing"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
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

type CartItem struct {
	SKU   uint32
	Count uint16
}

type ResponseListCart struct {
	Items []CartItem
}

func (r *Repository) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout/cartitems/add_to_cart")
	defer span.Finish()

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
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("build query AddToCart: %s", err))
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query AddToCart: %s", err))
	}

	return nil
}

func (r *Repository) DeleteFromCart(ctx context.Context, user int64, sku uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout/cartitems/delete_from_cart")
	defer span.Finish()

	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnUserId: user, ColumnSKU: sku})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("build query DeleteFromCart: %s", err))
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query DeleteFromCart: %s", err))
	}

	return nil
}

func (r *Repository) ListCart(ctx context.Context, user int64) (ResponseListCart, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout/cartitems/list_cart")
	defer span.Finish()

	query := sq.
		Select(ColumnUserId, ColumnSKU, ColumnCount).
		From(TableName).
		Where(sq.Eq{ColumnUserId: user}).
		Where(sq.Gt{ColumnCount: 0})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("build query ListCart: %s", err))
		return ResponseListCart{}, err
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	var result []schema.CartItem
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		err = tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query ListCart: %s", err))
		return ResponseListCart{}, err
	}

	return ResponseListCart{
		Items: ConvCartItemsSchemaDomain(result),
	}, nil
}

func (r *Repository) DeleteCart(ctx context.Context, user int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout/cartitems/delete_cart")
	defer span.Finish()

	query := sq.
		Delete(TableName).
		Where(sq.Eq{ColumnUserId: user})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("build query DeleteCart: %s", err))
	}

	logger.Debugf("SQL: %s", queryRaw)
	logger.Debugf("SQL: %+v", queryArgs)

	_, err = r.pool.Exec(ctx, queryRaw, queryArgs...)
	if err != nil {
		return tracing.MarkSpanWithError(ctx, fmt.Errorf("exec query DeleteCart: %s", err))
	}

	return nil
}
