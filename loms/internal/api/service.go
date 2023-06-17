package api

import (
	"route256/loms/internal/domain"
	"route256/loms/pkg/loms"
)

type Service struct {
	loms.UnimplementedLomsServer
	model *domain.Model
}

func NewService(model *domain.Model) *Service {
	return &Service{model: model}
}
