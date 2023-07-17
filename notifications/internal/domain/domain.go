package domain

import (
	"context"
	"route256/libs/logger"
	"time"
)

type Model struct {
	repo  NotificationRepository
	cache Cache[CacheKey, []Notification]
}

type Cache[K comparable, V any] interface {
	Add(key K, value V) bool
	Get(key K) (V, bool)
	Remove(key K) bool
	Contains(key K) bool
	Len() int
}

type CacheKey struct {
	userId int64
	tsFrom time.Time
	tsTill time.Time
}

func NewModel(repo NotificationRepository, cache Cache[CacheKey, []Notification]) *Model {
	return &Model{
		repo:  repo,
		cache: cache,
	}
}

type NotificationRepository interface {
	ListNotifications(ctx context.Context, userId int64, tsFrom time.Time, tsTill time.Time) ([]Notification, error)
	SaveNotification(ctx context.Context, recordId int64, userId int64, message string) error
}

type Notification struct {
	Message   string
	CreatedAt time.Time
}

func (m *Model) List(ctx context.Context, userId int64, tsFrom time.Time, tsTill time.Time) ([]Notification, error) {
	cacheKey := CacheKey{
		userId: userId,
		tsFrom: tsFrom,
		tsTill: tsTill,
	}

	if m.cache != nil {
		if items, exists := m.cache.Get(cacheKey); exists {
			logger.Infof("cache hit for userId=%d", userId)
			return items, nil
		}
	}

	items, err := m.repo.ListNotifications(ctx, userId, tsFrom, tsTill)
	if err != nil {
		return nil, err
	}

	if m.cache != nil {
		m.cache.Add(cacheKey, items)
	}

	return items, nil
}
