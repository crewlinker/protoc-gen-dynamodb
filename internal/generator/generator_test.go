package generator_test

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	modelv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v1"
	modelv1ddbpath "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v1/ddbpath"
	fuzz "github.com/google/gofuzz"
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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGenerator(t *testing.T) {
	format.MaxLength = 0 // we produce long diff messages
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/generator")
}

// test with messages defined in the example directory
var _ = Describe("handling example messages", func() {
	It("should (un)marshal engine example", func() {
		c1 := &modelv1.Car{Engine: &modelv1.Engine{Brand: "somedrain", Dirtyness: modelv1.Dirtyness_DIRTYNESS_CLEAN}, Name: "foo"}
		m1, err := c1.MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())

		Expect(m1).To(Equal(map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{Value: "somedrain"},
				"2": &types.AttributeValueMemberN{Value: "1"},
			}},
			"2": &types.AttributeValueMemberS{Value: "foo"},
		}))

		var c2 modelv1.Car
		Expect(c2.UnmarshalDynamoItem(m1)).To(Succeed())
		ExpectProtoEqual(&c2, c1)
	})

	It("should have generated key functions", func() {
		Expect((&modelv1.Car{}).DynamoKeyNames()).To(Equal([]string{"ws"}))
		Expect(modelv1ddbpath.CarKeyNames()).To(Equal([]string{"ws"}))
		Expect(modelv1ddbpath.CarPartitionKey()).To(Equal(expression.Key("ws")))
		Expect((&modelv1.Car{}).DynamoPartitionKey()).To(Equal(expression.Key("ws")))
		Expect(modelv1ddbpath.CarPartitionKeyName()).To(Equal(expression.Name("ws")))
		Expect((&modelv1.Car{}).DynamoPartitionKeyName()).To(Equal(expression.Name("ws")))

		Expect((&modelv1.Kitchen{}).DynamoKeyNames()).To(Equal([]string{"1", "3"}))
		Expect(modelv1ddbpath.KitchenKeyNames()).To(Equal([]string{"1", "3"}))
		Expect(modelv1ddbpath.KitchenSortKey()).To(Equal(expression.Key("3")))
		Expect((&modelv1.Kitchen{}).DynamoSortKey()).To(Equal(expression.Key("3")))
		Expect(modelv1ddbpath.KitchenSortKeyName()).To(Equal(expression.Name("3")))
		Expect((&modelv1.Kitchen{}).DynamoSortKeyName()).To(Equal(expression.Name("3")))
	})

	It("should handle omit tags correctly", func() {
		msgt := reflect.TypeOf(&modelv1.Ignored{})
		_, ok := msgt.MethodByName("SortKey")
		Expect(ok).To(Equal(false))
		_, ok = msgt.MethodByName("PartitionKey")
		Expect(ok).To(Equal(false))

		pt := reflect.TypeOf(&modelv1ddbpath.IgnoredPath{})
		_, ok = pt.MethodByName("Pk")
		Expect(ok).To(Equal(false))
		_, ok = pt.MethodByName("Sk")
		Expect(ok).To(Equal(false))
		_, ok = pt.MethodByName("Other")
		Expect(ok).To(Equal(false))

		By("not marshalling from struct event if fields are not empty")
		m1, err := (&modelv1.Ignored{Pk: "Pk", Sk: "Sk", Other: "other", Visible: "visible"}).MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())
		Expect(m1).To(Equal(
			map[string]types.AttributeValue{
				"4": &types.AttributeValueMemberS{Value: "visible"},
			},
		))

		By("not unmarshalling into struct event if fieldsa re provided")
		m2 := &modelv1.Ignored{}
		Expect(m2.UnmarshalDynamoItem(map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberS{Value: "pk"},
			"2": &types.AttributeValueMemberS{Value: "sk"},
			"3": &types.AttributeValueMemberS{Value: "other"},
			"4": &types.AttributeValueMemberS{Value: "visible"},
		})).To(Succeed())
		Expect(m2).To(Equal(&modelv1.Ignored{Visible: "visible"}))
	})
})

// assert unmarshalling of various attribute maps
var _ = DescribeTable("kitchen unmarshaling", func(m map[string]types.AttributeValue, exp *modelv1.Kitchen) {
	var msg modelv1.Kitchen
	Expect(msg.UnmarshalDynamoItem(m)).To(Succeed())
	ExpectProtoEqual(&msg, exp)
},
	Entry("empty", map[string]types.AttributeValue{}, &modelv1.Kitchen{}),
)

// assert unmarshalling of various attribute maps with complicated presence
var _ = DescribeTable("presence unmarshaling", func(jb string, m map[string]types.AttributeValue, exp *modelv1.FieldPresence) {
	var msg modelv1.FieldPresence
	Expect(msg.UnmarshalDynamoItem(m)).To(Succeed())
	ExpectProtoEqual(&msg, exp)

	var jmsg modelv1.FieldPresence
	Expect(protojson.Unmarshal([]byte(jb), &jmsg)).To(Succeed())
	ExpectProtoEqual(&msg, &jmsg)
},
	Entry("empty", `{}`, map[string]types.AttributeValue{}, &modelv1.FieldPresence{}),

	Entry("null map 1", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": nil,
	}, &modelv1.FieldPresence{}),
	Entry("null map 2", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberNULL{},
	}, &modelv1.FieldPresence{}),
	Entry("null map 3", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberNULL{Value: true},
	}, &modelv1.FieldPresence{}),

	Entry("empty map 1", `{"strMap":{}}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberM{},
	}, &modelv1.FieldPresence{}),
	Entry("empty map 2", `{"strMap":{}}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
	}, &modelv1.FieldPresence{}),

	Entry("null msg", `{"msg":null}`, map[string]types.AttributeValue{}, &modelv1.FieldPresence{}),
	Entry("str value", `{"strVal":null}`, map[string]types.AttributeValue{}, &modelv1.FieldPresence{}),
	Entry("str value 2", `{"strVal":""}`, map[string]types.AttributeValue{
		"strVal": &types.AttributeValueMemberS{},
	}, &modelv1.FieldPresence{
		StrVal: wrapperspb.String(""),
	}),
)

// Assert marshalling output of the kitchen message
var _ = DescribeTable("kitchen marshaling", func(k *modelv1.Kitchen, exp map[string]types.AttributeValue, expErr string) {
	m, err := k.MarshalDynamoItem()
	if expErr == "" {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(MatchError(expErr))
	}

	// check message equality
	Expect(m).To(Equal(exp))
},
	Entry("zero value",
		&modelv1.Kitchen{},
		map[string]types.AttributeValue{}, nil),
	Entry("with values",
		&modelv1.Kitchen{
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
			Dirtyness:         modelv1.Dirtyness_DIRTYNESS_CLEAN,
			Furniture:         map[int64]*modelv1.Appliance{100: {Brand: "Siemens"}},
			Calendar:          map[string]int64{"nov": 31},
			Timer:             durationpb.New((time.Second * 100) + 5),
			WallTime:          timestamppb.New(time.Unix(1678145849, 100)),
			ApplianceEngines:  []*modelv1.Engine{{Brand: "Kooks"}, {Brand: "Simens"}},
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

			StringSet: []string{"a", "b", "c"},
			NumberSet: []int64{1, 100, 2000},
			BytesSet:  [][]byte{{0x01}, {0x02}, {0x03}},
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
			"12": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", modelv1.Dirtyness_DIRTYNESS_CLEAN)},
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
			"22": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberSS{Value: []string{"extra_kitchen.extra_kitchen.brand", "brand"}},
			}},
			// value field
			"23": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"dar": &types.AttributeValueMemberN{Value: "1"},
				"foo": &types.AttributeValueMemberS{Value: "bar"},
			}},
			"28": &types.AttributeValueMemberSS{
				Value: []string{"a", "b", "c"},
			},
			"29": &types.AttributeValueMemberNS{
				Value: []string{"1", "100", "2000"},
			},
			"30": &types.AttributeValueMemberBS{
				Value: [][]byte{{0x01}, {0x02}, {0x03}},
			},
		}, nil),
)

// Assert marshalling of json embeddings
var _ = DescribeTable("json embed marshalling", func(k *modelv1.JsonFields, exp map[string]types.AttributeValue, expErr string) {
	m, err := k.MarshalDynamoItem()
	if expErr == "" {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(MatchError(expErr))
	}

	// check message equality
	Expect(m).To(Equal(exp))
},
	Entry("zero value",
		&modelv1.JsonFields{},
		map[string]types.AttributeValue{}, nil),
	Entry("some data",
		&modelv1.JsonFields{
			JsonStrList:    []string{"a", "b", "c"},
			JsonEngine:     &modelv1.Engine{Brand: "brand-a"},
			JsonIntMap:     map[int64]string{100: "foo", 200: "bar"},
			JsonEngineList: []*modelv1.Engine{{Brand: "bar"}, {Brand: "foo"}},
			JsonEngineMap:  map[bool]*modelv1.Engine{true: {Brand: "true"}, false: {Brand: "false"}},
			JsonNrSet:      []int64{math.MaxInt64, 10},
		},
		map[string]types.AttributeValue{
			"1":           &types.AttributeValueMemberS{Value: `["a","b","c"]`},
			"json_engine": &types.AttributeValueMemberS{Value: `{"brand":"brand-a"}`},
			"4":           &types.AttributeValueMemberS{Value: `{"100":"foo","200":"bar"}`},
			"2":           &types.AttributeValueMemberS{Value: `[{"brand":"bar"},{"brand":"foo"}]`},
			"5":           &types.AttributeValueMemberS{Value: `{"false":{"brand":"false"},"true":{"brand":"true"}}`},
			"6":           &types.AttributeValueMemberS{Value: `[9223372036854775807,10]`},
		}, nil),
)

// Assert marshalling of json embeddings
var _ = DescribeTable("json embed oneof", func(in *modelv1.JsonOneofs, exp map[string]types.AttributeValue, expErr string) {
	m, err := in.MarshalDynamoItem()
	if expErr == "" {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(MatchError(expErr))
	}
	Expect(m).To(Equal(exp))
	var out modelv1.JsonOneofs
	Expect(out.UnmarshalDynamoItem(m)).To(Succeed())
	ExpectProtoEqual(in, &out)
},
	Entry("zero value",
		&modelv1.JsonOneofs{},
		map[string]types.AttributeValue{}, nil),
	Entry("one part",
		&modelv1.JsonOneofs{JsonOo: &modelv1.JsonOneofs_OneofMsg{OneofMsg: &modelv1.Engine{Brand: "brand"}}},
		map[string]types.AttributeValue{
			"8": &types.AttributeValueMemberS{Value: `{"brand":"brand"}`},
		}, nil),
	Entry("one part",
		&modelv1.JsonOneofs{JsonOo: &modelv1.JsonOneofs_OneofStr{OneofStr: "foo"}},
		map[string]types.AttributeValue{
			"7": &types.AttributeValueMemberS{Value: `"foo"`},
		}, nil),
)

// We fuzz json embedding
var _ = DescribeTable("json embed fuzz", func(seed int64) {
	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	for i := 0; i < 10000; i++ {
		var in, out modelv1.JsonFields
		f.Funcs(ddbtest.PbDurationFuzz, ddbtest.PbTimestampFuzz, ddbtest.PbValueFuzz).Fuzz(&in)

		item, err := in.MarshalDynamoItem()
		if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
			continue // skip, unsupported variant
		}

		Expect(err).ToNot(HaveOccurred())
		Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
		ExpectProtoEqual(&in, &out)
	}
},
	Entry("now()", time.Now().UnixNano()),
)

// We fuzz our implementation by filling the kitchen message with random data, then marshal and unmarshal
// to check if it results in the output being equal to the input.
var _ = DescribeTable("kitchen fuzz", func(seed int64) {
	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	for i := 0; i < 10000; i++ {
		var in, out modelv1.Kitchen
		f.Funcs(ddbtest.PbDurationFuzz, ddbtest.PbTimestampFuzz, ddbtest.PbValueFuzz).Fuzz(&in)

		item, err := in.MarshalDynamoItem()
		if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
			continue // skip, unsupported variant
		}

		Expect(err).ToNot(HaveOccurred())
		Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
		ExpectProtoEqual(&in, &out)
	}
},
	// Table entries allow seeds that detected a regression to be used as future test cases
	Entry("1", int64(1)),
	Entry("panic marshal appliance", int64(1678177674234883000)),
	Entry("now()", time.Now().UnixNano()),
)

// We fuzz maps in particular because they are pretty finnicky to encode
var _ = DescribeTable("map galore fuzz", func(seed int64) {
	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	for i := 0; i < 10000; i++ {
		var in, out modelv1.MapGalore
		f.Funcs(ddbtest.PbDurationFuzz, ddbtest.PbTimestampFuzz, ddbtest.PbValueFuzz).Fuzz(&in)
		item, err := in.MarshalDynamoItem()
		if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
			continue // skip, unsupported variant
		}

		Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
		ExpectProtoEqual(&in, &out)
	}
},
	// Table entries allow seeds that detected a regression to be used as future test cases
	Entry("now()", time.Now().UnixNano()),
	Entry("duration map", int64(1678219381135764000)),
)

// We fuzz structpb values in particular
var _ = DescribeTable("value galore fuzz", func(seed int64) {
	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	for i := 0; i < 10000; i++ {
		var in, out modelv1.ValueGalore
		f.Funcs(ddbtest.PbDurationFuzz, ddbtest.PbTimestampFuzz, ddbtest.PbValueFuzz).Fuzz(&in)
		item, err := in.MarshalDynamoItem()
		if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
			continue // skip, unsupported variant
		}

		Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
		ExpectProtoEqual(&in, &out)
	}
},
	// Table entries allow seeds that detected a regression to be used as future test cases
	Entry("now()", time.Now().UnixNano()),
	Entry("first test", int64(1678219381135764000)),
)

// Assert that fields present in the resuling map are the same as present in canonical json encoding.
// This works if we have a message where dynamo attr names are explicitely set to match that of the
// the json encoding.
var _ = DescribeTable("field presence", func(m *modelv1.FieldPresence, e map[string]types.AttributeValue) {
	it, err := m.MarshalDynamoItem()
	Expect(err).ToNot(HaveOccurred())
	ExpectJSONFieldPresence(m, it)

	var m2 modelv1.FieldPresence
	Expect(m2.UnmarshalDynamoItem(it)).To(Succeed())

	ExpectProtoEqual(&m2, m)
	Expect(it).To(Equal(e))
},
	Entry("presence message zero",
		&modelv1.FieldPresence{}, map[string]types.AttributeValue{}),
	Entry("presence with field values at (non-nil) zero",
		&modelv1.FieldPresence{
			Str:     "",
			OptStr:  proto.String(""),
			Msg:     &modelv1.Engine{},
			OptMsg:  &modelv1.Engine{},
			StrList: []string{},
			MsgList: []*modelv1.Engine{},
			StrMap:  map[string]string{},
			MsgMap:  map[string]*modelv1.Engine{},
			Enum:    modelv1.Dirtyness_DIRTYNESS_UNSPECIFIED,
			OptEnum: modelv1.Dirtyness_DIRTYNESS_UNSPECIFIED.Enum(),
			Oo:      &modelv1.FieldPresence_OneofStr{},

			StrVal:    wrapperspb.String(""),
			BoolVal:   wrapperspb.Bool(false),
			BytesVal:  wrapperspb.Bytes(nil),
			DoubleVal: wrapperspb.Double(0),
			FloatVal:  wrapperspb.Float(0),
			Int32Val:  wrapperspb.Int32(0),
			Int64Val:  wrapperspb.Int64(0),
			Uint32Val: wrapperspb.UInt32(0),
			Uint64Val: wrapperspb.UInt64(0),
		}, map[string]types.AttributeValue{
			"optMsg":   &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)},
			"optEnum":  &types.AttributeValueMemberN{Value: "0"},
			"oneofStr": &types.AttributeValueMemberS{},
			"optStr":   &types.AttributeValueMemberS{},
			"msg":      &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},

			"strVal":    &types.AttributeValueMemberS{},
			"boolVal":   &types.AttributeValueMemberBOOL{},
			"bytesVal":  &types.AttributeValueMemberNULL{Value: true},
			"doubleVal": &types.AttributeValueMemberN{Value: "0"},
			"floatVal":  &types.AttributeValueMemberN{Value: "0"},
			"int32Val":  &types.AttributeValueMemberN{Value: "0"},
			"int64Val":  &types.AttributeValueMemberN{Value: "0"},
			"uint32Val": &types.AttributeValueMemberN{Value: "0"},
			"uint64Val": &types.AttributeValueMemberN{Value: "0"},
		}),
	Entry("just oneof mesage",
		&modelv1.FieldPresence{
			Oo: &modelv1.FieldPresence_OneofMsg{OneofMsg: &modelv1.Engine{Brand: "foo"}},
		}, map[string]types.AttributeValue{
			"oneofMsg": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"1": &types.AttributeValueMemberS{
						Value: "foo",
					},
				},
			},
		}),

	Entry("nil values for map entries and lists", &modelv1.FieldPresence{
		MsgMap:  map[string]*modelv1.Engine{"foo": nil},
		MsgList: []*modelv1.Engine{nil},
	}, map[string]types.AttributeValue{
		"msgList": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberNULL{
					Value: true,
				},
			},
		},
		"msgMap": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"foo": &types.AttributeValueMemberNULL{
					Value: true,
				},
			},
		},
	}),
)

// ExpectProtoEqual compares to proto messages while providing easier to debug output if
// it fails.
func ExpectProtoEqual(a, b proto.Message) {
	GinkgoHelper()
	atxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(a)
	Expect(err).ToNot(HaveOccurred())
	btxt, err := (prototext.MarshalOptions{Multiline: true}).Marshal(b)
	Expect(err).ToNot(HaveOccurred())

	Expect(string(atxt)).To(Equal(string(btxt))) // compare
	Expect(proto.Equal(a, b)).To(BeTrue(), fmt.Sprintf("%s \n!=\n %s", atxt, btxt))
}

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
