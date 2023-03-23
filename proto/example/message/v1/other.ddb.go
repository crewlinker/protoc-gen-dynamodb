// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

package messagev1

import (
	"fmt"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddb "github.com/crewlinker/protoc-gen-dynamodb/ddb"
	v1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// MarshalDynamoItem marshals data into a dynamodb attribute map
func (x *OtherKitchen) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.AnotherKitchen != nil {
		m16, err := ddb.MarshalMessage(x.GetAnotherKitchen(), ddb.Embed(v1.Encoding_ENCODING_UNSPECIFIED))
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'AnotherKitchen': %w", err)
		}
		m["16"] = m16
	}
	if x.OtherTimer != nil {
		m17, err := ddb.MarshalMessage(x.GetOtherTimer(), ddb.Embed(v1.Encoding_ENCODING_UNSPECIFIED))
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
		err = ddb.UnmarshalMessage(m["16"], x.AnotherKitchen, ddb.Embed(v1.Encoding_ENCODING_UNSPECIFIED))
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'AnotherKitchen': %w", err)
		}
	}
	if m["17"] != nil {
		x.OtherTimer = new(durationpb.Duration)
		err = ddb.UnmarshalMessage(m["17"], x.OtherTimer, ddb.Embed(v1.Encoding_ENCODING_UNSPECIFIED))
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OtherTimer': %w", err)
		}
	}
	return nil
}

// OtherKitchenP allows for constructing type-safe expression names
type OtherKitchenP struct {
	ddb.P
}

// Set allows generic list builder to replace the path value
func (p OtherKitchenP) Set(v string) OtherKitchenP {
	p.P = p.P.Set(v)
	return p
}

// OtherKitchenPath starts the building of an expression path into OtherKitchen
func OtherKitchenPath() OtherKitchenP {
	return OtherKitchenP{}
}

// AnotherKitchen returns 'p' with the attribute name appended and allow subselecting nested message
func (p OtherKitchenP) AnotherKitchen() KitchenP {
	return KitchenP{}.Set(p.Val() + ".16")
}

// OtherTimer returns 'p' with the attribute name appended
func (p OtherKitchenP) OtherTimer() ddb.P {
	return (ddb.P{}).Set(p.Val() + ".17")
}
