package domain

import (
	"context"
	"route256/libs/logger"
	"time"
)

type Model struct {
	repo  NotificationRepository
	cache Cache[int64, []Notification]
}

type Cache[K comparable, V any] interface {
	Add(key K, value V) bool
	Get(key K) (V, bool)
	Remove(key K) bool
	Contains(key K) bool
	Len() int
}

func NewModel(repo NotificationRepository, cache Cache[int64, []Notification]) *Model {
	return &Model{
		repo:  repo,
		cache: cache,
	}
}

type NotificationRepository interface {
	ListNotifications(ctx context.Context, userId int64) ([]Notification, error)
	SaveNotification(ctx context.Context, recordId int64, userId int64, message string) error
}

type Notification struct {
	Message   string
	CreatedAt time.Time
}

func (m *Model) List(ctx context.Context, userId int64) ([]Notification, error) {
	if m.cache != nil {
		if items, exists := m.cache.Get(userId); exists {
			logger.Infof("cache hit for userId=%d", userId)
			return items, nil
		}
	}

	items, err := m.repo.ListNotifications(ctx, userId)
	if err != nil {
		return nil, err
	}

	if m.cache != nil {
		m.cache.Add(userId, items)
	}

	return items, nil
}
