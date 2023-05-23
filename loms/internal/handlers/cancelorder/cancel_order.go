package cancelorder

import (
	"context"
	"errors"
	"log"
)

const Endpoint = "/cancelOrder"

type Handler struct {
}

type Request struct {
	OrderId int64 `json:"orderID"`
}

type Response struct {
}

var (
	ErrOrderNotFound = errors.New("order is not found")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v\n", req)

	if req.OrderId == 0 {
		return Response{}, ErrOrderNotFound
	}

	return Response{}, nil
}
