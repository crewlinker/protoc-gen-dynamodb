package messagev1

import (
	"fmt"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

// file_message_v1_other_proto_marshal_dynamo_item marshals into DynamoDB attribute value maps
func file_message_v1_other_proto_marshal_dynamo_item(x proto.Message) (types.AttributeValue, error) {
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
			return nil, fmt.Errorf("failed to unquote marshalled duration: %w", err)
		}
		return &types.AttributeValueMemberS{Value: xjsons}, nil
	default:
		return nil, fmt.Errorf("marshal of message type unsupported: %+T", xt)
	}
}

// file_message_v1_other_proto_marshal_dynamo_item unmarshals DynamoDB attribute value maps
func file_message_v1_other_proto_unmarshal_dynamo_item(m types.AttributeValue, x proto.Message) error {
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
	default:
		return fmt.Errorf("unmarshal of message type unsupported: %+T", xt)
	}
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *OtherKitchen) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.AnotherKitchen != nil {
		m16, err := file_message_v1_other_proto_marshal_dynamo_item(x.AnotherKitchen)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'AnotherKitchen': %w", err)
		}
		m["16"] = m16
	}
	if x.OtherTimer != nil {
		m17, err := file_message_v1_other_proto_marshal_dynamo_item(x.OtherTimer)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OtherTimer': %w", err)
		}
		m["17"] = m17
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *OtherKitchen) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	if m["16"] != nil {
		x.AnotherKitchen = new(Kitchen)
		err = file_message_v1_other_proto_unmarshal_dynamo_item(m["16"], x.AnotherKitchen)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'AnotherKitchen': %w", err)
		}
	}
	if m["17"] != nil {
		x.OtherTimer = new(durationpb.Duration)
		err = file_message_v1_other_proto_unmarshal_dynamo_item(m["17"], x.OtherTimer)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OtherTimer': %w", err)
		}
	}
	return nil
}
