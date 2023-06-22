package domain

import (
	"context"
	"log"
	cliloms "route256/checkout/internal/clients/loms"
)

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {
	cart, err := m.ListCart(ctx, user)
	log.Printf("Checkout.ListCart: %+v", cart)
	if err != nil {
		return 0, err
	}

	items := make([]cliloms.RequestCreateOrderItem, 0, len(cart))
	for _, item := range cart {
		items = append(items, cliloms.RequestCreateOrderItem{
			SKU:   item.SKU,
			Count: uint64(item.Count),
		})
	}

	res, err := m.loms.CreateOrder(ctx, user, items)
	log.Printf("LOMS.CreateOrder: %+v", res)
	if err != nil {
		return 0, err
	}

	err = m.cartItems.DeleteCart(ctx, user)
	if err != nil {
		return 0, err
	}

	return res.OrderId, nil
}
