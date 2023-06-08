package stocks

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewStocksRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const TableNameStocks = "stocks"
