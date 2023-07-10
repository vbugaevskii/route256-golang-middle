package domain

import "context"

type Model struct {
}

func NewModel() *Model {
	return &Model{}
}

type Item struct {
	OrderId   int64
	Status    string
	CreatedAt int64
}

func (m *Model) List(ctx context.Context, userId int64) ([]Item, error) {
	return nil, nil
}
