// Code generated by mockery v2.20.0. DO NOT EDIT.

package cms

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	proto "site/internal/external_api/cms/proto"
)

// mockCmsClient is an autogenerated mock type for the cmsClient type
type mockCmsClient struct {
	mock.Mock
}

// Article provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) Article(ctx context.Context, in *proto.ArticleRequest, opts ...grpc.CallOption) (*proto.ArticleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ArticleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleRequest, ...grpc.CallOption) (*proto.ArticleResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleRequest, ...grpc.CallOption) *proto.ArticleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleList provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) ArticleList(ctx context.Context, in *proto.ArticleListRequest, opts ...grpc.CallOption) (*proto.ArticleListResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ArticleListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListRequest, ...grpc.CallOption) (*proto.ArticleListResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListRequest, ...grpc.CallOption) *proto.ArticleListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForCategory provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) ArticleListForCategory(ctx context.Context, in *proto.ArticleListForCategoryRequest, opts ...grpc.CallOption) (*proto.ArticleListForCategoryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ArticleListForCategoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForCategoryRequest, ...grpc.CallOption) (*proto.ArticleListForCategoryResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForCategoryRequest, ...grpc.CallOption) *proto.ArticleListForCategoryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListForCategoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForCategoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForQuery provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) ArticleListForQuery(ctx context.Context, in *proto.ArticleListForQueryRequest, opts ...grpc.CallOption) (*proto.ArticleListResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ArticleListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForQueryRequest, ...grpc.CallOption) (*proto.ArticleListResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForQueryRequest, ...grpc.CallOption) *proto.ArticleListResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForQueryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ArticleListForTag provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) ArticleListForTag(ctx context.Context, in *proto.ArticleListForTagRequest, opts ...grpc.CallOption) (*proto.ArticleListForTagResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ArticleListForTagResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForTagRequest, ...grpc.CallOption) (*proto.ArticleListForTagResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ArticleListForTagRequest, ...grpc.CallOption) *proto.ArticleListForTagResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ArticleListForTagResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.ArticleListForTagRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Asset provides a mock function with given fields: ctx, in, opts
func (_m *mockCmsClient) Asset(ctx context.Context, in *proto.AssetRequest, opts ...grpc.CallOption) (*proto.AssetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.AssetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AssetRequest, ...grpc.CallOption) (*proto.AssetResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AssetRequest, ...grpc.CallOption) *proto.AssetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.AssetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.AssetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockCmsClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockCmsClient creates a new instance of mockCmsClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockCmsClient(t mockConstructorTestingTnewMockCmsClient) *mockCmsClient {
	mock := &mockCmsClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
