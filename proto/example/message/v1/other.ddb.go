// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

package messagev1

import (
	"fmt"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddb "github.com/crewlinker/protoc-gen-dynamodb/ddb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *OtherKitchen) MarshalDynamoItem(o ...ddb.EncodingOption) (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.AnotherKitchen != nil {
		m16, err := ddb.MarshalDynamoMessage(x.GetAnotherKitchen(), o...)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'AnotherKitchen': %w", err)
		}
		m["16"] = m16
	}
	if x.OtherTimer != nil {
		m17, err := ddb.MarshalDynamoMessage(x.GetOtherTimer(), o...)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OtherTimer': %w", err)
		}
		m["17"] = m17
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *OtherKitchen) UnmarshalDynamoItem(m map[string]types.AttributeValue, o ...ddb.DecodingOption) (err error) {
	if m["16"] != nil {
		x.AnotherKitchen = new(Kitchen)
		err = ddb.UnmarshalDynamoMessage(m["16"], x.AnotherKitchen, o...)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'AnotherKitchen': %w", err)
		}
	}
	if m["17"] != nil {
		x.OtherTimer = new(durationpb.Duration)
		err = ddb.UnmarshalDynamoMessage(m["17"], x.OtherTimer, o...)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OtherTimer': %w", err)
		}
	}
	return nil
}
