// Package ddb provides DynamoDB utility for Protobuf messages
package ddb

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// MarshalDynamoMessage will marshal a protobuf message 'm' into an attribute value. It supports several
// well-known Protobuf types and if 'x' implements its own MarshalDynamoItem method it will be called to
// delegate the marshalling.
func MarshalDynamoMessage(x proto.Message, opts ...EncodingOption) (a types.AttributeValue, err error) {
	if mx, ok := x.(interface {
		MarshalDynamoItem(...EncodingOption) (map[string]types.AttributeValue, error)
	}); ok {
		mm, err := mx.MarshalDynamoItem(opts...)
		return &types.AttributeValueMemberM{Value: mm}, err
	}

	switch xt := x.(type) {
	case *durationpb.Duration, *timestamppb.Timestamp:
		xjson, err := protojson.Marshal(xt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal duration: %w", err)
		}
		xjsons, err := strconv.Unquote(string(xjson))
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value: %w", err)
		}
		return &types.AttributeValueMemberS{Value: xjsons}, nil
	case *anypb.Any:
		mv := &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}}
		mv.Value["1"], err = attributevalue.MarshalWithOptions(xt.TypeUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's TypeURL field: %w", err)
		}
		mv.Value["2"], err = attributevalue.MarshalWithOptions(xt.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's Value field: %w", err)
		}
		return mv, nil
	case *fieldmaskpb.FieldMask:
		return &types.AttributeValueMemberSS{Value: xt.Paths}, nil
	case *structpb.Value:
		return attributevalue.MarshalWithOptions(xt.AsInterface())
	case *wrapperspb.StringValue:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.MarshalWithOptions(xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.MarshalWithOptions(xt.Value)
	default:
		return nil, fmt.Errorf("marshal of message type unsupported: %+T", xt)
	}
}

// UnmarshalDynamoMessage will attempt to unmarshal 'm' into a protobuf message 'x'. It provides special
// support for several well-known protobuf message types. If 'x' implements the MarshalDynamoItem method
// it will be called to delegate the unmarshalling.
func UnmarshalDynamoMessage(m types.AttributeValue, x proto.Message, opts ...DecodingOption) (err error) {
	if mx, ok := x.(interface {
		UnmarshalDynamoItem(map[string]types.AttributeValue, ...DecodingOption) error
	}); ok {
		mm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal: no map attribute provided")
		}
		return mx.UnmarshalDynamoItem(mm.Value, opts...)
	}

	switch xt := x.(type) {
	case *durationpb.Duration, *timestamppb.Timestamp:
		ms, ok := m.(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no string attribute provided")
		}
		return protojson.Unmarshal([]byte(strconv.Quote(ms.Value)), x)
	case *anypb.Any:
		mm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no map attribute provided")
		}
		err = attributevalue.UnmarshalWithOptions(mm.Value["1"], &xt.TypeUrl)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's TypeURL field: %w", err)
		}
		err = attributevalue.UnmarshalWithOptions(mm.Value["2"], &xt.Value)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's Value field: %w", err)
		}
		return nil
	case *fieldmaskpb.FieldMask:
		ss, ok := m.(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no string set attribute provided")
		}
		xt.Paths = ss.Value
		return nil
	case *structpb.Value:
		switch m.(type) {
		case *types.AttributeValueMemberL:
			vx := []any{}
			err = attributevalue.UnmarshalWithOptions(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			lv, err := structpb.NewList(vx)
			if err != nil {
				return fmt.Errorf("failed to init structpb.Value: %w", err)
			}
			xt.Kind = &structpb.Value_ListValue{ListValue: lv}
			return nil
		case *types.AttributeValueMemberM:
			vx := map[string]any{}
			err = attributevalue.UnmarshalWithOptions(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			lv, err := structpb.NewStruct(vx)
			if err != nil {
				return fmt.Errorf("failed to init structpb.Value: %w", err)
			}
			xt.Kind = &structpb.Value_StructValue{StructValue: lv}
			return nil
		case *types.AttributeValueMemberS:
			var vx string
			err = attributevalue.UnmarshalWithOptions(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			xt.Kind = &structpb.Value_StringValue{StringValue: vx}
			return nil
		case *types.AttributeValueMemberBOOL:
			var vx bool
			err = attributevalue.UnmarshalWithOptions(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			xt.Kind = &structpb.Value_BoolValue{BoolValue: vx}
			return nil
		case *types.AttributeValueMemberN:
			var vx float64
			err = attributevalue.UnmarshalWithOptions(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			xt.Kind = &structpb.Value_NumberValue{NumberValue: vx}
			return nil
		case *types.AttributeValueMemberNULL:
			xt.Kind = &structpb.Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}
			return nil
		default:
			return fmt.Errorf("failed to unmarshal struct value: unsupported attribute value")
		}
	// wrapper types can just call the sdk unmarshal on the wrapped value
	case *wrapperspb.StringValue:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.UnmarshalWithOptions(m, &xt.Value)
	default:
		return fmt.Errorf("unmarshal of message type unsupported: %+T", xt)
	}
}
