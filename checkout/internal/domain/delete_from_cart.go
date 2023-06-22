package domain

import (
	"context"
	"log"
)

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
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
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return err
	}

	for _, item := range cartItems.Items {
		err = m.cartItems.DeleteFromCart(ctx, user, item.SKU)
		if err != nil {
			return err
		}
	}

	return nil
}
