package purchase

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

	res, err := h.Model.Purchase(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	return Response{OrderId: res}, nil
}
