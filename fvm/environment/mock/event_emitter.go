// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	cadence "github.com/onflow/cadence"

	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// EventEmitter is an autogenerated mock type for the EventEmitter type
type EventEmitter struct {
	mock.Mock
}

// EmitEvent provides a mock function with given fields: event
func (_m *EventEmitter) EmitEvent(event cadence.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(cadence.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Events provides a mock function with given fields:
func (_m *EventEmitter) Events() []flow.Event {
	ret := _m.Called()

	var r0 []flow.Event
	if rf, ok := ret.Get(0).(func() []flow.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Event)
		}
	}

	return r0
}

// ServiceEvents provides a mock function with given fields:
func (_m *EventEmitter) ServiceEvents() []flow.Event {
	ret := _m.Called()

	var r0 []flow.Event
	if rf, ok := ret.Get(0).(func() []flow.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Event)
		}
	}

	return r0
}

type mockConstructorTestingTNewEventEmitter interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventEmitter creates a new instance of EventEmitter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventEmitter(t mockConstructorTestingTNewEventEmitter) *EventEmitter {
	mock := &EventEmitter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
