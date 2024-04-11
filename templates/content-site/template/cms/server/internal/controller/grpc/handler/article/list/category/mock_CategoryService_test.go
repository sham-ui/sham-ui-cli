// Code generated by mockery v2.20.0. DO NOT EDIT.

package category

import (
	model "cms/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCategoryService is an autogenerated mock type for the CategoryService type
type MockCategoryService struct {
	mock.Mock
}

// GetBySlug provides a mock function with given fields: ctx, slug
func (_m *MockCategoryService) GetBySlug(ctx context.Context, slug model.CategorySlug) (*model.Category, error) {
	ret := _m.Called(ctx, slug)

	var r0 *model.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.CategorySlug) (*model.Category, error)); ok {
		return rf(ctx, slug)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.CategorySlug) *model.Category); ok {
		r0 = rf(ctx, slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.CategorySlug) error); ok {
		r1 = rf(ctx, slug)
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