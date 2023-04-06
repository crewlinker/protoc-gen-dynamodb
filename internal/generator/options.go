package generator

import (
	"fmt"
	"strconv"

	ddbexpression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// MessageOptions returns our plugin specific options for a field. If the field has no options
// it returns nil.
func MessageOptions(m *protogen.Message) *ddbv1.MessageOptions {
	opts, ok := m.Desc.Options().(*descriptorpb.MessageOptions)
	if !ok {
		return nil
	}
	ext, ok := proto.GetExtension(opts, ddbv1.E_Msg).(*ddbv1.MessageOptions)
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

// attrKindForPkSk returns the attribute kind for a pk/sk (of a gsi/lsi)
func (tg *Target) attrKindForPkSk(f *protogen.Field) (ddbexpression.DynamoDBAttributeType, error) {
	if !tg.isValidKeyField(f) {
		return "", fmt.Errorf("field '%s' of '%s' marked as pk/sk for base table or lsi/gsi but it's not a basic type", f.GoName, f.Message.GoIdent)
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return ddbexpression.String, nil
	case protoreflect.BytesKind:
		return ddbexpression.Binary, nil
	case protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.DoubleKind,
		protoreflect.FloatKind:
		return ddbexpression.Number, nil
	default:
		return "", fmt.Errorf("unsupported kind for pk/sk: %v", f.Desc.Kind())
	}
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

// gsiPlacement contains gen data for fields of a message that describe a
// global secondary index
type gsiPlacement struct {
	pkField   *protogen.Field
	pkType    ddbexpression.DynamoDBAttributeType
	skField   *protogen.Field
	skType    ddbexpression.DynamoDBAttributeType
	projected []*protogen.Field
}

// lsiPlacement contains gen info of fields that make up a sk
type lsiPlacement struct {
	skField   *protogen.Field
	skType    ddbexpression.DynamoDBAttributeType
	projected []*protogen.Field
}

// holds information from Protobuf options that together describes in what
// table a message is placed.
type tablePlacementInfo struct {
	tableNames []string
	message    *protogen.Message
	pkField    *protogen.Field
	pkType     ddbexpression.DynamoDBAttributeType
	skField    *protogen.Field
	skType     ddbexpression.DynamoDBAttributeType
	gsis       map[string]gsiPlacement
	lsis       map[string]lsiPlacement
}

// tablePlacementOptions returns all options from the protobuf definition that influence
// how the message is placed in a dynamodb table
func (tg *Target) tablePlacementOptions(msg *protogen.Message) (mp *tablePlacementInfo, err error) {
	mopts := MessageOptions(msg)
	if mopts == nil || len(mopts.Table) < 1 {
		return nil, nil // no tables to placement configured
	}

	mp = &tablePlacementInfo{
		message:    msg,
		tableNames: mopts.Table,
		gsis:       map[string]gsiPlacement{},
		lsis:       map[string]lsiPlacement{}}

	for _, fld := range msg.Fields {
		fopts := FieldOptions(fld)
		if fopts == nil || (fopts.Omit != nil && *fopts.Omit) {
			continue // no options, or omitted
		}

		if fopts.Pk != nil && *fopts.Pk {
			mp.pkField = fld
			if mp.pkType, err = tg.attrKindForPkSk(fld); err != nil {
				return nil, err
			}
		}

		if fopts.Sk != nil && *fopts.Sk {
			mp.skField = fld
			if mp.skType, err = tg.attrKindForPkSk(fld); err != nil {
				return nil, err
			}
		}

		// @TODO error if sk and pk are the same field, for both indexes
		// @TODO error if pk/sk is set twice for a certain gsi/lsi
		// @TODO error if field is part of the projection and of a sk/pk

		// @TODO error if the message has no pk defined for base table (even though it wants to be placed in a table)

		for _, gsio := range fopts.Gsi {
			if gsio == nil {
				continue // not part of gsi in any way
			}
			curr := mp.gsis[*gsio.Name]
			if gsio.Pk != nil && *gsio.Pk {
				curr.pkField = fld
				if curr.pkType, err = tg.attrKindForPkSk(fld); err != nil {
					return nil, err
				}
			}
			if gsio.Sk != nil && *gsio.Sk {
				curr.skField = fld
				if curr.skType, err = tg.attrKindForPkSk(fld); err != nil {
					return nil, err
				}
			}
			if (gsio.Sk == nil || !*gsio.Sk) || (gsio.Pk == nil || !*gsio.Pk) {
				curr.projected = append(curr.projected, fld)
			}

			mp.gsis[*gsio.Name] = curr
		}

		for _, lsio := range fopts.Lsi {
			if lsio == nil {
				continue // not part of lsi in any way
			}
			curr := mp.lsis[*lsio.Name]
			if lsio.Sk != nil && *lsio.Sk {
				curr.skField = fld
				if curr.skType, err = tg.attrKindForPkSk(fld); err != nil {
					return nil, err
				}
			}
			if lsio.Sk == nil || !*lsio.Sk {
				curr.projected = append(curr.projected, fld)
			}

			mp.lsis[*lsio.Name] = curr
		}
	}

	return
}
