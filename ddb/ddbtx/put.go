package ddbtx

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// PutOption allows configuring put calls
type PutOption func(*putOpts)

// PutIfNotExists adds a condition check that checks if the item exists
func PutIfNotExists() PutOption {
	return func(po *putOpts) {
		po.ifNotExists = true
	}
}

// Put adds a PutItem to the transaction
func (wtx WriteTx) Put(it Item, os ...PutOption) WriteTx {
	var put types.Put
	var err error
	var exprb expression.Builder

	opts := applyPutOptions(os...)
	put.Item, err = it.MarshalDynamoItem()
	if err != nil {
		return wtx.errorf("failed to marshal: %w", err)
	}

	if opts.ifNotExists {
		exprb = withShouldNotExist(exprb, it)
	}

	expr, err := exprb.Build()
	if err != nil {
		return wtx.errorf("failed to build expression: %w", err)
	}

	put.ConditionExpression = expr.Condition()
	put.ExpressionAttributeNames = expr.Names()
	put.ExpressionAttributeValues = expr.Values()

	wtx.its = append(wtx.its, types.TransactWriteItem{Put: &put})
	return wtx
}

// putOptions holds optional configurations for put
type putOpts struct {
	ifNotExists bool
}

// applyPutOptions applies put option defaults and overwites whatever is configured
func applyPutOptions(pos ...PutOption) (po putOpts) {
	for _, pof := range pos {
		pof(&po)
	}
	return
}
