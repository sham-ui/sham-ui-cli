// Code generated by mockery v2.20.0. DO NOT EDIT.

package integration

import (
	context "context"
	proto "site/internal/external_api/cms/proto"

	mock "github.com/stretchr/testify/mock"
)

// mockCmsServer is an autogenerated mock type for the cmsServer type
type mockCmsServer struct {
	mock.Mock
}

// Article provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) Article(_a0 context.Context, _a1 *proto.ArticleRequest) (*proto.ArticleResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.ArticleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleRequest) (*proto.ArticleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleRequest) *proto.ArticleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleList provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) ArticleList(_a0 context.Context, _a1 *proto.ArticleListRequest) (*proto.ArticleListResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.ArticleListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListRequest) (*proto.ArticleListResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListRequest) *proto.ArticleListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForCategory provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) ArticleListForCategory(_a0 context.Context, _a1 *proto.ArticleListForCategoryRequest) (*proto.ArticleListForCategoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.ArticleListForCategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForCategoryRequest) (*proto.ArticleListForCategoryResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForCategoryRequest) *proto.ArticleListForCategoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListForCategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForCategoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForQuery provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) ArticleListForQuery(_a0 context.Context, _a1 *proto.ArticleListForQueryRequest) (*proto.ArticleListResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.ArticleListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForQueryRequest) (*proto.ArticleListResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForQueryRequest) *proto.ArticleListResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForQueryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForTag provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) ArticleListForTag(_a0 context.Context, _a1 *proto.ArticleListForTagRequest) (*proto.ArticleListForTagResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.ArticleListForTagResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForTagRequest) (*proto.ArticleListForTagResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForTagRequest) *proto.ArticleListForTagResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListForTagResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForTagRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Asset provides a mock function with given fields: _a0, _a1
func (_m *mockCmsServer) Asset(_a0 context.Context, _a1 *proto.AssetRequest) (*proto.AssetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.AssetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AssetRequest) (*proto.AssetResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AssetRequest) *proto.AssetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.AssetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.AssetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockCmsServer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockCmsServer creates a new instance of mockCmsServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockCmsServer(t mockConstructorTestingTnewMockCmsServer) *mockCmsServer {
	mock := &mockCmsServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
