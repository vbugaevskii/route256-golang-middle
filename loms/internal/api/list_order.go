package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"
)

func (s *Service) ListOrder(ctx context.Context, req *loms.RequestListOrder) (*loms.ResponseListOrder, error) {
	log.Printf("%+v\n", req)

	if req.OrderID == 0 {
		return nil, ErrOrderNotFound
	}

	order, err := s.model.ListOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	items := make([]*loms.ResponseListOrder_OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &loms.ResponseListOrder_OrderItem{
			Sku:   item.Sku,
			Count: uint64(item.Count),
		})
	}

	return &loms.ResponseListOrder{
		Status: string(order.Status),
		User:   order.User,
		Items:  items,
	}, nil
}
