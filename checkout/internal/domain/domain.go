package domain

import (
	cliloms "route256/checkout/internal/clients/loms"
	clips "route256/checkout/internal/clients/productservice"
)

type Model struct {
	Loms           cliloms.LomsClient
	ProductService clips.ProductServiceClient
}

func New(loms *cliloms.LomsService, productService *clips.ProductService) *Model {
	return &Model{
		Loms:           loms,
		ProductService: productService,
	}
}
