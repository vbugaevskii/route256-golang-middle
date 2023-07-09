package domain

import (
	"context"
	"route256/libs/logger"
)

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	logger.Infof("CartItems.ListCart: %+v", cartItems)
	if err != nil {
		return err
	}

	var countInCart uint16
	for _, item := range cartItems.Items {
		if item.SKU == sku {
			countInCart += item.Count
		}
	}

	// If count to be deleted is greater than count in cart, it's OK
	// We will remove all items from the cart
	if countInCart < count {
		countInCart = 0
	} else {
		countInCart -= count
	}

	if countInCart > 0 {
		err = m.cartItems.AddToCart(ctx, user, sku, countInCart)
	} else {
		err = m.cartItems.DeleteFromCart(ctx, user, sku)
	}
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) DeleteCart(ctx context.Context, user int64) error {
	err := m.cartItems.DeleteCart(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
