package ordersreservations

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewOrdersReservationsRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const TableNameOrdersReservations = "orders_reservations"
