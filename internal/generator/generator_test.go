package generator_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/generator")
}

var _ = BeforeSuite(func(ctx context.Context) {
	// cmd := exec.CommandContext(ctx, "buf", "generate")
	// cmd.Dir = filepath.Join("..", "..")
	// cmd.Stderr = GinkgoWriter
	// Expect(cmd.Run()).To(Succeed())
})

// ExpectJSONFieldPresense will assert that all fields are present in 'item' that are also present
// when encoding 'm' using canonical json.
func ExpectJSONFieldPresence(m proto.Message, it map[string]types.AttributeValue) {
	jm := map[string]any{}
	jb, err := protojson.Marshal(m)
	Expect(err).ToNot(HaveOccurred())
	Expect(json.Unmarshal(jb, &jm)).ToNot(HaveOccurred())

	jkeys, dkeys := lo.Keys(jm), lo.Keys(it)
	sort.Strings(jkeys)
	sort.Strings(dkeys)
	Expect(dkeys).To(Equal(jkeys))
}

// test with messages defined in the example directory
var _ = Describe("example generation", func() {
	It("should (un)marshal engine example", func() {
		k1 := &messagev1.Car{Engine: &messagev1.Engine{Brand: "somedrain", Dirtyness: messagev1.Dirtyness_DIRTYNESS_CLEAN}, Name: "foo"}
		m1, err := k1.MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())

		Expect(m1).To(Equal(map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{Value: "somedrain"},
				"2": &types.AttributeValueMemberN{Value: "1"},
			}},
			"2": &types.AttributeValueMemberS{Value: "foo"},
		}))

		var k2 messagev1.Car
		Expect(k2.UnmarshalDynamoItem(m1)).To(Succeed())
		Expect(proto.Equal(&k2, k1)).To(BeTrue())
	})

	// Assert that fields present in the resuling map are the same as present in canonical json encoding.
	// This works if we have a message where dynamo attr names are explicitely set to match that of the
	// the json encoding.
	DescribeTable("field presence", func(m *messagev1.FieldPresence, e map[string]types.AttributeValue) {
		it, err := m.MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())
		ExpectJSONFieldPresence(m, it)
	},
		Entry("presence message zero",
			&messagev1.FieldPresence{}, map[string]types.AttributeValue{}),
		Entry("presence with field values at (non-nil) zero",
			&messagev1.FieldPresence{
				Str:     "",
				OptStr:  proto.String(""),
				Msg:     &messagev1.Engine{},
				OptMsg:  &messagev1.Engine{},
				StrList: []string{},
				MsgList: []*messagev1.Engine{},
				StrMap:  map[string]string{},
				MsgMap:  map[string]*messagev1.Engine{},
				Enum:    messagev1.Dirtyness_DIRTYNESS_UNSPECIFIED,
				OptEnum: messagev1.Dirtyness_DIRTYNESS_UNSPECIFIED.Enum(),
			}, map[string]types.AttributeValue{}),
	)

	// Assert encoding of various kitchen messages
	DescribeTable("kitchen example", func(k *messagev1.Kitchen, exp map[string]types.AttributeValue, expErr string) {
		m, err := k.MarshalDynamoItem()
		if expErr == "" {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(MatchError(expErr))
		}

		// check message equality
		format.MaxLength = 0
		Expect(m).To(Equal(exp))
	},
		Entry("zero value",
			&messagev1.Kitchen{},
			map[string]types.AttributeValue{}, nil),
		Entry("with values",
			&messagev1.Kitchen{
				Brand:             "Siemens",
				IsRenovated:       true,
				QrCode:            []byte{'A'},
				NumSmallKnifes:    math.MaxInt32,
				NumSharpKnifes:    6,
				NumBluntKnifes:    math.MaxUint32,
				NumSmallForks:     math.MaxInt64,
				NumMediumForks:    20,
				NumLargeForks:     math.MaxUint64,
				PercentBlackTiles: math.MaxFloat32,
				PercentWhiteTiles: math.MaxFloat64,
				Dirtyness:         messagev1.Dirtyness_DIRTYNESS_CLEAN,
				Furniture:         map[int64]*messagev1.Appliance{100: {Brand: "Siemens"}},
				Calendar:          map[string]int64{"nov": 31},
				Timer:             durationpb.New((time.Second * 100) + 5),
				WallTime:          timestamppb.New(time.Unix(1678145849, 100)),
				ApplianceEngines:  []*messagev1.Engine{{Brand: "Kooks"}, {Brand: "Simens"}},
				OtherBrands:       []string{"Bosch", "Magimix"},
				SomeAny: &anypb.Any{
					TypeUrl: "type.googleapis.com/message.v1.Engine",
					Value:   []byte{10, 5, 75, 105, 107, 99, 104},
				},
				SomeMask: &fieldmaskpb.FieldMask{
					Paths: []string{"extra_kitchen.extra_kitchen.brand", "brand"},
				},
				SomeValue: func() *structpb.Value {
					v, _ := structpb.NewValue(map[string]any{"foo": "bar", "dar": 1})
					return v
				}(),

				// @TODO test with nil values for embedded messages
				// @TODO test with nil values for map entries
				// @TODO add test with nil value for map value, should marshal to NullValue
			},
			map[string]types.AttributeValue{
				// string/bool/bytes
				"1": &types.AttributeValueMemberS{Value: "Siemens"},
				"2": &types.AttributeValueMemberBOOL{Value: true},
				"3": &types.AttributeValueMemberB{Value: []byte{'A'}},
				// int32
				"4": &types.AttributeValueMemberN{Value: "2147483647"},
				"5": &types.AttributeValueMemberN{Value: "6"},
				"6": &types.AttributeValueMemberN{Value: "4294967295"},
				// int64
				"7": &types.AttributeValueMemberN{Value: "9223372036854775807"},
				"8": &types.AttributeValueMemberN{Value: "20"},
				"9": &types.AttributeValueMemberN{Value: "18446744073709551615"},
				// float/double
				"10": &types.AttributeValueMemberN{Value: "340282350000000000000000000000000000000"},
				"11": &types.AttributeValueMemberN{Value: "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
				// enum
				"12": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", messagev1.Dirtyness_DIRTYNESS_CLEAN)},
				// maps
				"13": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"100": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"1": &types.AttributeValueMemberS{Value: "Siemens"},
					}},
				}},
				"14": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"nov": &types.AttributeValueMemberN{Value: "31"},
				}},
				// duration uses protojson string encoding
				"17": &types.AttributeValueMemberS{Value: "100.000000005s"},
				// timestamp is encoded using protojson string encoding
				"18": &types.AttributeValueMemberS{Value: "2023-03-06T23:37:29.000000100Z"},
				// repeated nested message
				"19": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"1": &types.AttributeValueMemberS{Value: "Kooks"},
					}},
					&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"1": &types.AttributeValueMemberS{Value: "Simens"},
					}},
				}},
				// basic slice
				"20": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "Bosch"},
					&types.AttributeValueMemberS{Value: "Magimix"},
				}},
				// any message
				"21": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"1": &types.AttributeValueMemberS{Value: "type.googleapis.com/message.v1.Engine"},
					"2": &types.AttributeValueMemberB{Value: []byte{10, 5, 75, 105, 107, 99, 104}},
				}},
				// fieldmask message
				"22": &types.AttributeValueMemberSS{Value: []string{"extra_kitchen.extra_kitchen.brand", "brand"}},
				// value field
				"23": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"dar": &types.AttributeValueMemberN{Value: "1"},
					"foo": &types.AttributeValueMemberS{Value: "bar"},
				}},
			}, nil),
	)

	// We fuzz our implementation by filling the kitchen message with random data, then marshal and unmarshal
	// to check if it results in the output being equal to the input.
	DescribeTable("kitchen fuzz", func(seed int64) {
		f := fuzz.NewWithSeed(seed).NilChance(0.5)
		fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
		for i := 0; i < 10000; i++ {
			var in, out messagev1.Kitchen
			f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)

			item, err := in.MarshalDynamoItem()
			if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
				continue // skip, unsupported variant
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(out.UnmarshalDynamoItem(item)).To(Succeed())

			intxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&in)
			Expect(err).ToNot(HaveOccurred())
			outtxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&out)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(outtxt)).To(Equal(string(intxt))) // compare
			Expect(proto.Equal(&out, &in)).To(BeTrue(), fmt.Sprintf("%s \n!=\n %s", intxt, outtxt))
		}
	},
		// Table entries allow seeds that detected a regression to be used as future test cases
		Entry("1", int64(1)),
		Entry("panic marshal appliance", int64(1678177674234883000)),
		Entry("now()", time.Now().UnixNano()),
	)

	// We fuzz maps in particular because they are pretty finnicky to encode
	DescribeTable("map galore fuzz", func(seed int64) {
		f := fuzz.NewWithSeed(seed).NilChance(0.5)
		fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
		for i := 0; i < 10000; i++ {
			var in, out messagev1.MapGalore
			f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)
			item, err := in.MarshalDynamoItem()
			if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
				continue // skip, unsupported variant
			}

			Expect(out.UnmarshalDynamoItem(item)).To(Succeed())

			intxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&in)
			Expect(err).ToNot(HaveOccurred())
			outtxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&out)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(outtxt)).To(Equal(string(intxt))) // compare
			Expect(proto.Equal(&out, &in)).To(BeTrue(), fmt.Sprintf("%s \n!=\n %s", intxt, outtxt))
		}
	},
		// Table entries allow seeds that detected a regression to be used as future test cases
		Entry("now()", time.Now().UnixNano()),
		Entry("duration map", int64(1678219381135764000)),
	)

	// We fuzz maps in particular because they are pretty finnicky to encode
	DescribeTable("value galore fuzz", func(seed int64) {
		f := fuzz.NewWithSeed(seed).NilChance(0.5)
		fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
		for i := 0; i < 10000; i++ {
			var in, out messagev1.ValueGalore
			f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)
			item, err := in.MarshalDynamoItem()
			if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
				continue // skip, unsupported variant
			}

			Expect(out.UnmarshalDynamoItem(item)).To(Succeed())

			intxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&in)
			Expect(err).ToNot(HaveOccurred())
			outtxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(&out)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(outtxt)).To(Equal(string(intxt))) // compare
			Expect(proto.Equal(&out, &in)).To(BeTrue(), fmt.Sprintf("%s \n!=\n %s", intxt, outtxt))
		}
	},
		// Table entries allow seeds that detected a regression to be used as future test cases
		Entry("now()", time.Now().UnixNano()),
		Entry("first test", int64(1678219381135764000)),
	)
})

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
	var sf *structpb.Value
	switch c.Rand.Int63n(8) {
	case 0:
		sf = structpb.NewBoolValue(c.RandBool())
	case 1:
		sf = structpb.NewStringValue(c.RandString())
	case 2:
		sf = structpb.NewNullValue()
	case 3:
		sf = structpb.NewNumberValue(c.ExpFloat64())
	case 4:
		sf, _ = structpb.NewValue([]any{c.RandString()})
	case 5:
		sf, _ = structpb.NewValue(map[string]any{c.RandString(): c.RandString()})
	case 6:
		sf, _ = structpb.NewValue(c.RandUint64())
	case 7:
		p := make([]byte, 10)
		c.Read(p)
		sf, _ = structpb.NewValue(p)
	}
	*s = *sf
}
