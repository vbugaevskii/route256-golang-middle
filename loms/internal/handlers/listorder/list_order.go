package listorder

import (
	"context"
	"errors"
	"log"
)

type Handler struct {
}

type Request struct {
	OrderId int64 `json:"orderID"`
}

type StatusType string

const (
	New             StatusType = "new"
	AwaitingPayment StatusType = "awaiting payment"
	Failed          StatusType = "failed"
	Payed           StatusType = "payed"
	Cancelled       StatusType = "cancelled"
)

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint64 `json:"count"`
}

type Response struct {
	Status StatusType  `json:"status"`
	User   int64       `json:"user"`
	Items  []OrderItem `json:"items"`
}

var (
	ErrOrderNotFound = errors.New("order is not found")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v\n", req)

	if req.OrderId == 0 {
		return Response{}, ErrOrderNotFound
	}

	// TODO: add communication with product-service

	return Response{
		Status: New,
		User:   42,
		Items: []OrderItem{
			{SKU: 1, Count: 200},
		},
	}, nil
}
