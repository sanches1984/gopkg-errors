package middleware

import (
	"context"
	"github.com/sanches1984/gopkg-errors"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// NewConvertErrorsServerInterceptor converts errors
func NewConvertErrorsServerInterceptor(converters []errors.ErrorConverter, counter *prometheus.Counter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			appErr := (*errWithDetails)(errors.Convert(ctx, err, converters...))
			if counter != nil && appErr.Unwrap().IsTyped(errors.Internal) {
				(*counter).Inc()
			}
			return nil, appErr
		}

		return resp, nil
	}
}
