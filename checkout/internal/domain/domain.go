package domain

import (
	"context"
	"errors"
	"log"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
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
	ListCart(ctx context.Context, user int64) ([]CartItem, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
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

func (m *Model) ListCart(ctx context.Context, user int64) ([]CartItem, error) {
	cartItems, err := m.cartItems.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return nil, err
	}

	for _, item := range cartItems {
		product, err := m.product.GetProduct(ctx, item.SKU)
		log.Printf("Product.GetProduct: %+v\n", product)
		if err != nil {
			return nil, err
		}

		item.Name = product.Name
		item.Price = product.Price
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

	var countTotal uint64
	for _, stock := range stocks.Stocks {
		countTotal += stock.Count
		if countTotal >= uint64(count+countInCart) {
			break
		}
	}

	if countTotal < uint64(count) {
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

	if countInCart < count {
		countInCart = 0
	} else {
		countInCart -= count
	}

	err = m.cartItems.AddToCart(ctx, user, sku, countInCart)
	if err != nil {
		return err
	}

	return nil
}
