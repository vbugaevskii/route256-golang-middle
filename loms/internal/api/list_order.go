package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"
)

type StatusType string

const (
	New             StatusType = "new"
	AwaitingPayment StatusType = "awaiting payment"
	Failed          StatusType = "failed"
	Payed           StatusType = "payed"
	Cancelled       StatusType = "cancelled"
)

func (s *Service) ListOrder(ctx context.Context, req *loms.RequestListOrder) (*loms.ResponseListOrder, error) {
	log.Printf("%+v\n", req)

	if req.OrderID == 0 {
		return nil, ErrOrderNotFound
	}

	_, _ = s.model.ListOrder(ctx, req.OrderID)
	// if err != nil {
	// 	return nil, err
	// }

	return &loms.ResponseListOrder{
		Status: string(New),
		User:   42,
		Items: []*loms.ResponseListOrder_OrderItem{
			{Sku: 1, Count: 200},
		},
	}, nil
}
