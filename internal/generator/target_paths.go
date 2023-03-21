package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

// genBasicFieldPath implements the generation of method for building baths for basic type fields
func (tg *Target) genBasicFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	f.Commentf("%s returns 'p' with the attribute name appended", field.GoName)
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "P")).
		Block(
			Return(
				Call(Qual(tg.idents.ddb, "P").Values()).Dot("Set").Call(Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field)))),
			),
		)
	return nil
}

// genMessageFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genMessageFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// repeated message field in a different package we don't support any special type-safe path construciton
	// and instead return a basic path.
	if !tg.isSamePkgIdent(field.Message.GoIdent) {
		return tg.genBasicFieldPath(f, m, field) // NOTE: may support traversing of this in the future
	}

	// generate the method that append the path element
	f.Commentf("%s returns 'p' with the attribute name appended and allow subselecting nested message", field.GoName)
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Id(field.Message.GoIdent.GoName + "P")).
		Block(
			Return(Id(field.Message.GoIdent.GoName + "P").Values().Dot("Set").Call(
				Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field))),
			)),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// if it's a list of basic types, or the repeated message is not in the same package we return a
	// a basic list accessor. NOTE: we maybe support external messages in the future
	if field.Message == nil || !tg.isSamePkgIdent(field.Message.GoIdent) {
		// If its not a list of messages, the path will always end and the generated method will return
		// basic path builder that always returns a string with the final path
		f.Commentf("%s returns 'p' appended with the attribute name and allow indexing", field.GoName)
		f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddb, "BasicListP")).
			Block(
				Return(Call(Qual(tg.idents.ddb, "BasicListP").Values()).
					Dot("Set").Call(Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field))))),
			)
		return nil
	}

	// else, we assume its a message in the same package
	got := tg.fieldGoType(field, "P")
	f.Commentf("%s returns 'p' appended with the attribute while allow indexing a nested message", field.GoName)
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "ListP").Types(got)).
		Block(
			Return(Call(Qual(tg.idents.ddb, "ListP").Types(got).Values()).Dot("Set").Call(
				Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field))),
			)),
		)

	return nil
}

// genMapFieldPath implements the generation of method for building baths for field of a map type
func (tg *Target) genMapFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	val := field.Message.Fields[1] // value type of the message

	// if it's a list of basic types, or the repeated message is not in the same package we return a
	// a basic list accessor. NOTE: we maybe support external messages in the future
	if val.Message == nil || !tg.isSamePkgIdent(val.Message.GoIdent) {
		f.Commentf("%s returns 'p' appended with the attribute name and allow map keys to be specified", field.GoName)
		f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddb, "BasicMapP")).
			Block(
				Return(Call(Qual(tg.idents.ddb, "BasicMapP").Values()).
					Dot("Set").Call(Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field))))),
			)
		return nil
	}

	// else, we assume its a message in the same package
	got := tg.fieldGoType(val, "P")
	f.Commentf("%s returns 'p' appended with the attribute while allow map keys on a nested message", field.GoName)
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "MapP").Types(got)).
		Block(
			Return(Call(Qual(tg.idents.ddb, "MapP").Types(got).Values()).Dot("Set").Call(
				Id("p").Dot("Val").Call().Op("+").Lit(fmt.Sprintf(".%s", tg.attrName(field))),
			)),
		)

	return nil
}

// genMessagePaths k
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) error {
	f.Commentf("%sP allows for constructing type-safe expression names", m.GoIdent.GoName)
	f.Type().Id(m.GoIdent.GoName + "P").Struct(Qual(tg.idents.ddb, "P"))

	// generate the "Set" method for the path struct, required to make generic list builder work
	f.Commentf("Set allows generic list builder to replace the path value")
	f.Func().Params(Id("p").Id(m.GoIdent.GoName+"P")).Id("Set").
		Params(Id("v").String()).
		Params(Id(m.GoIdent.GoName+"P")).
		Block(
			Id("p").Dot("P").Op("=").Id("p").Dot("P").Dot("Set").Call(Id("v")),
			Return(Id("p")),
		)

	// Generate path function to make it more ergonomic to start a path
	f.Commentf("%sPath starts the building of an expression path into %s", m.GoIdent.GoName, m.GoIdent.GoName)
	f.Func().Id(m.GoIdent.GoName + "Path").
		Params().
		Params(Id(m.GoIdent.GoName + "P")).
		Block(
			Return(Id(m.GoIdent.GoName + "P").Values()),
		)

	// generate path building methods for each field
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // no path building for ignored fields
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

	return nil
}
