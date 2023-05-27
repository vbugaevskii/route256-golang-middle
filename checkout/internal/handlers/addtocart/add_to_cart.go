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
	ErrUserNotFound = errors.New("user not found")
)

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("%+v", req)

	if req.User == 0 {
		return Response{}, ErrUserNotFound
	}

	err := h.Model.AddToCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, err
}
