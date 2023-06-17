package schema

type OrdersReservationsItem struct {
	OrderId     int64  `db:"order_id"`
	WarehouseId int64  `db:"warehouse_id"`
	SKU         uint32 `db:"sku"`
	Count       uint16 `db:"count"`
}
