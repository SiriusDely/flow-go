// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

// Collections is an autogenerated mock type for the Collections type
type Collections struct {
	mock.Mock
}

// ByFingerprint provides a mock function with given fields: fingerprint
func (_m *Collections) ByFingerprint(fingerprint flow.Fingerprint) (*flow.Collection, error) {
	ret := _m.Called(fingerprint)

	var r0 *flow.Collection
	if rf, ok := ret.Get(0).(func(flow.Fingerprint) *flow.Collection); ok {
		r0 = rf(fingerprint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Collection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Fingerprint) error); ok {
		r1 = rf(fingerprint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: tx
func (_m *Collections) Insert(tx *flow.Collection) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Collection) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: fingerprint
func (_m *Collections) Remove(fingerprint flow.Fingerprint) error {
	ret := _m.Called(fingerprint)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Fingerprint) error); ok {
		r0 = rf(fingerprint)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
