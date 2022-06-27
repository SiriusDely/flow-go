// Code generated by mockery v2.13.0. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// IdentifierMap is an autogenerated mock type for the IdentifierMap type
type IdentifierMap struct {
	mock.Mock
}

// Append provides a mock function with given fields: key, id
func (_m *IdentifierMap) Append(key flow.Identifier, id flow.Identifier) error {
	ret := _m.Called(key, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) error); ok {
		r0 = rf(key, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *IdentifierMap) Get(key flow.Identifier) ([]flow.Identifier, bool) {
	ret := _m.Called(key)

	var r0 []flow.Identifier
	if rf, ok := ret.Get(0).(func(flow.Identifier) []flow.Identifier); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Identifier)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Has provides a mock function with given fields: key
func (_m *IdentifierMap) Has(key flow.Identifier) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Keys provides a mock function with given fields:
func (_m *IdentifierMap) Keys() ([]flow.Identifier, bool) {
	ret := _m.Called()

	var r0 []flow.Identifier
	if rf, ok := ret.Get(0).(func() []flow.Identifier); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Identifier)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Rem provides a mock function with given fields: key
func (_m *IdentifierMap) Rem(key flow.Identifier) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemIdFromKey provides a mock function with given fields: key, id
func (_m *IdentifierMap) RemIdFromKey(key flow.Identifier, id flow.Identifier) error {
	ret := _m.Called(key, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) error); ok {
		r0 = rf(key, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *IdentifierMap) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

type NewIdentifierMapT interface {
	mock.TestingT
	Cleanup(func())
}

// NewIdentifierMap creates a new instance of IdentifierMap. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIdentifierMap(t NewIdentifierMapT) *IdentifierMap {
	mock := &IdentifierMap{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
