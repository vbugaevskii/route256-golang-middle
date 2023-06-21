package tests

import (
	"context"
	"route256/checkout/internal/api"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"route256/checkout/pkg/checkout"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		sku1 := uint32(1)
		sku2 := uint32(2)

		item1 := domain.CartItem{
			SKU:   sku1,
			Count: 2,
		}
		item2 := domain.CartItem{
			SKU:   sku2,
			Count: 3,
		}

		resp1 := cliproduct.ResponseGetProduct{
			Name:  "product_1",
			Price: 100,
		}
		resp2 := cliproduct.ResponseGetProduct{
			Name:  "product_2",
			Price: 10,
		}

		cartItems := []*domain.CartItem{
			&item1,
			&item2,
		}

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, int64(1)).Return(cartItems, nil)

		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku1).Return(resp1, nil)
		productService.On("GetProduct", mock.Anything, sku2).Return(resp2, nil)

		lomsService := mocks.NewLomsClient(t)

		// Act
		model := domain.New(lomsService, productService, cartItemsRepo)

		service := api.NewService(model)
		result, err := service.ListCart(context.Background(), &checkout.RequestListCart{
			User: 1,
		})

		expected := &checkout.ResponseListCart{
			Items: []*checkout.ResponseListCart_CartItem{
				{
					Sku:   item1.SKU,
					Count: uint32(item1.Count),
					Name:  resp1.Name,
					Price: resp1.Price,
				},
				{
					Sku:   item2.SKU,
					Count: uint32(item2.Count),
					Name:  resp2.Name,
					Price: resp2.Price,
				},
			},
			TotalPrice: resp1.Price*uint32(item1.Count) + resp2.Price*uint32(item2.Count),
		}

		// Assert
		require.NoError(t, err)
		require.Len(t, result.Items, 2, "len(result.Items) != 2")
		require.Equal(t, expected, result)
	})
}
