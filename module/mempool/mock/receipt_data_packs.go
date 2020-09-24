// Code generated by mockery v1.0.0. DO NOT EDIT.

package mempool

import (
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	verification "github.com/onflow/flow-go/model/verification"
)

// ReceiptDataPacks is an autogenerated mock type for the ReceiptDataPacks type
type ReceiptDataPacks struct {
	mock.Mock
}

// Add provides a mock function with given fields: rdp
func (_m *ReceiptDataPacks) Add(rdp *verification.ReceiptDataPack) bool {
	ret := _m.Called(rdp)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*verification.ReceiptDataPack) bool); ok {
		r0 = rf(rdp)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// All provides a mock function with given fields:
func (_m *ReceiptDataPacks) All() []*verification.ReceiptDataPack {
	ret := _m.Called()

	var r0 []*verification.ReceiptDataPack
	if rf, ok := ret.Get(0).(func() []*verification.ReceiptDataPack); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*verification.ReceiptDataPack)
		}
	}

	return r0
}

// Get provides a mock function with given fields: rdpID
func (_m *ReceiptDataPacks) Get(rdpID flow.Identifier) (*verification.ReceiptDataPack, bool) {
	ret := _m.Called(rdpID)

	var r0 *verification.ReceiptDataPack
	if rf, ok := ret.Get(0).(func(flow.Identifier) *verification.ReceiptDataPack); ok {
		r0 = rf(rdpID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*verification.ReceiptDataPack)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(rdpID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Has provides a mock function with given fields: rdpID
func (_m *ReceiptDataPacks) Has(rdpID flow.Identifier) bool {
	ret := _m.Called(rdpID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(rdpID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Rem provides a mock function with given fields: rdpID
func (_m *ReceiptDataPacks) Rem(rdpID flow.Identifier) bool {
	ret := _m.Called(rdpID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(rdpID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *ReceiptDataPacks) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}
