package listcart

import (
	"context"
	"errors"
	"log"
)

type Handler struct {
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
	log.Printf("%+v", req)

	if req.User == 0 {
		return Response{}, ErrUserNotFound
	}

	// TODO: add communication with ProductService

	return Response{
		Items: []CartItem{
			{
				SKU:   12,
				Count: 2,
				Name:  "Молоко Домик в деревне",
				Price: 8500,
			},
		},
		TotalPrice: 17000,
	}, nil
}
