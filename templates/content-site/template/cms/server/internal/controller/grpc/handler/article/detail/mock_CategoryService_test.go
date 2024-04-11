// Code generated by mockery v2.20.0. DO NOT EDIT.

package detail

import (
	model "cms/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCategoryService is an autogenerated mock type for the CategoryService type
type MockCategoryService struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockCategoryService) GetByID(ctx context.Context, id model.CategoryID) (*model.Category, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.CategoryID) (*model.Category, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.CategoryID) *model.Category); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.CategoryID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockCategoryService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCategoryService creates a new instance of MockCategoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCategoryService(t mockConstructorTestingTNewMockCategoryService) *MockCategoryService {
	mock := &MockCategoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}