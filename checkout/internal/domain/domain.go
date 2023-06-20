package domain

import (
	"context"
	"errors"
	"log"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	wp "route256/libs/workerpool"
)

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) (cliloms.ResponseStocks, error)
	CreateOrder(ctx context.Context, user int64, items []cliloms.RequestCreateOrderItem) (cliloms.ResponseCreateOrder, error)
}

type ProductClient interface {
	GetProduct(ctx context.Context, sku uint32) (cliproduct.ResponseGetProduct, error)
	ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (cliproduct.ResponseListSkus, error)
}

type CartItemsRepository interface {
	ListCart(ctx context.Context, user int64) ([]*CartItem, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32) error
	DeleteCart(ctx context.Context, user int64) error
}

type Model struct {
	loms      LomsClient
	product   ProductClient
	cartItems CartItemsRepository
}

func New(loms LomsClient, product ProductClient, cartItems CartItemsRepository) *Model {
	return &Model{
		loms:      loms,
		product:   product,
		cartItems: cartItems,
	}
}

var (
	ErrProductInsufficient = errors.New("product insufficient")
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

const numProductWorkers = 5

func (m *Model) ListCart(ctx context.Context, user int64) ([]*CartItem, error) {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return nil, err
	}

	// prepare context for worker pool
	ctxPool, cancelPool := context.WithCancel(ctx)

	pool := wp.NewWorkerPool(
		ctxPool,
		numProductWorkers,
		// will change cartItems inplace
		func(item *CartItem) (struct{}, error) {
			product, err := m.product.GetProduct(ctx, item.SKU)
			log.Printf("Product.GetProduct: %+v\n", product)
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

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return err
	}

	var countInCart uint16
	for _, item := range cartItems {
		if item.SKU == sku {
			countInCart += item.Count
		}
	}

	stocks, err := m.loms.Stocks(ctx, sku)
	log.Printf("LOMS.Stocks: %+v", stocks)
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

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return err
	}

	var countInCart uint16
	for _, item := range cartItems {
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
