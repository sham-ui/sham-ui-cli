// Code generated by mockery v2.20.0. DO NOT EDIT.

package update

import (
	model "cms/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockSlugifyService is an autogenerated mock type for the SlugifyService type
type MockSlugifyService struct {
	mock.Mock
}

// SlugifyCategory provides a mock function with given fields: ctx, name
func (_m *MockSlugifyService) SlugifyCategory(ctx context.Context, name string) model.CategorySlug {
	ret := _m.Called(ctx, name)

	var r0 model.CategorySlug
	if rf, ok := ret.Get(0).(func(context.Context, string) model.CategorySlug); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(model.CategorySlug)
	}

	return r0
}

type mockConstructorTestingTNewMockSlugifyService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSlugifyService creates a new instance of MockSlugifyService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSlugifyService(t mockConstructorTestingTNewMockSlugifyService) *MockSlugifyService {
	mock := &MockSlugifyService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
