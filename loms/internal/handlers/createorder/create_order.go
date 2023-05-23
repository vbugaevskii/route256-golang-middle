package createorder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (r *Request) Prepare(ctx context.Context, netloc string) (*http.Request, error) {
	reqBytes, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	reqHttp, err := http.NewRequestWithContext(ctx, http.MethodPost, netloc+Endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	return reqHttp, nil
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
