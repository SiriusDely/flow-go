package environment_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/onflow/cadence/runtime/common"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/engine/execution/state/delta"
	"github.com/onflow/flow-go/fvm"
	"github.com/onflow/flow-go/fvm/environment"
	programsStorage "github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
)

func Test_Programs(t *testing.T) {

	addressA := flow.HexToAddress("0a")
	addressB := flow.HexToAddress("0b")
	addressC := flow.HexToAddress("0c")

	contractALocation := common.AddressLocation{
		Address: common.Address(addressA),
		Name:    "A",
	}

	contractBLocation := common.AddressLocation{
		Address: common.Address(addressB),
		Name:    "B",
	}

	contractCLocation := common.AddressLocation{
		Address: common.Address(addressC),
		Name:    "C",
	}

	contractA0Code := `
		pub contract A {
			pub fun hello(): String {
        		return "bad version"
    		}
		}
	`

	contractACode := `
		pub contract A {
			pub fun hello(): String {
        		return "hello from A"
    		}
		}
	`

	contractBCode := `
		import A from 0xa
	
		pub contract B {
			pub fun hello(): String {
       		return "hello from B but also ".concat(A.hello())
    		}
		}
	`

	contractCCode := `
		import B from 0xb
	
		pub contract C {
			pub fun hello(): String {
	   		return "hello from C, ".concat(B.hello())
			}
		}
	`

	callTx := func(name string, address flow.Address) *flow.TransactionBody {

		return flow.NewTransactionBody().SetScript([]byte(fmt.Sprintf(`
			import %s from %s
			transaction {
              prepare() {
                log(%s.hello())
              }
            }`, name, address.HexWithPrefix(), name)),
		)
	}

	contractDeployTx := func(name, code string, address flow.Address) *flow.TransactionBody {
		encoded := hex.EncodeToString([]byte(code))

		return flow.NewTransactionBody().SetScript([]byte(fmt.Sprintf(`transaction {
              prepare(signer: AuthAccount) {
                signer.contracts.add(name: "%s", code: "%s".decodeHex())
              }
            }`, name, encoded)),
		).AddAuthorizer(address)
	}

	updateContractTx := func(name, code string, address flow.Address) *flow.TransactionBody {
		encoded := hex.EncodeToString([]byte(code))

		return flow.NewTransactionBody().SetScript([]byte(fmt.Sprintf(`transaction {
             prepare(signer: AuthAccount) {
               signer.contracts.update__experimental(name: "%s", code: "%s".decodeHex())
             }
           }`, name, encoded)),
		).AddAuthorizer(address)
	}

	mainView := delta.NewView(func(_, _ string) (flow.RegisterValue, error) {
		return nil, nil
	})

	txnState := state.NewTransactionState(mainView, state.DefaultParameters())

	vm := fvm.NewVM()
	programs := programsStorage.NewEmptyBlockPrograms()

	accounts := environment.NewAccounts(txnState)

	err := accounts.Create(nil, addressA)
	require.NoError(t, err)

	err = accounts.Create(nil, addressB)
	require.NoError(t, err)

	err = accounts.Create(nil, addressC)
	require.NoError(t, err)

	//err = stm.
	require.NoError(t, err)

	fmt.Printf("Account created\n")

	context := fvm.NewContext(
		fvm.WithContractDeploymentRestricted(false),
		fvm.WithTransactionProcessors(fvm.NewTransactionInvoker()),
		fvm.WithCadenceLogging(true),
		fvm.WithBlockPrograms(programs))

	var contractAView *delta.View = nil
	var contractBView *delta.View = nil
	var txAView *delta.View = nil

	t.Run("contracts can be updated", func(t *testing.T) {
		retrievedContractA, err := accounts.GetContract("A", addressA)
		require.NoError(t, err)
		require.Empty(t, retrievedContractA)

		// deploy contract A0
		procContractA0 := fvm.Transaction(
			contractDeployTx("A", contractA0Code, addressA),
			programs.NextTxIndexForTestingOnly())
		err = vm.RunV2(context, procContractA0, mainView)
		require.NoError(t, err)

		retrievedContractA, err = accounts.GetContract("A", addressA)
		require.NoError(t, err)

		require.Equal(t, contractA0Code, string(retrievedContractA))

		// deploy contract A
		procContractA := fvm.Transaction(
			updateContractTx("A", contractACode, addressA),
			programs.NextTxIndexForTestingOnly())
		err = vm.RunV2(context, procContractA, mainView)
		require.NoError(t, err)
		require.NoError(t, procContractA.Err)

		retrievedContractA, err = accounts.GetContract("A", addressA)
		require.NoError(t, err)

		require.Equal(t, contractACode, string(retrievedContractA))

	})

	t.Run("register touches are captured for simple contract A", func(t *testing.T) {

		// deploy contract A
		procContractA := fvm.Transaction(
			contractDeployTx("A", contractACode, addressA),
			programs.NextTxIndexForTestingOnly())
		err := vm.RunV2(context, procContractA, mainView)
		require.NoError(t, err)

		fmt.Println("---------- Real transaction here ------------")

		// run a TX using contract A
		procCallA := fvm.Transaction(
			callTx("A", addressA),
			programs.NextTxIndexForTestingOnly())

		loadedCode := false
		viewExecA := delta.NewView(func(owner, key string) (flow.RegisterValue, error) {
			if key == environment.ContractKey("A") {
				loadedCode = true
			}

			return mainView.Peek(owner, key)
		})

		err = vm.RunV2(context, procCallA, viewExecA)
		require.NoError(t, err)

		// make sure tx was really run
		require.Contains(t, procCallA.Logs, "\"hello from A\"")

		// Make sure the code has been loaded from storage
		require.True(t, loadedCode)

		_, programState, has := programs.GetForTestingOnly(contractALocation)
		require.True(t, has)

		// type assertion for further inspections
		require.IsType(t, programState.View(), &delta.View{})

		// assert some reads were recorded (at least loading of code)
		deltaView := programState.View().(*delta.View)
		require.NotEmpty(t, deltaView.Interactions().Reads)

		contractAView = deltaView
		txAView = viewExecA

		// merge it back
		err = mainView.MergeView(viewExecA)
		require.NoError(t, err)

		// execute transaction again, this time make sure it doesn't load code
		viewExecA2 := delta.NewView(func(owner, key string) (flow.RegisterValue, error) {
			//this time we fail if a read of code occurs
			require.NotEqual(t, key, environment.ContractKey("A"))

			return mainView.Peek(owner, key)
		})

		procCallA = fvm.Transaction(
			callTx("A", addressA),
			programs.NextTxIndexForTestingOnly())

		err = vm.RunV2(context, procCallA, viewExecA2)
		require.NoError(t, err)

		require.Contains(t, procCallA.Logs, "\"hello from A\"")

		// same transaction should produce the exact same views
		// but only because we don't do any conditional update in a tx
		compareViews(t, viewExecA, viewExecA2)

		// merge it back
		err = mainView.MergeView(viewExecA2)
		require.NoError(t, err)
	})

	t.Run("deploying another contract cleans programs storage", func(t *testing.T) {

		// deploy contract B
		procContractB := fvm.Transaction(
			contractDeployTx("B", contractBCode, addressB),
			programs.NextTxIndexForTestingOnly())
		err := vm.RunV2(context, procContractB, mainView)
		require.NoError(t, err)

		_, _, hasA := programs.GetForTestingOnly(contractALocation)
		_, _, hasB := programs.GetForTestingOnly(contractBLocation)

		require.False(t, hasA)
		require.False(t, hasB)
	})

	var viewExecB *delta.View

	t.Run("contract B imports contract A", func(t *testing.T) {

		// programs should have no entries for A and B, as per previous test

		// run a TX using contract B
		procCallB := fvm.Transaction(
			callTx("B", addressB),
			programs.NextTxIndexForTestingOnly())

		viewExecB = delta.NewView(mainView.Peek)

		err = vm.RunV2(context, procCallB, viewExecB)
		require.NoError(t, err)

		require.Contains(t, procCallB.Logs, "\"hello from B but also hello from A\"")

		_, programAState, has := programs.GetForTestingOnly(contractALocation)
		require.True(t, has)

		// state should be essentially the same as one which we got in tx with contract A
		require.IsType(t, programAState.View(), &delta.View{})
		deltaA := programAState.View().(*delta.View)

		compareViews(t, contractAView, deltaA)

		_, programBState, has := programs.GetForTestingOnly(contractBLocation)
		require.True(t, has)

		// program B should contain all the registers used by program A, as it depends on it
		require.IsType(t, programBState.View(), &delta.View{})
		deltaB := programBState.View().(*delta.View)

		idsA, valuesA := deltaA.Delta().RegisterUpdates()
		for i, id := range idsA {
			v, has := deltaB.Delta().Get(id.Owner, id.Key)
			require.True(t, has)

			require.Equal(t, valuesA[i], v)
		}

		for id, registerA := range deltaA.Interactions().Reads {

			registerB, has := deltaB.Interactions().Reads[id]
			require.True(t, has)

			require.Equal(t, registerA, registerB)
		}

		contractBView = deltaB

		// merge it back
		err = mainView.MergeView(viewExecB)
		require.NoError(t, err)

		// rerun transaction

		// execute transaction again, this time make sure it doesn't load code
		viewExecB2 := delta.NewView(func(owner, key string) (flow.RegisterValue, error) {
			//this time we fail if a read of code occurs
			require.NotEqual(t, key, environment.ContractKey("A"))
			require.NotEqual(t, key, environment.ContractKey("B"))

			return mainView.Peek(owner, key)
		})

		procCallB = fvm.Transaction(
			callTx("B", addressB),
			programs.NextTxIndexForTestingOnly())

		err = vm.RunV2(context, procCallB, viewExecB2)
		require.NoError(t, err)

		require.Contains(t, procCallB.Logs, "\"hello from B but also hello from A\"")

		compareViews(t, viewExecB, viewExecB2)

		// merge it back
		err = mainView.MergeView(viewExecB2)
		require.NoError(t, err)
	})

	t.Run("contract A runs from cache after program B has been loaded", func(t *testing.T) {

		// at this point programs cache should contain data for contract A
		// only because contract B has been called

		viewExecA := delta.NewView(func(owner, key string) (flow.RegisterValue, error) {
			require.NotEqual(t, key, environment.ContractKey("A"))
			return mainView.Peek(owner, key)
		})

		// run a TX using contract A
		procCallA := fvm.Transaction(
			callTx("A", addressA),
			programs.NextTxIndexForTestingOnly())

		err = vm.RunV2(context, procCallA, viewExecA)
		require.NoError(t, err)

		require.Contains(t, procCallA.Logs, "\"hello from A\"")

		compareViews(t, txAView, viewExecA)

		// merge it back
		err = mainView.MergeView(viewExecA)
		require.NoError(t, err)
	})

	t.Run("deploying contract C cleans programs", func(t *testing.T) {
		require.NotNil(t, contractBView)

		// deploy contract C
		procContractC := fvm.Transaction(
			contractDeployTx("C", contractCCode, addressC),
			programs.NextTxIndexForTestingOnly())
		err := vm.RunV2(context, procContractC, mainView)
		require.NoError(t, err)

		_, _, hasA := programs.GetForTestingOnly(contractALocation)
		_, _, hasB := programs.GetForTestingOnly(contractBLocation)
		_, _, hasC := programs.GetForTestingOnly(contractCLocation)

		require.False(t, hasA)
		require.False(t, hasB)
		require.False(t, hasC)

	})

	t.Run("importing C should chain-import B and A", func(t *testing.T) {
		procCallC := fvm.Transaction(
			callTx("C", addressC),
			programs.NextTxIndexForTestingOnly())

		viewExecC := delta.NewView(mainView.Peek)

		err = vm.RunV2(context, procCallC, viewExecC)
		require.NoError(t, err)

		require.Contains(t, procCallC.Logs, "\"hello from C, hello from B but also hello from A\"")

		// program A is the same
		_, programAState, has := programs.GetForTestingOnly(contractALocation)
		require.True(t, has)

		require.IsType(t, programAState.View(), &delta.View{})
		deltaA := programAState.View().(*delta.View)
		compareViews(t, contractAView, deltaA)

		// program B is the same
		_, programBState, has := programs.GetForTestingOnly(contractBLocation)
		require.True(t, has)

		require.IsType(t, programBState.View(), &delta.View{})
		deltaB := programBState.View().(*delta.View)
		compareViews(t, contractBView, deltaB)
	})
}

// compareViews compares views using only data that matters (ie. two different hasher instances
// trips the library comparison, even if actual SPoCKs are the same)
func compareViews(t *testing.T, a, b *delta.View) {
	require.Equal(t, a.Delta(), b.Delta())
	require.Equal(t, a.Interactions(), b.Interactions())
	require.Equal(t, a.ReadsCount(), b.ReadsCount())
	require.Equal(t, a.SpockSecret(), b.SpockSecret())
}
