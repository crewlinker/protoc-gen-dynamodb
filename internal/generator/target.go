package generator

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	//lint:ignore ST1001 we want expressive meta code
	. "github.com/dave/jennifer/jen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	attributevalues = "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbtypes   = "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Target facilitates generation from a single protobuf file
type Target struct {
	src    *protogen.File
	logs   *zap.Logger
	idents struct {
		marshal   string
		unmarshal string
	}
}

// determine the dyanmodb attribute name given the field definition
func (tg *Target) attrName(f *protogen.Field) string {
	if fopts := FieldOptions(f); fopts != nil && fopts.Name != nil {
		return *fopts.Name // explicit name option
	}

	return strconv.FormatInt(int64(f.Desc.Number()), 10)
}

// returns true as the identifier is part of the package we're generating for
func (tg *Target) isSamePkgIdent(ident protogen.GoIdent) bool {
	return ident.GoImportPath == tg.src.GoImportPath
}

// fieldGoType turns a field protoreflect kind into a go type
func (tg *Target) fieldGoType(f *protogen.Field) *Statement {
	if f.Message != nil {
		// if the message is from the same package path as we're generating for, assume we refer
		// to it without fullq qualifier
		if tg.isSamePkgIdent(f.Message.GoIdent) {
			return Id(f.Message.GoIdent.GoName)
		}

		// else refer to it with qualifier
		return Qual(string(f.Message.GoIdent.GoImportPath), f.Message.GoIdent.GoName)
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind, protoreflect.BoolKind,
		protoreflect.Int64Kind, protoreflect.Uint64Kind:
		return Id(f.Desc.Kind().String())
	case protoreflect.BytesKind:
		return Id("[]byte")
	case protoreflect.Fixed64Kind:
		return Id("uint64")
	case protoreflect.Sint64Kind:
		return Id("int64")
	case protoreflect.Sfixed64Kind:
		return Id("int64")
	case protoreflect.Int32Kind, protoreflect.Uint32Kind:
		return Id(f.Desc.Kind().String())
	case protoreflect.Fixed32Kind:
		return Id("uint32")
	case protoreflect.Sint32Kind:
		return Id("int32")
	case protoreflect.Sfixed32Kind:
		return Id("int32")
	case protoreflect.DoubleKind:
		return Id("float64")
	case protoreflect.FloatKind:
		return Id("float32")
	default:
		panic("unsupported field type: " + f.Desc.Kind().String())
	}
}

// fieldZeroValue determines the literal that is asserted against to determine if
// the field should be added to the result attribute map
func (tg *Target) fieldZeroValue(f *protogen.Field) *Statement {
	if f.Message != nil ||
		f.Desc.IsList() ||
		f.Desc.HasPresence() {
		return Nil()
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return Lit("")
	case protoreflect.BoolKind:
		return Lit(false)
	case protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind, protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.DoubleKind,
		protoreflect.FloatKind:
		return Lit(0)
	case protoreflect.BytesKind:
		return Nil()
	case protoreflect.EnumKind:
		return Lit(0)
	default:
		panic("unsupported zero value: " + f.Desc.Kind().String())
	}
}

// marshalPresenceCond generates code for the if statements condition that checks if the field
// should be included in the marshalled attribute map.
func (tg *Target) marshalPresenceCond(f *protogen.Field) []Code {
	switch {
	case f.Oneof != nil && !f.Desc.HasOptionalKeyword():
		return []Code{
			List(Id("onev"), Id("ok")).Op(":=").Id("x").Dot(f.Oneof.GoName).Assert(Op("*").
				Id(fmt.Sprintf("%s_%s", f.Parent.GoIdent.GoName, f.GoName))),
			Id("ok").Op("&&").Id("onev").Op("!=").Add(tg.fieldZeroValue(f)),
		}
	case f.Desc.IsList(), f.Desc.IsMap():
		return []Code{Len(Id("x").Dot(f.GoName)).Op("!=").Lit(0)}
	default:
		return []Code{Id("x").Dot(f.GoName).Op("!=").Add(tg.fieldZeroValue(f))}
	}
}

// Generate peforms the actual code generation
func (tg *Target) Generate(w io.Writer) error {
	f := NewFile(string(tg.src.GoPackageName))

	tg.idents.marshal, tg.idents.unmarshal =
		fmt.Sprintf("%s_marshal_dynamo_item", strings.ToLower(tg.src.GoDescriptorIdent.GoName)),
		fmt.Sprintf("%s_unmarshal_dynamo_item", strings.ToLower(tg.src.GoDescriptorIdent.GoName))

	// generate a single marshal function for messages. This way we can handle externally included messages
	// and locally generated message in the same way.
	if err := tg.genCentralMarshal(f); err != nil {
		return fmt.Errorf("failed to generate central message marshal: %w", err)
	}

	// generate a single unmarshal function per proto file so we can unmarshal external messages and
	// messages local to the package in the same way.
	if err := tg.genCentralUnmarshal(f); err != nil {
		return fmt.Errorf("failed to generate central message unmarshal: %w", err)
	}

	// generate per message marshal/unmarshal code
	for _, m := range tg.src.Messages {
		if err := tg.genMessageMarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate marshal: %w", err)
		}
		if err := tg.genMessageUnmarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate unmarshal: %w", err)
		}
	}

	return f.Render(w)
}
