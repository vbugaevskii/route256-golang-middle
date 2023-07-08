package domain

import (
	"context"
	"route256/libs/logger"
)

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	logger.Infof("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return err
	}

	var countInCart uint16
	for _, item := range cartItems.Items {
		if item.SKU == sku {
			countInCart += item.Count
		}
	}

	stocks, err := m.loms.Stocks(ctx, sku)
	logger.Infof("LOMS.Stocks: %+v", stocks)
	if err != nil {
		return err
	}

	countTotalMax := uint64(count + countInCart)
	var countTotal uint64
	for _, stock := range stocks.Stocks {
		countTotal += stock.Count
		if countTotal >= countTotalMax {
			break
		}
	}

	if countTotal < countTotalMax {
		return ErrProductInsufficient
	}

	err = m.cartItems.AddToCart(ctx, user, sku, count+countInCart)
	if err != nil {
		return err
	}

	return nil
}
