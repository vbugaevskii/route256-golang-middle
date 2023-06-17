package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CancelOrder(ctx context.Context, req *loms.RequestCancelOrder) (*emptypb.Empty, error) {
	log.Printf("%+v\n", req)

	if req.OrderID == 0 {
		return &emptypb.Empty{}, ErrOrderNotFound
	}

	err := s.model.CancelOrder(ctx, req.OrderID)
	return &emptypb.Empty{}, err
}
