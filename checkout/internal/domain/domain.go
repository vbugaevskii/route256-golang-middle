package domain

import (
	cliLoms "route256/loms/external/client"
)

type Model struct {
	Loms cliLoms.Client
}

func New(loms *LomsClient) *Model {
	return &Model{
		Loms: loms,
	}
}
