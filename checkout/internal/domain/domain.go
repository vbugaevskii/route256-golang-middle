package domain

import (
	cliLoms "route256/loms/external/client"
)

type Model struct {
	Loms           cliLoms.Client
	ProductService *ProductServiceClient
}

func New(loms *LomsClient, productService *ProductServiceClient) *Model {
	return &Model{
		Loms:           loms,
		ProductService: productService,
	}
}
