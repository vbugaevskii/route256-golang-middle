package stocks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const Endpoint = "/stocks"

type Handler struct {
}

type Request struct {
	SKU uint32 `json:"sku"`
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

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Response struct {
	Stocks []StockItem `json:"stocks"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)

	return Response{
		Stocks: []StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}, nil
}
