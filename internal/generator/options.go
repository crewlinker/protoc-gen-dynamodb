package generator

import (
	"strconv"

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

// determine the dyanmodb attribute name given the field definition
func (tg *Target) attrName(f *protogen.Field) string {
	if fopts := FieldOptions(f); fopts != nil && fopts.Name != nil {
		return *fopts.Name // explicit name option
	}

	return strconv.FormatInt(int64(f.Desc.Number()), 10)
}

// determine if the field is marked as the partition/sk key
func (tg *Target) isKey(f *protogen.Field) (isPk bool, isSk bool) {
	if fopts := FieldOptions(f); fopts != nil {
		return (fopts.Pk != nil && *fopts.Pk), (fopts.Sk != nil && *fopts.Sk)
	}

	return false, false
}

// determine if the field is marked as the partition/sk key
func (tg *Target) isOmitted(f *protogen.Field) bool {
	if fopts := FieldOptions(f); fopts != nil && fopts.Omit != nil {
		return *fopts.Omit
	}

	return false
}
