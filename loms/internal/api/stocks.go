package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"
)

func (s *Service) Stocks(ctx context.Context, req *loms.RequestStocks) (*loms.ResponseStocks, error) {
	log.Printf("%+v", req)

	return &loms.ResponseStocks{
		Stocks: []*loms.StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}, nil
}
