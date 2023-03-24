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

	if pkf != nil {
		if !tg.isValidKeyField(pkf) {
			return fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a PK", pkf.GoName)
		}

		if err := tg.genPartitionKeyMethod(f, m, pkf); err != nil {
			return fmt.Errorf("failed to generate partition key method: %w", err)
		}
	}

	if skf != nil {
		if !tg.isValidKeyField(skf) {
			return fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a SK", skf.GoName)
		}

		if err := tg.genSortKeyMethod(f, m, skf); err != nil {
			return fmt.Errorf("failed to generate sort key method: %w", err)
		}
	}

	// if no key fields are configured, so we don't generate a MarshalDynamoKey at all
	if pkf == nil && skf == nil {
		return nil
	}

	body := []Code{Id("m").Op("=").Make(Map(String()).Qual(types, "AttributeValue"))}
	if pkf != nil {
		body = append(body,
			List(Id("pk"), Id("pkv")).Op(":=").Id("x").Dot("PartitionKey").Call(),
			List(Id("m").Index(Id("pk")), Err()).Op("=").Qual(tg.idents.ddb, "Marshal").Call(
				Id("pkv"),
				tg.genEmbedOption(pkf),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal partition key field: %w"), Err())),
			),
		)
	}

	if skf != nil {
		body = append(body,
			List(Id("sk"), Id("skv")).Op(":=").Id("x").Dot("SortKey").Call(),
			List(Id("m").Index(Id("sk")), Err()).Op("=").Qual(tg.idents.ddb, "Marshal").Call(
				Id("skv"),
				tg.genEmbedOption(skf),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal sort key field: %w"), Err())),
			),
		)
	}

	f.Comment(`MarshalDynamoKey encodes only the item's key attributes into a DynamoDB attribute map`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("MarshalDynamoKey").
		Params().
		Params(
			Id("m").Map(String()).Qual(types, "AttributeValue"),
			Id("err").Id("error"),
		).
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

// genPartitionKeyMethod generates a method for the message to return partition key information
func (tg *Target) genPartitionKeyMethod(f *File, m *protogen.Message, kf *protogen.Field) (err error) {
	got, attrn := tg.fieldGoType(kf), tg.attrName(kf)
	f.Comment(`PartitionKey returns the name of the Dynamo attribute that holds th partition key and the current value of that key in the struct`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("PartitionKey").
		Params().
		Params(Id("name").String(), Id("value").Add(got)).
		Block(
			Return(Lit(attrn), Id("x").Dot(kf.GoName)),
		)

	return nil
}

// genSortKeyMethod generates a method to return sort key information for a message
func (tg *Target) genSortKeyMethod(f *File, m *protogen.Message, kf *protogen.Field) (err error) {
	got, attrn := tg.fieldGoType(kf), tg.attrName(kf)
	f.Comment(`Sortkey returns the name of the Dynamo attribute that holds the sort key and the current value of that key in the struct`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("SortKey").
		Params().
		Params(Id("name").String(), Id("value").Add(got)).
		Block(
			Return(Lit(attrn), Id("x").Dot(kf.GoName)),
		)

	return nil
}
