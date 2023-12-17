// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Media is an autogenerated mock type for the Media type
type Media struct {
	mock.Mock
}

// ID provides a mock function with given fields:
func (_m *Media) ID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// URI provides a mock function with given fields:
func (_m *Media) URI() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for URI")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMedia creates a new instance of Media. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMedia(t interface {
	mock.TestingT
	Cleanup(func())
}) *Media {
	mock := &Media{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
