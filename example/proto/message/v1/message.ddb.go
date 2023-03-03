package messagev1

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

// MarshalDynamoItem encodes all fields of Engine as a DynamoDB attribute map.
func (x *Engine) MarshalDynamoItem(opts ...func(*attributevalue.EncoderOptions)) (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue, 2)

	m["1"], err = attributevalue.NewEncoder(opts...).Encode(x.Brand)
	if err != nil {
		return nil, fmt.Errorf("failed to encode field 'Brand': %w", err)
	}

	m["2"], err = attributevalue.NewEncoder(opts...).Encode(x.Dirtyness)
	if err != nil {
		return nil, fmt.Errorf("failed to encode field 'Dirtyness': %w", err)
	}

	return
}

// UnmarshalDynamoItem decodes all DynamoDB attributes into Engine
func (x *Engine) UnmarshalDynamoItem(m map[string]types.AttributeValue, opts ...func(*attributevalue.DecoderOptions)) (err error) {

	if err = attributevalue.NewDecoder(opts...).Decode(m["1"], &x.Brand); err != nil {
		return fmt.Errorf("failed to decode into field 'Brand': %w", err)
	}

	if err = attributevalue.NewDecoder(opts...).Decode(m["2"], &x.Dirtyness); err != nil {
		return fmt.Errorf("failed to decode into field 'Dirtyness': %w", err)
	}

	return nil
}

// MarshalDynamoItem encodes all fields of Car as a DynamoDB attribute map.
func (x *Car) MarshalDynamoItem(opts ...func(*attributevalue.EncoderOptions)) (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue, 2)

	mEngine, err := x.Engine.MarshalDynamoItem(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message field 'Engine': %w", err)
	}
	m["1"] = &types.AttributeValueMemberM{Value: mEngine}

	m["2"], err = attributevalue.NewEncoder(opts...).Encode(x.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to encode field 'Name': %w", err)
	}

	return
}

// UnmarshalDynamoItem decodes all DynamoDB attributes into Car
func (x *Car) UnmarshalDynamoItem(m map[string]types.AttributeValue, opts ...func(*attributevalue.DecoderOptions)) (err error) {

	mEngine, ok := m["1"].(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("no attribute for '1', or not of 'M' type")
	}
	x.Engine = new(Engine)
	x.Engine.UnmarshalDynamoItem(mEngine.Value, opts...)

	if err = attributevalue.NewDecoder(opts...).Decode(m["2"], &x.Name); err != nil {
		return fmt.Errorf("failed to decode into field 'Name': %w", err)
	}

	return nil
}

// MarshalDynamoItem encodes all fields of Kitchen as a DynamoDB attribute map.
func (x *Kitchen) MarshalDynamoItem(opts ...func(*attributevalue.EncoderOptions)) (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue, 1)

	m["1"], err = attributevalue.NewEncoder(opts...).Encode(x.Brand)
	if err != nil {
		return nil, fmt.Errorf("failed to encode field 'Brand': %w", err)
	}

	return
}

// UnmarshalDynamoItem decodes all DynamoDB attributes into Kitchen
func (x *Kitchen) UnmarshalDynamoItem(m map[string]types.AttributeValue, opts ...func(*attributevalue.DecoderOptions)) (err error) {

	if err = attributevalue.NewDecoder(opts...).Decode(m["1"], &x.Brand); err != nil {
		return fmt.Errorf("failed to decode into field 'Brand': %w", err)
	}

	return nil
}
