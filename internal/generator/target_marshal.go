package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
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
			List(Id("m").Index(Lit(tg.attrName(f))), Err()).Op("=").Qual(tg.idents.ddb, "MarshalMappedMessage").Call(
				Id("x").Dot(f.GoName),
				tg.genEmbedOption(f),
			),
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
			List(Id(fmt.Sprintf("m%d", f.Desc.Number())), Id("err")).Op(":=").Qual(tg.idents.ddb, "MarshalMessage").Call(
				Id("x").Dot("Get"+f.GoName).Call(),
				tg.genEmbedOption(f),
			),
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
				Qual(tg.idents.ddb, "Marshal").Call(
				Id("x").Dot("Get"+f.GoName).Call(),
				tg.genEmbedOption(f),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// genSetFieldMarshal generates code to marshal a field into a StringSet, NumberSet or BinarySet
func (tg *Target) genSetFieldMarshal(f *protogen.Field) []Code {
	return []Code{
		If(tg.marshalPresenceCond(f)...).Block(
			List(Id("m").Index(Lit(tg.attrName(f))), Err()).Op("=").Qual(tg.idents.ddb, "MarshalSet").Call(
				Id("x").Dot(f.GoName),
				tg.genEmbedOption(f),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal set item of field '"+f.GoName+"': %w"), Err())),
			),
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
			List(Id("m").Index(Lit(tg.attrName(f))), Err()).Op("=").Qual(tg.idents.ddb, "MarshalRepeatedMessage").Call(
				Id("x").Dot(f.GoName),
				tg.genEmbedOption(f),
			),
			If(Err().Op("!=").Nil()).Block(
				Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal repeated message field '"+f.GoName+"': %w"), Err())),
			),
		),
	}
}

// generateMessageMarshal generates the logic for marshalling messages into dynamo items
func (tg *Target) genMessageMarshal(f *File, m *protogen.Message) error {

	// render method body
	body := []Code{Id("m").Op("=").Make(Map(String()).Qual(types, "AttributeValue"))}

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
			// else, assume basic marshalling can handle it
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
			Id("m").Map(String()).Qual(types, "AttributeValue"),
			Id("err").Id("error"),
		).Block(body...)
	return nil
}
