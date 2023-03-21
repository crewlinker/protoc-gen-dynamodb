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
		Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(String()).
		Block(
			Return(
				Qual("strings", "TrimPrefix").Call(String().Call(Id("p")).Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())), Lit("."))),
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
		Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(Id(field.Message.GoIdent.GoName + "Path")).
		Block(
			Return(Id(field.Message.GoIdent.GoName + "Path").Call(
				Id("p").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())),
			)),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	if field.Message == nil {
		// if its not a list of messages, the path will always end and the generated method will return
		// basic path builder that always returns a string with the final path
		f.Commentf("%s returns 'p' appended with the attribute name and allow indexing", field.GoName)
		f.Func().Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
			Params().
			Params(Qual(tg.idents.ddb, "BasicListPath")).
			Block(
				Return(Qual(tg.idents.ddb, "BasicListPath").Call(
					Id("p").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())),
				)),
			)
		return nil
	}

	got := tg.fieldGoType(field, "Path")
	f.Commentf("%s returns 'p' appended with the attribute while allow indexing a nested message", field.GoName)
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "ListPath").Types(got)).
		Block(
			Return(Qual(tg.idents.ddb, "ListPath").Types(got).Call(
				Id("p").Op("+").Lit(fmt.Sprintf(".%d", field.Desc.Number())),
			)),
		)

	return nil
}

// genMessagePaths k
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) error {
	f.Commentf("%s allows for constructing type-safe expression names", m.GoIdent.GoName)
	f.Type().Id(m.GoIdent.GoName + "Path").String()

	// generate a function that makes it more ergonomic to start building a path
	f.Commentf("In%s starts the building of a path into a kitchen item", m.GoIdent.GoName)
	f.Func().Id("In" + m.GoIdent.GoName).Params().Params(Id("p").Id(m.GoIdent.GoName + "Path")).
		Block(Return(Id("p")))

	// generate the "Set" method for the path struct, required to make generic list builder work
	f.Commentf("Set allows generic list builder to replace the path value")
	f.Func().Params(Id("p").Id(m.GoIdent.GoName+"Path")).Id("Set").
		Params(Id("v").String()).
		Params(Id(m.GoIdent.GoName+"Path")).
		Block(
			Id("p").Op("=").Id(m.GoIdent.GoName+"Path").Call(Id("v")),
			Return(Id("p")),
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
