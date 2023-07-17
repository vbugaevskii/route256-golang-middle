package schema

import "time"

type Notification struct {
	RecordId  int64     `db:"record_id"`
	UserID    int64     `db:"user_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
