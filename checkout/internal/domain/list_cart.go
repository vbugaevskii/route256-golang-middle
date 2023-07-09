package domain

import (
	"context"
	"route256/libs/logger"
	wp "route256/libs/workerpool"
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

const numProductWorkers = 5

func (m *Model) ListCart(ctx context.Context, user int64) ([]*CartItem, error) {
	cartItemsRepo, err := m.cartItems.ListCart(ctx, user)
	logger.Infof("CartItems.ListCart: %+v", cartItemsRepo)
	if err != nil {
		return nil, err
	}

	cartItems := make([]*CartItem, 0, len(cartItemsRepo.Items))
	for _, item := range cartItemsRepo.Items {
		cartItems = append(cartItems, &CartItem{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	// prepare context for worker pool
	ctxPool, cancelPool := context.WithCancel(ctx)

	pool := wp.NewWorkerPool(
		ctxPool,
		numProductWorkers,
		// will change cartItems inplace
		func(item *CartItem) (struct{}, error) {
			product, err := m.product.GetProduct(ctx, item.SKU)
			logger.Infof("Product.GetProduct: %+v", product)
			if err != nil {
				return struct{}{}, err
			}
			item.Name = product.Name
			item.Price = product.Price
			return struct{}{}, err
		},
	)

	// run reader
	result := pool.Submit(cartItems)
	go func() {
		for r := range result {
			if r.Err != nil {
				err = r.Err
				break
			}
		}

		// tell the pool we are ready
		cancelPool()
	}()

	// waiting all the tasks to be done
	pool.Close()

	if err != nil {
		return nil, err
	}

	return cartItems, nil
}
