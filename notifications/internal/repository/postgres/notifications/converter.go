package notifications

import (
	"route256/notifications/internal/domain"
	"route256/notifications/internal/repository/schema"
)

func ConvNotificationSchemaDomain(itemsSchema []schema.Notification) []domain.Notification {
	itemsDomain := make([]domain.Notification, 0, len(itemsSchema))
	for _, item := range itemsSchema {
		itemsDomain = append(itemsDomain, domain.Notification{
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
		})
	}
	return itemsDomain
}
