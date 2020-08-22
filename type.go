package errors

import (
	"context"
	"net/http"
)

// ErrorCode describes different error codes
type ErrorCode uint32

// ToHTTPCode returns HTTP response code from the ErrorCode
func (code ErrorCode) ToHTTPCode() int {
	if code < 0 || int(code) >= len(biz2httpCode) {
		return http.StatusInternalServerError
	}
	return biz2httpCode[code]
}

var (
	biz2httpCode = [...]int{
		// old
		CodeInternal:        http.StatusInternalServerError,
		CodeNotFound:        http.StatusNotFound,
		CodeConflict:        http.StatusConflict,
		CodeUnauthorized:    http.StatusUnauthorized,
		CodeBadRequest:      http.StatusBadRequest,
		CodeAccessDenied:    http.StatusForbidden,
		CodeBadGateway:      http.StatusBadGateway,
		CodeTooManyRequests: http.StatusTooManyRequests,
		// new
		CodePaymentRequired:              http.StatusPaymentRequired,
		CodeMethodNotAllowed:             http.StatusMethodNotAllowed,
		CodeNotAcceptable:                http.StatusNotAcceptable,
		CodeProxyAuthRequired:            http.StatusProxyAuthRequired,
		CodeRequestTimeout:               http.StatusRequestTimeout,
		CodeAlreadyExists:                http.StatusConflict,
		CodeGone:                         http.StatusGone,
		CodeLengthRequired:               http.StatusLengthRequired,
		CodePreconditionFailed:           http.StatusPreconditionFailed,
		CodeRequestEntityTooLarge:        http.StatusRequestEntityTooLarge,
		CodeRequestURITooLong:            http.StatusRequestURITooLong,
		CodeUnsupportedMediaType:         http.StatusUnsupportedMediaType,
		CodeRequestedRangeNotSatisfiable: http.StatusRequestedRangeNotSatisfiable,
		CodeExpectationFailed:            http.StatusExpectationFailed,
		CodeTeapot:                       http.StatusTeapot,
		//CodeMisdirectedRequest:           http.StatusMisdirectedRequest, in go >1.11
		CodeUnprocessableEntity:         http.StatusUnprocessableEntity,
		CodeLocked:                      http.StatusLocked,
		CodeFailedDependency:            http.StatusFailedDependency,
		CodeUpgradeRequired:             http.StatusUpgradeRequired,
		CodePreconditionRequired:        http.StatusPreconditionRequired,
		CodeRequestHeaderFieldsTooLarge: http.StatusRequestHeaderFieldsTooLarge,
		CodeUnavailableForLegalReasons:  http.StatusUnavailableForLegalReasons,
	}
)

// ErrorCode consts
const (
	CodeInternal        ErrorCode = iota // 500
	CodeNotFound                         // 404
	CodeConflict                         // 409
	CodeUnauthorized                     // 401
	CodeBadRequest                       // 400
	CodeAccessDenied                     // 403
	CodeBadGateway                       // 502 // return if dependency fails
	CodeTooManyRequests                  // 429
	CodePaymentRequired
	CodeMethodNotAllowed
	CodeNotAcceptable
	CodeProxyAuthRequired
	CodeRequestTimeout
	CodeAlreadyExists
	CodeGone
	CodeLengthRequired
	CodePreconditionFailed
	CodeRequestEntityTooLarge
	CodeRequestURITooLong
	CodeUnsupportedMediaType
	CodeRequestedRangeNotSatisfiable
	CodeExpectationFailed
	CodeTeapot
	//CodeMisdirectedRequest // in go >1.11
	CodeUnprocessableEntity
	CodeLocked
	CodeFailedDependency
	CodeUpgradeRequired
	CodePreconditionRequired
	CodeRequestHeaderFieldsTooLarge
	CodeUnavailableForLegalReasons
)

// ErrorType ...
type ErrorType struct {
	errorCode ErrorCode
}

// Global errors types
var (
	Internal        = NewType(CodeInternal)
	NotFound        = NewType(CodeNotFound)
	Conflict        = NewType(CodeConflict)
	Unauthenticated = NewType(CodeUnauthorized)
	BadRequest      = NewType(CodeBadRequest)
	Unprocessable   = NewType(CodeUnprocessableEntity)
	AccessDenied    = NewType(CodeAccessDenied)
	Timeout         = NewType(CodeRequestTimeout)
	NullError       = NewType(CodeInternal)
)

// NewType ...
func NewType(errorCode ErrorCode) *ErrorType {
	return &ErrorType{errorCode: errorCode}
}

// Err ...
func (et *ErrorType) Err(ctx context.Context, msg string) *Error {
	return et.new(msg, false)
}

// ErrWithTag ...
func (et *ErrorType) ErrWithTag(ctx context.Context, msg string, tag *ErrorTag) *Error {
	return et.new(msg, false).WithTag(tag)
}

// ErrWithStack ...
func (et *ErrorType) ErrWithStack(ctx context.Context, msg string) *Error {
	return et.new(msg, true)
}

// ErrWrap ...
func (et *ErrorType) ErrWrap(ctx context.Context, msg string, err error) *Error {
	// prevent double wrap
	if typed, ok := err.(*Error); ok {
		return typed
	}
	wrappedErr := et.new(msg, false)
	wrappedErr.cause = err
	return wrappedErr
}

// IsTyped ...
func (et *ErrorType) IsTyped(err error) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*Error); ok {
		return appErr.IsTyped(et)
	}
	return false
}

func (et *ErrorType) new(msg string, stacktrace bool) *Error {
	err := &Error{
		errorType: et,
		message:   msg,
	}

	if stacktrace {
		err.stackTrace = getStackTrace()
	}

	return err
}
