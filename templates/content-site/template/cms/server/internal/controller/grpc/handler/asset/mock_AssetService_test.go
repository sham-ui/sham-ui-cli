// Code generated by mockery v2.20.0. DO NOT EDIT.

package asset

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAssetService is an autogenerated mock type for the AssetService type
type MockAssetService struct {
	mock.Mock
}

// ReadFile provides a mock function with given fields: ctx, path
func (_m *MockAssetService) ReadFile(ctx context.Context, path string) ([]byte, error) {
	ret := _m.Called(ctx, path)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, error)); ok {
		return rf(ctx, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockAssetService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAssetService creates a new instance of MockAssetService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAssetService(t mockConstructorTestingTNewMockAssetService) *MockAssetService {
	mock := &MockAssetService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
