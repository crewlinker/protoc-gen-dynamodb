package ddbv1

import (
	"fmt"
	attributevalue "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"
)

// file_ddb_v1_options_proto_marshal_dynamo_item marshals into DynamoDB attribute value maps
func file_ddb_v1_options_proto_marshal_dynamo_item(x proto.Message) (a types.AttributeValue, err error) {
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

// file_ddb_v1_options_proto_marshal_dynamo_item unmarshals DynamoDB attribute value maps
func file_ddb_v1_options_proto_unmarshal_dynamo_item(m types.AttributeValue, x proto.Message) (err error) {
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
		var vv any
		switch m.(type) {
		case *types.AttributeValueMemberL:
			vx := []any{}
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberM:
			vx := map[string]any{}
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberS:
			var vx string
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberBOOL:
			var vx bool
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberN:
			var vx float64
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberNULL:
			sv, _ := structpb.NewValue(nil)
			*xt = *sv
			return nil
		default:
			return fmt.Errorf("failed to unmarshal struct value: unsupported attribute value")
		}
		if err != nil {
			return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
		}
		sv, err := structpb.NewValue(vv)
		if err != nil {
			return fmt.Errorf("failed to init structpb value: %w", err)
		}
		*xt = *sv
		return nil
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

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *FieldOptions) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Name != nil {
		m["1"], err = attributevalue.Marshal(x.GetName())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Name': %w", err)
		}
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *FieldOptions) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["1"], &x.Name)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Name': %w", err)
	}
	return nil
}
