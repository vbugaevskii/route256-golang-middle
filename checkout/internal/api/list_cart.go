package api

import (
	"context"
	"route256/checkout/pkg/checkout"
	"route256/libs/logger"
)

func (s *Service) ListCart(ctx context.Context, req *checkout.RequestListCart) (*checkout.ResponseListCart, error) {
	logger.Infof("%+v", req)

	if req.User == 0 {
		return nil, ErrUserNotFound
	}

	items, err := s.model.ListCart(ctx, req.User)
	if err != nil {
		return nil, err
	}

	cart := checkout.ResponseListCart{}
	cart.Items = make([]*checkout.ResponseListCart_CartItem, 0, len(items))

	for _, item := range items {
		cart.Items = append(cart.Items, &checkout.ResponseListCart_CartItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
		cart.TotalPrice += item.Price * uint32(item.Count)
	}

	return &cart, nil
}
