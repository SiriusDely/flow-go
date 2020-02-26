package flow

import (
	"github.com/dapperlabs/flow-go/crypto"
)

type ResultApprovalBody struct {
	// CorrectnessAttestation
	ApproverID           Identifier       // TODO: actually fill in
	BlockID              Identifier       // TODO: actually fill in
	ExecutionResultID    Identifier       // hash of approved execution result
	ChunkIndex           uint64           // index of chunk that is approved
	AttestationSignature crypto.Signature // signature over ExecutionResultID and ChunkIndex, this has been separated for BLS aggregation
	Spock                Spock            // proof of re-computation, one per each chunk
}

type ResultApproval struct {
	ResultApprovalBody ResultApprovalBody
	VerifierSignature  crypto.Signature // signature over all above fields
}

func (ra ResultApproval) Body() interface{} {
	return ra.ResultApprovalBody
}

func (ra ResultApproval) ID() Identifier {
	return MakeID(ra.ResultApprovalBody)
}

func (ra ResultApproval) Checksum() Identifier {
	return MakeID(ra)
}
