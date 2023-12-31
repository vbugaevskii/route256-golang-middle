package api

import (
	"context"
	"route256/checkout/pkg/checkout"
	"route256/libs/logger"
)

func (s *Service) Purchase(ctx context.Context, req *checkout.RequestPurchase) (*checkout.ResponsePurchase, error) {
	logger.Infof("%+v", req)

	if req.User == 0 {
		return nil, ErrUserNotFound
	}

	res, err := s.model.Purchase(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &checkout.ResponsePurchase{OrderID: res}, nil
}
