package converter

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvStocksItemsSchemaDomain(itemsSchema []schema.StocksItem) []domain.StocksItem {
	itemsDomain := make([]domain.StocksItem, 0, len(itemsSchema))

	for _, item := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.StocksItem{
			WarehouseId: item.WarehouseId,
			Count:       item.Count,
		})
	}

	return itemsDomain
}
