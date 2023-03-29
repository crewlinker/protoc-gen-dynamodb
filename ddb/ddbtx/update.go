package ddbtx

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
)

// UpdateOption allows configuring update calls
type UpdateOption func(*updateOpts)

// UpdateMask allows updates to only set part of the item
func UpdateMask(v ...string) UpdateOption {
	return func(uo *updateOpts) { uo.mask = v }
}

// UpdateIfExists adds a condition check that checks if the item exists
func UpdateIfExists() UpdateOption {
	return func(uo *updateOpts) {
		uo.ifExists = true
	}
}

// Update adds an UpdateItem to the transaction. Edits an existing item's attributes, or adds a new item
// to the table if it does not already exist.
func (wtx WriteTx) Update(it Item, os ...UpdateOption) WriteTx {
	var upd types.Update
	var exprb expression.Builder

	opts := applyUpdateOptions(os...)
	avm, err := it.MarshalDynamoItem()
	if err != nil {
		return wtx.errorf("failed to marshal: %w", err)
	}

	upd.Key, err = ddbpath.SelectMapValues(avm, it.DynamoKeyNames()...)
	if err != nil {
		return wtx.errorf("failed to select key attribute values: %w", err)
	}

	// if the mask is set it will ensure partial updates of selected values
	if len(opts.mask) > 0 {
		// @TODO should validation the mask

		mavm, err := ddbpath.SelectMapValues(avm, opts.mask...)
		if err != nil {
			return wtx.errorf("failed to select mask attribute values: %w", err)
		}

		var sets expression.UpdateBuilder
		for name, val := range mavm {
			// @TODO could add an option to "REMOVE", "DEL", or "ADD" the masked values
			// @TODO could maybe even make it part of the path string (improved lexer)
			sets.Set(expression.Name(name), expression.Value(val))
		}

		exprb.WithUpdate(sets)
	}

	// if we only wanna update if it exists, add the condition
	if opts.ifExists {
		exprb = withShouldExist(exprb, it)
	}

	expr, err := exprb.Build()
	if err != nil {
		return wtx.errorf("failed to build expression: %w", err)
	}

	upd.ExpressionAttributeNames = expr.Names()
	upd.ExpressionAttributeValues = expr.Values()
	upd.UpdateExpression = expr.Update()
	upd.ConditionExpression = expr.Condition()

	wtx.its = append(wtx.its, types.TransactWriteItem{Update: &upd})
	return wtx
}

// updateOptions holds optional configurations for update
type updateOpts struct {
	mask     []string
	ifExists bool
}

// applyUpdateOptions applies update option defaults and overwites whatever is configured
func applyUpdateOptions(uos ...UpdateOption) (uo updateOpts) {
	for _, uof := range uos {
		uof(&uo)
	}
	return
}
