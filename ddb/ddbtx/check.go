package ddbtx

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
)

// CheckOption allows configuring check calls
type CheckOption func(*checkOpts)

// CheckIfNotExists adds a condition check that checks if the item exists
func CheckIfNotExists() CheckOption {
	return func(uo *checkOpts) {
		uo.ifNotExists = true
	}
}

// CheckIfExists adds a condition check that checks if the item exists
func CheckIfExists() CheckOption {
	return func(uo *checkOpts) {
		uo.ifExists = true
	}
}

// checkOptions holds optional configurations for check
type checkOpts struct {
	ifNotExists bool
	ifExists    bool
}

// applyCheckOptions applies check option defaults and overwites whatever is configured
func applyCheckOptions(uos ...CheckOption) (uo checkOpts) {
	for _, uof := range uos {
		uof(&uo)
	}
	return
}

// Check adds a check item operation to the tx
func (wtx WriteTx) Check(it Item, os ...CheckOption) WriteTx {
	var cch types.ConditionCheck
	var exprb expression.Builder

	opts := applyCheckOptions(os...)
	avm, err := it.MarshalDynamoItem()
	if err != nil {
		return wtx.errorf("failed to marshal: %w", err)
	}

	cch.Key, err = ddbpath.SelectMapValues(avm, it.DynamoKeyNames()...)
	if err != nil {
		return wtx.errorf("failed to select key attribute values: %w", err)
	}

	if opts.ifExists {
		exprb = withShouldExist(exprb, it)
	} else if opts.ifNotExists {
		exprb = withShouldNotExist(exprb, it)
	}

	expr, err := exprb.Build()
	if err != nil {
		return wtx.errorf("failed to build expression: %w", err)
	}

	cch.ExpressionAttributeNames = expr.Names()
	cch.ExpressionAttributeValues = expr.Values()
	cch.ConditionExpression = expr.Condition()

	wtx.its = append(wtx.its, types.TransactWriteItem{ConditionCheck: &cch})
	return wtx
}
