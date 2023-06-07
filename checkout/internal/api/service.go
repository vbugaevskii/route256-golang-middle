package api

import (
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout"
)

type Service struct {
	checkout.UnimplementedCheckoutServer
	model *domain.Model
}

func NewService(model *domain.Model) *Service {
	return &Service{model: model}
}
