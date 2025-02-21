package fvm

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/onflow/flow-go/fvm/crypto"
	"github.com/onflow/flow-go/fvm/environment"
	"github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/fvm/programs"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/trace"
)

// TransactionVerifier verifies the content of the transaction by
// checking accounts (authorizers, payer, proposer) are not frozen
// checking there is no double signature
// all signatures are valid
// all accounts provides enoguh weights
//
// if KeyWeightThreshold is set to a negative number, signature verification is skipped
type TransactionVerifier struct {
	KeyWeightThreshold int
}

func NewTransactionVerifier(keyWeightThreshold int) *TransactionVerifier {
	return &TransactionVerifier{
		KeyWeightThreshold: keyWeightThreshold,
	}
}

func (v *TransactionVerifier) Process(
	ctx Context,
	proc *TransactionProcedure,
	txnState *state.TransactionState,
	_ *programs.TransactionPrograms,
) error {
	return v.verifyTransaction(proc, ctx, txnState)
}

func newInvalidEnvelopeSignatureError(txSig flow.TransactionSignature, err error) error {
	return errors.NewInvalidEnvelopeSignatureError(txSig.Address, txSig.KeyIndex, err)
}

func newInvalidPayloadSignatureError(txSig flow.TransactionSignature, err error) error {
	return errors.NewInvalidPayloadSignatureError(txSig.Address, txSig.KeyIndex, err)
}

func (v *TransactionVerifier) verifyTransaction(
	proc *TransactionProcedure,
	ctx Context,
	txnState *state.TransactionState,
) error {
	if ctx.Tracer != nil && proc.TraceSpan != nil {
		span := ctx.Tracer.StartSpanFromParent(proc.TraceSpan, trace.FVMVerifyTransaction)
		span.SetAttributes(
			attribute.String("transaction.ID", proc.ID.String()),
		)
		defer span.End()
	}

	tx := proc.Transaction
	accounts := environment.NewAccounts(txnState)
	if tx.Payer == flow.EmptyAddress {
		err := errors.NewInvalidAddressErrorf(tx.Payer, "payer address is invalid")
		return fmt.Errorf("transaction verification failed: %w", err)
	}

	var err error
	var payloadWeights map[flow.Address]int
	var proposalKeyVerifiedInPayload bool

	err = v.checkSignatureDuplications(tx)
	if err != nil {
		return fmt.Errorf("transaction verification failed: %w", err)
	}

	// TODO(Janez): move disabling limits out of the verifier. Verifier should not be metered anyway.
	// TODO(Janez): verification is part of inclusion fees, not execution fees.

	// check accounts uses the state, but if the limits are too low, this might fail.
	// we shouldn't fail here if the limits are too low as fee deduction won't happen
	txnState.RunWithAllLimitsDisabled(
		func() {
			err = v.checkAccountsAreNotFrozen(tx, accounts)
		},
	)
	if err != nil {
		return fmt.Errorf("transaction verification failed: %w", err)
	}

	if v.KeyWeightThreshold < 0 {
		return nil
	}

	payloadWeights, proposalKeyVerifiedInPayload, err = v.verifyAccountSignatures(
		txnState,
		accounts,
		tx.PayloadSignatures,
		tx.PayloadMessage(),
		tx.ProposalKey,
		newInvalidPayloadSignatureError,
	)
	if err != nil {
		return errors.NewInvalidProposalSignatureError(tx.ProposalKey.Address, tx.ProposalKey.KeyIndex, err)
	}

	var envelopeWeights map[flow.Address]int
	var proposalKeyVerifiedInEnvelope bool

	envelopeWeights, proposalKeyVerifiedInEnvelope, err = v.verifyAccountSignatures(
		txnState,
		accounts,
		tx.EnvelopeSignatures,
		tx.EnvelopeMessage(),
		tx.ProposalKey,
		newInvalidEnvelopeSignatureError,
	)
	if err != nil {
		return errors.NewInvalidProposalSignatureError(tx.ProposalKey.Address, tx.ProposalKey.KeyIndex, err)

	}

	proposalKeyVerified := proposalKeyVerifiedInPayload || proposalKeyVerifiedInEnvelope
	if !proposalKeyVerified {
		err := fmt.Errorf("either the payload or the envelope should provide proposal signatures")
		return errors.NewInvalidProposalSignatureError(tx.ProposalKey.Address, tx.ProposalKey.KeyIndex, err)
	}

	for _, addr := range tx.Authorizers {
		// Skip this authorizer if it is also the payer. In the case where an account is
		// both a PAYER as well as an AUTHORIZER or PROPOSER, that account is required
		// to sign only the envelope.
		if addr == tx.Payer {
			continue
		}
		// hasSufficientKeyWeight
		if !v.hasSufficientKeyWeight(payloadWeights, addr) {
			msg := fmt.Sprintf("authorizer account does not have sufficient signatures (%d < %d)", payloadWeights[addr], v.KeyWeightThreshold)
			return errors.NewAccountAuthorizationErrorf(addr, msg)
		}
	}

	if !v.hasSufficientKeyWeight(envelopeWeights, tx.Payer) {
		// TODO change this to payer error (needed for fees)
		msg := fmt.Sprintf("payer account does not have sufficient signatures (%d < %d)", envelopeWeights[tx.Payer], v.KeyWeightThreshold)
		return errors.NewAccountAuthorizationErrorf(tx.Payer, msg)
	}

	return nil
}

// getPublicKey skips checking limits when getting the public key
func (v *TransactionVerifier) getPublicKey(
	txnState *state.TransactionState,
	accounts environment.Accounts,
	address flow.Address,
	keyIndex uint64,
) (pub flow.AccountPublicKey, err error) {
	txnState.RunWithAllLimitsDisabled(
		func() {
			pub, err = accounts.GetPublicKey(address, keyIndex)
		},
	)
	return
}

func (v *TransactionVerifier) verifyAccountSignatures(
	txnState *state.TransactionState,
	accounts environment.Accounts,
	signatures []flow.TransactionSignature,
	message []byte,
	proposalKey flow.ProposalKey,
	errorBuilder func(flow.TransactionSignature, error) error,
) (
	weights map[flow.Address]int,
	proposalKeyVerified bool,
	err error,
) {
	weights = make(map[flow.Address]int)

	for _, txSig := range signatures {

		accountKey, err := v.getPublicKey(txnState, accounts, txSig.Address, txSig.KeyIndex)
		if err != nil {
			return nil, false, errorBuilder(txSig, err)
		}
		err = v.verifyAccountSignature(accountKey, txSig, message, errorBuilder)
		if err != nil {
			return nil, false, err
		}
		if !proposalKeyVerified && v.sigIsForProposalKey(txSig, proposalKey) {
			proposalKeyVerified = true
		}

		weights[txSig.Address] += accountKey.Weight
	}

	return
}

// verifyAccountSignature verifies that an account signature is valid for the
// account and given message.
//
// If the signature is valid, this function returns the associated account key.
//
// An error is returned if the account does not contain a public key that
// correctly verifies the signature against the given message.
func (v *TransactionVerifier) verifyAccountSignature(
	accountKey flow.AccountPublicKey,
	txSig flow.TransactionSignature,
	message []byte,
	errorBuilder func(flow.TransactionSignature, error) error,
) error {

	if accountKey.Revoked {
		return errorBuilder(txSig, fmt.Errorf("account key has been revoked"))
	}

	valid, err := crypto.VerifySignatureFromTransaction(
		txSig.Signature,
		message,
		accountKey.PublicKey,
		accountKey.HashAlgo,
	)
	if err != nil {
		return errorBuilder(txSig, err)
	}

	if valid {
		return nil
	}

	return errorBuilder(txSig, fmt.Errorf("signature is not valid"))
}

func (v *TransactionVerifier) hasSufficientKeyWeight(
	weights map[flow.Address]int,
	address flow.Address,
) bool {
	return weights[address] >= v.KeyWeightThreshold
}

func (v *TransactionVerifier) sigIsForProposalKey(
	txSig flow.TransactionSignature,
	proposalKey flow.ProposalKey,
) bool {
	return txSig.Address == proposalKey.Address && txSig.KeyIndex == proposalKey.KeyIndex
}

func (v *TransactionVerifier) checkSignatureDuplications(tx *flow.TransactionBody) error {
	type uniqueKey struct {
		address flow.Address
		index   uint64
	}
	observedSigs := make(map[uniqueKey]bool)
	for _, sig := range tx.PayloadSignatures {
		if observedSigs[uniqueKey{sig.Address, sig.KeyIndex}] {
			err := fmt.Errorf("duplicate signatures are provided for the same key")
			return errors.NewInvalidPayloadSignatureError(sig.Address, sig.KeyIndex, err)
		}
		observedSigs[uniqueKey{sig.Address, sig.KeyIndex}] = true
	}

	for _, sig := range tx.EnvelopeSignatures {
		if observedSigs[uniqueKey{sig.Address, sig.KeyIndex}] {
			err := fmt.Errorf("duplicate signatures are provided for the same key")
			return errors.NewInvalidEnvelopeSignatureError(sig.Address, sig.KeyIndex, err)
		}
		observedSigs[uniqueKey{sig.Address, sig.KeyIndex}] = true
	}
	return nil
}

func (v *TransactionVerifier) checkAccountsAreNotFrozen(
	tx *flow.TransactionBody,
	accounts environment.Accounts,
) error {
	authorizers := make([]flow.Address, 0, len(tx.Authorizers)+2)
	authorizers = append(authorizers, tx.Authorizers...)
	authorizers = append(authorizers, tx.ProposalKey.Address, tx.Payer)

	for _, authorizer := range authorizers {
		err := accounts.CheckAccountNotFrozen(authorizer)
		if err != nil {
			return fmt.Errorf("checking frozen account failed: %w", err)
		}
	}

	return nil
}
