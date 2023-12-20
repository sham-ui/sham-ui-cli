// Code generated by mockery v2.20.0. DO NOT EDIT.

package tag

import (
	context "context"
	model "site/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// MockArticlesService is an autogenerated mock type for the ArticlesService type
type MockArticlesService struct {
	mock.Mock
}

// ArticleListForTag provides a mock function with given fields: ctx, tagSlug, offset, limit
func (_m *MockArticlesService) ArticleListForTag(ctx context.Context, tagSlug string, offset int64, limit int64) (*model.PaginatedArticleForTag, error) {
	ret := _m.Called(ctx, tagSlug, offset, limit)

	var r0 *model.PaginatedArticleForTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64) (*model.PaginatedArticleForTag, error)); ok {
		return rf(ctx, tagSlug, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64) *model.PaginatedArticleForTag); ok {
		r0 = rf(ctx, tagSlug, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PaginatedArticleForTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64, int64) error); ok {
		r1 = rf(ctx, tagSlug, offset, limit)
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
