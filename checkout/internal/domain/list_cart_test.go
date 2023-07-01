package domain

import (
	"context"
	"errors"
	"fmt"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/repository/postgres/cartitems"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			userId int64 = 1

			sku1 uint32 = 1
			sku2 uint32 = 2
		)

		item1 := cartitems.CartItem{
			SKU:   sku1,
			Count: 2,
		}
		item2 := cartitems.CartItem{
			SKU:   sku2,
			Count: 3,
		}

		resp1 := cliproduct.ResponseGetProduct{
			Name:  "product_1",
			Price: 10,
		}
		resp2 := cliproduct.ResponseGetProduct{
			Name:  "product_2",
			Price: 100,
		}

		itemList := []cartitems.CartItem{
			item1,
			item2,
		}

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)

		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku1).Return(resp1, nil)
		productService.On("GetProduct", mock.Anything, sku2).Return(resp2, nil)

		lomsService := mocks.NewLomsClient(t)

		expected := []*CartItem{
			{SKU: item1.SKU, Count: item1.Count, Name: resp1.Name, Price: resp1.Price},
			{SKU: item2.SKU, Count: item2.Count, Name: resp2.Name, Price: resp2.Price},
		}

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		result, err := model.ListCart(context.Background(), userId)

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2, "len(result) != 2")
		require.ElementsMatch(t, expected, result)
	})

	t.Run("success multiple", func(t *testing.T) {
		t.Parallel()

		var userId int64 = 1

		skuList := []uint32{6, 8, 7, 1, 3, 9, 10, 2, 4, 5}

		itemList := make([]cartitems.CartItem, 0, len(skuList))
		for i, sku := range skuList {
			item := cartitems.CartItem{
				SKU:   sku,
				Count: uint16(i + 1),
			}
			itemList = append(itemList, item)
		}

		respList := make([]cliproduct.ResponseGetProduct, 0, len(skuList))
		for i, sku := range skuList {
			resp := cliproduct.ResponseGetProduct{
				Name:  fmt.Sprintf("product_%d", sku),
				Price: uint32(i * 10),
			}
			respList = append(respList, resp)
		}

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)

		productService := mocks.NewProductClient(t)
		for i, sku := range skuList {
			productService.On("GetProduct", mock.Anything, sku).Return(respList[i], nil)
		}

		lomsService := mocks.NewLomsClient(t)

		expected := make([]*CartItem, 0, len(skuList))
		for i, item := range itemList {
			itemCopy := CartItem{
				SKU:   item.SKU,
				Count: item.Count,
				Name:  respList[i].Name,
				Price: respList[i].Price,
			}
			expected = append(expected, &itemCopy)
		}

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		result, err := model.ListCart(context.Background(), userId)

		// Assert
		require.NoError(t, err)
		require.Len(t, result, len(expected), "len(result) != len(expected)")
		require.ElementsMatch(t, expected, result)
	})

	t.Run("fail cartItems.ListCart", func(t *testing.T) {
		t.Parallel()

		var userId int64 = 1

		errExpected := errors.New("error expected")

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{}, errExpected)

		productService := mocks.NewProductClient(t)
		lomsService := mocks.NewLomsClient(t)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		_, err := model.ListCart(context.Background(), userId)

		// Assert
		require.Error(t, err, errExpected)
	})

	t.Run("fail productService.GetProduct", func(t *testing.T) {
		t.Parallel()

		var (
			userId int64  = 1
			sku    uint32 = 2
		)

		errExpected := errors.New("error expected")

		itemList := []cartitems.CartItem{
			{SKU: sku, Count: 10},
		}

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)

		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku).
			Return(cliproduct.ResponseGetProduct{}, errExpected)

		lomsService := mocks.NewLomsClient(t)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		_, err := model.ListCart(context.Background(), userId)

		// Assert
		require.Error(t, err, errExpected)
	})
}
