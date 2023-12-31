package api

import (
	"context"
	"errors"
	"route256/checkout/internal/api/mocks"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var userId int64 = 1

		cartItems := []*domain.CartItem{
			{SKU: 1, Count: 2, Name: "product_1", Price: 100},
			{SKU: 2, Count: 3, Name: "product_2", Price: 10},
		}

		model := mocks.NewImpl(t)
		model.On("ListCart", mock.Anything, userId).Return(cartItems, nil)

		// Act
		service := NewService(model)
		result, err := service.ListCart(context.Background(), &checkout.RequestListCart{
			User: userId,
		})

		expected := &checkout.ResponseListCart{
			Items:      make([]*checkout.ResponseListCart_CartItem, 0, len(cartItems)),
			TotalPrice: 0,
		}
		for _, item := range cartItems {
			itemPb := checkout.ResponseListCart_CartItem{
				Sku:   item.SKU,
				Count: uint32(item.Count),
				Name:  item.Name,
				Price: item.Price,
			}
			expected.Items = append(expected.Items, &itemPb)

			expected.TotalPrice += itemPb.Price * uint32(itemPb.Count)
		}

		// Assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("fail zero user", func(t *testing.T) {
		t.Parallel()

		model := mocks.NewImpl(t)

		// Act
		service := NewService(model)
		_, err := service.ListCart(context.Background(), &checkout.RequestListCart{
			User: 0,
		})

		// Assert
		require.Error(t, err, ErrUserNotFound)
	})

	t.Run("fail domain list cart", func(t *testing.T) {
		t.Parallel()

		var userId int64 = 1
		errExpected := errors.New("error expected")

		model := mocks.NewImpl(t)
		model.On("ListCart", mock.Anything, userId).Return(nil, errExpected)

		// Act
		service := NewService(model)
		_, err := service.ListCart(context.Background(), &checkout.RequestListCart{
			User: userId,
		})

		// Assert
		require.Error(t, err, errExpected)
	})
}
