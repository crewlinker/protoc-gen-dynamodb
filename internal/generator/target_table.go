package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// fieldByNumber returns a protogen field description by the field nr of a message
func fieldByNumber(m *protogen.Message, n uint32) *protogen.Field {
	for _, f := range m.Fields {
		if f.Desc.Number() == protowire.Number(n) {
			return f
		}
	}
	return nil
}

// fieldToTableAttr turns the field number into an attribute literal, for key or entity type field
func (tg *Target) genKeyAttrLitByNumber(nr uint32, m *protogen.Message, isKeyElseEntity bool) (Code, error) {
	f := fieldByNumber(m, nr)
	if f == nil {
		return nil, fmt.Errorf("field nr '%d' was referenced, but couldn't find it on message '%s'", nr, m.GoIdent.String())
	}

	if isKeyElseEntity && !tg.isValidKeyField(f) {
		return nil, fmt.Errorf("not a valid type for a key field: %s", f.Desc.Kind())
	} else if !isKeyElseEntity && f.Desc.Kind() != protoreflect.EnumKind {
		return nil, fmt.Errorf("entity type field must be an enum, got: %s", f.Desc.Kind())
	}

	var typ Code
	switch f.Desc.Kind() {
	case protoreflect.EnumKind:
		typ = Qual(expression, "String")
	case protoreflect.StringKind:
		typ = Qual(expression, "String")
	case protoreflect.BytesKind:
		typ = Qual(expression, "Binary")
	case protoreflect.Int64Kind,
		protoreflect.Uint64Kind,
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
		typ = Qual(expression, "Number")
	default:
		return nil, fmt.Errorf("unsupported kind for to reference: %v", f.Desc.Kind())
	}

	return Op("&").Qual(tg.idents.ddbtable, "Attribute").Values(Dict{
		Id("Name"): Lit(tg.attrName(f)),
		Id("Type"): typ,
	}), nil
}

// genTableRegistration generates path building types
func (tg *Target) genTableRegistration(f *File, m *protogen.Message) (err error) {
	topts := TableOptions(m)
	if topts == nil {
		return nil // no table options, nothing to register
	}

	tbl := Dict{Id("Name"): Lit(*topts.Name)}
	if tbl[Id("PartitionKey")], err = tg.genKeyAttrLitByNumber(*topts.Pk, m, true); err != nil {
		return fmt.Errorf("failed to determine attr qualifier for partition key: %w", err)
	}

	// optional sort key for table
	if topts.Sk != nil {
		if tbl[Id("SortKey")], err = tg.genKeyAttrLitByNumber(*topts.Sk, m, true); err != nil {
			return fmt.Errorf("failed to determinw attr qualifier for sort key: %w", err)
		}
	}

	// build global secondary index definitions
	var gsivals []Code
	for _, gsi := range topts.Gsi {
		gsid := Dict{Id("Name"): Lit(*gsi.Name)}
		if gsid[Id("PartitionKey")], err = tg.genKeyAttrLitByNumber(*gsi.Pk, m, true); err != nil {
			return fmt.Errorf("failed to determine attr qualifier for gsi pk: %w", err)
		}

		if gsi.Sk != nil {
			if gsid[Id("SortKey")], err = tg.genKeyAttrLitByNumber(*gsi.Sk, m, true); err != nil {
				return fmt.Errorf("failed to determine attr qualifier for gsi sk: %w", err)
			}
		}

		gsivals = append(gsivals, Values(gsid))
	}

	tbl[Id("GlobalIndexes")] = Index().Op("*").Qual(tg.idents.ddbtable, "GlobalIndex").Values(gsivals...)

	// if the tables storese multipe entity types we're looking for one-of with the
	// entity mapping field.
	entityFieldDetermined := false
	for _, oneof := range m.Oneofs {
		eopts := EntityOptions(oneof)
		if eopts == nil {
			continue
		}

		if entityFieldDetermined {
			return fmt.Errorf("can only be one oneOf entity field per table")
		}

		if tbl[Id("EntityType")], err = tg.genKeyAttrLitByNumber(*eopts.TypeAttr, m, false); err != nil {
			return fmt.Errorf("failed to determine entity type attribute: %w", err)
		}

		entityFieldDetermined = true
	}

	// run the actual init method
	f.Commentf("%sTableDefinition can be used to register the table in the ddbtable registry", m.GoIdent.GoName)
	f.Var().Id(fmt.Sprintf("%sTableDefinition", m.GoIdent.GoName)).Op("=").Qual(tg.idents.ddbtable, "Table").Values(tbl)
	f.Comment("register table in the default registry")
	f.Func().Id("init").Params().Block(
		Qual(tg.idents.ddbtable, "Register").Call(Op("&").Id(fmt.Sprintf("%sTableDefinition", m.GoIdent.GoName))),
	)

	return nil
}
