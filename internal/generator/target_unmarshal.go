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

	// we cannot solve key unmarshalling using type parameters so we determine
	// the correct function here.
	keyType := tg.fieldGoType(key)
	var keyFunc *Statement
	switch key.Desc.Kind() {
	case protoreflect.StringKind:
		keyFunc = Qual(tg.idents.ddb, "StringMapKey")
	case protoreflect.BoolKind:
		keyFunc = Qual(tg.idents.ddb, "BoolMapKey")
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		keyFunc = Qual(tg.idents.ddb, "UintMapKey").Types(keyType)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		keyFunc = Qual(tg.idents.ddb, "IntMapKey").Types(keyType)
	default:
		panic("unsupported map key type: " + key.Desc.Kind().String())
	}

	// defer to the generic unmarshal implementation
	return []Code{
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			List(Id("x").Dot(f.GoName),
				Err()).Op("=").Qual(tg.idents.ddb, "UnmarshalMappedMessage").Types(keyType, tg.fieldGoType(val)).Call(
				Id("m").Index(Lit(tg.attrName(f))),
				keyFunc,
				tg.genEmbedOption(f),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal repeated message field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// generate nested message marshalling
func (tg *Target) genMessageFieldUnmarshal(f *protogen.Field) []Code {
	return []Code{
		// only unmarshal map, if the attribute is not nil
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			Id("x").Dot(f.GoName).Op("=").New(tg.fieldGoType(f)),
			Err().Op("=").Qual(tg.idents.ddb, "UnmarshalMessage").Call(
				Id("m").Index(Lit(tg.attrName(f))),
				Id("x").Dot(f.GoName),
				tg.genEmbedOption(f),
			),
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
			Qual(tg.idents.ddb, "Unmarshal").Call(
			Id("m").Index(Lit(tg.attrName(f))),
			Op("&").Id("x").Dot(f.GoName),
			tg.genEmbedOption(f),
		),
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': %w"), Err())),
		),
	}
}

// genListFieldUnmarshal generates Unmarshal code for a repeated field
func (tg *Target) genListFieldUnmarshal(f *protogen.Field) []Code {
	// if its not a list of messages, no recursing is necessary and we can just
	// unmarshal like a basic type
	if f.Message == nil {
		return tg.genBasicFieldUnmarshal(f)
	}

	// for messages we loop over each item and unmarshal them one by one
	return []Code{
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(
			List(Id("x").Dot(f.GoName),
				Err()).Op("=").Qual(tg.idents.ddb, "UnmarshalRepeatedMessage").Types(tg.fieldGoType(f)).Call(
				Id("m").Index(Lit(tg.attrName(f))),
				tg.genEmbedOption(f),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal repeated message field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// genOneOfFieldUnmarshal generates unmarshal code for one-of fields. This needs special care because
// the optional value is held in a special "FieldPresence_" type.
func (tg *Target) genOneOfFieldUnmarshal(f *protogen.Field) []Code {
	unmarshal := []Code{
		Var().Id("mo").Id(fmt.Sprintf("%s_%s", f.Parent.GoIdent.GoName, f.GoName)),
	}

	switch {
	case f.Message != nil:
		// oneof field is a message
		unmarshal = append(unmarshal,
			Id("mo").Dot(f.GoName).Op("=").New(tg.fieldGoType(f)),
			Err().Op("=").Qual(tg.idents.ddb, "UnmarshalMessage").Call(
				Id("m").Index(Lit(tg.attrName(f))), Id("mo").Dot(f.GoName),
				tg.genEmbedOption(f),
			),
		)
	default:
		// else, assume the oneof field is a basic type
		unmarshal = append(unmarshal,
			Err().Op("=").
				Qual(tg.idents.ddb, "Unmarshal").Call(
				Id("m").Index(Lit(tg.attrName(f))),
				Op("&").Id("mo").Dot(f.GoName),
				tg.genEmbedOption(f),
			))
	}

	// handle error for either case
	unmarshal = append(unmarshal,
		If(Err().Op("!=").Nil()).Block(
			Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal field '"+f.GoName+"': %w"), Err())),
		),
		Id("x").Dot(f.Oneof.GoName).Op("=").Op("&").Id("mo"),
	)

	return []Code{
		If(Id("m").Index(Lit(tg.attrName(f))).Op("!=").Nil()).Block(unmarshal...),
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

	f.Comment(`UnmarshalDynamoItem unmarshals data from a dynamodb attribute map`)
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("UnmarshalDynamoItem").
		Params(Id("m").Map(String()).Qual(types, "AttributeValue")).
		Params(Id("err").Id("error")).Block(append(body, Return(Nil()))...)

	return nil
}
