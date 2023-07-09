package converter

import (
	"encoding/json"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ConvNotificationsOutboxSchemaDomain(itemsSchema []schema.Notification) ([]domain.Notification, error) {
	itemsDomain := make([]domain.Notification, 0, len(itemsSchema))

	for _, item := range itemsSchema {
		orderMsg := domain.Notification{}
		err := json.Unmarshal([]byte(item.Value), &orderMsg)
		if err != nil {
			return nil, err
		}
		orderMsg.RecordId = item.RecordId
		orderMsg.CreatedAt = item.CreatedAt

		itemsDomain = append(itemsDomain, orderMsg)
	}

	return itemsDomain, nil
}
