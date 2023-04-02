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
		Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
		Params().
		Params(Qual(expression, "NameBuilder")).
		Block(
			Return(Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field))))),
		)

	return nil
}

// genMessageFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genMessageFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// don't generate recursive path methods if it is an external message, and not embedded
	if tg.notSupportPathing(field) {
		return tg.genBasicFieldPath(f, m, field) // NOTE: may support traversing of this in the future
	}

	// generate the method that append the path element
	f.Commentf("%s returns 'p' with the attribute name appended and allow subselecting nested message", field.GoName)
	f.Func().
		Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
		Params().
		Params(tg.pathStructType(field.Message)).
		Block(
			Return(tg.pathStructType(field.Message).Values(Dict{
				Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
			})),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// if it's a list of basic types, or the repeated message is not in the same package and not well-known.
	// Or if the message is a embedded, then it is also a basic path
	if tg.notSupportPathing(field) {

		// If its not a list of messages, the path will always end and the generated method will return
		// basic path builder that always returns a string with the final path
		f.Commentf("%s returns 'p' appended with the attribute name and allow indexing", field.GoName)
		f.Func().Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
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
	got := tg.pathStructType(field.Message)
	f.Commentf("%s returns 'p' appended with the attribute while allow indexing a nested message", field.GoName)
	f.Func().Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
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
	if tg.notSupportPathing(val) {
		f.Commentf("%s returns 'p' appended with the attribute name and allow map keys to be specified", field.GoName)
		f.Func().Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
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
	got := tg.pathStructType(val.Message)
	f.Commentf("%s returns 'p' appended with the attribute while allow map keys on a nested message", field.GoName)
	f.Func().Params(Id("p").Add(tg.pathStructType(m))).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddbpath, "ItemMap").Types(got)).
		Block(
			Return(Qual(tg.idents.ddbpath, "ItemMap").Types(got).Values(Dict{
				Id("NameBuilder"): Id("p").Dot("AppendName").Call(Qual(expression, "Name").Call(Lit(tg.attrName(field)))),
			})),
		)

	return nil
}

// path struct type name
func (tg *Target) pathStructIdentName(m *protogen.Message) string {
	return m.GoIdent.GoName + "Path"
}

// isWellKnownPathSupported returns true if a message is a well-known message and we support
// generating type-safe path accessors for it
func (tg *Target) isWellKnownPathSupported(m *protogen.Message) bool {
	switch m.GoIdent.String() {
	case `"\"google.golang.org/protobuf/types/known/anypb\"".Any`:
		return true
	case `"\"google.golang.org/protobuf/types/known/structpb\"".Value`:
		return true
	case `"\"google.golang.org/protobuf/types/known/fieldmaskpb\"".FieldMask`:
		return true
	}

	return false
}

// pathStructType returns an identifier or qualifier statement for a path struct.
func (tg *Target) pathStructType(m *protogen.Message) *Statement {
	switch m.GoIdent.String() {
	case `"\"google.golang.org/protobuf/types/known/anypb\"".Any`:
		return Qual(tg.idents.ddbpath, "AnyPath")
	case `"\"google.golang.org/protobuf/types/known/structpb\"".Value`:
		return Qual(tg.idents.ddbpath, "ValuePath")
	case `"\"google.golang.org/protobuf/types/known/fieldmaskpb\"".FieldMask`:
		return Qual(tg.idents.ddbpath, "FieldMaskPath")
	}

	return Id(tg.pathStructIdentName(m))
}

// genFieldRegistration generatiosn the registration code for a field
func (tg *Target) genFieldRegistration(field *protogen.Field) (Code, error) {

	// reflect on fields message for registration
	genFieldMsgReflect := func(m *protogen.Message) *Statement {
		return Qual("reflect", "TypeOf").Call(Add(tg.pathStructType(m)).Values())
	}

	d := Dict{}
	switch {
	case field.Desc.IsList():
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "FieldKindList")
		if !tg.notSupportPathing(field) {
			d[Id("Message")] = genFieldMsgReflect(field.Message)
		}
	case field.Desc.IsMap():
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "FieldKindMap")
		val := field.Message.Fields[1] // value type of the message
		if !tg.notSupportPathing(val) {
			d[Id("Message")] = genFieldMsgReflect(val.Message)
		}
	case field.Message != nil:
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "FieldKindSingle")
		if !tg.notSupportPathing(field) {
			d[Id("Message")] = genFieldMsgReflect(field.Message)
		}
	default:
		d[Id("Kind")] = Qual(tg.idents.ddbpath, "FieldKindSingle")
	}

	return Values(d), nil
}

// genMessagePaths generates path building types
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) (err error) {
	f.Commentf("%s allows for constructing type-safe expression names", tg.pathStructIdentName(m))
	f.Type().Add(tg.pathStructType(m)).Struct(Qual(expression, "NameBuilder"))

	// generate the "WithDynamoNameBuilder" to allow generic types to set the namebuilder is its decending
	f.Commentf("WithDynamoNameBuilder allows generic types to overwrite the path")
	f.Func().Params(Id("p").Add(tg.pathStructType(m))).Id("WithDynamoNameBuilder").
		Params(Id("n").Qual(expression, "NameBuilder")).
		Params(Add(tg.pathStructType(m))).
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
		Qual(tg.idents.ddbpath, "Register").Call(
			tg.pathStructType(m).Values(),
			Map(String()).Qual(tg.idents.ddbpath, "FieldInfo").Values(regFields),
		),
	)

	return nil
}
