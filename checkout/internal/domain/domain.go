package domain

import (
	cliLoms "route256/loms/external/client"
)

type Model struct {
	Loms *cliLoms.Client
}

func New(netlocLoms string) *Model {
	return &Model{
		Loms: cliLoms.New(netlocLoms),
	}
}
