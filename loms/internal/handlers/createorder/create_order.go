package createorder

import (
	"context"
	"errors"
	"log"
)

const Endpoint = "/createOrder"

type Handler struct {
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint64 `json:"count"`
}

type Request struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

type Response struct {
	OrderId int64 `json:"orderID"`
}

var (
	ErrEmptyOrder = errors.New("empty order")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v\n", req)

	if len(req.Items) == 0 {
		return Response{}, ErrEmptyOrder
	}

	return Response{
		OrderId: 42,
	}, nil
}
