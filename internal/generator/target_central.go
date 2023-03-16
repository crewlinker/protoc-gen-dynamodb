package generator

import (
	. "github.com/dave/jennifer/jen"
)

// genCentralMarshal generates the central marshalling code. The generated code takes an Any value and
// asserts if the it implements the DynamoItemUnmarshaller interface. If not it also supports well-known types
// that we allow to be unmarshalled.
func (tg *Target) genCentralMarshal(f *File) error {
	f.Commentf("%s marshals into DynamoDB attribute value maps", tg.idents.marshal)
	f.Func().Id(tg.idents.marshal).
		Params(Id("x").Qual("google.golang.org/protobuf/proto", "Message")).
		Params(Id("a").Qual(dynamodbtypes, "AttributeValue"), Id("err").Error()).
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
						Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to unquote value: %w"), Err())),
					),
					Return(Op("&").Qual(dynamodbtypes, "AttributeValueMemberS").Values(Dict{Id("Value"): Id("xjsons")}), Nil()),
				),

				// Any type is marshalled using field numbers, into a map
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/anypb", "Any"),
				).Block(
					Id("mv").Op(":=").Op("&").Qual(dynamodbtypes, "AttributeValueMemberM").Values(
						Dict{Id("Value"): Map(String()).Qual(dynamodbtypes, "AttributeValue").Values()}),
					List(Id("mv").Dot("Value").Index(Lit("1")), Err()).Op("=").Qual(attributevalues, "Marshal").Call(Id("xt").Dot("TypeUrl")),
					If(Err().Op("!=").Nil()).Block(
						Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal Any's TypeURL field: %w"), Err())),
					),
					List(Id("mv").Dot("Value").Index(Lit("2")), Err()).Op("=").Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")),
					If(Err().Op("!=").Nil()).Block(
						Return(Nil(), Qual("fmt", "Errorf").Call(Lit("failed to marshal Any's Value field: %w"), Err())),
					),
					Return(Id("mv"), Nil()),
				),

				// FieldMask type is marshalled as a string set
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/fieldmaskpb", "FieldMask"),
				).Block(
					Return(Op("&").Qual(dynamodbtypes, "AttributeValueMemberSS").Values(
						Dict{Id("Value"): Id("xt").Dot("Paths")}), Nil()),
				),

				// Structpb Value are marshalled using sdk attribute marshalling on their interface{} value
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/structpb", "Value"),
				).Block(
					Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("AsInterface").Call())),
				),

				// WrappersPb values just marshal the wrapped value
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "StringValue")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "BoolValue")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "BytesValue")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "DoubleValue")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "FloatValue")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "Int32Value")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "Int64Value")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "UInt32Value")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "UInt64Value")).
					Block(Return(Qual(attributevalues, "Marshal").Call(Id("xt").Dot("Value")))),

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
		Params(Id("err").Error()).
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

				// durationpb.Duration, timestamppb.Timestamp are unmarshalled using protojson
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

				// anypb.Any type is unmarshalled using field numbers
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/anypb", "Any"),
				).Block(
					List(Id("mm"), Id("ok")).Op(":=").Id("m").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberM")),
					If(Op("!").Id("ok")).Block(
						Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal duration: no map attribute provided"))),
					),

					Err().Op("=").Qual(attributevalues, "Unmarshal").Call(
						Id("mm").Dot("Value").Index(Lit("1")),
						Op("&").Id("xt").Dot("TypeUrl")),
					If(Err().Op("!=").Nil()).Block(
						Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal Any's TypeURL field: %w"), Err())),
					),

					Err().Op("=").Qual(attributevalues, "Unmarshal").Call(
						Id("mm").Dot("Value").Index(Lit("2")),
						Op("&").Id("xt").Dot("Value")),
					If(Err().Op("!=").Nil()).Block(
						Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal Any's Value field: %w"), Err())),
					),

					Return(Nil()),
				),

				// fieldmaskpb.FieldMask type is unmarshalled from a stringset
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/fieldmaskpb", "FieldMask"),
				).Block(
					List(Id("ss"), Id("ok")).Op(":=").Id("m").Assert(Op("*").Qual(dynamodbtypes, "AttributeValueMemberSS")),
					If(Op("!").Id("ok")).Block(
						Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal duration: no string set attribute provided"))),
					),

					Id("xt").Dot("Paths").Op("=").Id("ss").Dot("Value"),
					Return(Nil()),
				),

				// structpb.Value types are not self describing. So we have to unmarshal based on the
				// dynamo representation. This is a bit painfull and luckily dynamo types map nicely on
				// json-oriented nature of structpb.Value
				Case(
					Op("*").Qual("google.golang.org/protobuf/types/known/structpb", "Value"),
				).Block(
					Switch(Id("m").Assert(Type())).Block(
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberL")).Block(
							Id("vx").Op(":=").Index().Any().Values(),
							Err().Op("=").Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal structpb Value field: %w"), Err())),
							),
							List(Id("lv"), Err()).Op(":=").Qual("google.golang.org/protobuf/types/known/structpb", "NewList").Call(Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to init structpb.Value: %w"), Err())),
							),
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_ListValue").
								Values(Dict{Id("ListValue"): Id("lv")}),
							Return(Nil()),
						),
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberM")).Block(
							Id("vx").Op(":=").Map(String()).Any().Values(),
							Err().Op("=").Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal structpb Value field: %w"), Err())),
							),
							List(Id("lv"), Err()).Op(":=").Qual("google.golang.org/protobuf/types/known/structpb", "NewStruct").Call(Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to init structpb.Value: %w"), Err())),
							),
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_StructValue").
								Values(Dict{Id("StructValue"): Id("lv")}),
							Return(Nil()),
						),
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberS")).Block(
							Var().Id("vx").String(),
							Err().Op("=").Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal structpb Value field: %w"), Err())),
							),
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_StringValue").
								Values(Dict{Id("StringValue"): Id("vx")}),
							Return(Nil()),
						),
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberBOOL")).Block(
							Var().Id("vx").Bool(),
							Err().Op("=").Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal structpb Value field: %w"), Err())),
							),
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_BoolValue").
								Values(Dict{Id("BoolValue"): Id("vx")}),
							Return(Nil()),
						),
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberN")).Block(
							Var().Id("vx").Float64(),
							Err().Op("=").Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("vx")),
							If(Err().Op("!=").Nil()).Block(
								Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal structpb Value field: %w"), Err())),
							),
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_NumberValue").
								Values(Dict{Id("NumberValue"): Id("vx")}),
							Return(Nil()),
						),
						// in case of a nill value, we don't have to unmarshal and just assign it
						Case(Op("*").Qual(dynamodbtypes, "AttributeValueMemberNULL")).Block(
							Id("xt").Dot("Kind").Op("=").Op("&").
								Qual("google.golang.org/protobuf/types/known/structpb", "Value_NullValue").
								Values(Dict{Id("NullValue"): Qual("google.golang.org/protobuf/types/known/structpb", "NullValue_NULL_VALUE")}),
							Return(Nil()),
						),
						Default().Block(
							Return(Qual("fmt", "Errorf").Call(Lit("failed to unmarshal struct value: unsupported attribute value"))),
						),
					),
				),

				// unmarshal all wrapper types by using the value member
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "StringValue")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "BoolValue")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "BytesValue")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "DoubleValue")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "FloatValue")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "Int32Value")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "Int64Value")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "UInt32Value")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),
				Case(Op("*").Qual("google.golang.org/protobuf/types/known/wrapperspb", "UInt64Value")).
					Block(Return(Qual(attributevalues, "Unmarshal").Call(Id("m"), Op("&").Id("xt").Dot("Value")))),

				// explicitely don't support anything else
				Default().Block(
					Return(Qual("fmt", "Errorf").Call(Lit("unmarshal of message type unsupported: %+T"), Id("xt"))),
				),
			),
		)

	return nil
}
