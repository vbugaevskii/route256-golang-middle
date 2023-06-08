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

const TableNameOrders = "orders"
