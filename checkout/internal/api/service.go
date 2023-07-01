package api

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout"
)

//go:generate ${GOBIN}/mockery --filename impl_mock.go --name Impl
type Impl interface {
	ListCart(ctx context.Context, user int64) ([]*domain.CartItem, error)
	Purchase(ctx context.Context, user int64) (int64, error)
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
}

type Service struct {
	checkout.UnimplementedCheckoutServer
	model Impl
}

func NewService(model Impl) *Service {
	return &Service{model: model}
}
