// Code generated by mockery v2.20.0. DO NOT EDIT.

package reset_password

import (
	model "{{ shortName }}/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockPasswordService is an autogenerated mock type for the PasswordService type
type MockPasswordService struct {
	mock.Mock
}

// Hash provides a mock function with given fields: ctx, raw
func (_m *MockPasswordService) Hash(ctx context.Context, raw string) (model.MemberHashedPassword, error) {
	ret := _m.Called(ctx, raw)

	var r0 model.MemberHashedPassword
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.MemberHashedPassword, error)); ok {
		return rf(ctx, raw)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.MemberHashedPassword); ok {
		r0 = rf(ctx, raw)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.MemberHashedPassword)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, raw)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockPasswordService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPasswordService creates a new instance of MockPasswordService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPasswordService(t mockConstructorTestingTNewMockPasswordService) *MockPasswordService {
	mock := &MockPasswordService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
