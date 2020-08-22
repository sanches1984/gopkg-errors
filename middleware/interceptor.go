package middleware

import (
	"context"
	"fmt"

	"github.com/severgroup-tt/gopkg-errors"
	"github.com/severgroup-tt/gopkg-logger"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// NewConvertErrorsServerInterceptor converts errors
func NewConvertErrorsServerInterceptor(converters []errors.ErrorConverter, counter *prometheus.Counter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			appErr := (*errWithDetails)(errors.Convert(ctx, err, converters...))
			logger.Error(ctx, fmt.Sprintf("error: %+v\n", appErr.Unwrap()))
			if counter != nil && appErr.Unwrap().IsTyped(errors.Internal) {
				(*counter).Inc()
			}
			return nil, appErr
		}

		return resp, nil
	}
}
