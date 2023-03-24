package generator_test

import (
	"encoding/base64"
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
	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
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
		c1 := &messagev1.Car{Engine: &messagev1.Engine{Brand: "somedrain", Dirtyness: messagev1.Dirtyness_DIRTYNESS_CLEAN}, Name: "foo"}
		m1, err := c1.MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())

		Expect(m1).To(Equal(map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{Value: "somedrain"},
				"2": &types.AttributeValueMemberN{Value: "1"},
			}},
			"2": &types.AttributeValueMemberS{Value: "foo"},
		}))

		var c2 messagev1.Car
		Expect(c2.UnmarshalDynamoItem(m1)).To(Succeed())
		ExpectProtoEqual(&c2, c1)
	})

	It("should have generated pk/sk methods", func() {
		c1 := &messagev1.Car{Name: "foo", NrOfWheels: 4}
		pk, pkv := c1.PartitionKey()
		Expect(pk).To(Equal("ws"))
		Expect(pkv).To(Equal(int64(4)))

		sk, skv := c1.SortKey()
		Expect(sk).To(Equal("2"))
		Expect(skv).To(Equal("foo"))
	})

	It("marshal key should work as expected", func() {
		c1 := &messagev1.Car{Name: "foo", NrOfWheels: 4}
		k1, err := c1.MarshalDynamoKey()
		Expect(err).ToNot(HaveOccurred())
		Expect(k1).To(Equal(map[string]types.AttributeValue{
			"ws": &types.AttributeValueMemberN{Value: "4"},
			"2":  &types.AttributeValueMemberS{Value: "foo"},
		}))

		c2 := &messagev1.Car{}
		k2, err := c2.MarshalDynamoKey()
		Expect(err).ToNot(HaveOccurred())
		Expect(k2).To(Equal(map[string]types.AttributeValue{
			"ws": &types.AttributeValueMemberN{Value: "0"},
			"2":  &types.AttributeValueMemberS{Value: ""},
		}))
	})

	It("should handle omit tags correctly", func() {
		msgt := reflect.TypeOf(&messagev1.Ignored{})
		_, ok := msgt.MethodByName("SortKey")
		Expect(ok).To(Equal(false))
		_, ok = msgt.MethodByName("PartitionKey")
		Expect(ok).To(Equal(false))

		pt := reflect.TypeOf(&messagev1.IgnoredP{})
		_, ok = pt.MethodByName("Pk")
		Expect(ok).To(Equal(false))
		_, ok = pt.MethodByName("Sk")
		Expect(ok).To(Equal(false))
		_, ok = pt.MethodByName("Other")
		Expect(ok).To(Equal(false))

		By("not marshalling from struct event if fields are not empty")
		m1, err := (&messagev1.Ignored{Pk: "Pk", Sk: "Sk", Other: "other", Visible: "visible"}).MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())
		Expect(m1).To(Equal(
			map[string]types.AttributeValue{
				"4": &types.AttributeValueMemberS{Value: "visible"},
			},
		))

		By("not unmarshalling into struct event if fieldsa re provided")
		m2 := &messagev1.Ignored{}
		Expect(m2.UnmarshalDynamoItem(map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberS{Value: "pk"},
			"2": &types.AttributeValueMemberS{Value: "sk"},
			"3": &types.AttributeValueMemberS{Value: "other"},
			"4": &types.AttributeValueMemberS{Value: "visible"},
		})).To(Succeed())
		Expect(m2).To(Equal(&messagev1.Ignored{Visible: "visible"}))
	})
})

// test the building of paths
var _ = DescribeTable("path building", func(s interface {
	fmt.Stringer
	N() expression.NameBuilder
}, exp string) {
	Expect(s.String()).To(Equal(exp))

	// check that expression builder accepts the paths
	_, err := expression.NewBuilder().WithUpdate(
		expression.Set(s.N(), expression.Value("foo")),
	).Build()
	Expect(err).ToNot(HaveOccurred())
},
	Entry("basic type field", messagev1.KitchenPath().Brand(), "1"),
	Entry("message type fields", messagev1.KitchenPath().ExtraKitchen().Brand(), "16.1"),
	Entry("lists of basic types itself", messagev1.KitchenPath().OtherBrands().At(10), "20[10]"),
	Entry("through list of messages", messagev1.KitchenPath().ApplianceEngines().At(3).Brand(), "19[3].1"),
	Entry("to list of messages itself", messagev1.KitchenPath().ApplianceEngines(), "19"),
	Entry("to message field itself", messagev1.KitchenPath().ExtraKitchen(), "16"),
	Entry("to field with renamed attr", messagev1.FieldPresencePath().Str(), "str"),

	// message not in the same package only support direct path building, not "Through". Including
	// well-known types.
	Entry("well-known message field", messagev1.KitchenPath().Timer(), "17"),
	Entry("list of well-known message field", messagev1.KitchenPath().ListOfTs().At(4), "27[4]"),

	// map access, als has limit on messages outside of the package
	Entry("to map of basic types itself", messagev1.MapGalorePath().Int64Int64().Key("a"), "1.a"),
	Entry("to map of well-known itself", messagev1.MapGalorePath().Stringtimestamp().Key("b"), "17.b"),
	Entry("through map of messages", messagev1.KitchenPath().ExtraKitchen().Furniture().Key("foo").Brand(), "16.13.foo.1"),
)

// assert unmarshalling of various attribute maps
var _ = DescribeTable("kitchen unmarshaling", func(m map[string]types.AttributeValue, exp *messagev1.Kitchen) {
	var msg messagev1.Kitchen
	Expect(msg.UnmarshalDynamoItem(m)).To(Succeed())
	ExpectProtoEqual(&msg, exp)
},
	Entry("empty", map[string]types.AttributeValue{}, &messagev1.Kitchen{}),
)

// assert unmarshalling of various attribute maps with complicated presence
var _ = DescribeTable("presence unmarshaling", func(jb string, m map[string]types.AttributeValue, exp *messagev1.FieldPresence) {
	var msg messagev1.FieldPresence
	Expect(msg.UnmarshalDynamoItem(m)).To(Succeed())
	ExpectProtoEqual(&msg, exp)

	var jmsg messagev1.FieldPresence
	Expect(protojson.Unmarshal([]byte(jb), &jmsg)).To(Succeed())
	ExpectProtoEqual(&msg, &jmsg)
},
	Entry("empty", `{}`, map[string]types.AttributeValue{}, &messagev1.FieldPresence{}),

	Entry("null map 1", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": nil,
	}, &messagev1.FieldPresence{}),
	Entry("null map 2", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberNULL{},
	}, &messagev1.FieldPresence{}),
	Entry("null map 3", `{"strMap":null}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberNULL{Value: true},
	}, &messagev1.FieldPresence{}),

	Entry("empty map 1", `{"strMap":{}}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberM{},
	}, &messagev1.FieldPresence{}),
	Entry("empty map 2", `{"strMap":{}}`, map[string]types.AttributeValue{
		"strMap": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
	}, &messagev1.FieldPresence{}),

	Entry("null msg", `{"msg":null}`, map[string]types.AttributeValue{}, &messagev1.FieldPresence{}),
	Entry("str value", `{"strVal":null}`, map[string]types.AttributeValue{}, &messagev1.FieldPresence{}),
	Entry("str value 2", `{"strVal":""}`, map[string]types.AttributeValue{
		"strVal": &types.AttributeValueMemberS{},
	}, &messagev1.FieldPresence{
		StrVal: wrapperspb.String(""),
	}),
)

// Assert marshalling output of the kitchen message
var _ = DescribeTable("kitchen marshaling", func(k *messagev1.Kitchen, exp map[string]types.AttributeValue, expErr string) {
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
var _ = DescribeTable("json embed marshalling", func(k *messagev1.JsonFields, exp map[string]types.AttributeValue, expErr string) {
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
		&messagev1.JsonFields{},
		map[string]types.AttributeValue{}, nil),
	Entry("some data",
		&messagev1.JsonFields{
			JsonStrList:    []string{"a", "b", "c"},
			JsonEngine:     &messagev1.Engine{Brand: "brand-a"},
			JsonIntMap:     map[int64]string{100: "foo", 200: "bar"},
			JsonEngineList: []*messagev1.Engine{{Brand: "bar"}, {Brand: "foo"}},
			JsonEngineMap:  map[bool]*messagev1.Engine{true: {Brand: "true"}, false: {Brand: "false"}},
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

// We fuzz json embedding
var _ = DescribeTable("json embed fuzz", func(seed int64) {
	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	for i := 0; i < 10000; i++ {
		var in, out messagev1.JsonFields
		f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)

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
		var in, out messagev1.Kitchen
		f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)

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
		var in, out messagev1.MapGalore
		f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)
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
		var in, out messagev1.ValueGalore
		f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)
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
var _ = DescribeTable("field presence", func(m *messagev1.FieldPresence, e map[string]types.AttributeValue) {
	it, err := m.MarshalDynamoItem()
	Expect(err).ToNot(HaveOccurred())
	ExpectJSONFieldPresence(m, it)

	var m2 messagev1.FieldPresence
	Expect(m2.UnmarshalDynamoItem(it)).To(Succeed())

	ExpectProtoEqual(&m2, m)
	Expect(it).To(Equal(e))
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
			Oo:      &messagev1.FieldPresence_OneofStr{},

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
		&messagev1.FieldPresence{
			Oo: &messagev1.FieldPresence_OneofMsg{OneofMsg: &messagev1.Engine{Brand: "foo"}},
		}, map[string]types.AttributeValue{
			"oneofMsg": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"1": &types.AttributeValueMemberS{
						Value: "foo",
					},
				},
			},
		}),

	Entry("nil values for map entries and lists", &messagev1.FieldPresence{
		MsgMap:  map[string]*messagev1.Engine{"foo": nil},
		MsgList: []*messagev1.Engine{nil},
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
	switch c.Rand.Int63n(8) {
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
