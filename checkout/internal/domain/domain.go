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
}

type Model struct {
	loms          LomsClient
	product       ProductClient
	cartItemsRepo CartItemsRepository
}

func New(loms LomsClient, product ProductClient, cartItemsRepo CartItemsRepository) *Model {
	return &Model{
		loms:          loms,
		product:       product,
		cartItemsRepo: cartItemsRepo,
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
	cartItems, err := m.cartItemsRepo.ListCart(ctx, user)
	log.Printf("CartItems.ListCart: %+v\n", cartItems)
	if err != nil {
		return nil, err
	}

	product, err := m.product.GetProduct(ctx, 773297411)
	log.Printf("Product.GetProduct: %+v\n", product)
	if err != nil {
		return nil, err
	}

	cartItems = append(cartItems, CartItem{
		SKU:   773297411,
		Count: 2,
		Name:  product.Name,
		Price: product.Price,
	})

	return cartItems, nil
}

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {
	cart, err := m.ListCart(ctx, user)
	log.Printf("Checkout.listcart: %+v", cart)
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
	log.Printf("LOMS.createOrder: %+v", res)
	if err != nil {
		return 0, err
	}

	return res.OrderId, nil
}

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := m.loms.Stocks(ctx, sku)
	log.Printf("LOMS.stocks: %+v", stocks)
	if err != nil {
		return err
	}

	var countTotal uint64
	for _, stock := range stocks.Stocks {
		countTotal += stock.Count
		if countTotal >= uint64(count) {
			break
		}
	}

	if countTotal < uint64(count) {
		return ErrProductInsufficient
	}

	return nil
}
