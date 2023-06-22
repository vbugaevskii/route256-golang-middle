package domain

import (
	"context"
	"errors"
	cliloms "route256/checkout/internal/clients/loms"
	cliproduct "route256/checkout/internal/clients/product"
	"route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/repository/postgres/cartitems"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	t.Parallel()

	var (
		userId  int64 = 1
		orderId int64 = 42

		sku1 uint32 = 1
		sku2 uint32 = 2
	)

	resp1 := cliproduct.ResponseGetProduct{
		Name:  "product_1",
		Price: 10,
	}
	resp2 := cliproduct.ResponseGetProduct{
		Name:  "product_2",
		Price: 100,
	}

	itemList := []cartitems.CartItem{
		{SKU: sku1, Count: 2},
		{SKU: sku2, Count: 3},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)
		cartItemsRepo.On("DeleteCart", mock.Anything, userId).Return(nil)

		// NOTE: I can't mock ListCart for Model, so I have to mock these function calls :(
		// I have mocks in api package, but then I have cycle dependencies
		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku1).Return(resp1, nil)
		productService.On("GetProduct", mock.Anything, sku2).Return(resp2, nil)

		// NOTE: I have to use "loms" instead of "cliloms" due to autogenerated code
		lomsService := mocks.NewLomsClient(t)
		lomsService.On("CreateOrder", mock.Anything, userId, mock.AnythingOfType("[]loms.RequestCreateOrderItem")).
			Return(cliloms.ResponseCreateOrder{OrderId: orderId}, nil)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		result, err := model.Purchase(context.Background(), userId)

		// Assert
		require.NoError(t, err)
		require.Equal(t, result, orderId)
	})

	t.Run("fail lomsService.CreateOrder", func(t *testing.T) {
		t.Parallel()

		errExpected := errors.New("error expected")

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)

		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku1).Return(resp1, nil)
		productService.On("GetProduct", mock.Anything, sku2).Return(resp2, nil)

		lomsService := mocks.NewLomsClient(t)
		lomsService.On("CreateOrder", mock.Anything, userId, mock.AnythingOfType("[]loms.RequestCreateOrderItem")).
			Return(cliloms.ResponseCreateOrder{}, errExpected)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		_, err := model.Purchase(context.Background(), userId)

		// Assert
		require.Error(t, err, errExpected)
	})

	t.Run("fail cartItems.DeleteCart", func(t *testing.T) {
		t.Parallel()

		errExpected := errors.New("error expected")

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{Items: itemList}, nil)
		cartItemsRepo.On("DeleteCart", mock.Anything, userId).Return(errExpected)

		productService := mocks.NewProductClient(t)
		productService.On("GetProduct", mock.Anything, sku1).Return(resp1, nil)
		productService.On("GetProduct", mock.Anything, sku2).Return(resp2, nil)

		lomsService := mocks.NewLomsClient(t)
		lomsService.On("CreateOrder", mock.Anything, userId, mock.AnythingOfType("[]loms.RequestCreateOrderItem")).
			Return(cliloms.ResponseCreateOrder{OrderId: orderId}, nil)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		_, err := model.Purchase(context.Background(), userId)

		// Assert
		require.Error(t, err, errExpected)
	})

	t.Run("fail cartItems.ListCart", func(t *testing.T) {
		t.Parallel()

		errExpected := errors.New("error expected")

		cartItemsRepo := mocks.NewCartItemsRepository(t)
		cartItemsRepo.On("ListCart", mock.Anything, userId).Return(cartitems.ResponseListCart{}, errExpected)

		productService := mocks.NewProductClient(t)

		lomsService := mocks.NewLomsClient(t)

		// Act
		model := New(lomsService, productService, cartItemsRepo)
		_, err := model.Purchase(context.Background(), userId)

		// Assert
		require.Error(t, err, errExpected)
	})
}
