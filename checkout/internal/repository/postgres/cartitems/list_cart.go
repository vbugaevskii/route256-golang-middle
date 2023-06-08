package cartitems

import (
	"context"
	"fmt"
	"log"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) ListCart(ctx context.Context, user int64) ([]domain.CartItem, error) {
	query := sq.
		Select("user_id", "sku", "count", "created_at", "updated_at", "deleted_at").
		From(TableNameCartItems).
		Where(sq.Eq{"user_id": user, "deleted_at": nil})

	queryRaw, queryArgs, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for filter: %s", err)
	}

	log.Printf("SQL: %s\n", queryRaw)
	log.Printf("SQL: %+v\n", queryArgs)

	var result []schema.CartItem
	err = pgxscan.Select(ctx, r.pool, &result, queryRaw, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec query for filter: %s", err)
	}
	return ConvertCartItems(result), nil
}

func ConvertCartItems(itemsSchema []schema.CartItem) []domain.CartItem {
	itemsDomain := make([]domain.CartItem, 0, len(itemsSchema))

	for _, v := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.CartItem{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}

	return itemsDomain
}
