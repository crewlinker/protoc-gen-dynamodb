package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// genMessageKeying generates partition/sort key methods
func (tg *Target) genMessageKeying(f *File, m *protogen.Message) (err error) {
	var pkf, skf *protogen.Field
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // omitted, don't try to turn it into a key
		}

		isPk, isSk := tg.isKey(field)
		if isPk && isSk {
			// check that field cannot be marked as both sk and pk
			return fmt.Errorf("field '%s' is both marked as PK and as SK", field.GoName)
		}

		if isPk {
			if pkf != nil { // only one field can be marked as PK
				return fmt.Errorf("field '%s' is already marked as PK", pkf.GoName)
			}

			pkf = field
		}

		if isSk {
			if skf != nil { // only one field can be marked as SK
				return fmt.Errorf("field '%s' is already marked as SK", skf.GoName)
			}

			skf = field
		}
	}

	// if the message has a sort key, but not a partition key that doesn't make sense. The other
	// way around is ok.
	if pkf == nil && skf != nil {
		return fmt.Errorf("message '%s' has a sort key, but not a partition key", m.GoIdent.GoName)
	}

	// if no key fields are configured, so we don't generate a MarshalDynamoKey at all
	if pkf == nil && skf == nil {
		return nil
	}

	var body []Code
	if pkf != nil {
		if !tg.isValidKeyField(pkf) {
			return fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a PK", pkf.GoName)
		}

		body = append(body,
			Id("v").Op("=").Append(Id("v"), Lit(tg.attrName(pkf))),
		)
	}

	if skf != nil {
		if !tg.isValidKeyField(skf) {
			return fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a SK", skf.GoName)
		}

		body = append(body,
			Id("v").Op("=").Append(Id("v"), Lit(tg.attrName(skf))),
		)
	}

	f.Comment("DynamoKeyNames returns the attribute names of the partition and sort keys respectively")
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName)).Id("DynamoKeyNames").
		Params().
		Params(Id("v").Index().String()).
		Block(append(body, Return())...)

	return nil
}

// isValidKeyField returns whether a protobuf field can be a valid key
func (tg *Target) isValidKeyField(f *protogen.Field) bool {
	if f.Message != nil {
		return false // only basic types can be keys
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind,
		protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.BytesKind,
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
		return true // what encodes to dynamo Number, String or Bytes
	default:
		return false
	}
}
