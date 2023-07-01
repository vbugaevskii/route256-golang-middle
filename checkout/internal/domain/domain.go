package domain

import (
	"context"
	"errors"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/repository/postgres/cartitems"
)

//go:generate ${GOBIN}/mockery --filename loms_mock.go --name LomsClient
type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) (cliloms.ResponseStocks, error)
	CreateOrder(ctx context.Context, user int64, items []cliloms.RequestCreateOrderItem) (cliloms.ResponseCreateOrder, error)
}

//go:generate ${GOBIN}/mockery --filename product_mock.go --name ProductClient
type ProductClient interface {
	GetProduct(ctx context.Context, sku uint32) (cliproduct.ResponseGetProduct, error)
	ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (cliproduct.ResponseListSkus, error)
}

//go:generate ${GOBIN}/mockery --filename cart_items_mock.go --name CartItemsRepository
type CartItemsRepository interface {
	ListCart(ctx context.Context, user int64) (cartitems.ResponseListCart, error)
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
