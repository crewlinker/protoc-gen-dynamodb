package generator

import (
	"strconv"

	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// EntityOptions on oneoff for single-table design
func EntityOptions(of *protogen.Oneof) *ddbv1.EntityOptions {
	opts, ok := of.Desc.Options().(*descriptorpb.OneofOptions)
	if !ok {
		return nil
	}
	ext, ok := proto.GetExtension(opts, ddbv1.E_Entity).(*ddbv1.EntityOptions)
	if !ok {
		return nil
	}
	if ext == nil {
		return nil
	}
	return ext
}

// TableOptions returns our plugin specific options for a field. If the field has no options
// it returns nil.
func TableOptions(m *protogen.Message) *ddbv1.TableOptions {
	opts, ok := m.Desc.Options().(*descriptorpb.MessageOptions)
	if !ok {
		return nil
	}
	ext, ok := proto.GetExtension(opts, ddbv1.E_Table).(*ddbv1.TableOptions)
	if !ok {
		return nil
	}
	if ext == nil {
		return nil
	}
	return ext
}

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

// // attrKindForPkSk returns the attribute kind for a pk/sk (of a gsi/lsi)
// func (tg *Target) attrKindForPkSk(f *protogen.Field) (ddbexpression.DynamoDBAttributeType, error) {
// 	if !tg.isValidKeyField(f) {
// 		return "", fmt.Errorf("field '%s' of '%s' marked as pk/sk for base table or lsi/gsi but it's not a basic type", f.GoName, f.Message.GoIdent)
// 	}

// 	switch f.Desc.Kind() {
// 	case protoreflect.StringKind:
// 		return ddbexpression.String, nil
// 	case protoreflect.BytesKind:
// 		return ddbexpression.Binary, nil
// 	case protoreflect.Int64Kind,
// 		protoreflect.Uint64Kind,
// 		protoreflect.Fixed64Kind,
// 		protoreflect.Sint64Kind,
// 		protoreflect.Sfixed64Kind,
// 		protoreflect.Int32Kind,
// 		protoreflect.Uint32Kind,
// 		protoreflect.Fixed32Kind,
// 		protoreflect.Sint32Kind,
// 		protoreflect.Sfixed32Kind,
// 		protoreflect.DoubleKind,
// 		protoreflect.FloatKind:
// 		return ddbexpression.Number, nil
// 	default:
// 		return "", fmt.Errorf("unsupported kind for pk/sk: %v", f.Desc.Kind())
// 	}
// }

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

// determine if the field is marked as the a set of strings,numbers or bytes
func (tg *Target) isSet(f *protogen.Field) bool {
	if fopts := FieldOptions(f); fopts != nil && fopts.Set != nil {
		return *fopts.Set
	}

	return false
}

// returns the embedding encoding
func (tg *Target) embedEncoding(f *protogen.Field) ddbv1.Encoding {
	if fopts := FieldOptions(f); fopts != nil && fopts.Embed != nil {
		return *fopts.Embed
	}
	return ddbv1.Encoding_ENCODING_UNSPECIFIED
}

// notSupportPathing returns wether a field doesn't support deep pathing
func (tg *Target) notSupportPathing(field *protogen.Field) bool {
	return field.Message == nil || // if field is not a message, never support pathing
		(!tg.isSamePkgIdent(field.Message.GoIdent) && !tg.isWellKnownPathSupported(field.Message)) ||
		(tg.embedEncoding(field) != ddbv1.Encoding_ENCODING_DYNAMO &&
			tg.embedEncoding(field) != ddbv1.Encoding_ENCODING_UNSPECIFIED)
}
