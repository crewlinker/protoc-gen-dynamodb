package messagev1

import (
	"fmt"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func file_message_v1_other_proto_marshal_dynamo_item(x any) (map[string]types.AttributeValue, error) {
	if mx, ok := x.(interface {
		MarshalDynamoItem() (map[string]types.AttributeValue, error)
	}); ok {
		return mx.MarshalDynamoItem()
	}
	return nil, nil
}
func file_message_v1_other_proto_unmarshal_dynamo_item(m map[string]types.AttributeValue, x any) error {
	if mx, ok := x.(interface {
		UnmarshalDynamoItem(map[string]types.AttributeValue) error
	}); ok {
		return mx.UnmarshalDynamoItem(m)
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *OtherKitchen) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.AnotherKitchen != nil {
		m16, err := file_message_v1_other_proto_marshal_dynamo_item(x.AnotherKitchen)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'AnotherKitchen': %w", err)
		}
		m["16"] = &types.AttributeValueMemberM{Value: m16}
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *OtherKitchen) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	if m["16"] != nil {
		m16, ok := m["16"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'AnotherKitchen': no map attribute provided")
		}
		x.AnotherKitchen = new(Kitchen)
		err = file_message_v1_other_proto_unmarshal_dynamo_item(m16.Value, x.AnotherKitchen)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'AnotherKitchen': %w", err)
		}
	}
	return nil
}
