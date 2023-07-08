package domain

import (
	"context"
	cliloms "route256/checkout/internal/clients/loms"
	"route256/libs/logger"
)

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {
	cart, err := m.cartItems.ListCart(ctx, user)
	logger.Infof("Checkout.ListCart: %+v", cart)
	if err != nil {
		return 0, err
	}

	items := make([]cliloms.RequestCreateOrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, cliloms.RequestCreateOrderItem{
			SKU:   item.SKU,
			Count: uint64(item.Count),
		})
	}

	res, err := m.loms.CreateOrder(ctx, user, items)
	logger.Infof("LOMS.CreateOrder: %+v", res)
	if err != nil {
		return 0, err
	}

	err = m.cartItems.DeleteCart(ctx, user)
	if err != nil {
		return 0, err
	}

	return res.OrderId, nil
}
