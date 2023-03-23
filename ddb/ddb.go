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

// ProtoMessage is a constraint to a protobuf message pointer.
type ProtoMessage[T any] interface {
	proto.Message
	*T
}

// UintMapKey parses 's' as an unsigned integer value
func UintMapKey[K ~uint32 | ~uint64](s string) (K, error) {
	k, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return K(k), nil
}

// IntMapKey parses 's' as a signed integer value
func IntMapKey[K ~int32 | ~int64](s string) (K, error) {
	k, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return K(k), nil
}

// BoolMapKey parses 's' as a boolean 'true' or 'false' value
func BoolMapKey(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool key: %v", s)
	}
}

// StringMapKey parses 's' as a string map key
func StringMapKey(s string) (string, error) {
	return s, nil
}

// UnmarshalMappedMessage decodes the dynamodb representation of a map of messages
func UnmarshalMappedMessage[K comparable, T any, TP ProtoMessage[T]](m types.AttributeValue, fv func(s string) (K, error)) (xm map[K]TP, err error) {
	xm = make(map[K]TP)
	mm, ok := m.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshal mapped field: no map attribute provided")
	}
	for k, v := range mm.Value {
		kv, err := fv(k)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal map key: %w", err)
		}
		if _, ok := v.(*types.AttributeValueMemberNULL); ok {
			xm[kv] = nil // set explicit nil
			continue
		}
		var mv TP = new(T)
		if err = UnmarshalMessage(v, mv); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message map value: %w", err)
		}
		xm[kv] = mv
	}
	return
}

// MarshalMappedMessage takes a map of messages and marshals it to a dynamodb representation
func MarshalMappedMessage[K comparable, T any, TP ProtoMessage[T]](x map[K]TP) (types.AttributeValue, error) {
	m := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
	for k, v := range x {
		var kv string
		switch kt := any(k).(type) {
		case string:
			kv = kt
		case bool:
			if kt {
				kv = "true"
			} else {
				kv = "false"
			}
		case int32, int64, uint32, uint64:
			kv = fmt.Sprintf("%d", kt)
		default:
			return nil, fmt.Errorf("unsupported map key type: %T", k)
		}
		if kv == "" {
			return nil, fmt.Errorf("failed to marshal map key: map key cannot be empty")
		}
		if v == nil {
			m.Value[kv] = &types.AttributeValueMemberNULL{Value: true}
			continue
		}
		mv, err := MarshalMessage(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal mapped message: %w", err)
		}
		m.Value[kv] = mv
	}
	return m, nil
}

// UnmarshalRepeatedMessage provides a generic function for unmarshalling a repeated field of messages from
// the DynamoDB representation.
func UnmarshalRepeatedMessage[T any, TP ProtoMessage[T]](m types.AttributeValue) (xl []TP, err error) {
	ml, ok := m.(*types.AttributeValueMemberL)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshal repeated field: dynamo value is not a list")
	}

	for i, v := range ml.Value {
		if _, ok := v.(*types.AttributeValueMemberNULL); ok {
			xl = append(xl, nil) // append explicit nil
			continue
		}

		var mv TP = new(T)
		if err = UnmarshalMessage(v, mv); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message item '%d' of field: %w", i, err)
		}
		xl = append(xl, mv)
	}
	return
}

// MarshalSet will marshal a slice of 'T' to a dynamo set.
func MarshalSet[T ~uint64 | ~uint32 | ~int32 | ~int64 | string | []byte](s []T) (types.AttributeValue, error) {
	switch st := any(s).(type) {
	case []string:
		a := &types.AttributeValueMemberSS{}
		a.Value = append(a.Value, st...)
		return a, nil
	case [][]byte:
		a := &types.AttributeValueMemberBS{}
		a.Value = append(a.Value, st...)
		return a, nil
	case []uint64, []uint32, []int32, []int64:
		a := &types.AttributeValueMemberNS{}
		for _, v := range s {
			av, err := attributevalue.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal numeric set item: %w", err)
			}

			avn, ok := av.(*types.AttributeValueMemberN)
			if !ok {
				return nil, fmt.Errorf("expected N member encoding for numeric set item, got: %T", av)
			}
			a.Value = append(a.Value, avn.Value)
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unsupported set item encoding: %T", st)
	}
}

// MarshalRepeatedMessage provides a generic function for marshalling a repeated field as long as the
// generated code provides the concrete type as the Type parameter.
func MarshalRepeatedMessage[T any, TP ProtoMessage[T]](x []TP) (types.AttributeValue, error) {
	a := &types.AttributeValueMemberL{}
	for i, m := range x {
		if m == nil {
			a.Value = append(a.Value, &types.AttributeValueMemberNULL{Value: true})
			continue
		}
		v, err := MarshalMessage(m)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item '%d' of repeated message field': %w", i, err)
		}
		a.Value = append(a.Value, v)
	}
	return a, nil
}

// MarshalMessage will marshal a protobuf message 'm' into an attribute value. It supports several
// well-known Protobuf types and if 'x' implements its own MarshalDynamoItem method it will be called to
// delegate the marshalling.
func MarshalMessage(x proto.Message) (a types.AttributeValue, err error) {
	if mx, ok := x.(interface {
		MarshalDynamoItem() (map[string]types.AttributeValue, error)
	}); ok {
		mm, err := mx.MarshalDynamoItem()
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
		return &types.AttributeValueMemberSS{Value: xt.Paths}, nil
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
func UnmarshalMessage(m types.AttributeValue, x proto.Message) (err error) {
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

// Marshal will marshal basic types, and composite types that only hold basic types. It defers to the
// offical AWS sdk but is still put here to make it easier to change behaviour in the future.
func Marshal(in any) (types.AttributeValue, error) {
	return attributevalue.Marshal(in)
}

// Unmarshal will marshal basic types, and composite types that only hold basic types. It defers to the
// offical AWS sdk but is still put here to make it easier to change behaviour in the future.
func Unmarshal(av types.AttributeValue, out any) error {
	return attributevalue.Unmarshal(av, out)
}
