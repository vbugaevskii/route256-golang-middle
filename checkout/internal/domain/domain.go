package domain

import (
	cliloms "route256/checkout/internal/clients/loms"
	clips "route256/checkout/internal/clients/product"
)

type Model struct {
	Loms    cliloms.LomsClient
	Product clips.ProductClient
}

func New(loms *cliloms.LomsService, productService *clips.ProductService) *Model {
	return &Model{
		Loms:    loms,
		Product: productService,
	}
}
