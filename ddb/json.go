package ddb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// jsonUnmarshal unmarshals 'av' into 'out'. In case 'out' is a proto.Message it will use
// protojson encoding.
func jsonUnmarshal(av types.AttributeValue, out any) (err error) {
	if av == nil {
		return nil // nothing to decode
	}

	sav, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		return fmt.Errorf("expected json encoded embed in S attribute value, got: %T", av)
	}

	switch out := out.(type) {
	case proto.Message:
		err = protojson.Unmarshal([]byte(sav.Value), out)
	default:
		err = json.Unmarshal([]byte(sav.Value), out)
	}

	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

// jsonMarshal marshals 'in' to a dynamo S attribute. In case 'in' is a proto.Message it will
// use protojson encoding, else it will use the stdlib json encoding.
func jsonMarshal(in any) (av types.AttributeValue, err error) {
	var b []byte
	switch in := in.(type) {
	case proto.Message:
		b, err = protojson.Marshal(in)
	default:
		b, err = json.Marshal(in)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to json marshal: %w", err)
	}

	return &types.AttributeValueMemberS{Value: string(b)}, nil
}
