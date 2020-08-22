package errors

import (
	"context"
)

// ErrorConverter ...
type ErrorConverter func(context.Context, error) (*Error, bool)

// Convert ...
func Convert(ctx context.Context, err error, converters ...ErrorConverter) *Error {
	if err == nil {
		return NullError.Err(ctx, "")
	}
	// prevent double convert
	if pkgErr, ok := Unwrap(err); ok {
		return pkgErr
	}

	for _, converter := range converters {
		if result, ok := converter(ctx, err); ok {
			return result
		}
	}

	return Internal.ErrWrap(ctx, MessageInternalError, err)
}
