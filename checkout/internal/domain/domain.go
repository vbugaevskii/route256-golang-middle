package domain

import (
	"context"
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

func New(loms LomsClient, productService ProductClient) *Model {
	return &Model{
		Loms:    loms,
		Product: productService,
	}
}
