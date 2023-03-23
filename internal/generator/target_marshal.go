package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// generate marshalling code for a map field
func (tg *Target) genMapFieldMarshal(f *protogen.Field) (c []Code) {
	val := f.Message.Fields[1]

	// if the map value is not a message. We don't need to faciliate recursing so
	// we can just unmarshal it as a basic value.
	if val.Message == nil {
		return tg.genBasicFieldMarshal(f)
	}

	// for messages we loop over each item and marshal them one by one
	return []Code{
		If(tg.marshalPresenceCond(f)...).Block(
			List(Id("m").Index(Lit(tg.attrName(f))), Err()).Op("=").Qual(tg.idents.ddb, "MarshalMappedMessage").Call(Id("x").Dot(f.GoName)),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal mapped message field '"+f.GoName+"': %w"), Err())),
			),
		),
	}

}

// generate nested message marshalling
func (tg *Target) genMessageFieldMarshal(f *protogen.Field) []Code {
	return []Code{
		// only marshal message field if the value is not nil at runtime
		If(tg.marshalPresenceCond(f)...).Block(
			List(Id(fmt.Sprintf("m%d", f.Desc.Number())), Id("err")).Op(":=").Qual(tg.idents.ddb, "MarshalMessage").Call(Id("x").Dot("Get"+f.GoName).Call()),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal field '"+f.GoName+"': %w"), Err())),
			),
			Id("m").Index(Lit(tg.attrName(f))).Op("=").Id(fmt.Sprintf("m%d", f.Desc.Number())),
		),
	}
}

// basic field generates code for marshaling a regular field with basic types
func (tg *Target) genBasicFieldMarshal(f *protogen.Field) []Code {
	return []Code{
		If(tg.marshalPresenceCond(f)...).Block(
			List(
				Id("m").Index(Lit(tg.attrName(f))),
				Id("err"),
			).Op("=").
				Qual(attributevalues, "Marshal").Call(Id("x").Dot("Get"+f.GoName).Call()),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// genSetFieldMarshal generates code to marshal a field into a StringSet, NumberSet or BinarySet
func (tg *Target) genSetFieldMarshal(f *protogen.Field) []Code {
	mid := fmt.Sprintf("m%d", f.Desc.Number())
	var attrValMemberType string
	var loop []Code
	switch f.Desc.Kind() {
	case protoreflect.BytesKind:
		attrValMemberType = "AttributeValueMemberBS"
		loop = append(loop,
			Id(mid).Dot("Value").Op("=").Append(Id(mid).Dot("Value"), Id("v")))

	case protoreflect.StringKind:
		attrValMemberType = "AttributeValueMemberSS"
		loop = append(loop,
			Id(mid).Dot("Value").Op("=").Append(Id(mid).Dot("Value"), Id("v")))
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
		attrValMemberType = "AttributeValueMemberNS"

		// the looping logic is more difficult for number sets because we wanna use
		// the official library to encode numbers. Which may error, and need to be type asserted.
		loop = append(loop,
			List(Id("av"), Err()).Op(":=").Qual(attributevalues, "Marshal").Call(Id("v")),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal set item of field '"+f.GoName+"': %w"), Err())),
			),
			List(Id("avn"), Id("ok")).Op(":=").Id("av").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberN")),
			If(Op("!").Id("ok")).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("set item of field '"+f.GoName+"' dit not marshal to a N value"))),
			),
			Id(mid).Dot("Value").Op("=").Append(Id(mid).Dot("Value"), Id("avn").Dot("Value")),
		)

	default:
		panic("unsupported set field: " + f.Desc.Kind().String())
	}

	return []Code{
		If(tg.marshalPresenceCond(f)...).Block(
			Id(mid).Op(":=").Op("&").Qual(dynamodbtypes, attrValMemberType).Values(),
			For(List(Id("_"), Id("v")).Op(":=").Range().Id("x").Dot(f.GoName)).Block(loop...),
			Id("m").Index(Lit(tg.attrName(f))).Op("=").Id(mid),
		),
	}
}

// genListFieldMarshal generates marshal code for a repeated field
func (tg *Target) genListFieldMarshal(f *protogen.Field) []Code {

	// if its not a list of messages, it could be marked as as a set
	if f.Message == nil && tg.isSet(f) {
		return tg.genSetFieldMarshal(f)
	} else if f.Message == nil {
		// else, it must be a repeated field of basic type, we can just marshal as usual
		return tg.genBasicFieldMarshal(f)
	}

	// for messages we loop over each item and marshal them one by one
	return []Code{
		// only marshal if its not the zero value
		If(tg.marshalPresenceCond(f)...).Block(
			List(Id("m").Index(Lit(tg.attrName(f))), Err()).Op("=").Qual(tg.idents.ddb, "MarshalRepeatedMessage").Call(Id("x").Dot(f.GoName)),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal repeated message field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// generateMessageMarshal generates the logic for marshalling messages into dynamo items
func (tg *Target) genMessageMarshal(f *File, m *protogen.Message) error {

	// render method body
	body := []Code{Id("m").Op("=").Make(Map(String()).Qual(dynamodbtypes, "AttributeValue"))}

	// generate field marshalling code
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // generate no marshallling code for omitted fields
		}

		switch {
		case field.Desc.IsList():
			// lists are repeated fields
			body = append(body, tg.genListFieldMarshal(field)...)
		case field.Desc.IsMap():
			// field map is technically a message but we marshal it differently
			body = append(body, tg.genMapFieldMarshal(field)...)
		case field.Message != nil:
			// nested message, not part of a one-of
			body = append(body, tg.genMessageFieldMarshal(field)...)
		default:
			// else, assume attributevalue package can handle it
			body = append(body, tg.genBasicFieldMarshal(field)...)
		}
	}

	body = append(body,
		Return(Id("m"), Nil()))

	// render function with body
	f.Comment(`MarshalDynamoItem marshals data into a dynamodb attribute map`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("MarshalDynamoItem").
		Params().
		Params(
			Id("m").Map(String()).Qual(dynamodbtypes, "AttributeValue"),
			Id("err").Id("error"),
		).Block(body...)
	return nil
}
