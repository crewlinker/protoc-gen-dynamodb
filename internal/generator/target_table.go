package generator

import (
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
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

	var keyStructFields []Code
	keyStructName := entityKeyStructName(m)
	keyMapperName := entityKeyMapperName(m)

	tbl := Dict{Id("Name"): Lit(*topts.Name)}
	if tbl[Id("PartitionKey")], err = tg.genKeyAttrLitByNumber(*topts.Pk, m, true); err != nil {
		return fmt.Errorf("failed to determine attr qualifier for partition key: %w", err)
	}

	keyStructFields = append(keyStructFields, Id("Pk").Add(tg.fieldGoType(fieldByNumber(m, *topts.Pk))))

	// optional sort key for table
	if topts.Sk != nil {
		if tbl[Id("SortKey")], err = tg.genKeyAttrLitByNumber(*topts.Sk, m, true); err != nil {
			return fmt.Errorf("failed to determinw attr qualifier for sort key: %w", err)
		}

		keyStructFields = append(keyStructFields, Id("Sk").Add(tg.fieldGoType(fieldByNumber(m, *topts.Sk))))
	}

	// build global secondary index definitions
	var gsivals []Code
	for _, gsi := range topts.Gsi {
		gsid := Dict{Id("Name"): Lit(*gsi.Name)}
		if gsid[Id("PartitionKey")], err = tg.genKeyAttrLitByNumber(*gsi.Pk, m, true); err != nil {
			return fmt.Errorf("failed to determine attr qualifier for gsi pk: %w", err)
		}

		keyStructFields = append(keyStructFields, Id(strcase.ToCamel(*gsi.Name)+"Pk").Op("*").Add(tg.fieldGoType(fieldByNumber(m, *gsi.Pk))))

		if gsi.Sk != nil {
			if gsid[Id("SortKey")], err = tg.genKeyAttrLitByNumber(*gsi.Sk, m, true); err != nil {
				return fmt.Errorf("failed to determine attr qualifier for gsi sk: %w", err)
			}
			keyStructFields = append(keyStructFields, Id(strcase.ToCamel(*gsi.Name)+"Sk").Op("*").Add(tg.fieldGoType(fieldByNumber(m, *gsi.Sk))))
		}

		gsivals = append(gsivals, Values(gsid))
	}

	tbl[Id("GlobalIndexes")] = Index().Op("*").Qual(tg.idents.ddbtable, "GlobalIndex").Values(gsivals...)

	var keyMapperMethods []Code

	// if the tables storese multipe entity types we're looking for one-of with the
	// entity mapping field.
	var entityOneof *protogen.Oneof
	for _, oneof := range m.Oneofs {
		eopts := EntityOptions(oneof)
		if eopts == nil {
			continue
		}

		if entityOneof != nil {
			return fmt.Errorf("can only be one oneOf entity field per table, already saw: %s", entityOneof.GoIdent)
		}

		if tbl[Id("EntityType")], err = tg.genKeyAttrLitByNumber(*eopts.TypeAttr, m, false); err != nil {
			return fmt.Errorf("failed to determine entity type attribute: %w", err)
		}

		// generate key mapper interface methods
		for _, oof := range oneof.Fields {
			keyMapperMethods = append(keyMapperMethods,
				Id(fmt.Sprintf("Map%s", oof.Message.GoIdent.GoName)).Params(
					Op("*").Id(oof.Message.GoIdent.GoName),
				).Params(Id(keyStructName), Error()))
		}

		entityOneof = oneof
	}

	// generate key struct type definition
	f.Commentf("%s is populated by a key mapper to construct index values", keyStructName)
	f.Type().Id(keyStructName).Struct(keyStructFields...)

	// generate key mapper interface
	f.Commentf("%s interface can be implemented to customize how index attributes are build", keyMapperName)
	f.Type().Id(keyMapperName).Interface(keyMapperMethods...)

	// generate FromDynamoEntity method if has an entity defined
	if entityOneof != nil {
		if err = tg.genFromDynamoEntityMethod(f, m, entityOneof); err != nil {
			return fmt.Errorf("failed to generate DynamoFromEntity method: %w", err)
		}
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

func entityKeyStructName(m *protogen.Message) string {
	return fmt.Sprintf("%sKeys", m.GoIdent.GoName)
}

func entityKeyMapperName(m *protogen.Message) string {
	return fmt.Sprintf("%sKeyMapper", m.GoIdent.GoName)
}

func typeFieldEnumValue(enum *protogen.Enum, oof *protogen.Field) *protogen.EnumValue {
	for _, enumv := range enum.Values {
		parts := strings.Split(enumv.GoIdent.GoName, "_")
		if parts[len(parts)-1] == strings.ToUpper(oof.GoName) {
			return enumv
		}
	}

	return nil
}

func (tg *Target) genFromDynamoEntityMethod(f *File, m *protogen.Message, oo *protogen.Oneof) (err error) {
	cases := []Code{Default().Block(
		Return(Qual("fmt", "Errorf").Call(Lit("unsupported entity: %T"), Id("et"))),
	)}

	// generate key assignment code after the switch statement
	var keyAssign []Code
	typeField := fieldByNumber(m, EntityOptions(oo).GetTypeAttr())
	topts := TableOptions(m)
	pkf := fieldByNumber(m, *topts.Pk)
	keyAssign = append(keyAssign, Id("x").Dot(pkf.GoName).Op("=").Id("keys").Dot("Pk"))
	var skf *protogen.Field
	if topts.Sk != nil {
		skf = fieldByNumber(m, *topts.Sk)
		keyAssign = append(keyAssign, Id("x").Dot(skf.GoName).Op("=").Id("keys").Dot("Sk"))
	}

	for _, gsi := range topts.Gsi {
		gsipkf := fieldByNumber(m, *gsi.Pk)
		keyAssign = append(keyAssign, If(Id("keys").Dot(fmt.Sprintf("%sPk", strcase.ToCamel(*gsi.Name))).Op("!=").Nil()).Block(
			Id("x").Dot(gsipkf.GoName).Op("=").Op("*").Id("keys").Dot(fmt.Sprintf("%sPk", strcase.ToCamel(*gsi.Name))),
		))
		if gsi.Sk != nil {
			gsiskf := fieldByNumber(m, *gsi.Sk)
			keyAssign = append(keyAssign, If(Id("keys").Dot(fmt.Sprintf("%sSk", strcase.ToCamel(*gsi.Name))).Op("!=").Nil()).Block(
				Id("x").Dot(gsiskf.GoName).Op("=").Op("*").Id("keys").Dot(fmt.Sprintf("%sSk", strcase.ToCamel(*gsi.Name))),
			))
		}
	}

	// generate switch case code for key mapping caller code
	for _, ef := range oo.Fields {
		if ef.Message == nil {
			return fmt.Errorf("entity oneof can only containe message fields, got: %v", ef.Desc.Kind())
		}

		if !tg.isSamePkgIdent(ef.Message.GoIdent) {
			return fmt.Errorf("entity oneof messages must be in the same package, got: %s", ef.Message.GoIdent)
		}

		typev := typeFieldEnumValue(typeField.Enum, ef)
		if typev == nil {
			return fmt.Errorf("failed to find entity type enum value for oneof field '%s'", ef.GoName)
		}

		cases = append(cases, Case(Op("*").Id(fmt.Sprintf("%s_%s", m.GoIdent.GoName, ef.Message.GoIdent.GoName))).Block(
			Id("x").Dot(typeField.GoName).Op("=").Id(typev.GoIdent.GoName),
			Id("x").Dot(oo.GoName).Op("=").Id("et"),
			List(Id("keys"), Err()).Op("=").Id("m").Dot(fmt.Sprintf("Map%s", ef.Message.GoIdent.GoName)).Call(Id("et").Dot(ef.Message.GoIdent.GoName)),
		))
	}

	// finally, generate the method block
	oneOfIfaceName := fmt.Sprintf("is%s", oo.GoIdent.GoName)
	f.Commentf("FromDynamoEntity propulates the table message from an entity message")
	f.Func().Params(Id("x").Op("*").Id(m.GoIdent.GoName)).Id("FromDynamoEntity").
		Params(Id("e").Id(oneOfIfaceName), Id("m").Id(entityKeyMapperName(m))).
		Params(Err().Error()).
		Block(append([]Code{
			Var().Id("keys").Id(entityKeyStructName(m)),
			Switch(Id("et").Op(":=").Id("e").Assert(Id("type"))).Block(cases...),
			If(Err().Op("!=").Nil()).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("failed to map keys: %w"), Err())),
			),
		}, append(keyAssign, Return())...)...)
	return nil
}
