package orders

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewOrdersRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const (
	TableName = "orders"

	ColumnOrderId = "order_id"
	ColumnUserId  = "user_id"
	ColumnStatus  = "status"
)
