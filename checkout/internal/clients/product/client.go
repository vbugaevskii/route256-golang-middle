package product

import (
	"context"
	pbproduct "route256/checkout/pkg/product"
	"route256/libs/cache/ttlcache"
	"route256/libs/ratelimiter"

	"google.golang.org/grpc"
)

type RequestGetProduct struct {
	Token string
	SKU   uint32
}

type ResponseGetProduct struct {
	Name  string
	Price uint32
}

type RequestListSkus struct {
	Token         string
	StartAfterSku uint32
	Count         uint32
}

type ResponseListSkus struct {
	SKUList []uint32
}

type RateLimiter interface {
	Acquire()
}

type Cache[K comparable, V any] interface {
	Add(key K, value V) bool
	Get(key K) (V, bool)
	Remove(key K) bool
	Contains(key K) bool
	Len() int
}

type ProductService struct {
	token  string
	client pbproduct.ProductServiceClient
	rate   RateLimiter
	cache  Cache[uint32, ResponseGetProduct]
}

func NewProductClient(con grpc.ClientConnInterface, token string, rps int, cacheSize, cacheTtl int) *ProductService {
	return &ProductService{
		token:  token,
		client: pbproduct.NewProductServiceClient(con),
		rate:   ratelimiter.NewRateLimiter(rps),
		cache:  ttlcache.NewTTLCache[uint32, ResponseGetProduct](cacheSize, cacheTtl),
	}
}

func (cli *ProductService) GetProduct(ctx context.Context, sku uint32) (ResponseGetProduct, error) {
	if cli.cache != nil {
		res, exists := cli.cache.Get(sku)
		if exists {
			return res, nil
		}
	}

	cli.rate.Acquire()

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

	if cli.cache != nil {
		cli.cache.Add(sku, res)
	}

	return res, nil
}

func (cli *ProductService) ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (ResponseListSkus, error) {
	cli.rate.Acquire()

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
