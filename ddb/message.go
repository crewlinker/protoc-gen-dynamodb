package ddb

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// MarshalMessage will marshal a protobuf message 'm' into an attribute value. It supports several
// well-known Protobuf types and if 'x' implements its own MarshalDynamoItem method it will be called to
// delegate the marshalling.
func MarshalMessage(x proto.Message, os ...Option) (a types.AttributeValue, err error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		return jsonMarshal(x)
	case ddbv1.Encoding_ENCODING_DYNAMO:
	default:
		return nil, errEmbedEncoding()
	}

	// check if the message implements its own marshalling, if so defer to that
	if mx, ok := x.(interface {
		MarshalDynamoItem() (map[string]types.AttributeValue, error)
	}); ok {
		mm, err := mx.MarshalDynamoItem()
		return &types.AttributeValueMemberM{Value: mm}, err
	}

	// else, check for some special well-known types and handle these cases specifically
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
		mv.Value["1"], err = attributevalue.Marshal(xt.TypeUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's TypeURL field: %w", err)
		}
		mv.Value["2"], err = attributevalue.Marshal(xt.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's Value field: %w", err)
		}
		return mv, nil
	case *fieldmaskpb.FieldMask:
		// we encode the fieldmask as a map with a single key. Because if we encode it as a set directly
		// it causes trouble when building paths. As repeated fieldmask fields would require list of list
		// indexing. Which is blocked by: https://github.com/crewlinker/protoc-gen-dynamodb/issues/45
		return &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberSS{Value: xt.Paths},
		}}, nil
	case *structpb.Value:
		return attributevalue.Marshal(xt.AsInterface())
	case *wrapperspb.StringValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.Marshal(xt.Value)
	default:
		return nil, fmt.Errorf("marshal of message type unsupported: %+T", xt)
	}
}

// UnmarshalMessage will attempt to unmarshal 'm' into a protobuf message 'x'. It provides special
// support for several well-known protobuf message types. If 'x' implements the MarshalDynamoItem method
// it will be called to delegate the unmarshalling.
func UnmarshalMessage(m types.AttributeValue, x proto.Message, os ...Option) (err error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		return jsonUnmarshal(m, x)
	case ddbv1.Encoding_ENCODING_DYNAMO:
	default:
		return errEmbedEncoding()
	}

	if mx, ok := x.(interface {
		UnmarshalDynamoItem(map[string]types.AttributeValue) error
	}); ok {
		mm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal: no map attribute provided")
		}
		return mx.UnmarshalDynamoItem(mm.Value)
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
		err = attributevalue.Unmarshal(mm.Value["1"], &xt.TypeUrl)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's TypeURL field: %w", err)
		}
		err = attributevalue.Unmarshal(mm.Value["2"], &xt.Value)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's Value field: %w", err)
		}
		return nil
	case *fieldmaskpb.FieldMask:
		fmm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no map attribute provided")
		}

		ss, ok := fmm.Value["1"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no string set attribute provided")
		}
		xt.Paths = ss.Value
		return nil
	case *structpb.Value:
		switch m.(type) {
		case *types.AttributeValueMemberL:
			vx := []any{}
			err = attributevalue.Unmarshal(m, &vx)
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
			err = attributevalue.Unmarshal(m, &vx)
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
			err = attributevalue.Unmarshal(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			xt.Kind = &structpb.Value_StringValue{StringValue: vx}
			return nil
		case *types.AttributeValueMemberBOOL:
			var vx bool
			err = attributevalue.Unmarshal(m, &vx)
			if err != nil {
				return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
			}
			xt.Kind = &structpb.Value_BoolValue{BoolValue: vx}
			return nil
		case *types.AttributeValueMemberN:
			var vx float64
			err = attributevalue.Unmarshal(m, &vx)
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
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	default:
		return fmt.Errorf("unmarshal of message type unsupported: %+T", xt)
	}
}
