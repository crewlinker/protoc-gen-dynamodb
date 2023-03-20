package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

// genBasicFieldPath implements the generation of method for building baths for basic type fields
func (tg *Target) genBasicFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(Qual(tg.idents.ddb, "Path")).
		Block(
			Return(
				Id("p").Dot("Append").Call(Lit(fmt.Sprintf("%d", field.Desc.Number()))),
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
	f.Func().
		Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(Id(field.Message.GoIdent.GoName + "Path")).
		Block(
			Return(Id(field.Message.GoIdent.GoName + "Path").Values(
				Id("p").Dot("Append").Call(Lit(fmt.Sprintf("%d", field.Desc.Number()))),
			)),
		)

	return nil
}

// genListFieldPath implements the generation of method for building baths for message type fields
func (tg *Target) genListFieldPath(f *File, m *protogen.Message, field *protogen.Field) error {
	if field.Message == nil {
		return nil // @TODO skip generating lists of basic messages, does't work with generic ddb type
	}

	got := tg.fieldGoType(field, "Path")
	f.Func().Params(Id("p").Id(m.GoIdent.GoName + "Path")).Id(field.GoName).
		Params().
		Params(
			Qual(tg.idents.ddb, "ListPath").Types(got, Op("*").Add(got)),
		).
		Block(
			Return(Qual(tg.idents.ddb, "ListPath").Types(got, Op("*").Add(got)).Values(Dict{
				Id("Path"): Id("p").Dot("Append").Call(Lit(fmt.Sprintf("%d", field.Desc.Number()))),
			})),
		)

	return nil
}

// genMessagePaths k
func (tg *Target) genMessagePaths(f *File, m *protogen.Message) error {
	f.Commentf("//%sPath allows for constructing type-safe expression names", m.GoIdent.GoName)
	f.Type().Id(m.GoIdent.GoName + "Path").Struct(
		Qual(tg.idents.ddb, "Path"),
	)

	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("DynamoPath").
		Params().
		Params(Id(m.GoIdent.GoName + "Path")).
		Block(
			Return(Id(m.GoIdent.GoName + "Path").Values(Qual(tg.idents.ddb, "Path").Values())),
		)

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
