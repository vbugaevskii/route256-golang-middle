package api

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain"
	"route256/loms/pkg/loms"
)

func (s *Service) CreateOrder(ctx context.Context, req *loms.RequestCreateOrder) (*loms.ResponseCreateOrder, error) {
	logger.Infof("%+v", req)

	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			Sku:   item.Sku,
			Count: int32(item.Count),
		})
	}

	orderId, err := s.model.CreateOrder(ctx, req.User, items)
	if err != nil {
		return nil, err
	}

	return &loms.ResponseCreateOrder{
		OrderID: orderId,
	}, nil
}
