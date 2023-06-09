package cartitems

import (
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
