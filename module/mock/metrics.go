// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

import time "time"

// Metrics is an autogenerated mock type for the Metrics type
type Metrics struct {
	mock.Mock
}

// BadgerDBSize provides a mock function with given fields: sizeBytes
func (_m *Metrics) BadgerDBSize(sizeBytes int64) {
	_m.Called(sizeBytes)
}

// CollectionGuaranteed provides a mock function with given fields: collection
func (_m *Metrics) CollectionGuaranteed(collection flow.LightCollection) {
	_m.Called(collection)
}

// CollectionProposed provides a mock function with given fields: collection
func (_m *Metrics) CollectionProposed(collection flow.LightCollection) {
	_m.Called(collection)
}

// CollectionsInFinalizedBlock provides a mock function with given fields: count
func (_m *Metrics) CollectionsInFinalizedBlock(count int) {
	_m.Called(count)
}

// CollectionsPerBlock provides a mock function with given fields: count
func (_m *Metrics) CollectionsPerBlock(count int) {
	_m.Called(count)
}

// ExecutionGasUsedPerBlock provides a mock function with given fields: gas
func (_m *Metrics) ExecutionGasUsedPerBlock(gas uint64) {
	_m.Called(gas)
}

// ExecutionStateReadsPerBlock provides a mock function with given fields: reads
func (_m *Metrics) ExecutionStateReadsPerBlock(reads uint64) {
	_m.Called(reads)
}

// ExecutionStateStorageDiskTotal provides a mock function with given fields: bytes
func (_m *Metrics) ExecutionStateStorageDiskTotal(bytes int64) {
	_m.Called(bytes)
}

// ExecutionStorageStateCommitment provides a mock function with given fields: bytes
func (_m *Metrics) ExecutionStorageStateCommitment(bytes int64) {
	_m.Called(bytes)
}

// FinalizedBlocks provides a mock function with given fields: count
func (_m *Metrics) FinalizedBlocks(count int) {
	_m.Called(count)
}

// FinishBlockReceivedToExecuted provides a mock function with given fields: blockID
func (_m *Metrics) FinishBlockReceivedToExecuted(blockID flow.Identifier) {
	_m.Called(blockID)
}

// FinishBlockToSeal provides a mock function with given fields: blockID
func (_m *Metrics) FinishBlockToSeal(blockID flow.Identifier) {
	_m.Called(blockID)
}

// FinishCollectionToFinalized provides a mock function with given fields: collectionID
func (_m *Metrics) FinishCollectionToFinalized(collectionID flow.Identifier) {
	_m.Called(collectionID)
}

// HotStuffBusyDuration provides a mock function with given fields: duration, event
func (_m *Metrics) HotStuffBusyDuration(duration time.Duration, event string) {
	_m.Called(duration, event)
}

func (_m *Metrics) HotStuffBusySecondsTotalAdd(duration time.Duration, event string) {
	_m.Called(duration, event)
}

// HotStuffIdleDuration provides a mock function with given fields: duration
func (_m *Metrics) HotStuffIdleDuration(duration time.Duration) {
	_m.Called(duration)
}

// HotStuffWaitDuration provides a mock function with given fields: duration, event
func (_m *Metrics) HotStuffWaitDuration(duration time.Duration, event string) {
	_m.Called(duration, event)
}

// NetworkMessageSent provides a mock function with given fields: sizeBytes
func (_m *Metrics) NetworkMessageSent(sizeBytes int) {
	_m.Called(sizeBytes)
}

// NewestKnownQC provides a mock function with given fields: view
func (_m *Metrics) NewestKnownQC(view uint64) {
	_m.Called(view)
}

// OnChunkDataAdded provides a mock function with given fields: chunkID, size
func (_m *Metrics) OnChunkDataAdded(chunkID flow.Identifier, size float64) {
	_m.Called(chunkID, size)
}

// OnChunkDataRemoved provides a mock function with given fields: chunkID, size
func (_m *Metrics) OnChunkDataRemoved(chunkID flow.Identifier, size float64) {
	_m.Called(chunkID, size)
}

// OnChunkVerificationFinished provides a mock function with given fields: chunkID
func (_m *Metrics) OnChunkVerificationFinished(chunkID flow.Identifier) {
	_m.Called(chunkID)
}

// OnChunkVerificationStarted provides a mock function with given fields: chunkID
func (_m *Metrics) OnChunkVerificationStarted(chunkID flow.Identifier) {
	_m.Called(chunkID)
}

// OnResultApproval provides a mock function with given fields:
func (_m *Metrics) OnResultApproval() {
	_m.Called()
}

// SealsInFinalizedBlock provides a mock function with given fields: count
func (_m *Metrics) SealsInFinalizedBlock(count int) {
	_m.Called(count)
}

// StartBlockReceivedToExecuted provides a mock function with given fields: blockID
func (_m *Metrics) StartBlockReceivedToExecuted(blockID flow.Identifier) {
	_m.Called(blockID)
}

// StartBlockToSeal provides a mock function with given fields: blockID
func (_m *Metrics) StartBlockToSeal(blockID flow.Identifier) {
	_m.Called(blockID)
}

// StartCollectionToFinalized provides a mock function with given fields: collectionID
func (_m *Metrics) StartCollectionToFinalized(collectionID flow.Identifier) {
	_m.Called(collectionID)
}

// StartNewView provides a mock function with given fields: view
func (_m *Metrics) StartNewView(view uint64) {
	_m.Called(view)
}

// TransactionReceived provides a mock function with given fields: txID
func (_m *Metrics) TransactionReceived(txID flow.Identifier) {
	_m.Called(txID)
}
