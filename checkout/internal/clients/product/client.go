package product

import (
	"context"
	pbproduct "route256/checkout/pkg/product"

	"google.golang.org/grpc"
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
	token  string
	client pbproduct.ProductServiceClient
}

func NewProductClient(con grpc.ClientConnInterface, token string) *ProductService {
	return &ProductService{
		token:  token,
		client: pbproduct.NewProductServiceClient(con),
	}
}

func (cli *ProductService) GetProduct(ctx context.Context, sku uint32) (ResponseGetProduct, error) {
	reqProto := pbproduct.RequestGetProduct{
		Token: cli.token,
		Sku:   sku,
	}

	resProto, err := cli.client.GetProduct(ctx, &reqProto)
	if err != nil {
		return ResponseGetProduct{}, err
	}

	res := ResponseGetProduct{
		Name:  resProto.Name,
		Price: resProto.Price,
	}
	return res, nil
}

func (cli *ProductService) ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (ResponseListSkus, error) {
	reqProto := pbproduct.RequestListSkus{
		Token:         cli.token,
		StartAfterSku: startAfterSku,
		Count:         count,
	}

	resProto, err := cli.client.ListSkus(ctx, &reqProto)
	if err != nil {
		return ResponseListSkus{}, err
	}

	res := ResponseListSkus{
		SKUList: make([]uint32, 0, len(resProto.Skus)),
	}
	res.SKUList = append(res.SKUList, resProto.Skus...)
	return res, nil
}
