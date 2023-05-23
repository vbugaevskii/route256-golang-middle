package purchase

import (
	"context"
	"errors"
	"log"
)

const Endpoint = "/purchase"

type Handler struct {
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

	// TODO: add communication with LOMS

	return Response{}, nil
}
