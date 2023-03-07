package generator

import (
	//lint:ignore ST1001 we want expressive meta code
	. "github.com/dave/jennifer/jen"
)

// genCentralMarshal generates the central marshalling code. The generated code takes an Any value and
// asserts if the it implements the DynamoItemUnmarshaller interface. If not it also supports well-known types
// that we allow to be unmarshalled.
func (tg *Target) genCentralMarshal(f *File) error {
	f.Commentf("%s marshals into DynamoDB attribute value maps", tg.idents.marshal)
	f.Func().Id(tg.idents.marshal).
		Params(Id("x").Qual("google.golang.org/protobuf/proto", "Message")).
		Params(Qual(dynamodbtypes, "AttributeValue"), Error()).
		Block(
			// if the passed in type implements the marshal interface we can call it directly
			If(List(Id("mx"), Id("ok")).Op(":=").Id("x").Assert(Interface(
				Id("MarshalDynamoItem").Params().Params(Map(String()).Qual(dynamodbtypes, "AttributeValue"), Error()),
			)), Id("ok")).Block(
				List(Id("mm"), Err()).Op(":=").Id("mx").Dot("MarshalDynamoItem").Call(),
				Return(Op("&").Qual(dynamodbtypes, "AttributeValueMemberM").Values(Dict{Id("Value"): Id("mm")}), Err()),
			),

			Switch(Id("xt").Op(":=").Id("x").Assert(Type())).Block(
				// Duration, Timestamp are encoded as strings using protojson
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/durationpb", "Duration"),
					Op("*").Qual("google.golang.org/protobuf/types/known/timestamppb", "Timestamp"),
				).Block(
					List(Id("xjson"), Err()).Op(":=").Qual("google.golang.org/protobuf/encoding/protojson", "Marshal").Call(Id("xt")),
					If(Err().Op("!=").Nil()).Block(
						Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal duration: %w"), Err())),
					),
					List(Id("xjsons"), Err()).Op(":=").Qual("strconv", "Unquote").Call(String().Call(Id("xjson"))),
					If(Err().Op("!=").Nil()).Block(
						Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to unquote marshalled duration: %w"), Err())),
					),
					Return(Op("&").Qual(dynamodbtypes, "AttributeValueMemberS").Values(Dict{Id("Value"): Id("xjsons")}), Nil()),
				),

				// or, any other type, return unsupported message
				Default().Block(
					Return(Nil(), Qual("fmt", "Errorf").Call(Lit("marshal of message type unsupported: %+T"), Id("xt"))),
				),
			),
		)
	return nil
}

// genCentralUnmarshal generates the central Unmarshal function. The generated code takes an Any value and
// asserts if the it implements the DynamoItemUnmarshaller interface. If not it also supports well-known types
// that we allow to be unmarshalled.
func (tg *Target) genCentralUnmarshal(f *File) error {
	f.Commentf("%s unmarshals DynamoDB attribute value maps", tg.idents.marshal)
	f.Func().Id(tg.idents.unmarshal).
		Params(Id("m").Qual(dynamodbtypes, "AttributeValue"), Id("x").Qual("google.golang.org/protobuf/proto", "Message")).
		Params(Error()).
		Block(
			// assert to Unmarshaller interface, if so hand off unmarshalling
			If(List(Id("mx"), Id("ok")).Op(":=").Id("x").Assert(Interface(
				Id("UnmarshalDynamoItem").Params(Map(String()).Qual(dynamodbtypes, "AttributeValue")).Params(Error()),
			)), Id("ok")).Block(
				List(Id("mm"), Id("ok")).Op(":=").Id("m").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberM")),
				If(Op("!").Id("ok")).Block(
					Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal: no map attribute provided"))),
				),

				Return(Id("mx").Dot("UnmarshalDynamoItem").Call(Id("mm").Dot("Value"))),
			),

			// Else, take care of some special cases for well-known/common types
			Switch(Id("xt").Op(":=").Id("x").Assert(Type())).Block(
				// Duration, Timestamp are unmarshalled using protojson
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/durationpb", "Duration"),
					Op("*").Qual("google.golang.org/protobuf/types/known/timestamppb", "Timestamp"),
				).Block(
					List(Id("ms"), Id("ok")).Op(":=").Id("m").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberS")),
					If(Op("!").Id("ok")).Block(
						Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal duration: no string attribute provided"))),
					),
					Return(Qual("google.golang.org/protobuf/encoding/protojson", "Unmarshal").Call(
						Index().Byte().Call(Qual("strconv", "Quote").Call(Id("ms").Dot("Value"))),
						Id("x"))),
				),

				// or, any other type, return unsupported message
				Default().Block(
					Return(Qual("fmt", "Errorf").Call(Lit("unmarshal of message type unsupported: %+T"), Id("xt"))),
				),
			),
		)

	return nil
}
