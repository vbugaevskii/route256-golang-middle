package api

import "errors"

var (
	ErrEmptyOrder    = errors.New("empty order")
	ErrOrderNotFound = errors.New("order is not found")
)
