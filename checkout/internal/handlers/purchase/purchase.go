package purchase

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
	"route256/loms/external/client"
)

const Endpoint = "/purchase"

type Handler struct {
	Model *domain.Model
}

type Request struct {
	User int64 `json:"user"`
}

type Response struct {
	OrderId int64 `json:"orderID"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCartIsEmpty  = errors.New("cart is empty")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)

	if req.User == 0 {
		return Response{}, ErrUserNotFound
	}

	// TODO: Go to listCart, then pass items to createOrder
	items := make([]client.CreateOrderItem, 0)
	items = append(items, client.CreateOrderItem{
		SKU:   12,
		Count: 33,
	})

	res, err := h.Model.Loms.CreateOrder(ctx, req.User, items)
	log.Printf("LOMS.createOrder: %+v", res)
	if err != nil {
		return Response{}, nil
	}

	return Response{OrderId: res.OrderId}, nil
}
