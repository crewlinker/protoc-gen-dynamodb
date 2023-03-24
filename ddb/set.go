package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
)

// MarshalSet will marshal a slice of 'T' to a dynamo set.
func MarshalSet[T ~uint64 | ~uint32 | ~int32 | ~int64 | string | []byte](s []T, os ...Option) (types.AttributeValue, error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		return jsonMarshal(s)
	case ddbv1.Encoding_ENCODING_DYNAMO:
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
	default:
		return nil, errEmbedEncoding()
	}
}
