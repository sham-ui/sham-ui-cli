// Code generated by mockery v2.20.0. DO NOT EDIT.

package reset_password

import (
	model "{{ shortName }}/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockMemberService is an autogenerated mock type for the MemberService type
type MockMemberService struct {
	mock.Mock
}

// UpdatePassword provides a mock function with given fields: ctx, id, password
func (_m *MockMemberService) UpdatePassword(ctx context.Context, id model.MemberID, password model.MemberHashedPassword) error {
	ret := _m.Called(ctx, id, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.MemberID, model.MemberHashedPassword) error); ok {
		r0 = rf(ctx, id, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
