package generator

import (
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FieldOptions returns our plugin specific options for a field. If the field has no options
// it returns nil.
func FieldOptions(f *protogen.Field) *ddbv1.FieldOptions {
	opts, ok := f.Desc.Options().(*descriptorpb.FieldOptions)
	if !ok {
		return nil
	}
	ext, ok := proto.GetExtension(opts, ddbv1.E_Field).(*ddbv1.FieldOptions)
	if !ok {
		return nil
	}
	if ext == nil {
		return nil
	}
	return ext
}
