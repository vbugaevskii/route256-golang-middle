package api

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrProductNotFound     = errors.New("product not found")
	ErrProductInsufficient = errors.New("product insufficient")
	ErrCartIsEmpty         = errors.New("cart is empty")
)
