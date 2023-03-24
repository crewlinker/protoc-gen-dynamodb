package ddb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// UnmarshalMappedMessage decodes the dynamodb representation of a map of messages
func UnmarshalMappedMessage[K comparable, T any, TP ProtoMessage[T]](m types.AttributeValue, fv func(s string) (K, error), os ...Option) (xm map[K]TP, err error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		var outer map[string]json.RawMessage
		if err := jsonUnmarshal(m, &outer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal outer map: %w", err)
		}
		xm = make(map[K]TP)
		for k, b := range outer {
			kv, err := fv(k)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal map key: %w", err)
			}
			var mv TP = new(T)
			if err := protojson.Unmarshal(b, mv); err != nil {
				return nil, fmt.Errorf("failed to unmarshal mapped message '%s': %w", k, err)
			}
			xm[kv] = mv
		}
		return
	case ddbv1.Encoding_ENCODING_DYNAMO:
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
			if err = UnmarshalMessage(v, mv, os...); err != nil {
				return nil, fmt.Errorf("failed to unmarshal message map value: %w", err)
			}
			xm[kv] = mv
		}
		return
	default:
		return nil, errEmbedEncoding()
	}
}

// MarshalMappedMessage takes a map of messages and marshals it to a dynamodb representation
func MarshalMappedMessage[K comparable, T any, TP ProtoMessage[T]](x map[K]TP, os ...Option) (types.AttributeValue, error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		outer := make(map[string]json.RawMessage, len(x))
		for k, v := range x {
			kv, err := marshalMapKey(k)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map key: %w", err)
			}
			if outer[kv], err = protojson.Marshal(v); err != nil {
				return nil, fmt.Errorf("failed to marshal mapped message '%s': %w", kv, err)
			}
		}
		return jsonMarshal(outer)
	case ddbv1.Encoding_ENCODING_DYNAMO:
		m := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x {
			kv, err := marshalMapKey(k)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map key: %w", err)
			}
			if v == nil {
				m.Value[kv] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := MarshalMessage(v, os...)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal mapped message: %w", err)
			}
			m.Value[kv] = mv
		}
		return m, nil
	default:
		return nil, errEmbedEncoding()
	}
}
