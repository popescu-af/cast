// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	library "model/library"

	mock "github.com/stretchr/testify/mock"
)

// Library is an autogenerated mock type for the Library type
type Library struct {
	mock.Mock
}

// ID provides a mock function with given fields:
func (_m *Library) ID() string {
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

// List provides a mock function with given fields:
func (_m *Library) List() ([]library.Media, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []library.Media
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]library.Media, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []library.Media); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]library.Media)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLibrary creates a new instance of Library. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLibrary(t interface {
	mock.TestingT
	Cleanup(func())
}) *Library {
	mock := &Library{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}