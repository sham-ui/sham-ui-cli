// Code generated by mockery v2.20.0. DO NOT EDIT.

package detail

import (
	model "cms/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockArticleService is an autogenerated mock type for the ArticleService type
type MockArticleService struct {
	mock.Mock
}

type MockArticleService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockArticleService) EXPECT() *MockArticleService_Expecter {
	return &MockArticleService_Expecter{mock: &_m.Mock}
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *MockArticleService) FindByID(ctx context.Context, id model.ArticleID) (model.Article, error) {
	ret := _m.Called(ctx, id)

	var r0 model.Article
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ArticleID) (model.Article, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ArticleID) model.Article); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Article)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ArticleID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockArticleService_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type MockArticleService_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id model.ArticleID
func (_e *MockArticleService_Expecter) FindByID(ctx interface{}, id interface{}) *MockArticleService_FindByID_Call {
	return &MockArticleService_FindByID_Call{Call: _e.mock.On("FindByID", ctx, id)}
}

func (_c *MockArticleService_FindByID_Call) Run(run func(ctx context.Context, id model.ArticleID)) *MockArticleService_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ArticleID))
	})
	return _c
}

func (_c *MockArticleService_FindByID_Call) Return(_a0 model.Article, _a1 error) *MockArticleService_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockArticleService_FindByID_Call) RunAndReturn(run func(context.Context, model.ArticleID) (model.Article, error)) *MockArticleService_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetTags provides a mock function with given fields: ctx, articleID
func (_m *MockArticleService) GetTags(ctx context.Context, articleID model.ArticleID) ([]model.Tag, error) {
	ret := _m.Called(ctx, articleID)

	var r0 []model.Tag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ArticleID) ([]model.Tag, error)); ok {
		return rf(ctx, articleID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ArticleID) []model.Tag); ok {
		r0 = rf(ctx, articleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Tag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ArticleID) error); ok {
		r1 = rf(ctx, articleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockArticleService_GetTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTags'
type MockArticleService_GetTags_Call struct {
	*mock.Call
}

// GetTags is a helper method to define mock.On call
//   - ctx context.Context
//   - articleID model.ArticleID
func (_e *MockArticleService_Expecter) GetTags(ctx interface{}, articleID interface{}) *MockArticleService_GetTags_Call {
	return &MockArticleService_GetTags_Call{Call: _e.mock.On("GetTags", ctx, articleID)}
}

func (_c *MockArticleService_GetTags_Call) Run(run func(ctx context.Context, articleID model.ArticleID)) *MockArticleService_GetTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ArticleID))
	})
	return _c
}

func (_c *MockArticleService_GetTags_Call) Return(_a0 []model.Tag, _a1 error) *MockArticleService_GetTags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockArticleService_GetTags_Call) RunAndReturn(run func(context.Context, model.ArticleID) ([]model.Tag, error)) *MockArticleService_GetTags_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockArticleService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockArticleService creates a new instance of MockArticleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockArticleService(t mockConstructorTestingTNewMockArticleService) *MockArticleService {
	mock := &MockArticleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
