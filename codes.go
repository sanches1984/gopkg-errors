package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Message codes
const (
	MessageInternalError  = "INTERNAL_ERROR"
	MessageNotFound       = "NOT_FOUND"
	MessageConflict       = "CONFLICT"
	MessageUnauthorized   = "UNAUTHORIZED"
	MessageBadRequest     = "BAD_REQUEST"
	MessageAccessDenied   = "ACCESS_DENIED"
	MessageRequestTimeout = "REQUEST_TIMEOUT"
	MessageUnprocessable  = "UNPROCESSABLE"
)

var (
	httpCodesMap = map[int]string{
		http.StatusNotFound:            MessageNotFound,
		http.StatusConflict:            MessageConflict,
		http.StatusUnauthorized:        MessageUnauthorized,
		http.StatusBadRequest:          MessageBadRequest,
		http.StatusForbidden:           MessageAccessDenied,
		http.StatusRequestTimeout:      MessageRequestTimeout,
		http.StatusUnprocessableEntity: MessageUnprocessable,
	}

	grpcCodesMap = map[codes.Code]string{
		codes.NotFound:           MessageNotFound,
		codes.FailedPrecondition: MessageConflict,
		codes.Unauthenticated:    MessageUnauthorized,
		codes.InvalidArgument:    MessageBadRequest,
		codes.PermissionDenied:   MessageAccessDenied,
		codes.DeadlineExceeded:   MessageRequestTimeout,
	}

	appCodeMap = map[ErrorCode]string{
		CodeNotFound:            MessageNotFound,
		CodeConflict:            MessageConflict,
		CodeUnauthorized:        MessageUnauthorized,
		CodeBadRequest:          MessageBadRequest,
		CodeAccessDenied:        MessageAccessDenied,
		CodeRequestTimeout:      MessageRequestTimeout,
		CodeUnprocessableEntity: MessageUnprocessable,
	}
)

// PrintHTTPCode responds with string code
func PrintHTTPCode(httpCode int) string {
	if code, ok := httpCodesMap[httpCode]; ok {
		return code
	}
	return MessageInternalError
}

// PrintGRPCCode responds with string code
func PrintGRPCCode(grpcCode codes.Code) string {
	if code, ok := grpcCodesMap[grpcCode]; ok {
		return code
	}
	return MessageInternalError
}

// PrintAppCode responds with string code
func PrintAppCode(code ErrorCode) string {
	if code, ok := appCodeMap[code]; ok {
		return code
	}
	return MessageInternalError
}
