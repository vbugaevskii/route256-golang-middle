package api

import (
	"context"
	"log"
	"route256/checkout/pkg/checkout"
)

func (s *Service) ListCart(ctx context.Context, req *checkout.RequestListCart) (*checkout.ResponseListCart, error) {
	log.Printf("%+v\n", req)

	if req.User == 0 {
		return nil, ErrUserNotFound
	}

	items, err := s.model.ListCart(ctx, req.User)
	if err != nil {
		return nil, err
	}

	cart := checkout.ResponseListCart{}
	cart.Items = make([]*checkout.CartItem, 0, len(items))

	for _, item := range items {
		cart.Items = append(cart.Items, &checkout.CartItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
		cart.TotalPrice += item.Price * uint32(item.Count)
	}

	return &cart, nil
}
