package transport

import (
	"fmt"

	"github.com/sanches1984/gopkg-errors"
	"github.com/sanches1984/gopkg-errors/pb"
)

// GetProtoMessage ...
func GetProtoMessage(err *errors.Error) *pb.Error {
	kv := err.GetPayloadKV()
	data := make([]*pb.Error_Entry, 0, len(kv)/2)
	for i := 1; i < len(kv); i += 2 {
		data = append(data, &pb.Error_Entry{
			Key:   fmt.Sprint(kv[i-1]),
			Value: fmt.Sprint(kv[i]),
		})
	}

	return &pb.Error{
		Code:    errors.PrintAppCode(err.GetScratchCode()),
		Message: err.Error(),
		Data:    data,
	}
}
