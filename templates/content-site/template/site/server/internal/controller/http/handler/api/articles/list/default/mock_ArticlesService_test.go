// Code generated by mockery v2.20.0. DO NOT EDIT.

package _default

import (
	context "context"

	model "site/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// MockArticlesService is an autogenerated mock type for the ArticlesService type
type MockArticlesService struct {
	mock.Mock
}

// Articles provides a mock function with given fields: ctx, offset, limit
func (_m *MockArticlesService) Articles(ctx context.Context, offset int64, limit int64) (*model.PaginatedArticles, error) {
	ret := _m.Called(ctx, offset, limit)

	var r0 *model.PaginatedArticles
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (*model.PaginatedArticles, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) *model.PaginatedArticles); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PaginatedArticles)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockArticlesService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockArticlesService creates a new instance of MockArticlesService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockArticlesService(t mockConstructorTestingTNewMockArticlesService) *MockArticlesService {
	mock := &MockArticlesService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
