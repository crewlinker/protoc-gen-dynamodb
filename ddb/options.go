package ddb

import (
	"strings"
)

// EncodingOptions configure the marshalling process
type EncodingOptions struct {
	mask map[string][]string
}

// IsMasked returns true if the provided field is selected at the front of each mask entry.
func (o EncodingOptions) IsMasked(attrname string) bool {
	if len(o.mask) < 1 {
		return true
	}
	for _, p := range o.mask {
		if p[0] == attrname {
			return true
		}
	}
	return false
}

// SubMask will remove the head of each mask in the encoding option if it matches the attrname. If it
// does not match the attrname it will be removed from the mask. This is usefull during marshalling to
// recurse into a subfield and only pass along the masks that are relevant.
func (o EncodingOptions) SubMask(attrname string) EncodingOptions {
	if !o.IsMasked(attrname) {
		o.mask = nil
		return o
	}

	newm := map[string][]string{}
	for _, p := range o.mask {
		if p[0] != attrname {
			continue
		}
		pp := p[1:]
		if len(pp) > 0 {
			newm[strings.Join(pp, ".")] = pp
		}
	}

	o.mask = newm
	return o
}

// EncodingOption types an option that configures the marshalling
type EncodingOption func(*EncodingOptions)

// WithEncodingOptions is an option that provides all encoding options directly
func WithEncodingOptions(v EncodingOptions) EncodingOption {
	return func(eo *EncodingOptions) {
		*eo = v
	}
}

// WithMask configures the marshalling process to ONLY encode fields that match the mask. This is usefull for
// example to only marshal the key attributes, or to generate expression values for updates, etc. Duplicate
// entries en empty strings will be ignored.
func WithMask(v ...string) EncodingOption {
	return func(eo *EncodingOptions) {
		eo.mask = map[string][]string{}
		for _, vv := range v {
			if vv == "" {
				continue
			}

			eo.mask[vv] = strings.Split(vv, ".")
		}
	}
}

// ApplyEncodingOptions merges the options into the options struct
func ApplyEncodingOptions(os ...EncodingOption) (opts EncodingOptions) {
	for _, o := range os {
		o(&opts)
	}
	return
}

// DecodingOptions configure the unmarshalling process
type DecodingOptions struct{}

// DecodingOption types an option that configures the marshalling
type DecodingOption func(*DecodingOptions)

// ApplyDecodingOptions merges the options into the options struct
func ApplyDecodingOptions(os ...DecodingOption) (opts DecodingOptions) {
	for _, o := range os {
		o(&opts)
	}
	return
}
