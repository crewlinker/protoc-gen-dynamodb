package generator

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

// genMessageKeying generates partition/sort key methods on the messages itself
func (tg *Target) genMessageKeying(f *File, m *protogen.Message) (err error) {
	pkf, skf, err := tg.keyFields(m)
	if err != nil {
		return fmt.Errorf("failed to determine key fields: %w", err)
	}

	// if no key fields are configured, so we don't generate a MarshalDynamoKey at all
	if pkf == nil && skf == nil {
		return nil
	}

	if pkf != nil {
		// generate function that returns they partition key as a KeyBuilder for easy key conditions
		f.Commentf("DynamoPartitionKey returns a key builder for the partition key")
		f.Func().
			Params(Id("x").Op("*").Id(m.GoIdent.GoName)).
			Id("DynamoPartitionKey").
			Params().
			Params(Id("v").Qual(expression, "KeyBuilder")).
			Block(Return(Qual(tg.idents.ddbimp, m.GoIdent.GoName+"PartitionKey").Call()))

		// generate function that returns they partition key as a NameBuilder for easy conditions
		f.Commentf("DynamoPartitionKeyName returns a key builder for the partition key")
		f.Func().
			Params(Id("x").Op("*").Id(m.GoIdent.GoName)).
			Id("DynamoPartitionKeyName").
			Params().
			Params(Id("v").Qual(expression, "NameBuilder")).
			Block(Return(Qual(tg.idents.ddbimp, m.GoIdent.GoName+"PartitionKeyName").Call()))
	}

	if skf != nil {
		// generate function that returns they sort key as a KeyBuilder for easy key conditions
		f.Commentf("DynamoSortKey returns a key builder for the sort key")
		f.Func().
			Params(Id("x").Op("*").Id(m.GoIdent.GoName)).
			Id("DynamoSortKey").
			Params().
			Params(Id("v").Qual(expression, "KeyBuilder")).
			Block(Return(Qual(tg.idents.ddbimp, m.GoIdent.GoName+"SortKey").Call()))

		// generate function that returns they sort key as a NameBuilder for easy conditions
		f.Commentf("DynamoSortKeyName returns a key builder for the sort key")
		f.Func().
			Params(Id("x").Op("*").Id(m.GoIdent.GoName)).
			Id("DynamoSortKeyName").
			Params().
			Params(Id("v").Qual(expression, "NameBuilder")).
			Block(Return(Qual(tg.idents.ddbimp, m.GoIdent.GoName+"SortKeyName").Call()))
	}

	// Generate method that returns the key names as a string slice, usefull for masking attribute value maps
	f.Commentf("DynamoKeyNames returns the attribute names of the partition and sort keys respectively")
	f.Func().
		Params(Id("x").Op("*").Id(m.GoIdent.GoName)).
		Id("DynamoKeyNames").
		Params().
		Params(Id("v").Index().String()).
		Block(Return(
			Qual(tg.idents.ddbimp, m.GoIdent.GoName+"KeyNames").Call(),
		))

	// @TODO .DynamoPartitionKey()
	// @TODO .DynamoSortKey()
	// @TODO .DynamoPartitionKeyName()
	// @TODO .DynamoSortKeyName()

	return nil
}

// genDdbKeying generates static partition/sort key functions in the ddb package
func (tg *Target) genDdbKeying(f *File, m *protogen.Message) (err error) {
	pkf, skf, err := tg.keyFields(m)
	if err != nil {
		return fmt.Errorf("failed to determine key fields: %w", err)
	}

	// if no key fields are configured, so we don't generate a MarshalDynamoKey at all
	if pkf == nil && skf == nil {
		return nil
	}

	var body []Code
	if pkf != nil {

		// append to names slice
		body = append(body, Id("v").Op("=").Append(Id("v"), Lit(tg.attrName(pkf))))

		// generate function that returns they partition key as a KeyBuilder
		f.Commentf("%sPartitionKey returns a key builder for the partition key", m.GoIdent.GoName)
		f.Func().
			Id(fmt.Sprintf("%sPartitionKey", m.GoIdent.GoName)).
			Params().
			Params(Id("v").Qual(expression, "KeyBuilder")).
			Block(Return(Qual(expression, "Key").Call(Lit(tg.attrName(pkf)))))

		// generate function that returns they partition key as a NameGuilder
		f.Commentf("%sPartitionKeyName returns a name builder for the partition key", m.GoIdent.GoName)
		f.Func().
			Id(fmt.Sprintf("%sPartitionKeyName", m.GoIdent.GoName)).
			Params().
			Params(Id("v").Qual(expression, "NameBuilder")).
			Block(Return(Qual(expression, "Name").Call(Lit(tg.attrName(pkf)))))
	}

	if skf != nil {

		// append to names slice
		body = append(body, Id("v").Op("=").Append(Id("v"), Lit(tg.attrName(skf))))

		// generate function that returns they sort key as a KeyBuilder
		f.Commentf("%sSortKey returns a key builder for the sort key", m.GoIdent.GoName)
		f.Func().
			Id(fmt.Sprintf("%sSortKey", m.GoIdent.GoName)).
			Params().
			Params(Id("v").Qual(expression, "KeyBuilder")).
			Block(Return(Qual(expression, "Key").Call(Lit(tg.attrName(skf)))))

		// generate function that returns they sort key as a NameBuilder
		f.Commentf("%sSortKeyName returns a name builder for the sort key", m.GoIdent.GoName)
		f.Func().
			Id(fmt.Sprintf("%sSortKeyName", m.GoIdent.GoName)).
			Params().
			Params(Id("v").Qual(expression, "NameBuilder")).
			Block(Return(Qual(expression, "Name").Call(Lit(tg.attrName(skf)))))
	}

	// static function that returns the key names for a certain message
	f.Commentf("%sKeyNames returns the attribute names of the partition and sort keys respectively", m.GoIdent.GoName)
	f.Func().
		Id(fmt.Sprintf("%sKeyNames", m.GoIdent.GoName)).
		Params().
		Params(Id("v").Index().String()).
		Block(append(body, Return())...)

	return nil
}
