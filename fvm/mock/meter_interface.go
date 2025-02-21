// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	common "github.com/onflow/cadence/runtime/common"

	mock "github.com/stretchr/testify/mock"
)

// MeterInterface is an autogenerated mock type for the MeterInterface type
type MeterInterface struct {
	mock.Mock
}

// Meter provides a mock function with given fields: _a0, _a1
func (_m *MeterInterface) Meter(_a0 common.ComputationKind, _a1 uint) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.ComputationKind, uint) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMeterInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMeterInterface creates a new instance of MeterInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMeterInterface(t mockConstructorTestingTNewMeterInterface) *MeterInterface {
	mock := &MeterInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
