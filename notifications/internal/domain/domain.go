package domain

import (
	"context"
	"time"
)

type Model struct {
	repo NotificationRepository
}

func NewModel(repo NotificationRepository) *Model {
	return &Model{
		repo: repo,
	}
}

type NotificationRepository interface {
	ListNotifications(ctx context.Context, userId int64) ([]Notification, error)
}

type Notification struct {
	Message   string
	CreatedAt time.Time
}

func (m *Model) List(ctx context.Context, userId int64) ([]Notification, error) {
	items, err := m.repo.ListNotifications(ctx, userId)
	if err != nil {
		return nil, err
	}

	return items, nil
}
