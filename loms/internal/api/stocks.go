package api

import (
	"context"
	"log"
	"route256/loms/pkg/loms"
)

func (s *Service) Stocks(ctx context.Context, req *loms.RequestStocks) (*loms.ResponseStocks, error) {
	log.Printf("%+v", req)

	_, _ = s.model.Stocks(ctx, req.Sku)
	// if err != nil {
	// 	return nil, err
	// }

	return &loms.ResponseStocks{
		Stocks: []*loms.ResponseStocks_StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}, nil
}
