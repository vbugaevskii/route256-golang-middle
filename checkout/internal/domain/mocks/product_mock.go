// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	product "route256/checkout/internal/clients/product"
)

// ProductClient is an autogenerated mock type for the ProductClient type
type ProductClient struct {
	mock.Mock
}

type ProductClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductClient) EXPECT() *ProductClient_Expecter {
	return &ProductClient_Expecter{mock: &_m.Mock}
}

// GetProduct provides a mock function with given fields: ctx, sku
func (_m *ProductClient) GetProduct(ctx context.Context, sku uint32) (product.ResponseGetProduct, error) {
	ret := _m.Called(ctx, sku)

	var r0 product.ResponseGetProduct
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (product.ResponseGetProduct, error)); ok {
		return rf(ctx, sku)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) product.ResponseGetProduct); ok {
		r0 = rf(ctx, sku)
	} else {
		r0 = ret.Get(0).(product.ResponseGetProduct)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, sku)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductClient_GetProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProduct'
type ProductClient_GetProduct_Call struct {
	*mock.Call
}

// GetProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - sku uint32
func (_e *ProductClient_Expecter) GetProduct(ctx interface{}, sku interface{}) *ProductClient_GetProduct_Call {
	return &ProductClient_GetProduct_Call{Call: _e.mock.On("GetProduct", ctx, sku)}
}

func (_c *ProductClient_GetProduct_Call) Run(run func(ctx context.Context, sku uint32)) *ProductClient_GetProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32))
	})
	return _c
}

func (_c *ProductClient_GetProduct_Call) Return(_a0 product.ResponseGetProduct, _a1 error) *ProductClient_GetProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductClient_GetProduct_Call) RunAndReturn(run func(context.Context, uint32) (product.ResponseGetProduct, error)) *ProductClient_GetProduct_Call {
	_c.Call.Return(run)
	return _c
}

// ListSkus provides a mock function with given fields: ctx, startAfterSku, count
func (_m *ProductClient) ListSkus(ctx context.Context, startAfterSku uint32, count uint32) (product.ResponseListSkus, error) {
	ret := _m.Called(ctx, startAfterSku, count)

	var r0 product.ResponseListSkus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (product.ResponseListSkus, error)); ok {
		return rf(ctx, startAfterSku, count)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) product.ResponseListSkus); ok {
		r0 = rf(ctx, startAfterSku, count)
	} else {
		r0 = ret.Get(0).(product.ResponseListSkus)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(ctx, startAfterSku, count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductClient_ListSkus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSkus'
type ProductClient_ListSkus_Call struct {
	*mock.Call
}

// ListSkus is a helper method to define mock.On call
//   - ctx context.Context
//   - startAfterSku uint32
//   - count uint32
func (_e *ProductClient_Expecter) ListSkus(ctx interface{}, startAfterSku interface{}, count interface{}) *ProductClient_ListSkus_Call {
	return &ProductClient_ListSkus_Call{Call: _e.mock.On("ListSkus", ctx, startAfterSku, count)}
}

func (_c *ProductClient_ListSkus_Call) Run(run func(ctx context.Context, startAfterSku uint32, count uint32)) *ProductClient_ListSkus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32), args[2].(uint32))
	})
	return _c
}

func (_c *ProductClient_ListSkus_Call) Return(_a0 product.ResponseListSkus, _a1 error) *ProductClient_ListSkus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductClient_ListSkus_Call) RunAndReturn(run func(context.Context, uint32, uint32) (product.ResponseListSkus, error)) *ProductClient_ListSkus_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductClient creates a new instance of ProductClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductClient {
	mock := &ProductClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}