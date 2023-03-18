package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// generate marshalling code for a map field
func (tg *Target) genMapFieldMarshal(f *protogen.Field) (c []Code) {
	key := f.Message.Fields[0]
	val := f.Message.Fields[1]
	mid := fmt.Sprintf("m%d", f.Desc.Number())

	// if the map value is not a message. We don't need to faciliate recursing so
	// we can just unmarshal it as a basic value.
	if val.Message == nil {
		return tg.genBasicFieldMarshal(f)
	}

	// marshal the map key into a string (only map key type supported id Dynamo). For numbers it will be
	// base-10 encoded, bool values will be "true" or "false" and strings are just used as keys directly.
	var loop []Code
	switch key.Desc.Kind() {
	case protoreflect.StringKind:
		loop = append(loop, Id("mk").Op(":=").Id("k"))
	case protoreflect.BoolKind:
		loop = append(loop, Id("mk").Op(":=").Qual("fmt", "Sprintf").Call(Lit("%t"), Id("k")))
	case protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind:
		loop = append(loop, Id("mk").Op(":=").Qual("fmt", "Sprintf").Call(Lit("%d"), Id("k")))
	default:
		panic("unsupported map key: " + key.Desc.Kind().String())
	}

	// error when the resulting map key is an empty string, this also errors the marshalling logic
	// in the official Go SDK
	loop = append(loop, If(Id("mk").Op("==").Lit("")).Block(
		Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal map key of field '"+f.GoName+"': map key cannot be empty"))),
	))

	// marshal the map value. It makes a special case to check for null attribute value since
	// this needs to just be decoded as nil.
	loop = append(loop,
		// if the map value is nill, we encode as a null attribute value
		If(Id("v").Op("==").Nil()).Block(
			Id(mid).Dot("Value").Index(Id("mk")).Op("=").Op("&").Qual(dynamodbtypes, "AttributeValueMemberNULL").Values(
				Dict{Id("Value"): Lit(true)}),
			Continue(), // next map item
		),

		// marshal non-nil map value by calling the centrally generated function
		List(Id("mv"), Err()).Op(":=").Add(tg.idents.marshal).Call(Id("v"), Id("o").Op("...")),
		If(Err().Op("!=").Nil()).Block(
			Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal map value of field '"+f.GoName+"': %w"), Err())),
		),
		Id(mid).Dot("Value").Index(Id("mk")).Op("=").Id("mv"),
	)

	return []Code{
		// only marshal the map at all if it's not nil at runtime
		If(tg.marshalPresenceCond(f)...).Block(
			Id(mid).Op(":=").Op("&").Qual(dynamodbtypes, "AttributeValueMemberM").Values(
				Dict{Id("Value"): Make(Map(String()).Qual(dynamodbtypes, "AttributeValue"))}),

			For(List(Id("k"), Id("v")).Op(":=").Range().Id("x").Dot(f.GoName)).Block(loop...),
			Id("m").Index(Lit(tg.attrName(f))).Op("=").Id(mid),
		),
	}
}

// generate nested message marshalling
func (tg *Target) genMessageFieldMarshal(f *protogen.Field) []Code {
	return []Code{
		// only marshal message field if the value is not nil at runtime
		If(tg.marshalPresenceCond(f)...).Block(
			List(Id(fmt.Sprintf("m%d", f.Desc.Number())), Id("err")).Op(":=").Add(tg.idents.marshal).Call(
				Id("x").Dot("Get"+f.GoName).Call(), Id("o").Op("...")),
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

// genListFieldMarshal generates marshal code for a repeated field
func (tg *Target) genListFieldMarshal(f *protogen.Field) []Code {
	mid := fmt.Sprintf("m%d", f.Desc.Number())

	// if its not a list of messages, they are basic types and we can just
	// marshal it as normal
	if f.Message == nil {
		return tg.genBasicFieldMarshal(f)
	}

	// for messages we loop over each item and marshal them one by one
	return []Code{
		// only marshal if its not the zero  value
		If(tg.marshalPresenceCond(f)...).Block(
			Id(mid).Op(":=").Op("&").Qual(dynamodbtypes, "AttributeValueMemberL").Values(),
			For(List(Id("k"), Id("v")).Op(":=").Range().Id("x").Dot(f.GoName)).Block(

				// the value can also be nil, append null attribute and continue
				If(Id("v").Op("==").Nil()).Block(

					Id(mid).Dot("Value").Op("=").
						Append(Id(mid).Dot("Value"), Op("&").Qual(dynamodbtypes, "AttributeValueMemberNULL").Values(
							Dict{Id("Value"): Lit(true)})),

					Continue(), // next map item
				),

				// else, marshal the item
				List(Id("mv"), Err()).Op(":=").Add(tg.idents.marshal).Call(Id("v"), Id("o").Op("...")),
				If(Err().Op("!=").Nil()).Block(
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal item '%d' of field '"+f.GoName+"': %w"), Id("k"), Err())),
				),
				Id(mid).Dot("Value").Op("=").
					Append(Id(mid).Dot("Value"), Id("mv")),
			),
			Id("m").Index(Lit(tg.attrName(f))).Op("=").Id(mid),
		),
	}
}

// generateMessageMarshal generates the logic for marshalling messages into dynamo items
func (tg *Target) genMessageMarshal(f *File, m *protogen.Message) error {

	// render method body
	body := []Code{Id("m").Op("=").Make(Map(String()).Qual(dynamodbtypes, "AttributeValue"))}

	// generate field marshalling code
	for _, field := range m.Fields {
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
	f.Comment(`MarshalDynamoItem marshals dat into a dynamodb attribute map`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("MarshalDynamoItem").
		Params(Id("o").Op("...").Add(tg.idents.encopt)).
		Params(
			Id("m").Map(String()).Qual(dynamodbtypes, "AttributeValue"),
			Id("err").Id("error"),
		).Block(body...)
	return nil
}
