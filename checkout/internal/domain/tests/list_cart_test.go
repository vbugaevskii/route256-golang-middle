package tests

import (
	"context"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
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
			Price: 10,
		}
		resp2 := cliproduct.ResponseGetProduct{
			Name:  "product_2",
			Price: 10,
		}

		itemWithFields1 := item1
		itemWithFields1.Name = resp1.Name
		itemWithFields1.Price = resp1.Price

		itemWithFields2 := item2
		itemWithFields2.Name = resp2.Name
		itemWithFields2.Price = resp2.Price

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

		expected := []*domain.CartItem{
			&itemWithFields1,
			&itemWithFields2,
		}

		// Act
		model := domain.New(lomsService, productService, cartItemsRepo)
		result, err := model.ListCart(context.Background(), int64(1))

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2, "len(result) != 2")
		require.ElementsMatch(t, expected, result)
	})
}
