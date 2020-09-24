package chunks

import (
	"bytes"
	"fmt"

	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/engine/verification"
	"github.com/onflow/flow-go/fvm"
	"github.com/onflow/flow-go/fvm/state"
	chmodels "github.com/onflow/flow-go/model/chunks"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/storage/ledger"
	"github.com/onflow/flow-go/storage/ledger/ptrie"
)

type VirtualMachine interface {
	Run(fvm.Context, fvm.Procedure, state.Ledger) error
}

// ChunkVerifier is a verifier based on the current definitions of the flow network
type ChunkVerifier struct {
	vm    VirtualMachine
	vmCtx fvm.Context
}

// NewChunkVerifier creates a chunk verifier containing a flow virtual machine
func NewChunkVerifier(vm VirtualMachine, vmCtx fvm.Context) *ChunkVerifier {
	return &ChunkVerifier{
		vm:    vm,
		vmCtx: vmCtx,
	}
}

// Verify verifies a given VerifiableChunk corresponding to a non-system chunk.
// by executing it and checking the final state commitment
// It returns a Spock Secret as a byte array, verification fault of the chunk, and an error.
// Note: Verify should only be executed on non-system chunks. It returns an error if it is invoked on
// system chunks.
func (fcv *ChunkVerifier) Verify(vc *verification.VerifiableChunkData) ([]byte, chmodels.ChunkFault, error) {
	if vc.IsSystemChunk {
		return nil, nil, fmt.Errorf("wrong method invoked for verifying system chunk")
	}

	transactions := make([]*fvm.TransactionProcedure, 0)
	for _, txBody := range vc.Collection.Transactions {
		tx := fvm.Transaction(txBody)
		transactions = append(transactions, tx)
	}

	return fcv.verifyTransactions(vc.Chunk, vc.ChunkDataPack, vc.Result, vc.Header, transactions, vc.EndState)
}

// VerifySystemChunk verifies a given VerifiableChunk corresponding to a system chunk.
// by executing it and checking the final state commitment
// It returns a Spock Secret as a byte array, verification fault of the chunk, and an error.
// Note: SystemChunkVerify should only be executed on system chunks. It returns an error if it is invoked on
// non-system chunks.
func (fcv *ChunkVerifier) SystemChunkVerify(vc *verification.VerifiableChunkData) ([]byte, chmodels.ChunkFault, error) {
	if !vc.IsSystemChunk {
		return nil, nil, fmt.Errorf("wrong method invoked for verifying non-system chunk")
	}

	// transaction body of system chunk
	txBody := fvm.SystemChunkTransaction(fcv.vmCtx.Chain.ServiceAddress())
	tx := fvm.Transaction(txBody)
	transactions := []*fvm.TransactionProcedure{tx}

	return fcv.verifyTransactions(vc.Chunk, vc.ChunkDataPack, vc.Result, vc.Header, transactions, vc.EndState)
}

func (fcv *ChunkVerifier) verifyTransactions(chunk *flow.Chunk,
	chunkDataPack *flow.ChunkDataPack,
	result *flow.ExecutionResult,
	header *flow.Header,
	transactions []*fvm.TransactionProcedure,
	endState flow.StateCommitment) ([]byte, chmodels.ChunkFault, error) {

	// TODO check collection hash to match
	// TODO check datapack hash to match
	// TODO check the number of transactions and computation used

	chIndex := chunk.Index
	execResID := result.ID()

	// build a block context
	blockCtx := fvm.NewContextFromParent(fcv.vmCtx, fvm.WithBlockHeader(header))

	if chunkDataPack == nil {
		return nil, nil, fmt.Errorf("missing chunk data pack")
	}

	// constructing a partial trie given chunk data package
	psmt, err := ptrie.NewPSMT(chunkDataPack.StartState,
		ledger.RegisterKeySize,
		chunkDataPack.Registers(),
		chunkDataPack.Values(),
		chunkDataPack.Proofs(),
	)
	if err != nil {
		// TODO provide more details based on the error type
		return nil, chmodels.NewCFInvalidVerifiableChunk("error constructing partial trie: ", err, chIndex, execResID),
			nil
	}

	// chunk view construction
	// unknown register tracks access to parts of the partial trie which
	// are not expanded and values are unknown.
	unknownRegTouch := make(map[string]bool)
	regMap := chunkDataPack.GetRegisterValues()
	getRegister := func(owner, controller, key string) (flow.RegisterValue, error) {
		// check if register has been provided in the chunk data pack
		k := state.RegisterID(owner, controller, key)

		val, ok := regMap[string(k)]
		if !ok {
			unknownRegTouch[string(k)] = true
			return nil, fmt.Errorf("missing register")
		}
		return val, nil
	}

	chunkView := delta.NewView(getRegister)

	// executes all transactions in this chunk
	for i, tx := range transactions {
		txView := chunkView.NewChild()

		// tx := fvm.Transaction(txBody)

		err := fcv.vm.Run(blockCtx, tx, txView)
		if err != nil {
			// this covers unexpected and very rare cases (e.g. system memory issues...),
			// so we shouldn't be here even if transaction naturally fails (e.g. permission, runtime ... )
			return nil, nil, fmt.Errorf("failed to execute transaction: %d (%w)", i, err)
		}

		if tx.Err == nil {
			// if tx is successful, we apply changes to the chunk view by merging the txView into chunk view
			chunkView.MergeView(txView)
		}
	}

	// check read access to unknown registers
	if len(unknownRegTouch) > 0 {
		var missingRegs []string
		for key := range unknownRegTouch {
			missingRegs = append(missingRegs, key)
		}
		return nil, chmodels.NewCFMissingRegisterTouch(missingRegs, chIndex, execResID), nil
	}

	// applying chunk delta (register updates at chunk level) to the partial trie
	// this returns the expected end state commitment after updates and the list of
	// register keys that was not provided by the chunk data package (err).
	regs, values := chunkView.Delta().RegisterUpdates()
	expEndStateComm, failedKeys, err := psmt.Update(regs, values)
	if err != nil {
		return nil, chmodels.NewCFMissingRegisterTouch(failedKeys, chIndex, execResID), nil
	}

	// TODO check if exec node provided register touches that was not used (no read and no update)
	// check if the end state commitment mentioned in the chunk matches
	// what the partial trie is providing.
	if !bytes.Equal(expEndStateComm, endState) {
		return nil, chmodels.NewCFNonMatchingFinalState(expEndStateComm, endState, chIndex, execResID), nil
	}
	return chunkView.SpockSecret(), nil, nil
}
