package generator

import (
	"fmt"
	"io"

	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	. "github.com/dave/jennifer/jen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// we refer to this in the code in all sorts of places so lets setup a handy shortcut
	types = "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	// expression package is used a lot as well
	expression = "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

// Target facilitates generation from a single protobuf file
type Target struct {
	src    *protogen.File
	logs   *zap.Logger
	idents struct {
		ddb     string
		ddbv1   string
		ddbpath string
		ddbimp  string
	}
}

// genEmbedOption generates the statement to configure embed encoding
func (tg *Target) genEmbedOption(f *protogen.Field) *Statement {
	switch tg.embedEncoding(f) {
	case ddbv1.Encoding_ENCODING_JSON:
		return Qual(tg.idents.ddb, "Embed").Call(Qual(tg.idents.ddbv1, "Encoding_ENCODING_JSON"))
	case ddbv1.Encoding_ENCODING_DYNAMO:
		return Qual(tg.idents.ddb, "Embed").Call(Qual(tg.idents.ddbv1, "Encoding_ENCODING_DYNAMO"))
	default:
		return Qual(tg.idents.ddb, "Embed").Call(Qual(tg.idents.ddbv1, "Encoding_ENCODING_DYNAMO"))
	}
}

// returns true as the identifier is part of the package we're generating for
func (tg *Target) isSamePkgIdent(ident protogen.GoIdent) bool {
	return ident.GoImportPath == tg.src.GoImportPath
}

// keyFields consults the the fields of the message and return the fields that describe pk/sks
func (tg *Target) keyFields(m *protogen.Message) (pkf, skf *protogen.Field, err error) {
	for _, field := range m.Fields {
		if tg.isOmitted(field) {
			continue // omitted, don't try to turn it into a key
		}

		isPk, isSk := tg.isKey(field)
		if isPk && isSk {
			// check that field cannot be marked as both sk and pk
			return nil, nil, fmt.Errorf("field '%s' is both marked as PK and as SK", field.GoName)
		}

		if isPk {
			if pkf != nil { // only one field can be marked as PK
				return nil, nil, fmt.Errorf("field '%s' is already marked as PK", pkf.GoName)
			}

			pkf = field
			if !tg.isValidKeyField(pkf) {
				return nil, nil, fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a PK", pkf.GoName)
			}
		}

		if isSk {
			if skf != nil { // only one field can be marked as SK
				return nil, nil, fmt.Errorf("field '%s' is already marked as SK", skf.GoName)
			}

			skf = field
			if !tg.isValidKeyField(skf) {
				return nil, nil, fmt.Errorf("field '%s' must be a basic type that marshals to Number,String or Bytes to be a SK", skf.GoName)
			}
		}
	}

	// if the message has a sort key, but not a partition key that doesn't make sense. The other
	// way around is ok.
	if pkf == nil && skf != nil {
		return nil, nil, fmt.Errorf("message '%s' has a sort key, but not a partition key", m.GoIdent.GoName)
	}

	return
}

// isValidKeyField returns whether a protobuf field can be a valid key
func (tg *Target) isValidKeyField(f *protogen.Field) bool {
	if f.Message != nil {
		return false // only basic types can be keys
	}

	switch f.Desc.Kind() {
	case protoreflect.StringKind,
		protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
		protoreflect.BytesKind,
		protoreflect.Fixed64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind,
		protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.DoubleKind,
		protoreflect.FloatKind:
		return true // what encodes to dynamo Number, String or Bytes
	default:
		return false
	}
}

// fieldGoType turns a field protoreflect kind into a go type
func (tg *Target) fieldGoType(f *protogen.Field, msgSuffix ...string) *Statement {
	if f.Message != nil {
		identName := f.Message.GoIdent.GoName
		if len(msgSuffix) > 0 {
			identName = identName + msgSuffix[0]
		}

		// if the message is from the same package path as we're generating for, assume we refer
		// to it without fullq qualifier
		if tg.isSamePkgIdent(f.Message.GoIdent) {
			return Id(identName)
		}

		// else refer to it with qualifier
		return Qual(string(f.Message.GoIdent.GoImportPath), identName)
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

// GeneratePathBuilding generates code for type-safe document pathing building
func (tg *Target) GeneratePathBuilding(w io.Writer) error {
	pkgname := string(tg.src.GoPackageName + "ddb")
	f := NewFile(pkgname)
	f.PackageComment(fmt.Sprintf("Package %s holds generated schema structure", pkgname))
	f.HeaderComment("Code generated by protoc-gen-dynamodb. DO NOT EDIT.")

	// generate per message dynamo logic
	for _, m := range tg.src.Messages {
		// generate the message paths
		if err := tg.genMessagePaths(f, m); err != nil {
			return fmt.Errorf("failed to generate message path building: %w", err)
		}

		// generate pk/sk methods
		if err := tg.genDdbKeying(f, m); err != nil {
			return fmt.Errorf("failed to generate keying: %w", err)
		}
	}

	return f.Render(w)
}

// GenerateMessageLogic peforms the actual code generation
func (tg *Target) GenerateMessageLogic(w io.Writer) error {
	f := NewFile(string(tg.src.GoPackageName))
	f.HeaderComment("Code generated by protoc-gen-dynamodb. DO NOT EDIT.")

	// generate per message marshal/unmarshal code
	for _, m := range tg.src.Messages {

		// generate the marshal method
		if err := tg.genMessageMarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate marshal: %w", err)
		}

		// generate the unmarshal method
		if err := tg.genMessageUnmarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate unmarshal: %w", err)
		}

		// generate methods that return key info
		if err := tg.genMessageKeying(f, m); err != nil {
			return fmt.Errorf("failed to generate key methods: %w", err)
		}

	}

	return f.Render(w)
}
