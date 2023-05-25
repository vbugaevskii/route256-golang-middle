package product

import (
	"context"
	"net/http"
	"route256/checkout/internal/config"
	"route256/libs/cliwrapper"
)

type RequestGetProduct struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ResponseGetProduct struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type RequestListSkus struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

type ResponseListSkus struct {
	SKUList []uint32 `json:"skus"`
}

type ProductService struct {
	Token string

	GetProductHandler *cliwrapper.Wrapper[*RequestGetProduct, ResponseGetProduct]
	ListSkusHandler   *cliwrapper.Wrapper[*RequestListSkus, ResponseListSkus]
}

func NewProduct(cfg config.ConfigService) *ProductService {
	return &ProductService{
		Token: cfg.Token,
		GetProductHandler: cliwrapper.New[*RequestGetProduct, ResponseGetProduct](
			cfg.Netloc,
			"/get_product",
			http.MethodPost,
		),
		ListSkusHandler: cliwrapper.New[*RequestListSkus, ResponseListSkus](
			cfg.Netloc,
			"/list_skus",
			http.MethodPost,
		),
	}
}

func (cli *ProductService) GetProduct(ctx context.Context, sku uint32) (ResponseGetProduct, error) {
	req := RequestGetProduct{
		Token: cli.Token,
		SKU:   sku,
	}
	return cli.GetProductHandler.Retrieve(ctx, &req)
}

func (cli *ProductService) ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (ResponseListSkus, error) {
	req := RequestListSkus{
		Token:         cli.Token,
		StartAfterSku: startAfterSku,
		Count:         count,
	}
	return cli.ListSkusHandler.Retrieve(ctx, &req)
}
