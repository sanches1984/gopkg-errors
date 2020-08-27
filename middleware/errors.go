package middleware

import (
	"github.com/sanches1984/gopkg-errors"
	"github.com/sanches1984/gopkg-errors/transport"

	"github.com/golang/protobuf/proto"
)

type errWithDetails errors.Error

// Error ...
func (e *errWithDetails) Error() string {
	return e.Unwrap().Error()
}

// Details is implementation of platform/errors.errDetails interface
func (e *errWithDetails) Details() interface{} {
	return []proto.Message{transport.GetProtoMessage(e.Unwrap())}
}

// Unwrap returns original *error.Error
func (e *errWithDetails) Unwrap() *errors.Error {
	return (*errors.Error)(e)
}

// Cause ...
func (e *errWithDetails) Cause() error {
	return e.Unwrap()
}
