package grpc

import (
	"context"
	"fmt"
	"regexp"

	"github.com/severgroup-tt/gopkg-errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rpc error: code = InvalidArgument desc = already exists
var errorRe = regexp.MustCompile(`^rpc error: code = (\S+) desc = (.+)$`)

const (
	msgInternal        = "Internal error"
	msgNotFound        = "Entity not found"
	msgConflict        = "Entity already exists"
	msgBadRequest      = "Invalid request"
	msgUnprocessable   = "Unprocessable entity"
	msgTimeout         = "Request timeout"
	msgUnauthenticated = "Authentication required"
	msgAccessDenied    = "Access denied"
)

type grpcStatus interface {
	GRPCStatus() *status.Status
}

type named interface {
	GetName() string
}

type stringable interface {
	String() string
}

// Converter ...
func Converter(service string) errors.ErrorConverter {
	return func(ctx context.Context, err error) (*errors.Error, bool) {
		for {
			if errTyped, ok := err.(grpcStatus); ok {
				return convert(ctx, service, errTyped.GRPCStatus()), true
			}

			errC, ok := err.(errors.Causer)
			if !ok {
				return nil, false
			}

			err = errC.Cause()
		}
	}
}

func convert(ctx context.Context, service string, err *status.Status) *errors.Error {
	if err.Code() == codes.OK {
		return errors.NullError.Err(ctx, "")
	}

	result := fromCode(ctx, err.Code(), err.Message(), err.Err()).WithLogKV(
		"service", service,
		"grpc.code", err.Code().String(),
		"msg", err.Message(),
	)
	return result.WithLogKV(convertPayload(err.Details()...)...)
}

func fromCode(ctx context.Context, code codes.Code, msg string, err error) *errors.Error {
	switch code {
	case codes.InvalidArgument, codes.OutOfRange:
		return errors.BadRequest.ErrWrap(ctx, msgBadRequest, err)

	case codes.NotFound:
		return errors.NotFound.ErrWrap(ctx, msgNotFound, err)

	case codes.AlreadyExists, codes.FailedPrecondition:
		return errors.Conflict.ErrWrap(ctx, msgConflict, err)

	case codes.Unauthenticated:
		return errors.Unauthenticated.ErrWrap(ctx, msgUnauthenticated, err)

	case codes.PermissionDenied:
		return errors.AccessDenied.ErrWrap(ctx, msgAccessDenied, err)

	case codes.DataLoss, codes.DeadlineExceeded:
		return errors.Timeout.ErrWrap(ctx, msgTimeout, err)

	case codes.Unavailable:
		return errors.Unprocessable.ErrWrap(ctx, msgUnprocessable, err)

	case codes.Internal:
		matches := errorRe.FindStringSubmatch(msg)
		if len(matches) == 3 {
			pkgErr := fromTextCode(ctx, matches[1], matches[2], err)
			if !pkgErr.IsNil() {
				return pkgErr
			}
		}
	}

	return errors.Internal.ErrWrap(ctx, msgInternal, err)
}

func fromTextCode(ctx context.Context, code string, msg string, err error) *errors.Error {
	switch code {
	case "InvalidArgument", "OutOfRange":
		return errors.BadRequest.ErrWrap(ctx, msg, err)

	case "NotFound":
		return errors.NotFound.ErrWrap(ctx, msg, err)

	case "AlreadyExists", "FailedPrecondition":
		return errors.Conflict.ErrWrap(ctx, msg, err)

	case "Unauthenticated":
		return errors.Unauthenticated.ErrWrap(ctx, msg, err)

	case "PermissionDenied":
		return errors.AccessDenied.ErrWrap(ctx, msg, err)

	case "DataLoss", "DeadlineExceeded":
		return errors.Timeout.ErrWrap(ctx, msg, err)

	case "Internal":
		// we need to go deeper, real life error:
		// rpc error: code = Internal desc = rpc error: code = InvalidArgument desc = already exists
		matches := errorRe.FindStringSubmatch(msg)
		if len(matches) == 3 {
			return fromTextCode(ctx, matches[1], matches[2], err)
		}

		return errors.Internal.ErrWrap(ctx, msg, err)
	}

	return errors.NullError.Err(ctx, "")
}

func convertPayload(payload ...interface{}) []interface{} {
	result := make([]interface{}, 0, len(payload)*2)
	for i, entity := range payload {
		result = append(result, extractKey(i, entity), extractValue(entity))
	}

	return result
}

func extractKey(index int, entity interface{}) string {
	switch typed := entity.(type) {
	case named:
		return typed.GetName()
	default:
		return fmt.Sprintf("key_%v", index)
	}
}

func extractValue(entity interface{}) interface{} {
	switch typed := entity.(type) {
	case stringable:
		return typed.String()
	default:
		return entity
	}
}
