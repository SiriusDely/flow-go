package computation

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/engine/execution"
	state2 "github.com/onflow/flow-go/engine/execution/state"
	unittest2 "github.com/onflow/flow-go/engine/execution/state/unittest"
	"github.com/onflow/flow-go/ledger/complete"
	"github.com/onflow/flow-go/ledger/complete/wal/fixtures"
	requesterunit "github.com/onflow/flow-go/module/state_synchronization/requester/unittest"

	"github.com/onflow/flow-go/engine/execution/computation/committer"
	"github.com/onflow/flow-go/engine/execution/computation/computer"
	"github.com/onflow/flow-go/engine/execution/computation/computer/uploader"
	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/engine/execution/testutil"
	"github.com/onflow/flow-go/fvm"
	fvmErrors "github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/executiondatasync/execution_data"
	"github.com/onflow/flow-go/module/executiondatasync/provider"
	"github.com/onflow/flow-go/module/executiondatasync/tracker"
	mocktracker "github.com/onflow/flow-go/module/executiondatasync/tracker/mock"
	"github.com/onflow/flow-go/module/mempool/entity"
	"github.com/onflow/flow-go/module/metrics"
	module "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/module/trace"
	"github.com/onflow/flow-go/utils/unittest"
)

var scriptLogThreshold = 1 * time.Second

func TestComputeBlockWithStorage(t *testing.T) {
	chain := flow.Mainnet.Chain()

	vm := fvm.NewVM()
	execCtx := fvm.NewContext(fvm.WithChain(chain))

	privateKeys, err := testutil.GenerateAccountPrivateKeys(2)
	require.NoError(t, err)

	ledger := testutil.RootBootstrappedLedger(vm, execCtx)
	accounts, err := testutil.CreateAccounts(vm, ledger, programs.NewEmptyBlockPrograms(), privateKeys, chain)
	require.NoError(t, err)

	tx1 := testutil.DeployCounterContractTransaction(accounts[0], chain)
	tx1.SetProposalKey(chain.ServiceAddress(), 0, 0).
		SetGasLimit(1000).
		SetPayer(chain.ServiceAddress())

	err = testutil.SignPayload(tx1, accounts[0], privateKeys[0])
	require.NoError(t, err)

	err = testutil.SignEnvelope(tx1, chain.ServiceAddress(), unittest.ServiceAccountPrivateKey)
	require.NoError(t, err)

	tx2 := testutil.CreateCounterTransaction(accounts[0], accounts[1])
	tx2.SetProposalKey(chain.ServiceAddress(), 0, 0).
		SetGasLimit(1000).
		SetPayer(chain.ServiceAddress())

	err = testutil.SignPayload(tx2, accounts[1], privateKeys[1])
	require.NoError(t, err)

	err = testutil.SignEnvelope(tx2, chain.ServiceAddress(), unittest.ServiceAccountPrivateKey)
	require.NoError(t, err)

	transactions := []*flow.TransactionBody{tx1, tx2}

	col := flow.Collection{Transactions: transactions}

	guarantee := flow.CollectionGuarantee{
		CollectionID: col.ID(),
		Signature:    nil,
	}

	block := flow.Block{
		Header: &flow.Header{
			View: 42,
		},
		Payload: &flow.Payload{
			Guarantees: []*flow.CollectionGuarantee{&guarantee},
		},
	}

	executableBlock := &entity.ExecutableBlock{
		Block: &block,
		CompleteCollections: map[flow.Identifier]*entity.CompleteCollection{
			guarantee.ID(): {
				Guarantee:    &guarantee,
				Transactions: transactions,
			},
		},
		StartState: unittest.StateCommitmentPointerFixture(),
	}

	me := new(module.Local)
	me.On("NodeID").Return(flow.ZeroID)

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	blockComputer, err := computer.NewBlockComputer(vm, execCtx, metrics.NewNoopCollector(), trace.NewNoopTracer(), zerolog.Nop(), committer.NewNoopViewCommitter(), prov)
	require.NoError(t, err)

	programsCache, err := programs.NewChainPrograms(10)
	require.NoError(t, err)

	engine := &Manager{
		blockComputer: blockComputer,
		me:            me,
		programsCache: programsCache,
		tracer:        trace.NewNoopTracer(),
	}

	view := delta.NewView(ledger.Get)
	blockView := view.NewChild()

	returnedComputationResult, err := engine.ComputeBlock(context.Background(), executableBlock, blockView)
	require.NoError(t, err)

	require.NotEmpty(t, blockView.(*delta.View).Delta())
	require.Len(t, returnedComputationResult.StateSnapshots, 1+1) // 1 coll + 1 system chunk
	assert.NotEmpty(t, returnedComputationResult.StateSnapshots[0].Delta)
	assert.True(t, returnedComputationResult.ComputationUsed > 0)
}

func TestComputeBlock_Uploader(t *testing.T) {

	noopCollector := &metrics.NoopCollector{}

	ledger, err := complete.NewLedger(&fixtures.NoopWAL{}, 10, noopCollector, zerolog.Nop(), complete.DefaultPathFinderVersion)
	require.NoError(t, err)

	compactor := fixtures.NewNoopCompactor(ledger)
	<-compactor.Ready()
	defer func() {
		<-ledger.Done()
		<-compactor.Done()
	}()

	me := new(module.Local)
	me.On("NodeID").Return(flow.ZeroID)

	computationResult := unittest2.ComputationResultFixture([][]flow.Identifier{
		{unittest.IdentifierFixture()},
		{unittest.IdentifierFixture()},
	})

	blockComputer := &FakeBlockComputer{
		computationResult: computationResult,
	}

	programsCache, err := programs.NewChainPrograms(10)
	require.NoError(t, err)

	fakeUploader := &FakeUploader{}

	manager := &Manager{
		blockComputer: blockComputer,
		me:            me,
		programsCache: programsCache,
		uploaders:     []uploader.Uploader{fakeUploader},
		tracer:        trace.NewNoopTracer(),
	}

	view := delta.NewView(state2.LedgerGetRegister(ledger, flow.StateCommitment(ledger.InitialState())))
	blockView := view.NewChild()

	_, err = manager.ComputeBlock(context.Background(), computationResult.ExecutableBlock, blockView)
	require.NoError(t, err)

	retrievedResult, has := fakeUploader.data[computationResult.ExecutableBlock.ID()]
	require.True(t, has)

	assert.Equal(t, computationResult, retrievedResult)
}

func TestExecuteScript(t *testing.T) {

	logger := zerolog.Nop()

	execCtx := fvm.NewContext(fvm.WithLogger(logger))

	me := new(module.Local)
	me.On("NodeID").Return(flow.ZeroID)

	vm := fvm.NewVM()

	ledger := testutil.RootBootstrappedLedger(vm, execCtx, fvm.WithExecutionMemoryLimit(math.MaxUint64))

	view := delta.NewView(ledger.Get)

	scriptView := view.NewChild()

	script := []byte(fmt.Sprintf(
		`
			import FungibleToken from %s

			pub fun main() {}
		`,
		fvm.FungibleTokenAddress(execCtx.Chain).HexWithPrefix(),
	))

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	engine, err := New(logger,
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		me,
		nil,
		execCtx,
		committer.NewNoopViewCommitter(),
		nil,
		prov,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       scriptLogThreshold,
			ScriptExecutionTimeLimit: DefaultScriptExecutionTimeLimit,
		},
	)
	require.NoError(t, err)

	header := unittest.BlockHeaderFixture()
	_, err = engine.ExecuteScript(context.Background(), script, nil, header, scriptView)
	require.NoError(t, err)
}

// Balance script used to swallow errors, which meant that even if the view was empty, a script that did nothing but get
// the balance of an account would succeed and return 0.
func TestExecuteScript_BalanceScriptFailsIfViewIsEmpty(t *testing.T) {

	logger := zerolog.Nop()

	execCtx := fvm.NewContext(fvm.WithLogger(logger))

	me := new(module.Local)
	me.On("NodeID").Return(flow.ZeroID)

	view := delta.NewView(func(owner, key string) (flow.RegisterValue, error) {
		return nil, fmt.Errorf("error getting register")
	})

	scriptView := view.NewChild()

	script := []byte(fmt.Sprintf(
		`
			pub fun main(): UFix64 {
				return getAccount(%s).balance
			}
		`,
		fvm.FungibleTokenAddress(execCtx.Chain).HexWithPrefix(),
	))

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	engine, err := New(logger,
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		me,
		nil,
		execCtx,
		committer.NewNoopViewCommitter(),
		nil,
		prov,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       scriptLogThreshold,
			ScriptExecutionTimeLimit: DefaultScriptExecutionTimeLimit,
		},
	)
	require.NoError(t, err)

	header := unittest.BlockHeaderFixture()
	_, err = engine.ExecuteScript(context.Background(), script, nil, header, scriptView)
	require.ErrorContains(t, err, "error getting register")
}

func TestExecuteScripPanicsAreHandled(t *testing.T) {

	ctx := fvm.NewContext()

	buffer := &bytes.Buffer{}
	log := zerolog.New(buffer)

	header := unittest.BlockHeaderFixture()

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	manager, err := New(log,
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		nil,
		nil,
		ctx,
		committer.NewNoopViewCommitter(),
		nil,
		prov,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       scriptLogThreshold,
			ScriptExecutionTimeLimit: DefaultScriptExecutionTimeLimit,
			NewCustomVirtualMachine: func() computer.VirtualMachine {
				return &PanickingVM{}
			},
		},
	)
	require.NoError(t, err)

	_, err = manager.ExecuteScript(context.Background(), []byte("whatever"), nil, header, noopView())

	require.Error(t, err)

	require.Contains(t, buffer.String(), "Verunsicherung")
}

func TestExecuteScript_LongScriptsAreLogged(t *testing.T) {

	ctx := fvm.NewContext()

	buffer := &bytes.Buffer{}
	log := zerolog.New(buffer)

	header := unittest.BlockHeaderFixture()

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	manager, err := New(log,
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		nil,
		nil,
		ctx,
		committer.NewNoopViewCommitter(),
		nil,
		prov,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       1 * time.Millisecond,
			ScriptExecutionTimeLimit: DefaultScriptExecutionTimeLimit,
			NewCustomVirtualMachine: func() computer.VirtualMachine {
				return &LongRunningVM{duration: 2 * time.Millisecond}
			},
		},
	)
	require.NoError(t, err)

	_, err = manager.ExecuteScript(context.Background(), []byte("whatever"), nil, header, noopView())

	require.NoError(t, err)

	require.Contains(t, buffer.String(), "exceeded threshold")
}

func TestExecuteScript_ShortScriptsAreNotLogged(t *testing.T) {

	ctx := fvm.NewContext()

	buffer := &bytes.Buffer{}
	log := zerolog.New(buffer)

	header := unittest.BlockHeaderFixture()

	bservice := requesterunit.MockBlobService(blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore())))
	trackerStorage := new(mocktracker.Storage)
	trackerStorage.On("Update", mock.Anything).Return(func(fn tracker.UpdateFn) error {
		return fn(func(uint64, ...cid.Cid) error { return nil })
	})

	prov := provider.NewProvider(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		execution_data.DefaultSerializer,
		bservice,
		trackerStorage,
	)

	manager, err := New(log,
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		nil,
		nil,
		ctx,
		committer.NewNoopViewCommitter(),
		nil,
		prov,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       1 * time.Second,
			ScriptExecutionTimeLimit: DefaultScriptExecutionTimeLimit,
			NewCustomVirtualMachine: func() computer.VirtualMachine {
				return &LongRunningVM{duration: 0}
			},
		},
	)
	require.NoError(t, err)

	_, err = manager.ExecuteScript(context.Background(), []byte("whatever"), nil, header, noopView())

	require.NoError(t, err)

	require.NotContains(t, buffer.String(), "exceeded threshold")
}

type PanickingVM struct{}

func (p *PanickingVM) Run(f fvm.Context, procedure fvm.Procedure, view state.View, p2 *programs.Programs) error {
	return p.RunV2(f, procedure, view)
}

func (p *PanickingVM) RunV2(f fvm.Context, procedure fvm.Procedure, view state.View) error {
	panic("panic, but expected with sentinel for test: Verunsicherung ")
}

func (p *PanickingVM) GetAccount(f fvm.Context, address flow.Address, view state.View, p2 *programs.Programs) (*flow.Account, error) {
	panic("not expected")
}

func (p *PanickingVM) GetAccountV2(f fvm.Context, address flow.Address, view state.View) (*flow.Account, error) {
	panic("not expected")
}

type LongRunningVM struct {
	duration time.Duration
}

func (l *LongRunningVM) Run(f fvm.Context, procedure fvm.Procedure, view state.View, p2 *programs.Programs) error {
	return l.RunV2(f, procedure, view)
}

func (l *LongRunningVM) RunV2(f fvm.Context, procedure fvm.Procedure, view state.View) error {
	time.Sleep(l.duration)
	// satisfy value marshaller
	if scriptProcedure, is := procedure.(*fvm.ScriptProcedure); is {
		scriptProcedure.Value = cadence.NewVoid()
	}

	return nil
}

func (l *LongRunningVM) GetAccount(f fvm.Context, address flow.Address, view state.View, p2 *programs.Programs) (*flow.Account, error) {
	panic("not expected")
}

func (l *LongRunningVM) GetAccountV2(f fvm.Context, address flow.Address, view state.View) (*flow.Account, error) {
	panic("not expected")
}

type FakeBlockComputer struct {
	computationResult *execution.ComputationResult
}

func (f *FakeBlockComputer) ExecuteBlock(context.Context, *entity.ExecutableBlock, state.View, *programs.BlockPrograms) (*execution.ComputationResult, error) {
	return f.computationResult, nil
}

type FakeUploader struct {
	data map[flow.Identifier]*execution.ComputationResult
}

func (f *FakeUploader) Upload(computationResult *execution.ComputationResult) error {
	if f.data == nil {
		f.data = make(map[flow.Identifier]*execution.ComputationResult)
	}
	f.data[computationResult.ExecutableBlock.ID()] = computationResult
	return nil
}

func noopView() *delta.View {
	return delta.NewView(func(_, _ string) (flow.RegisterValue, error) {
		return nil, nil
	})
}

func TestExecuteScriptTimeout(t *testing.T) {

	timeout := 1 * time.Millisecond
	manager, err := New(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		nil,
		nil,
		fvm.NewContext(),
		committer.NewNoopViewCommitter(),
		nil,
		nil,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       DefaultScriptLogThreshold,
			ScriptExecutionTimeLimit: timeout,
		},
	)

	require.NoError(t, err)

	script := []byte(`
	pub fun main(): Int {
		var i = 0
		while i < 10000 {
			i = i + 1
		}
		return i
	}
	`)

	header := unittest.BlockHeaderFixture()
	value, err := manager.ExecuteScript(context.Background(), script, nil, header, noopView())

	require.Error(t, err)
	require.Nil(t, value)
	require.Contains(t, err.Error(), fvmErrors.ErrCodeScriptExecutionTimedOutError.String())
}

func TestExecuteScriptCancelled(t *testing.T) {

	timeout := 30 * time.Second
	manager, err := New(
		zerolog.Nop(),
		metrics.NewNoopCollector(),
		trace.NewNoopTracer(),
		nil,
		nil,
		fvm.NewContext(),
		committer.NewNoopViewCommitter(),
		nil,
		nil,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       DefaultScriptLogThreshold,
			ScriptExecutionTimeLimit: timeout,
		},
	)

	require.NoError(t, err)

	script := []byte(`
	pub fun main(): Int {
		var i = 0
		var j = 0 
		while i < 10000000 {
			i = i + 1
			j = i + j
		}
		return i
	}
	`)

	var value []byte
	var wg sync.WaitGroup
	reqCtx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		header := unittest.BlockHeaderFixture()
		value, err = manager.ExecuteScript(reqCtx, script, nil, header, noopView())
		wg.Done()
	}()
	cancel()
	wg.Wait()
	require.Nil(t, value)
	require.Contains(t, err.Error(), fvmErrors.ErrCodeScriptExecutionCancelledError.String())
}

func TestScriptStorageMutationsDiscarded(t *testing.T) {

	timeout := 10 * time.Second
	chain := flow.Mainnet.Chain()
	ctx := fvm.NewContext(fvm.WithChain(chain))
	manager, _ := New(
		zerolog.Nop(),
		metrics.NewExecutionCollector(ctx.Tracer),
		trace.NewNoopTracer(),
		nil,
		nil,
		ctx,
		committer.NewNoopViewCommitter(),
		nil,
		nil,
		ComputationConfig{
			ProgramsCacheSize:        programs.DefaultProgramsCacheSize,
			ScriptLogThreshold:       DefaultScriptLogThreshold,
			ScriptExecutionTimeLimit: timeout,
		},
	)
	vm := manager.vm.(*fvm.VirtualMachine)
	view := testutil.RootBootstrappedLedger(vm, ctx)

	programs := programs.NewEmptyBlockPrograms()
	txnPrograms, err := programs.NewTransactionPrograms(0, 0)
	require.NoError(t, err)

	txnState := state.NewTransactionState(view, state.DefaultParameters())
	env := fvm.NewScriptEnv(context.Background(), ctx, txnState, txnPrograms)

	// Create an account private key.
	privateKeys, err := testutil.GenerateAccountPrivateKeys(1)
	require.NoError(t, err)

	// Bootstrap a ledger, creating accounts with the provided private keys and the root account.
	accounts, err := testutil.CreateAccounts(vm, view, programs, privateKeys, chain)
	require.NoError(t, err)
	account := accounts[0]
	address := cadence.NewAddress(account)
	commonAddress, _ := common.HexToAddress(address.Hex())

	script := []byte(`
	pub fun main(account: Address) {
		let acc = getAuthAccount(account)
		acc.save(3, to: /storage/x)
	}
	`)

	header := unittest.BlockHeaderFixture()
	scriptView := view.NewChild()
	_, err = manager.ExecuteScript(context.Background(), script, [][]byte{jsoncdc.MustEncode(address)}, header, scriptView)

	require.NoError(t, err)

	rt := env.BorrowCadenceRuntime()
	defer env.ReturnCadenceRuntime(rt)

	v, err := rt.ReadStored(
		commonAddress,
		cadence.NewPath("storage", "x"),
	)

	// the save should not update account storage by writing the delta from the child view back to the parent
	require.NoError(t, err)
	require.Equal(t, nil, v)
}
