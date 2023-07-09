package schema

import "time"

type NotificationState string

const (
	StateWaiting   NotificationState = "Waiting"
	StateDelivered NotificationState = "Delivered"
)

type Notification struct {
	RecordId  int64             `db:"record_id"`
	Key       string            `db:"key"`
	Value     string            `db:"value"`
	State     NotificationState `db:"state"`
	CreatedAt time.Time         `db:"created_at"`
}
