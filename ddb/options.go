package ddb

import ddbv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/ddb/v1"

// opts holds the options
type opts struct {
	embedEncoding ddbv1.Encoding
}

// applyOptions merges the options together into a single struct
func applyOptions(os ...Option) (o opts) {
	o.embedEncoding = ddbv1.Encoding_ENCODING_DYNAMO
	for _, f := range os {
		f(&o)
	}
	return
}

// Option configures the shared logic
type Option func(*opts)

// Embed option will signal to the marshalling/unmarshalling logic that the field
// id embedded in the Dynamo item and should be decoded.
func Embed(v ddbv1.Encoding) Option {
	return func(o *opts) {
		o.embedEncoding = v
	}
}
