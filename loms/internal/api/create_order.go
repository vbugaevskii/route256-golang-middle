package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"
)

func (s *Service) CreateOrder(ctx context.Context, req *loms.RequestCreateOrder) (*loms.ResponseCreateOrder, error) {
	log.Printf("%+v\n", req)

	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	return &loms.ResponseCreateOrder{
		OrderID: 42,
	}, nil
}
