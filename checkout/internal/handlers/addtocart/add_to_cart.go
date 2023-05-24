package addtocart

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
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Response struct {
}

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrProductInsufficient = errors.New("product insufficient")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)

	if req.User == 0 {
		return Response{}, ErrUserNotFound
	}

	stocks, err := h.Model.Loms.Stocks(ctx, req.SKU)
	log.Printf("LOMS.stocks: %+v", stocks)
	if err != nil {
		return Response{}, err
	}

	var count uint64
	for _, stock := range stocks.Stocks {
		count += stock.Count
		if count >= uint64(req.Count) {
			break
		}
	}

	if count < uint64(req.Count) {
		return Response{}, ErrProductInsufficient
	}

	return Response{}, nil
}
