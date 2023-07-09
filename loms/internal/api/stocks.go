package api

import (
	"context"
	"route256/libs/logger"
	"route256/loms/pkg/loms"
)

func (s *Service) Stocks(ctx context.Context, req *loms.RequestStocks) (*loms.ResponseStocks, error) {
	logger.Infof("%+v", req)

	stocks, err := s.model.Stocks(ctx, req.Sku)
	if err != nil {
		return nil, err
	}

	result := make([]*loms.ResponseStocks_StockItem, 0, len(stocks))
	for _, item := range stocks {
		result = append(result, &loms.ResponseStocks_StockItem{
			WarehouseID: item.WarehouseId,
			Count:       uint64(item.Count),
		})
	}

	return &loms.ResponseStocks{Stocks: result}, nil
}
