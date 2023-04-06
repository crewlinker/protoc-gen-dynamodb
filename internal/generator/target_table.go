package generator

import (
	"fmt"

	ddbexpression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	. "github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/compiler/protogen"
)

func genPkSkAttrType(typ ddbexpression.DynamoDBAttributeType) Code {
	switch typ {
	case ddbexpression.String:
		return Qual(expression, "String")
	case ddbexpression.Number:
		return Qual(expression, "Number")
	case ddbexpression.Binary:
		return Qual(expression, "Binary")
	default:
		panic(fmt.Sprintf("unsupported pk/sk attr type: %T", typ))
	}
}

// genRegisterTablePlacement generates the init function that registers the message's table placement
func (tg *Target) genRegisterTablePlacement(f *File, m *protogen.Message) (err error) {
	placement, err := tg.tablePlacementOptions(m)
	if err != nil {
		return fmt.Errorf("failed to determine table placement options: %w", err)
	}

	if placement == nil {
		return nil // message is not placed in any table
	}

	var tnames []Code
	for _, tname := range placement.tableNames {
		tnames = append(tnames, Lit(tname))
	}

	gsidict := Dict{}
	for name, gsi := range placement.gsis {
		gsidict[Lit(name)] = Values(Dict{
			Id("PartitionKey"): Qual(tg.idents.ddbtable, "KeyDescriptor").Values(Dict{
				Id("AttrName"): Lit(tg.attrName(gsi.pkField)),
				Id("AttrType"): genPkSkAttrType(gsi.pkType),
			}),
			Id("SortKey"): Op("&").Qual(tg.idents.ddbtable, "KeyDescriptor").Values(Dict{
				Id("AttrName"): Lit(tg.attrName(gsi.skField)),
				Id("AttrType"): genPkSkAttrType(gsi.skType),
			}),
			Id("OtherAttrNames"): Index().String().ValuesFunc(func(g *Group) {
				for _, f := range gsi.projected {
					g.Add(Lit(tg.attrName(f)))
				}
			}),
		})
	}

	lsidict := Dict{}
	for name, lsi := range placement.lsis {
		lsidict[Lit(name)] = Values(Dict{
			Id("SortKey"): Qual(tg.idents.ddbtable, "KeyDescriptor").Values(Dict{
				Id("AttrName"): Lit(tg.attrName(lsi.skField)),
				Id("AttrType"): genPkSkAttrType(lsi.skType),
			}),
			Id("OtherAttrNames"): Index().String().ValuesFunc(func(g *Group) {
				for _, f := range lsi.projected {
					g.Add(Lit(tg.attrName(f)))
				}
			}),
		})
	}

	// generate init functions that will register the types for path validation
	f.Func().Id("init").Params().Block(
		Qual(tg.idents.ddbtable, "MustRegister").Call(
			Op("&").Id(m.GoIdent.GoName).Values(),
			Op("&").Qual(tg.idents.ddbtable, "TablePlacement").Values(Dict{
				Id("TableNames"): Index().String().Values(tnames...),
				Id("PartitionKey"): Qual(tg.idents.ddbtable, "KeyDescriptor").Values(Dict{
					Id("AttrName"): Lit(tg.attrName(placement.pkField)),
					Id("AttrType"): genPkSkAttrType(placement.pkType),
				}),
				Id("SortKey"): Op("&").Qual(tg.idents.ddbtable, "KeyDescriptor").Values(Dict{
					Id("AttrName"): Lit(tg.attrName(placement.skField)),
					Id("AttrType"): genPkSkAttrType(placement.skType),
				}),
				Id("GlobalSecondaryIdxs"): Map(String()).Op("*").Qual(tg.idents.ddbtable, "GlobalSecondaryIndexPlacement").
					Values(gsidict),
				Id("LocalSecondaryIdxs"): Map(String()).Op("*").Qual(tg.idents.ddbtable, "LocalSecondaryIndexPlacement").
					Values(lsidict),
			}),
		),
	)

	return
}
