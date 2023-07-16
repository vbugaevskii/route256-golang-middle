package api

import (
	"context"
	"errors"
	"route256/notifications/internal/domain"
	nofity "route256/notifications/pkg/notifications"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Impl interface {
	List(ctx context.Context, user int64, tsFrom time.Time, tsTill time.Time) ([]domain.Notification, error)
}

type Service struct {
	nofity.UnimplementedNotificationsServer
	model Impl
}

func NewService(model Impl) *Service {
	return &Service{model: model}
}

var (
	ErrInvalidPeriod = errors.New("invalid period: `tsTill` should be greater than `tsFrom`")
)

func (s *Service) List(ctx context.Context, req *nofity.RequestList) (*nofity.ResponseList, error) {
	var tsFrom, tsTill time.Time

	if req.TsFrom != nil {
		tsFrom = req.TsFrom.AsTime()
	}

	if req.TsTill != nil {
		tsTill = req.TsTill.AsTime()
	} else {
		tsTill = time.Now()
	}

	if req.TsFrom != nil && req.TsTill != nil && tsFrom.After(tsTill) {
		return nil, ErrInvalidPeriod
	}

	resp, err := s.model.List(ctx, req.User, tsFrom, tsTill)
	if err != nil {
		return nil, err
	}

	respPb := nofity.ResponseList{
		User: req.User,
	}
	for _, item := range resp {
		respPb.Items = append(respPb.Items, &nofity.ResponseList_Notification{
			Message:   item.Message,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}
	return &respPb, nil
}
