package schema

import (
	"database/sql"
	"time"
)

type StocksItem struct {
	WarehouseId int64        `db:"warehouse_id"`
	SKU         uint32       `db:"sku"`
	Count       uint16       `db:"count"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}
