package ddb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
)

type Mask struct {
	paths map[string]struct{}
}

func NewMask(m proto.Message, p []string) (*Mask, error) {
	// @TODO error if duplicate
	// @TODO validate against message

	return &Mask{}, nil
}

// NameValue is returned from a masked item to facilitate conditions,updates and projections
type NameValue struct {
	Name  expression.NameBuilder
	Value expression.ValueBuilder
}

// NameValues return a flatted sllice of expression names and values. It returns only names and values
// that are defined in the mask and exist in the attribute map.
func (m *Mask) NameValues(it map[string]types.AttributeValue) []NameValue {
	return nil
}

// Names returns a slice of expression names. This is usefull for building projections
func (m *Mask) Names() []expression.NameBuilder {
	return nil
}
