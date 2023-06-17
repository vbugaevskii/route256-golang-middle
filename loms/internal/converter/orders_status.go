package converter

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvStatusDomainSchema(statusDomain domain.StatusType) schema.StatusType {
	var statusSchema schema.StatusType

	switch statusDomain {
	case domain.StatusNew:
		statusSchema = schema.StatusNew
	case domain.StatusAwaitingPayment:
		statusSchema = schema.StatusAwaitingPayment
	case domain.StatusFailed:
		statusSchema = schema.StatusFailed
	case domain.StatusPayed:
		statusSchema = schema.StatusPayed
	case domain.StatusCancelled:
		statusSchema = schema.StatusCancelled
	}

	return statusSchema
}

func ConvStatusSchemaDomain(statusSchema schema.StatusType) domain.StatusType {
	var statusDomain domain.StatusType

	switch statusSchema {
	case schema.StatusNew:
		statusDomain = domain.StatusNew
	case schema.StatusAwaitingPayment:
		statusDomain = domain.StatusAwaitingPayment
	case schema.StatusFailed:
		statusDomain = domain.StatusFailed
	case schema.StatusPayed:
		statusDomain = domain.StatusPayed
	case schema.StatusCancelled:
		statusDomain = domain.StatusCancelled
	}

	return statusDomain
}
