// Code generated by mockery v2.20.0. DO NOT EDIT.

package list

import (
	model "{{ shortName }}/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockMemberService is an autogenerated mock type for the MemberService type
type MockMemberService struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, offset, limit
func (_m *MockMemberService) Find(ctx context.Context, offset int64, limit int64) ([]model.Member, error) {
	ret := _m.Called(ctx, offset, limit)

	var r0 []model.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) ([]model.Member, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []model.Member); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Total provides a mock function with given fields: ctx
func (_m *MockMemberService) Total(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockMemberService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMemberService creates a new instance of MockMemberService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMemberService(t mockConstructorTestingTNewMockMemberService) *MockMemberService {
	mock := &MockMemberService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
