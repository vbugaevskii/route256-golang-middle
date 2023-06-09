package converter

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvOrderSchemaDomain(orderSchema schema.Order) domain.Order {
	return domain.Order{
		Status: ConvStatusSchemaDomain(orderSchema.Status),
		User:   orderSchema.UserId,
	}
}
