package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// generate marshalling code for a map field.
func (tg *Target) genMapFieldUnmarshal(f *protogen.Field) (c []Code) {
	key := f.Message.Fields[0]
	val := f.Message.Fields[1]

	// if the map value is not a message. We don't need to faciliate recursing so
	// we can just unmarshal it as a basic value.
	if val.Message == nil {
		return tg.genBasicFieldUnmarshal(f)
	}

	// generate loop code for unmarshalling the map key
	var loop []Code
	var addUnmarshalKeyErr bool
	switch key.Desc.Kind() {
	case protoreflect.StringKind:
		loop = append(loop, Id("mk").Op(":=").
			Id("k"))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		loop = append(loop, List(Id("mk"), Err()).Op(":=").
			Qual("strconv", "ParseInt").Call(Id("k"), Lit(10), Lit(64)))
		addUnmarshalKeyErr = true
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		loop = append(loop, List(Id("mk"), Err()).Op(":=").
			Qual("strconv", "ParseUint").Call(Id("k"), Lit(10), Lit(64)))
		addUnmarshalKeyErr = true
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		loop = append(loop, List(Id("mk"), Err()).Op(":=").
			Qual("strconv", "ParseInt").Call(Id("k"), Lit(10), Lit(32)))
		addUnmarshalKeyErr = true
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		loop = append(loop, List(Id("mk"), Err()).Op(":=").
			Qual("strconv", "ParseUint").Call(Id("k"), Lit(10), Lit(32)))
		addUnmarshalKeyErr = true
	case protoreflect.BoolKind:
		loop = append(loop,
			Var().Id("mk").Id("bool"),
			Switch(Id("k")).Block(
				Case(Lit("true")).Block(
					Id("mk").Op("=").Lit(true),
				),
				Case(Lit("false")).Block(
					Id("mk").Op("=").Lit(false),
				),
				Default().Block(
					Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal map key for field '"+f.GoName+"': not 'true' or 'false' value"))),
				),
			))
	default:
		panic("unsupported key type for map: " + key.Desc.Kind().String())
	}

	if addUnmarshalKeyErr {
		loop = append(loop, If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal map key for field '"+f.GoName+"': %w"), Err())),
		))
	}

	// unmarshal the map values
	loop = append(loop,
		// if the map value is NULL, we just assign nil and don't attempt to unmarshal it
		If(List(Id("_"), Id("ok")).Op(":=").Id("v").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberNULL")), Id("ok")).Block(
			Id("x").Dot(f.GoName).Index(tg.fieldGoType(key).Call(Id("mk"))).Op("=").Nil(),
			Continue(),
		),

		// else, we unmarshal into a not-nil message
		Var().Id("mv").Add(tg.fieldGoType(val)),
		Err().Op("=").Add(tg.idents.unmarshal).Call(Id("v"), Op("&").Id("mv")),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal map value for field '"+f.GoName+"': %w"), Err())),
		),

		// assign map value while type casting to map key, the cast is only necessary because ParseInt always
		// returns 64-bit values while we sometime wanna assign 32-bit values. It should downcast at worst
		Id("x").Dot(f.GoName).Index(tg.fieldGoType(key).Call(Id("mk"))).Op("=").Op("&").Id("mv"),
	)

	mid := fmt.Sprintf("m%d", f.Desc.Number())
	valtypid := tg.fieldGoType(val)
	if val.Message != nil {
		valtypid = Op("*").Add(valtypid) // in case it's a message, pointer ref
	}

	return []Code{

		// only unmarshal map, if the attribute is not nil
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			Id("x").Dot(f.GoName).
				Op("=").Make(Map(tg.fieldGoType(key)).Add(valtypid)),
			List(Id(mid), Id("ok")).Op(":=").
				Id("m").Index(Lit(tg.attrName(f))).Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberM")),
			If(Op("!").Id("ok")).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': no map attribute provided"))),
			),
			For(List(Id("k"), Id("v")).Op(":=").Range().Id(mid).Dot("Value")).Block(loop...),
		),
	}
}

// generate nested message marshalling
func (tg *Target) genMessageFieldUnmarshal(f *protogen.Field) []Code {
	return []Code{
		// only unmarshal map, if the attribute is not nil
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			Id("x").Dot(f.GoName).Op("=").New(tg.fieldGoType(f)),
			Err().Op("=").Add(tg.idents.unmarshal).Call(Id("m").Index(Lit(tg.attrName(f))), Id("x").Dot(f.GoName)),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// basic field generates code for marshaling a regular field with basic types
func (tg *Target) genBasicFieldUnmarshal(f *protogen.Field) []Code {
	return []Code{
		Err().Op("=").
			Qual(attributevalues, "Unmarshal").Call(
			Id("m").Index(Lit(tg.attrName(f))),
			Op("&").Id("x").Dot(f.GoName),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': %w"), Err())),
		),
	}
}

// genListFieldUnmarshal generates Unmarshal code for a repeated field
func (tg *Target) genListFieldUnmarshal(f *protogen.Field) []Code {
	mid := fmt.Sprintf("m%d", f.Desc.Number())

	// if its not a list of messages, no recursing is necessary and we can just
	// unmarshal like a basic type
	if f.Message == nil {
		return tg.genBasicFieldUnmarshal(f)
	}

	// for messages we loop over each item and unmarshal them one by one
	return []Code{
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			List(Id(mid), Id("ok")).Op(":=").Id("m").Index(Lit(tg.attrName(f))).Assert(Op("*").Qual(dynamodbtypes, " AttributeValueMemberL")),
			If(Op("!").Id("ok")).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': no list attribute provided"))),
			),
			For(List(Id("k"), Id("v")).Op(":=").Range().Id(mid).Dot("Value")).Block(
				// for list items, they can also be NULL attributes, so we take special
				// care of that scenario.
				If(List(Id("_"), Id("ok")).Op(":=").Id("v").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberNULL")), Id("ok")).Block(
					Id("x").Dot(f.GoName).Op("=").Append(Id("x").Dot(f.GoName), Nil()),
					Continue(),
				),
				// else, init empty message and ummarshal into it
				Var().Id("mv").Add(tg.fieldGoType(f)),
				Err().Op("=").Add(tg.idents.unmarshal).Call(
					Id("v"),
					Op("&").Id("mv"),
				),
				If(Err().Op("!=").Nil()).Block(
					Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal item '%d' of field '"+f.GoName+"': %w"), Id("k"), Err())),
				),
				Id("x").Dot(f.GoName).Op("=").Append(Id("x").Dot(f.GoName), Op("&").Id("mv")),
			),
		),
	}
}

// genOneOfFieldUnmarshal generates unmarshal code for one-of fields
func (tg *Target) genOneOfFieldUnmarshal(f *protogen.Field) []Code {
	marshal := []Code{
		Var().Id("mo").Id(fmt.Sprintf("%s_%s", f.Parent.GoIdent.GoName, f.GoName)),
	}

	switch {
	case f.Message != nil:
		// oneof field is a message
		marshal = append(marshal,
			Id("mo").Dot(f.GoName).Op("=").New(tg.fieldGoType(f)),
			Err().Op("=").Add(tg.idents.unmarshal).Call(Id("m").Index(Lit(tg.attrName(f))), Id("mo").Dot(f.GoName)),
		)
	default:
		// else, assume the oneof field is a basic type
		marshal = append(marshal,
			Err().Op("=").
				Qual(attributevalues, "Unmarshal").Call(
				Id("m").Index(Lit(tg.attrName(f))),
				Op("&").Id("mo").Dot(f.GoName),
			))
	}

	// handle error
	marshal = append(marshal,
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': %w"), Err())),
		),
		Id("x").Dot(f.Oneof.GoName).Op("=").Op("&").Id("mo"),
	)

	return []Code{
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(marshal...),
	}
}

// genMessageUnmarshal generates the unmarshaling logic
func (tg *Target) genMessageUnmarshal(f *File, m *protogen.Message) error {
	var body []Code

	// generate unmarschalling code per field kind
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // don't generate unmarshal code for omitted field
		}

		switch {
		case field.Oneof != nil && !field.Desc.HasOptionalKeyword():
			// special case are explicit oneOf fields (not optional fields)
			body = append(body, tg.genOneOfFieldUnmarshal(field)...)
		case field.Desc.IsList(): // repeated fields
			body = append(body, tg.genListFieldUnmarshal(field)...)
		case field.Desc.IsMap(): // map
			body = append(body, tg.genMapFieldUnmarshal(field)...)
		case field.Message != nil: // (nested) message
			body = append(body, tg.genMessageFieldUnmarshal(field)...)
		default: // other, basic types
			body = append(body, tg.genBasicFieldUnmarshal(field)...)
		}
	}

	body = append(body, Return(Nil()))

	f.Comment(`UnmarshalDynamoItem unmarshals data from a dynamodb attribute map`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("UnmarshalDynamoItem").
		Params(Id("m").Map(String()).Qual(dynamodbtypes, "AttributeValue")).
		Params(Id("err").Id("error")).Block(body...)

	return nil
}
