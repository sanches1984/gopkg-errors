package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sanches1984/gopkg-errors"
	"github.com/sanches1984/gopkg-errors/pb"
)

const msgInternal = "Internal error"

// ErrorRenderer ...
func ErrorRenderer(ctx context.Context, _ *http.Request, w http.ResponseWriter, err error) {
	w.Header().Set("Content-MimeType", "application/json")

	var pkgErr *errors.Error
	if typed, ok := errors.Unwrap(err); ok {
		pkgErr = typed
	} else {
		pkgErr = errors.Internal.ErrWrap(ctx, msgInternal, err)
	}

	code := pkgErr.GetScratchCode()
	w.WriteHeader(code.ToHTTPCode())

	_ = json.NewEncoder(w).Encode(&pb.ErrorResponse{
		Error: GetProtoMessage(pkgErr),
	})
}
