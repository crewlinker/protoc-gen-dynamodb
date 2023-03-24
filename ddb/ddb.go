// Package ddb provides DynamoDB utility for Protobuf messages
package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"
	"google.golang.org/protobuf/proto"
)

// ProtoMessage is a constraint to a protobuf message pointer.
type ProtoMessage[T any] interface {
	proto.Message
	*T
}

// Marshal will marshal basic types, and composite types that only hold basic types. It defers to the
// offical AWS sdk but is still put here to make it easier to change behaviour in the future.
func Marshal(in any, os ...Option) (types.AttributeValue, error) {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		return jsonMarshal(in)
	case ddbv1.Encoding_ENCODING_DYNAMO:
		return attributevalue.Marshal(in)
	default:
		return nil, errEmbedEncoding()
	}
}

// Unmarshal will marshal basic types, and composite types that only hold basic types. It takes into
// account the embed encoding option.
func Unmarshal(av types.AttributeValue, out any, os ...Option) error {
	opts := applyOptions(os...)
	switch opts.embedEncoding {
	case ddbv1.Encoding_ENCODING_JSON:
		return jsonUnmarshal(av, out)
	case ddbv1.Encoding_ENCODING_DYNAMO:
		return attributevalue.Unmarshal(av, out)
	default:
		return errEmbedEncoding()
	}
}

var (
	// ErrUnsupportedEmbedEncoding is returned when a unsupported embed encoding is used
	ErrUnsupportedEmbedEncoding = fmt.Errorf("unsupported embed encoding, supports: %s %s",
		ddbv1.Encoding_ENCODING_JSON, ddbv1.Encoding_ENCODING_DYNAMO)
)

// errEmbedEncoding returns an error that forces comparing with errors.Is instead of "=="
func errEmbedEncoding() error {
	return fmt.Errorf("%w", ErrUnsupportedEmbedEncoding)
}
