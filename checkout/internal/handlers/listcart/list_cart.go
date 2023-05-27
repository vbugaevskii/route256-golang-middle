package listcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Request struct {
	User int64 `json:"user"`
}

type CartItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"uint32"`
}

type Response struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"totalPrice"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v\n", req)

	if req.User == 0 {
		return Response{}, ErrUserNotFound
	}

	items, err := h.Model.ListCart(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	cart := Response{}
	cart.Items = make([]CartItem, 0, len(items))

	for _, item := range items {
		cart.Items = append(cart.Items, CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
		cart.TotalPrice += item.Price * uint32(item.Count)
	}

	return cart, nil
}
