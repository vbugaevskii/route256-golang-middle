package converter

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvOrdersReservationsSchemaDomain(itemsSchema []schema.OrdersReservationsItem) []domain.OrdersReservationsItem {
	itemsDomain := make([]domain.OrdersReservationsItem, 0, len(itemsSchema))

	for _, item := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.OrdersReservationsItem{
			WarehouseId: item.WarehouseId,
			Sku:         item.SKU,
			Count:       item.Count,
		})
	}

	return itemsDomain
}
