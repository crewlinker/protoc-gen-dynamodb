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
				Call(Qual(tg.idents.ddb, "P").Values()).Dot("Set").Call(Id("p").Dot("v").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number()))),
			),
		)
	return nil
}

// genMessageFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genMessageFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {

	// for well-known types we need to return specially crafted paths
	switch field.Message.GoIdent.GoImportPath {
	case "google.golang.org/protobuf/types/known/durationpb", // @TODO return ddb.Path
		"google.golang.org/protobuf/types/known/timestamppb", // @TODO return ddb.Path
		"google.golang.org/protobuf/types/known/anypb",       // @TODO return path that allows selecting 0, or 1
		"google.golang.org/protobuf/types/known/structpb",    // @TODO return map select
		"google.golang.org/protobuf/types/known/wrapperspb",  // @TODO return ddb.Path
		"google.golang.org/protobuf/types/known/fieldmaskpb": // @TODO return ddb.IndexPath
		return nil // @TODO implement
	}

	// generate the method that append the path element
	f.Commentf("%s returns 'p' with the attribute name appended and allow subselecting nested message", field.GoName)
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Id(field.Message.GoIdent.GoName + "P")).
		Block(
			Return(Id(field.Message.GoIdent.GoName + "P").Values(Dict{
				Id("v"): Id("p").Dot("v").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())),
			})),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	if field.Message == nil {
		// if its not a list of messages, the path will always end and the generated method will return
		// basic path builder that always returns a string with the final path
		f.Commentf("%s returns 'p' appended with the attribute name and allow indexing", field.GoName)
		f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddb, "BasicListP")).
			Block(
				Return(Call(Qual(tg.idents.ddb, "BasicListP").Values()).
					Dot("Set").Call(Id("p").Dot("v").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())))),
			)
		return nil
	}

	got := tg.fieldGoType(field, "P")
	f.Commentf("%s returns 'p' appended with the attribute while allow indexing a nested message", field.GoName)
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "ListP").Types(got)).
		Block(
			Return(Call(Qual(tg.idents.ddb, "ListP").Types(got).Values()).Dot("Set").Call(
				Id("p").Dot("v").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())),
			)),
		)

	return nil
}

// genMessagePaths k
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) error {
	f.Commentf("%sP allows for constructing type-safe expression names", m.GoIdent.GoName)
	f.Type().Id(m.GoIdent.GoName + "P").Struct(Id("v").String())

	// generate the "Set" method for the path struct, required to make generic list builder work
	f.Commentf("Set allows generic list builder to replace the path value")
	f.Func().Params(Id("p").Id(m.GoIdent.GoName+"P")).Id("Set").
		Params(Id("v").String()).
		Params(Id(m.GoIdent.GoName+"P")).
		Block(
			Id("p").Dot("v").Op("=").Id("v"),
			Return(Id("p")),
		)

	// generate the "String" method for the path struct, required to allow paths to be formatted correctly
	f.Commentf("String formats the path and returns it")
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "P")).Id("String").
		Params().
		Params(String()).
		Block(
			Return(Qual("strings", "TrimPrefix").Call(Id("p").Dot("v"), Lit("."))),
		)

	// generate path building methods for each field
	for _, field := range m.Fields {
		switch {
		case field.Desc.IsList():
			tg.genListFieldPath(f, m, field)
		case field.Desc.IsMap():
			// @TODO implement
		case field.Message != nil:
			tg.genMessageFieldPath(f, m, field)
		default:
			tg.genBasicFieldPath(f, m, field)
		}
	}

	return nil
}
