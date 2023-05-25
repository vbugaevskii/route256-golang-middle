package purchase

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/listcart"

	cliloms "route256/checkout/internal/clients/loms"
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

	handListCart := listcart.Handler{
		Model: h.Model,
	}
	cart, err := handListCart.Handle(ctx, listcart.Request{User: req.User})
	log.Printf("Checkout.listcart: %+v", cart)
	if err != nil {
		return Response{}, err
	}

	items := make([]cliloms.RequestCreateOrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, cliloms.RequestCreateOrderItem{
			SKU:   item.SKU,
			Count: uint64(item.Count),
		})
	}

	res, err := h.Model.Loms.CreateOrder(ctx, req.User, items)
	log.Printf("LOMS.createOrder: %+v", res)
	if err != nil {
		return Response{}, err
	}

	return Response{OrderId: res.OrderId}, nil
}
