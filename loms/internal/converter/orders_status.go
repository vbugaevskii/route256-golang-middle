package converter

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvStatusDomainSchema(statusDomain domain.StatusType) schema.StatusType {
	var statusSchema schema.StatusType

	switch statusDomain {
	case domain.New:
		statusSchema = schema.New
	case domain.AwaitingPayment:
		statusSchema = schema.AwaitingPayment
	case domain.Failed:
		statusSchema = schema.Failed
	case domain.Payed:
		statusSchema = schema.Payed
	case domain.Cancelled:
		statusSchema = schema.Cancelled
	}

	return statusSchema
}

func ConvStatusSchemaDomain(statusSchema schema.StatusType) domain.StatusType {
	var statusDomain domain.StatusType

	switch statusSchema {
	case schema.New:
		statusDomain = domain.New
	case schema.AwaitingPayment:
		statusDomain = domain.AwaitingPayment
	case schema.Failed:
		statusDomain = domain.Failed
	case schema.Payed:
		statusDomain = domain.Payed
	case schema.Cancelled:
		statusDomain = domain.Cancelled
	}

	return statusDomain
}
