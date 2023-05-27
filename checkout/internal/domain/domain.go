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

type Model struct {
	Loms    LomsClient
	Product ProductClient
}

func New(loms LomsClient, product ProductClient) *Model {
	return &Model{
		Loms:    loms,
		Product: product,
	}
}

var (
	ErrProductInsufficient = errors.New("product insufficient")
)

type CartItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"uint32"`
}

func (m *Model) ListCart(ctx context.Context, user int64) ([]CartItem, error) {
	// TODO: There should be a call to DB to retrieve items from the cart

	product, err := m.Product.GetProduct(ctx, 773297411)
	log.Printf("Product.GetProduct: %+v\n", product)
	if err != nil {
		return nil, err
	}

	return []CartItem{
		{
			SKU:   773297411,
			Count: 2,
			Name:  product.Name,
			Price: product.Price,
		},
	}, nil
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

	res, err := m.Loms.CreateOrder(ctx, user, items)
	log.Printf("LOMS.createOrder: %+v", res)
	if err != nil {
		return 0, err
	}

	return res.OrderId, nil
}

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := m.Loms.Stocks(ctx, sku)
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
