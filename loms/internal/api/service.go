package api

import (
	"route256/loms/pkg/loms"
)

type Service struct {
	loms.UnimplementedLomsServer
}

func NewService() *Service {
	return &Service{}
}
