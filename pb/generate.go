package pb

//go:generate go build github.com/gogo/protobuf/protoc-gen-gofast
//go:generate protoc --plugin=protoc-gen-gofast=./protoc-gen-gofast -I. --gofast_out=:. messages.proto
//go:generate rm protoc-gen-gofast
