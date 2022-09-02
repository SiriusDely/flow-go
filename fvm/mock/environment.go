// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	atree "github.com/onflow/atree"
	ast "github.com/onflow/cadence/runtime/ast"

	attribute "go.opentelemetry.io/otel/attribute"

	cadence "github.com/onflow/cadence"

	common "github.com/onflow/cadence/runtime/common"

	flow "github.com/onflow/flow-go/model/flow"

	fvm "github.com/onflow/flow-go/fvm"

	interpreter "github.com/onflow/cadence/runtime/interpreter"

	mock "github.com/stretchr/testify/mock"

	oteltrace "go.opentelemetry.io/otel/trace"

	sema "github.com/onflow/cadence/runtime/sema"

	stdlib "github.com/onflow/cadence/runtime/stdlib"

	time "time"

	trace "github.com/onflow/flow-go/module/trace"

	zerolog "github.com/rs/zerolog"
)

// Environment is an autogenerated mock type for the Environment type
type Environment struct {
	mock.Mock
}

// AddAccountKey provides a mock function with given fields: address, publicKey, hashAlgo, weight
func (_m *Environment) AddAccountKey(address common.Address, publicKey *stdlib.PublicKey, hashAlgo sema.HashAlgorithm, weight int) (*stdlib.AccountKey, error) {
	ret := _m.Called(address, publicKey, hashAlgo, weight)

	var r0 *stdlib.AccountKey
	if rf, ok := ret.Get(0).(func(common.Address, *stdlib.PublicKey, sema.HashAlgorithm, int) *stdlib.AccountKey); ok {
		r0 = rf(address, publicKey, hashAlgo, weight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stdlib.AccountKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, *stdlib.PublicKey, sema.HashAlgorithm, int) error); ok {
		r1 = rf(address, publicKey, hashAlgo, weight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddEncodedAccountKey provides a mock function with given fields: address, publicKey
func (_m *Environment) AddEncodedAccountKey(address common.Address, publicKey []byte) error {
	ret := _m.Called(address, publicKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, []byte) error); ok {
		r0 = rf(address, publicKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AllocateStorageIndex provides a mock function with given fields: owner
func (_m *Environment) AllocateStorageIndex(owner []byte) (atree.StorageIndex, error) {
	ret := _m.Called(owner)

	var r0 atree.StorageIndex
	if rf, ok := ret.Get(0).(func([]byte) atree.StorageIndex); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(atree.StorageIndex)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BLSAggregatePublicKeys provides a mock function with given fields: keys
func (_m *Environment) BLSAggregatePublicKeys(keys []*stdlib.PublicKey) (*stdlib.PublicKey, error) {
	ret := _m.Called(keys)

	var r0 *stdlib.PublicKey
	if rf, ok := ret.Get(0).(func([]*stdlib.PublicKey) *stdlib.PublicKey); ok {
		r0 = rf(keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stdlib.PublicKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*stdlib.PublicKey) error); ok {
		r1 = rf(keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BLSAggregateSignatures provides a mock function with given fields: sigs
func (_m *Environment) BLSAggregateSignatures(sigs [][]byte) ([]byte, error) {
	ret := _m.Called(sigs)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([][]byte) []byte); ok {
		r0 = rf(sigs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([][]byte) error); ok {
		r1 = rf(sigs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BLSVerifyPOP provides a mock function with given fields: pk, s
func (_m *Environment) BLSVerifyPOP(pk *stdlib.PublicKey, s []byte) (bool, error) {
	ret := _m.Called(pk, s)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*stdlib.PublicKey, []byte) bool); ok {
		r0 = rf(pk, s)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*stdlib.PublicKey, []byte) error); ok {
		r1 = rf(pk, s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BorrowCadenceRuntime provides a mock function with given fields:
func (_m *Environment) BorrowCadenceRuntime() *fvm.ReusableCadenceRuntime {
	ret := _m.Called()

	var r0 *fvm.ReusableCadenceRuntime
	if rf, ok := ret.Get(0).(func() *fvm.ReusableCadenceRuntime); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fvm.ReusableCadenceRuntime)
		}
	}

	return r0
}

// Chain provides a mock function with given fields:
func (_m *Environment) Chain() flow.Chain {
	ret := _m.Called()

	var r0 flow.Chain
	if rf, ok := ret.Get(0).(func() flow.Chain); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Chain)
		}
	}

	return r0
}

// CreateAccount provides a mock function with given fields: payer
func (_m *Environment) CreateAccount(payer common.Address) (common.Address, error) {
	ret := _m.Called(payer)

	var r0 common.Address
	if rf, ok := ret.Get(0).(func(common.Address) common.Address); ok {
		r0 = rf(payer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(payer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecodeArgument provides a mock function with given fields: argument, argumentType
func (_m *Environment) DecodeArgument(argument []byte, argumentType cadence.Type) (cadence.Value, error) {
	ret := _m.Called(argument, argumentType)

	var r0 cadence.Value
	if rf, ok := ret.Get(0).(func([]byte, cadence.Type) cadence.Value); ok {
		r0 = rf(argument, argumentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cadence.Value)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, cadence.Type) error); ok {
		r1 = rf(argument, argumentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EmitEvent provides a mock function with given fields: _a0
func (_m *Environment) EmitEvent(_a0 cadence.Event) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(cadence.Event) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateUUID provides a mock function with given fields:
func (_m *Environment) GenerateUUID() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAvailableBalance provides a mock function with given fields: address
func (_m *Environment) GetAccountAvailableBalance(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountBalance provides a mock function with given fields: address
func (_m *Environment) GetAccountBalance(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountContractCode provides a mock function with given fields: address, name
func (_m *Environment) GetAccountContractCode(address common.Address, name string) ([]byte, error) {
	ret := _m.Called(address, name)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(common.Address, string) []byte); ok {
		r0 = rf(address, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, string) error); ok {
		r1 = rf(address, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountContractNames provides a mock function with given fields: address
func (_m *Environment) GetAccountContractNames(address common.Address) ([]string, error) {
	ret := _m.Called(address)

	var r0 []string
	if rf, ok := ret.Get(0).(func(common.Address) []string); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountKey provides a mock function with given fields: address, index
func (_m *Environment) GetAccountKey(address common.Address, index int) (*stdlib.AccountKey, error) {
	ret := _m.Called(address, index)

	var r0 *stdlib.AccountKey
	if rf, ok := ret.Get(0).(func(common.Address, int) *stdlib.AccountKey); ok {
		r0 = rf(address, index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stdlib.AccountKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, int) error); ok {
		r1 = rf(address, index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockAtHeight provides a mock function with given fields: height
func (_m *Environment) GetBlockAtHeight(height uint64) (stdlib.Block, bool, error) {
	ret := _m.Called(height)

	var r0 stdlib.Block
	if rf, ok := ret.Get(0).(func(uint64) stdlib.Block); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Get(0).(stdlib.Block)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(uint64) bool); ok {
		r1 = rf(height)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(uint64) error); ok {
		r2 = rf(height)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetCode provides a mock function with given fields: location
func (_m *Environment) GetCode(location common.Location) ([]byte, error) {
	ret := _m.Called(location)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(common.Location) []byte); ok {
		r0 = rf(location)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Location) error); ok {
		r1 = rf(location)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentBlockHeight provides a mock function with given fields:
func (_m *Environment) GetCurrentBlockHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProgram provides a mock function with given fields: _a0
func (_m *Environment) GetProgram(_a0 common.Location) (*interpreter.Program, error) {
	ret := _m.Called(_a0)

	var r0 *interpreter.Program
	if rf, ok := ret.Get(0).(func(common.Location) *interpreter.Program); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*interpreter.Program)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Location) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSigningAccounts provides a mock function with given fields:
func (_m *Environment) GetSigningAccounts() ([]common.Address, error) {
	ret := _m.Called()

	var r0 []common.Address
	if rf, ok := ret.Get(0).(func() []common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
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

// GetStorageCapacity provides a mock function with given fields: address
func (_m *Environment) GetStorageCapacity(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStorageUsed provides a mock function with given fields: address
func (_m *Environment) GetStorageUsed(address common.Address) (uint64, error) {
	ret := _m.Called(address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(common.Address) uint64); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValue provides a mock function with given fields: owner, key
func (_m *Environment) GetValue(owner []byte, key []byte) ([]byte, error) {
	ret := _m.Called(owner, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, []byte) []byte); ok {
		r0 = rf(owner, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(owner, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hash provides a mock function with given fields: data, tag, hashAlgorithm
func (_m *Environment) Hash(data []byte, tag string, hashAlgorithm sema.HashAlgorithm) ([]byte, error) {
	ret := _m.Called(data, tag, hashAlgorithm)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, string, sema.HashAlgorithm) []byte); ok {
		r0 = rf(data, tag, hashAlgorithm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, string, sema.HashAlgorithm) error); ok {
		r1 = rf(data, tag, hashAlgorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImplementationDebugLog provides a mock function with given fields: message
func (_m *Environment) ImplementationDebugLog(message string) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LimitAccountStorage provides a mock function with given fields:
func (_m *Environment) LimitAccountStorage() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Logger provides a mock function with given fields:
func (_m *Environment) Logger() *zerolog.Logger {
	ret := _m.Called()

	var r0 *zerolog.Logger
	if rf, ok := ret.Get(0).(func() *zerolog.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*zerolog.Logger)
		}
	}

	return r0
}

// MeterComputation provides a mock function with given fields: operationType, intensity
func (_m *Environment) MeterComputation(operationType common.ComputationKind, intensity uint) error {
	ret := _m.Called(operationType, intensity)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.ComputationKind, uint) error); ok {
		r0 = rf(operationType, intensity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MeterMemory provides a mock function with given fields: usage
func (_m *Environment) MeterMemory(usage common.MemoryUsage) error {
	ret := _m.Called(usage)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.MemoryUsage) error); ok {
		r0 = rf(usage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProgramLog provides a mock function with given fields: _a0
func (_m *Environment) ProgramLog(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecordTrace provides a mock function with given fields: operation, location, duration, attrs
func (_m *Environment) RecordTrace(operation string, location common.Location, duration time.Duration, attrs []attribute.KeyValue) {
	_m.Called(operation, location, duration, attrs)
}

// RemoveAccountContractCode provides a mock function with given fields: address, name
func (_m *Environment) RemoveAccountContractCode(address common.Address, name string) error {
	ret := _m.Called(address, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, string) error); ok {
		r0 = rf(address, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResolveLocation provides a mock function with given fields: identifiers, location
func (_m *Environment) ResolveLocation(identifiers []ast.Identifier, location common.Location) ([]sema.ResolvedLocation, error) {
	ret := _m.Called(identifiers, location)

	var r0 []sema.ResolvedLocation
	if rf, ok := ret.Get(0).(func([]ast.Identifier, common.Location) []sema.ResolvedLocation); ok {
		r0 = rf(identifiers, location)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]sema.ResolvedLocation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]ast.Identifier, common.Location) error); ok {
		r1 = rf(identifiers, location)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceOwnerChanged provides a mock function with given fields: _a0, resource, oldOwner, newOwner
func (_m *Environment) ResourceOwnerChanged(_a0 *interpreter.Interpreter, resource *interpreter.CompositeValue, oldOwner common.Address, newOwner common.Address) {
	_m.Called(_a0, resource, oldOwner, newOwner)
}

// ReturnCadenceRuntime provides a mock function with given fields: _a0
func (_m *Environment) ReturnCadenceRuntime(_a0 *fvm.ReusableCadenceRuntime) {
	_m.Called(_a0)
}

// RevokeAccountKey provides a mock function with given fields: address, index
func (_m *Environment) RevokeAccountKey(address common.Address, index int) (*stdlib.AccountKey, error) {
	ret := _m.Called(address, index)

	var r0 *stdlib.AccountKey
	if rf, ok := ret.Get(0).(func(common.Address, int) *stdlib.AccountKey); ok {
		r0 = rf(address, index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stdlib.AccountKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, int) error); ok {
		r1 = rf(address, index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RevokeEncodedAccountKey provides a mock function with given fields: address, index
func (_m *Environment) RevokeEncodedAccountKey(address common.Address, index int) ([]byte, error) {
	ret := _m.Called(address, index)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(common.Address, int) []byte); ok {
		r0 = rf(address, index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, int) error); ok {
		r1 = rf(address, index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAccountFrozen provides a mock function with given fields: address, frozen
func (_m *Environment) SetAccountFrozen(address common.Address, frozen bool) error {
	ret := _m.Called(address, frozen)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, bool) error); ok {
		r0 = rf(address, frozen)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetProgram provides a mock function with given fields: _a0, _a1
func (_m *Environment) SetProgram(_a0 common.Location, _a1 *interpreter.Program) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Location, *interpreter.Program) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetValue provides a mock function with given fields: owner, key, value
func (_m *Environment) SetValue(owner []byte, key []byte, value []byte) error {
	ret := _m.Called(owner, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte, []byte) error); ok {
		r0 = rf(owner, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartExtensiveTracingSpanFromRoot provides a mock function with given fields: name
func (_m *Environment) StartExtensiveTracingSpanFromRoot(name trace.SpanName) oteltrace.Span {
	ret := _m.Called(name)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(trace.SpanName) oteltrace.Span); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	return r0
}

// StartSpanFromRoot provides a mock function with given fields: name
func (_m *Environment) StartSpanFromRoot(name trace.SpanName) oteltrace.Span {
	ret := _m.Called(name)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(trace.SpanName) oteltrace.Span); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	return r0
}

// UnsafeRandom provides a mock function with given fields:
func (_m *Environment) UnsafeRandom() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAccountContractCode provides a mock function with given fields: address, name, code
func (_m *Environment) UpdateAccountContractCode(address common.Address, name string, code []byte) error {
	ret := _m.Called(address, name, code)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, string, []byte) error); ok {
		r0 = rf(address, name, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VM provides a mock function with given fields:
func (_m *Environment) VM() *fvm.VirtualMachine {
	ret := _m.Called()

	var r0 *fvm.VirtualMachine
	if rf, ok := ret.Get(0).(func() *fvm.VirtualMachine); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fvm.VirtualMachine)
		}
	}

	return r0
}

// ValidatePublicKey provides a mock function with given fields: key
func (_m *Environment) ValidatePublicKey(key *stdlib.PublicKey) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*stdlib.PublicKey) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValueExists provides a mock function with given fields: owner, key
func (_m *Environment) ValueExists(owner []byte, key []byte) (bool, error) {
	ret := _m.Called(owner, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, []byte) bool); ok {
		r0 = rf(owner, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(owner, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifySignature provides a mock function with given fields: signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm
func (_m *Environment) VerifySignature(signature []byte, tag string, signedData []byte, publicKey []byte, signatureAlgorithm sema.SignatureAlgorithm, hashAlgorithm sema.HashAlgorithm) (bool, error) {
	ret := _m.Called(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, string, []byte, []byte, sema.SignatureAlgorithm, sema.HashAlgorithm) bool); ok {
		r0 = rf(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, string, []byte, []byte, sema.SignatureAlgorithm, sema.HashAlgorithm) error); ok {
		r1 = rf(signature, tag, signedData, publicKey, signatureAlgorithm, hashAlgorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEnvironment interface {
	mock.TestingT
	Cleanup(func())
}

// NewEnvironment creates a new instance of Environment. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEnvironment(t mockConstructorTestingTNewEnvironment) *Environment {
	mock := &Environment{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
