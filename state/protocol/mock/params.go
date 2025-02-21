// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Params is an autogenerated mock type for the Params type
type Params struct {
	mock.Mock
}

// ChainID provides a mock function with given fields:
func (_m *Params) ChainID() (flow.ChainID, error) {
	ret := _m.Called()

	var r0 flow.ChainID
	if rf, ok := ret.Get(0).(func() flow.ChainID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(flow.ChainID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProtocolVersion provides a mock function with given fields:
func (_m *Params) ProtocolVersion() (uint, error) {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Root provides a mock function with given fields:
func (_m *Params) Root() (*flow.Header, error) {
	ret := _m.Called()

	var r0 *flow.Header
	if rf, ok := ret.Get(0).(func() *flow.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Seal provides a mock function with given fields:
func (_m *Params) Seal() (*flow.Seal, error) {
	ret := _m.Called()

	var r0 *flow.Seal
	if rf, ok := ret.Get(0).(func() *flow.Seal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Seal)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SporkID provides a mock function with given fields:
func (_m *Params) SporkID() (flow.Identifier, error) {
	ret := _m.Called()

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func() flow.Identifier); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewParams interface {
	mock.TestingT
	Cleanup(func())
}

// NewParams creates a new instance of Params. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewParams(t mockConstructorTestingTNewParams) *Params {
	mock := &Params{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
