package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

// genBasicFieldPath implements the generation of method for building baths for basic type fields
func (tg *Target) genBasicFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	f.Commentf("%s appends the path being build", field.GoName)
	f.Func().
		Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
		Params().
		Params(Qual(expression, "NameBuilder")).
		Block(
			Return(Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field))))),
		)

	return nil
}

// genMessageFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genMessageFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// (repeated) message field in a different package we don't support any special type-safe path construciton
	// and instead return a basic path that the user can build on further by themselves.
	if !tg.isSamePkgIdent(field.Message.GoIdent) {
		return tg.genBasicFieldPath(f, m, field) // NOTE: may support traversing of this in the future
	}

	// generate the method that append the path element
	f.Commentf("%s returns 'p' with the attribute name appended and allow subselecting nested message", field.GoName)
	f.Func().
		Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
		Params().
		Params(Id(tg.pathIdentName(field.Message))).
		Block(
			Return(Id(tg.pathIdentName(field.Message)).Values(
				Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
			)),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// if it's a list of basic types, or the repeated message is not in the same package we return a
	// a generic basic list type. NOTE: we maybe support external messages in the future
	if field.Message == nil || !tg.isSamePkgIdent(field.Message.GoIdent) {
		// If its not a list of messages, the path will always end and the generated method will return
		// basic path builder that always returns a string with the final path
		f.Commentf("%s returns 'p' appended with the attribute name and allow indexing", field.GoName)
		f.Func().Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddbpath, "List")).
			Block(
				Return(Qual(tg.idents.ddbpath, "List").Values(
					Dict{
						Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
					},
				)),
			)
		return nil
	}

	// else, we assume its a message in the same package
	got := Id(tg.pathIdentName(field.Message))
	f.Commentf("%s returns 'p' appended with the attribute while allow indexing a nested message", field.GoName)
	f.Func().Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddbpath, "ItemList").Types(got)).
		Block(
			Return(Qual(tg.idents.ddbpath, "ItemList").Types(got).Values(Dict{
				Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
			})),
		)

	return nil
}

// genMapFieldPath implements the generation of method for building baths for field of a map type
func (tg *Target) genMapFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	val := field.Message.Fields[1] // value type of the message

	// if it's a map of basic types, or the repeated message is not in the same package we return a
	// a basic list accessor. NOTE: we maybe support external messages in the future
	if val.Message == nil || !tg.isSamePkgIdent(val.Message.GoIdent) {
		f.Commentf("%s returns 'p' appended with the attribute name and allow map keys to be specified", field.GoName)
		f.Func().Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddbpath, "Map")).
			Block(
				Return(Qual(tg.idents.ddbpath, "Map").Values(
					Dict{
						Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
					},
				)),
			)
		return nil
	}

	// else, we assume its a message in the same package
	got := Id(tg.pathIdentName(val.Message))
	f.Commentf("%s returns 'p' appended with the attribute while allow map keys on a nested message", field.GoName)
	f.Func().Params(Id("p").Id(tg.pathIdentName(m))).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddbpath, "ItemMap").Types(got)).
		Block(
			Return(Qual(tg.idents.ddbpath, "ItemMap").Types(got).Values(Dict{
				Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
			})),
		)

	return nil
}

func (tg *Target) pathIdentName(m *protogen.Message) string {
	return m.GoIdent.GoName + "Path"
}

// genFieldRegistration generatiosn the registration code for a field
func (tg *Target) genFieldRegistration(field *protogen.Field) (Code, error) {

	// reflect on fields message for registration
	genFieldMsgReflect := func(m *protogen.Message) *Statement {
		return Qual("reflect", "TypeOf").Call(Id(tg.pathIdentName(m)).Values())
	}

	d := Dict{}
	switch {
	case field.Desc.IsList():
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "ListKind")
		if field.Message != nil && tg.isSamePkgIdent(field.Message.GoIdent) {
			d[Id("Ref")] = genFieldMsgReflect(field.Message)
		}
	case field.Desc.IsMap():
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "MapKind")
		val := field.Message.Fields[1] // value type of the message
		if val.Message != nil && tg.isSamePkgIdent(val.Message.GoIdent) {
			d[Id("Ref")] = genFieldMsgReflect(val.Message)
		}
	case field.Message != nil:
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "BasicKind")
		if field.Message != nil && tg.isSamePkgIdent(field.Message.GoIdent) {
			d[Id("Ref")] = genFieldMsgReflect(field.Message)
		}
	default:
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "BasicKind")
	}

	return Values(d), nil
}

// genMessagePaths generates path building types
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) (err error) {
	f.Commentf("%s allows for constructing type-safe expression names", tg.pathIdentName(m))
	f.Type().Id(tg.pathIdentName(m)).Struct(Qual(expression, "NameBuilder"))

	// generate the "WithDynamoNameBuilder" to allow generic types to set the namebuilder is its decending
	f.Commentf("WithDynamoNameBuilder allows generic types to overwrite the path")
	f.Func().Params(Id("p").Id(tg.pathIdentName(m))).Id("WithDynamoNameBuilder").
		Params(Id("n").Qual(expression, "NameBuilder")).
		Params(Id(tg.pathIdentName(m))).
		Block(
			Id("p").Dot("NameBuilder").Op("=").Id("n"),
			Return(Id("p")),
		)

	// generate path building and field registration
	regFields := Dict{}
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // no path building for ignored fields
		}

		// add each field to the register call in the generated init
		regFields[Lit(tg.attrName(field))], err = tg.genFieldRegistration(field)
		if err != nil {
			return fmt.Errorf("failed to generate field registration: %w", err)
		}

		switch {
		case field.Desc.IsList():
			tg.genListFieldPath(f, m, field)
		case field.Desc.IsMap():
			tg.genMapFieldPath(f, m, field)
		case field.Message != nil:
			tg.genMessageFieldPath(f, m, field)
		default:
			tg.genBasicFieldPath(f, m, field)
		}
	}

	// generate init functions that will register the types for path validation
	f.Func().Id("init").Params().Block(
		Qual(tg.idents.ddbpath, "RegisterMessage").Call(
			Qual("reflect", "TypeOf").Call(Id(tg.pathIdentName(m)).Values()),
			Map(String()).Qual(tg.idents.ddbpath, "FieldInfo").Values(regFields),
		),
	)

	return nil
}
