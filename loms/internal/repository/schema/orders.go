package schema

type StatusType string

const (
	New             StatusType = "New"
	AwaitingPayment StatusType = "AwaitingPayment"
	Failed          StatusType = "Failed"
	Payed           StatusType = "Payed"
	Cancelled       StatusType = "Cancelled"
)

type Order struct {
	OrderId int64      `db:"order_id"`
	UserId  int64      `db:"user_id"`
	Status  StatusType `db:"status"`
}
