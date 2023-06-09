package converter

import (
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"
)

func ConvCartItemsSchemaDomain(itemsSchema []schema.CartItem) []domain.CartItem {
	itemsDomain := make([]domain.CartItem, 0, len(itemsSchema))

	for _, v := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.CartItem{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}

	return itemsDomain
}
