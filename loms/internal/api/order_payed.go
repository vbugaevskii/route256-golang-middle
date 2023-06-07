package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) OrderPayed(ctx context.Context, req *loms.RequestOrderPayed) (*emptypb.Empty, error) {
	log.Printf("%+v\n", req)

	if req.OrderID == 0 {
		return &emptypb.Empty{}, ErrOrderNotFound
	}

	return &emptypb.Empty{}, nil
}
