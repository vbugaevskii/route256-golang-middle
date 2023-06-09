package schema

type StocksItem struct {
	WarehouseId int64  `db:"warehouse_id"`
	SKU         uint32 `db:"sku"`
	Count       uint16 `db:"count"`
}
