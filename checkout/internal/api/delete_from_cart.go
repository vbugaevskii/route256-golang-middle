package api

import (
	"context"
	"log"
	"route256/checkout/pkg/checkout"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) DeleteFromCart(ctx context.Context, req *checkout.RequestDeleteFromCart) (*emptypb.Empty, error) {
	log.Printf("%+v", req)

	if req.User == 0 {
		return &emptypb.Empty{}, ErrUserNotFound
	}

	err := s.model.DeleteFromCart(ctx, req.User, req.Sku, uint16(req.Count))
	return &emptypb.Empty{}, err
}
