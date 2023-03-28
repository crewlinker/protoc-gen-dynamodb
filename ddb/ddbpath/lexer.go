// Package ddbpath parses Dynamo paths efficiently
package ddbpath

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// lexer data that changes as it progresses through the input.
type lexer struct {
	pos   int
	width int
	start int
	err   error
	input string
	res   []PathElement
}

// eof rune
var eof rune = -1

// ignore the current rune by advancing the start position
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup the last rune that was seen
func (l *lexer) backup() {
	l.pos -= l.width
}

// next consumes a rune from the input
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// set the error on the lexer and return the nil (end) state
func (l *lexer) errorf(format string, v ...any) stateFn {
	l.err = fmt.Errorf(format, v...)
	return nil
}

// stateFn represents a lexer state. It returns the mutated lexer and a new state.
type stateFn func(l lexer) (lexer, stateFn)

// append the captured field key into the result
func emitField(l lexer) (lexer, stateFn) {
	l.res = append(l.res, PathElement{l.input[l.start:l.pos], -1})
	return l, lexStart
}

// append the captured index to the result
func emitIndex(l lexer) (lexer, stateFn) {
	var part PathElement
	part.Field = strings.TrimSuffix(l.input[l.start:l.pos], "]")
	if len(part.Field) < 1 {
		return l, l.errorf("expected at least 1 digit for index")
	}

	var err error
	part.Index, err = strconv.Atoi(part.Field)
	if err != nil {
		return l, l.errorf("failed to parse '%s' as int: %w", part.Field, err)
	}

	l.res = append(l.res, part)
	return l, lexStart
}

// parse a string key into a map
func lexField(l lexer) (lexer, stateFn) {
	for {
		switch l.next() {
		case eof:
			return l, emitField
		case '.', '[':
			l.backup()
			return l, emitField
		}
	}
}

// parse an index into a list or set
func lexIndex(l lexer) (lexer, stateFn) {
	for {
		r := l.next()
		switch r {
		case eof, '.', '[':
			return l, l.errorf("unexpected end of index, got '%s' expected ']", string(r))
		case ']':
			return l, emitIndex
		}

		if !unicode.IsDigit(r) {
			return l, l.errorf("indexing requires digit, got: '%s'", string(r))
		}
	}
}

// start state
func lexStart(l lexer) (lexer, stateFn) {
	switch l.next() {
	case '.':
		l.ignore()
		return l, lexField
	case '[':
		l.ignore()
		return l, lexIndex
	case eof:
		return l, nil // done
	default:
		return l, l.errorf("expected dot or bracket")
	}
}

// PathElement describes a part of the path. It is either an numeric index into a list (or set)
// or it is a string key into the field.
type PathElement struct {
	Field string
	Index int
}

// ParsePath will parse 'p' in its elements and return them.
func ParsePath(p string) ([]PathElement, error) {
	return AppendParsePath(p, nil)
}

// AppendParsePath will parse path 'p' and append elements to 'r'. If 'r' is already allocated with enough
// space to hold all the parts of 'p' it will not allocate any additional memory on the heap.
func AppendParsePath(p string, r []PathElement) ([]PathElement, error) {
	l := lexer{input: p, res: r}
	for sf := lexStart; sf != nil; {
		l, sf = sf(l)
	}
	return l.res, l.err
}

func selectValue(in types.AttributeValue, els []PathElement) (out types.AttributeValue, err error) {
	var field string
	var index int
	for i := 0; i < len(els); i++ {
		field, index = els[i].Field, els[i].Index
		if index >= 0 {
			switch tin := in.(type) {
			case *types.AttributeValueMemberL:
				in = tin.Value[index]
			case *types.AttributeValueMemberSS:
				in = &types.AttributeValueMemberS{Value: tin.Value[index]}
			case *types.AttributeValueMemberBS:
				in = &types.AttributeValueMemberB{Value: tin.Value[index]}
			case *types.AttributeValueMemberNS:
				in = &types.AttributeValueMemberN{Value: tin.Value[index]}
			default:
				return nil, fmt.Errorf("expected L, SS, BS, or NS, got: %T", tin)
			}
		} else if m, ok := in.(*types.AttributeValueMemberM); ok {
			in, ok = m.Value[field]
			if !ok {
				return nil, nil
			}
		} else {
			return nil, fmt.Errorf("unsupported select %v/%v for: %T", field, index, in)
		}
	}
	return in, nil
}

// SelectValues will return a subset of a composite attribute value 'v' (maps or sets), specified by 'paths'.
// This is usefull when only part of a DynamoDB item is allowed or desired for an operation. For example
// when only the keys need to be selected, or a partial update is performed using a mask.
func SelectValues(v types.AttributeValue, paths ...string) (res map[string]types.AttributeValue, err error) {
	res = make(map[string]types.AttributeValue, len(paths))
	els := make([]PathElement, 10) // allocate space for upto 10 element deep paths
	for _, p := range paths {
		if els, err = AppendParsePath(p, els[:0]); err != nil {
			return nil, fmt.Errorf("failed to parse plath '%s': %w", p, err)
		}

		r, err := selectValue(v, els)
		switch {
		case err != nil:
			return nil, fmt.Errorf("failed to select values: %w", err)
		case r == nil:
			continue
		default:
			res[p] = r
		}
	}
	return
}
