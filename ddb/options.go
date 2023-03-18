package ddb

type encodingOpts struct{}

// EncodingOption types an option that configures the marshalling
type EncodingOption func(encodingOpts)

type decodingOpts struct{}

// DecodingOption types an option that configures the marshalling
type DecodingOption func(decodingOpts)
