package ddb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// UnmarshalRepeatedMessage provides a generic function for unmarshalling a repeated field of messages from
// the DynamoDB representation.
func UnmarshalRepeatedMessage[T any, TP ProtoMessage[T]](m types.AttributeValue, os ...Option) (xl []TP, err error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		var outer []json.RawMessage
		if err := jsonUnmarshal(m, &outer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal outer slice: %w", err)
		}
		for i, b := range outer {
			var mv TP = new(T)
			if err := protojson.Unmarshal(b, mv); err != nil {
				return nil, fmt.Errorf("failed to unmarshal message item '%d': %w", i, err)
			}
			xl = append(xl, mv)
		}
		return
	case ddbv1.Encoding_ENCODING_DYNAMO:
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
			if err = UnmarshalMessage(v, mv, os...); err != nil {
				return nil, fmt.Errorf("failed to unmarshal message item '%d' of field: %w", i, err)
			}
			xl = append(xl, mv)
		}
		return
	default:
		return nil, errEmbedEncoding()
	}
}

// MarshalRepeatedMessage provides a generic function for marshalling a repeated field as long as the
// generated code provides the concrete type as the Type parameter.
func MarshalRepeatedMessage[T any, TP ProtoMessage[T]](x []TP, os ...Option) (av types.AttributeValue, err error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		outer := make([]json.RawMessage, len(x))
		for i, m := range x {
			if outer[i], err = protojson.Marshal(m); err != nil {
				return nil, fmt.Errorf("failed to marshal repeated message '%d': %w", i, err)
			}
		}
		return jsonMarshal(outer)
	case ddbv1.Encoding_ENCODING_DYNAMO:
		a := &types.AttributeValueMemberL{}
		for i, m := range x {
			if m == nil {
				a.Value = append(a.Value, &types.AttributeValueMemberNULL{Value: true})
				continue
			}
			v, err := MarshalMessage(m, os...)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal item '%d' of repeated message field': %w", i, err)
			}
			a.Value = append(a.Value, v)
		}
		return a, nil
	default:
		return nil, errEmbedEncoding()
	}
}
