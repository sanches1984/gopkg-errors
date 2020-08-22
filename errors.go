package errors

import (
	"fmt"
	"io"
)

// Error ...
type Error struct {
	message    string
	errorType  *ErrorType
	cause      error
	stackTrace []string
	tags       []*ErrorTag
	payload    DataBagList
	log        DataBagList
}

// Causer is an interface
type Causer interface {
	Cause() error
}

// WithTag ...
func (e *Error) WithTag(et *ErrorTag) *Error {
	if !e.IsNil() {
		e.tags = append(e.tags, et)
	}
	return e
}

// WithCause ...
func (e *Error) WithCause(cause error) *Error {
	if !e.IsNil() {
		e.cause = cause
	}
	return e
}

// WithMessage ...
func (e *Error) WithMessage(message string) *Error {
	if !e.IsNil() {
		e.message = message
	}
	return e
}

// WithPayloadKV ...
func (e *Error) WithPayloadKV(items ...interface{}) *Error {
	if !e.IsNil() {
		dataBags := dataBagsFromSlice(items...)
		e.payload = append(e.payload, dataBags...)
		e.log = append(e.log, dataBags...)
	}
	return e
}

// WithLogKV ...
func (e *Error) WithLogKV(items ...interface{}) *Error {
	if !e.IsNil() {
		e.log = append(e.log, dataBagsFromSlice(items...)...)
	}
	return e
}

// Cause ...
func (e *Error) Cause() error {
	if !e.IsNil() {
		return e.cause
	}
	return nil
}

// GetScratchCode ...
func (e *Error) GetScratchCode() ErrorCode {
	if !e.IsNil() {
		return ErrorCode(e.errorType.errorCode)
	}

	return CodeTeapot
}

// GetLogKV ...
func (e *Error) GetLogKV() []interface{} {
	if e.IsNil() {
		return nil
	}
	result := make([]interface{}, 0, len(e.log)*2)
	for _, bag := range e.log {
		result = append(result, bag.Key, bag.Value)
	}
	return result
}

// GetPayloadKV ...
func (e *Error) GetPayloadKV() []interface{} {
	if e.IsNil() {
		return nil
	}
	result := make([]interface{}, 0, len(e.payload)*2)
	for _, bag := range e.payload {
		result = append(result, bag.Key, bag.Value)
	}
	return result
}

// IsTyped ...
func (e *Error) IsTyped(et *ErrorType) bool {
	if !e.IsNil() {
		return e.errorType == et
	}
	return false
}

// IsTagged ...
func (e *Error) IsTagged(et *ErrorTag) bool {
	if !e.IsNil() && len(e.tags) > 0 {
		for _, tag := range e.tags {
			if tag == et {
				return true
			}
		}
	}
	return false
}

// IsNil ...
func (e *Error) IsNil() bool {
	return e == nil || e.errorType == NullError
}

// Error ...
func (e *Error) Error() string {
	if !e.IsNil() {
		return e.message
	}
	return "<nil>"
}

// Format ...
func (e *Error) Format(s fmt.State, verb rune) {
	if e.IsNil() {
		_, _ = io.WriteString(s, "<nil>")
		return
	}

	_, _ = io.WriteString(s, e.message+"; payload: "+e.log.String())

	if verb == 'v' && s.Flag('+') {
		if len(e.stackTrace) > 0 {
			_, _ = io.WriteString(s, "\n ----------------------------------")
		}
		for _, call := range e.stackTrace {
			_, _ = io.WriteString(s, "\n")
			_, _ = io.WriteString(s, call)
		}

		if e.cause != nil {
			_, _ = io.WriteString(s, fmt.Sprintf("\ncause:\n%+v", e.cause))
		}
	}
}
