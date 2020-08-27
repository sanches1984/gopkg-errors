package http

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/sanches1984/gopkg-errors"

	"github.com/go-openapi/runtime"
)

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

// Converter ...
func Converter(service string) errors.ErrorConverter {
	return func(ctx context.Context, err error) (*errors.Error, bool) {
		for {
			if errTyped, ok := err.(*runtime.APIError); ok {
				return convert(ctx, service, errTyped), true
			}

			errC, ok := err.(errors.Causer)
			if !ok {
				return nil, false
			}

			err = errC.Cause()
		}
	}
}

func convert(ctx context.Context, service string, err *runtime.APIError) *errors.Error {
	result := fromCode(ctx, err.Code, err).WithLogKV(
		"service", service,
		"operation", err.OperationName,
		"http.code", err.Code,
	)

	if resp, ok := err.Response.(runtime.ClientResponse); ok {
		rawBody, _ := ioutil.ReadAll(resp.Body())
		result.WithLogKV(
			"status", resp.Message(),
			"body", string(rawBody),
		)
	}

	return result
}

func fromCode(ctx context.Context, code int, err error) *errors.Error {
	switch code {
	case http.StatusMovedPermanently,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
		http.StatusNotFound,
		http.StatusGone:
		return errors.NotFound.ErrWrap(ctx, msgNotFound, err)

	case http.StatusConflict:
		return errors.Conflict.ErrWrap(ctx, msgConflict, err)

	case http.StatusBadRequest, http.StatusMethodNotAllowed:
		return errors.BadRequest.ErrWrap(ctx, msgBadRequest, err)

	case http.StatusUnprocessableEntity:
		return errors.Unprocessable.ErrWrap(ctx, msgUnprocessable, err)

	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return errors.Timeout.ErrWrap(ctx, msgTimeout, err)

	case http.StatusUnauthorized:
		return errors.Unauthenticated.ErrWrap(ctx, msgUnauthenticated, err)

	case http.StatusForbidden:
		return errors.AccessDenied.ErrWrap(ctx, msgAccessDenied, err)
	}

	return errors.Internal.ErrWrap(ctx, msgInternal, err)
}
