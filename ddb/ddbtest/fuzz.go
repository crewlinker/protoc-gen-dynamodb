package ddbtest

import (
	"encoding/base64"
	"math"
	"time"

	fuzz "github.com/google/gofuzz"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PbDurationFuzz fuzzes with some bounds on the duration as specified here
// https://pkg.go.dev/google.golang.org/protobuf/types/known/durationpb#Duration
func PbDurationFuzz(s *durationpb.Duration, c fuzz.Continue) {
	max := int64(math.MaxInt64)
	*s = *durationpb.New(time.Duration(c.Rand.Int63n(max) - (max / 2)))
}

// PbTimestampFuzz fuzzes with some bounds on the timestamp as specified here
// https://pkg.go.dev/google.golang.org/protobuf/types/known/durationpb#Duration
func PbTimestampFuzz(s *timestamppb.Timestamp, c fuzz.Continue) {
	max := int64(99999999999)
	*s = *timestamppb.New(time.Unix(c.Rand.Int63n(max)-(max/2), int64(c.RandUint64())))
}

// PbValueFuzz fuzzes code for structpb value. It doesn't recurse because go fuzz can't handle
// maps or lists with interface values.
func PbValueFuzz(s *structpb.Value, c fuzz.Continue) {
	switch c.Int63n(8) {
	case 0:
		s.Kind = &structpb.Value_BoolValue{BoolValue: c.RandBool()}
		return
	case 1:
		s.Kind = &structpb.Value_StringValue{StringValue: c.RandString()}
		return
	case 2:
		s.Kind = &structpb.Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}
		return
	case 3:
		s.Kind = &structpb.Value_NumberValue{NumberValue: c.ExpFloat64()}
		return
	case 4:
		lv := &structpb.Value_ListValue{}
		lv.ListValue, _ = structpb.NewList([]any{c.RandString()})
		s.Kind = lv
		return
	case 5:
		lv := &structpb.Value_StructValue{}
		lv.StructValue, _ = structpb.NewStruct(map[string]any{c.RandString(): c.RandString()})
		s.Kind = lv
		return
	case 6:
		s.Kind = &structpb.Value_NumberValue{NumberValue: float64(c.RandUint64())}
		return
	case 7:
		p := make([]byte, 10)
		c.Read(p)
		s.Kind = &structpb.Value_StringValue{StringValue: base64.StdEncoding.EncodeToString(p)}
		return
	default:
		panic("unsupported")
	}
}
