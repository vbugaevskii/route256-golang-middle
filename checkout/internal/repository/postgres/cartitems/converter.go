package cartitems

import "route256/checkout/internal/repository/schema"

func ConvCartItemsSchemaDomain(itemsSchema []schema.CartItem) []CartItem {
	itemsDomain := make([]CartItem, 0, len(itemsSchema))

	for _, v := range itemsSchema {
		itemsDomain = append(itemsDomain, CartItem{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}

	return itemsDomain
}
