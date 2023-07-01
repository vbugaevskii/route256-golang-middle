package api

import (
	"context"
	"errors"
	"route256/checkout/internal/api/mocks"
	"route256/checkout/pkg/checkout"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			userId  int64 = 1
			orderId int64 = 5
		)

		model := mocks.NewImpl(t)
		model.On("Purchase", mock.Anything, userId).Return(orderId, nil)

		// Act
		service := NewService(model)
		result, err := service.Purchase(context.Background(), &checkout.RequestPurchase{
			User: userId,
		})

		expected := &checkout.ResponsePurchase{
			OrderID: orderId,
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
		_, err := service.Purchase(context.Background(), &checkout.RequestPurchase{
			User: 0,
		})

		// Assert
		require.Error(t, err, ErrUserNotFound)
	})

	t.Run("fail domain purchase", func(t *testing.T) {
		t.Parallel()

		var userId int64 = 1

		errExpected := errors.New("error expected")

		model := mocks.NewImpl(t)
		model.On("Purchase", mock.Anything, userId).Return(int64(0), errExpected)

		// Act
		service := NewService(model)
		_, err := service.Purchase(context.Background(), &checkout.RequestPurchase{
			User: userId,
		})

		// Assert
		require.Error(t, err, errExpected)
	})
}
