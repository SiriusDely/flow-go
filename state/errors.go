package state

import (
	"errors"
	"fmt"

	"github.com/onflow/flow-go/model/flow"
)

// InvalidExtensionError is an error for invalid extension of the state
type InvalidExtensionError struct {
	err error
}

func NewInvalidExtensionError(msg string) error {
	return NewInvalidExtensionErrorf(msg)
}

func NewInvalidExtensionErrorf(msg string, args ...interface{}) error {
	return InvalidExtensionError{
		err: fmt.Errorf(msg, args...),
	}
}

func (e InvalidExtensionError) Unwrap() error {
	return e.err
}

func (e InvalidExtensionError) Error() string {
	return e.err.Error()
}

// IsInvalidExtensionError returns whether the given error is an InvalidExtensionError error
func IsInvalidExtensionError(err error) bool {
	return errors.As(err, &InvalidExtensionError{})
}

// OutdatedExtensionError is an error for the extension of the state being outdated.
// Being outdated doesn't mean it's invalid or not.
// Knowing whether an outdated extension is an invalid extension or not would
// take more state queries.
type OutdatedExtensionError struct {
	err error
}

func NewOutdatedExtensionError(msg string) error {
	return NewOutdatedExtensionErrorf(msg)
}

func NewOutdatedExtensionErrorf(msg string, args ...interface{}) error {
	return OutdatedExtensionError{
		err: fmt.Errorf(msg, args...),
	}
}

func (e OutdatedExtensionError) Unwrap() error {
	return e.err
}

func (e OutdatedExtensionError) Error() string {
	return e.err.Error()
}

func IsOutdatedExtensionError(err error) bool {
	return errors.As(err, &OutdatedExtensionError{})
}

// NoValidChildBlockError is a sentinel error when the case where a certain block has
// no valid child.
type NoValidChildBlockError struct {
	err error
}

func NewNoValidChildBlockError(msg string) error {
	return NoValidChildBlockError{
		err: fmt.Errorf(msg),
	}
}

func NewNoValidChildBlockErrorf(msg string, args ...interface{}) error {
	return NewNoValidChildBlockError(fmt.Sprintf(msg, args...))
}

func (e NoValidChildBlockError) Unwrap() error {
	return e.err
}

func (e NoValidChildBlockError) Error() string {
	return e.err.Error()
}

func IsNoValidChildBlockError(err error) bool {
	return errors.As(err, &NoValidChildBlockError{})
}

// UnknownBlockError is a sentinel error indicating that a certain block
// has not been ingested yet.
type UnknownBlockError struct {
	blockID flow.Identifier
	err     error
}

// WrapAsUnknownBlockError wraps a given error as UnknownBlockError
func WrapAsUnknownBlockError(blockID flow.Identifier, err error) error {
	return UnknownBlockError{
		blockID: blockID,
		err:     fmt.Errorf("block %v has not been processed yet: %w", blockID, err),
	}
}

func NewUnknownBlockError(blockID flow.Identifier) error {
	return UnknownBlockError{
		blockID: blockID,
		err:     fmt.Errorf("block %v has not been processed yet", blockID),
	}
}

func (e UnknownBlockError) Unwrap() error { return e.err }
func (e UnknownBlockError) Error() string { return e.err.Error() }

func IsUnknownBlockError(err error) bool {
	var e UnknownBlockError
	return errors.As(err, &e)
}
