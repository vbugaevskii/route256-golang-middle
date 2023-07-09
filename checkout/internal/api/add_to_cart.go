package api

import (
	"context"
	"route256/checkout/pkg/checkout"
	"route256/libs/logger"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AddToCart(ctx context.Context, req *checkout.RequestAddToCart) (*emptypb.Empty, error) {
	logger.Infof("%+v", req)

	if req.User == 0 {
		return &emptypb.Empty{}, ErrUserNotFound
	}

	err := s.model.AddToCart(ctx, req.User, req.Sku, uint16(req.Count))
	return &emptypb.Empty{}, err
}
