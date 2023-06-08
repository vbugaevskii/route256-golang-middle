package schema

import (
	"database/sql"
	"time"
)

type CartItem struct {
	UserId    int64        `db:"user_id"`
	SKU       uint32       `db:"sku"`
	Count     int16        `db:"count"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
