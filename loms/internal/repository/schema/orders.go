package schema

import "time"

type StatusType string

const (
	StatusNew             StatusType = "New"
	StatusAwaitingPayment StatusType = "AwaitingPayment"
	StatusFailed          StatusType = "Failed"
	StatusPayed           StatusType = "Payed"
	StatusCancelled       StatusType = "Cancelled"
)

type Order struct {
	OrderId   int64      `db:"order_id"`
	UserId    int64      `db:"user_id"`
	Status    StatusType `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
}
