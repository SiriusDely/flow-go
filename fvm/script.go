package fvm

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/common"

	"github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/hash"
)

type ScriptProcedure struct {
	ID             flow.Identifier
	Script         []byte
	Arguments      [][]byte
	RequestContext context.Context
	Value          cadence.Value
	Logs           []string
	Events         []flow.Event
	GasUsed        uint64
	MemoryEstimate uint64
	Err            errors.Error
}

type ScriptProcessor interface {
	Process(
		Context,
		*ScriptProcedure,
		*state.TransactionState,
		*programs.TransactionPrograms,
	) error
}

func Script(code []byte) *ScriptProcedure {
	scriptHash := hash.DefaultHasher.ComputeHash(code)

	return &ScriptProcedure{
		Script:         code,
		ID:             flow.HashToID(scriptHash),
		RequestContext: context.Background(),
	}
}

func (proc *ScriptProcedure) WithArguments(args ...[]byte) *ScriptProcedure {
	return &ScriptProcedure{
		ID:             proc.ID,
		Script:         proc.Script,
		RequestContext: proc.RequestContext,
		Arguments:      args,
	}
}

func (proc *ScriptProcedure) WithRequestContext(
	reqContext context.Context,
) *ScriptProcedure {
	return &ScriptProcedure{
		ID:             proc.ID,
		Script:         proc.Script,
		RequestContext: reqContext,
		Arguments:      proc.Arguments,
	}
}

func NewScriptWithContextAndArgs(
	code []byte,
	reqContext context.Context,
	args ...[]byte,
) *ScriptProcedure {
	scriptHash := hash.DefaultHasher.ComputeHash(code)
	return &ScriptProcedure{
		ID:             flow.HashToID(scriptHash),
		Script:         code,
		RequestContext: reqContext,
		Arguments:      args,
	}
}

func (proc *ScriptProcedure) Run(
	ctx Context,
	txnState *state.TransactionState,
	programs *programs.TransactionPrograms,
) error {
	for _, p := range ctx.ScriptProcessors {
		err := p.Process(ctx, proc, txnState, programs)
		txError, failure := errors.SplitErrorTypes(err)
		if failure != nil {
			if errors.IsALedgerFailure(failure) {
				return fmt.Errorf("cannot execute the script, this error usually happens if the reference block for this script is not set to a recent block: %w", failure)
			}
			return failure
		}
		if txError != nil {
			proc.Err = txError
			return nil
		}
	}

	return nil
}

func (proc *ScriptProcedure) ComputationLimit(ctx Context) uint64 {
	computationLimit := ctx.ComputationLimit
	// if ctx.ComputationLimit is also zero, fallback to the default computation limit
	if computationLimit == 0 {
		computationLimit = DefaultComputationLimit
	}
	return computationLimit
}

func (proc *ScriptProcedure) MemoryLimit(ctx Context) uint64 {
	memoryLimit := ctx.MemoryLimit
	// if ctx.MemoryLimit is also zero, fallback to the default memory limit
	if memoryLimit == 0 {
		memoryLimit = DefaultMemoryLimit
	}
	return memoryLimit
}

func (proc *ScriptProcedure) ShouldDisableMemoryAndInteractionLimits(
	ctx Context,
) bool {
	return ctx.DisableMemoryAndInteractionLimits
}

func (ScriptProcedure) Type() ProcedureType {
	return ScriptProcedureType
}

func (proc *ScriptProcedure) InitialSnapshotTime() programs.LogicalTime {
	return programs.EndOfBlockExecutionTime
}

func (proc *ScriptProcedure) ExecutionTime() programs.LogicalTime {
	return programs.EndOfBlockExecutionTime
}

type ScriptInvoker struct{}

func NewScriptInvoker() ScriptInvoker {
	return ScriptInvoker{}
}

func (i ScriptInvoker) Process(
	ctx Context,
	proc *ScriptProcedure,
	txnState *state.TransactionState,
	programs *programs.TransactionPrograms,
) error {
	env := NewScriptEnv(proc.RequestContext, ctx, txnState, programs)

	rt := env.BorrowCadenceRuntime()
	defer env.ReturnCadenceRuntime(rt)

	value, err := rt.ExecuteScript(
		runtime.Script{
			Source:    proc.Script,
			Arguments: proc.Arguments,
		},
		common.ScriptLocation(proc.ID))

	if err != nil {
		return err
	}

	proc.Value = value
	proc.Logs = env.Logs()
	proc.Events = env.Events()
	proc.GasUsed = env.ComputationUsed()
	proc.MemoryEstimate = env.MemoryEstimate()
	return nil
}
