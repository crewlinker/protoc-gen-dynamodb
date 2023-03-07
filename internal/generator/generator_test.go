package generator_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/example/proto/message/v1"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/generator")
}

var _ = BeforeSuite(func(ctx context.Context) {
	if os.Getenv("PROTOC_GEN_DYNAMODB_TEST_NO_GENERATE") != "" {
		return // disable when env asks to
	}

	cmd := exec.CommandContext(ctx, "buf", "generate")
	cmd.Dir = filepath.Join("..", "..", "example")
	cmd.Stderr = GinkgoWriter
	Expect(cmd.Run()).To(Succeed())
})

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
		Expect(&k2).To(Equal(k1))
	})

	// Assert encoding of various kitchen messages
	DescribeTable("kitchen example", func(k *messagev1.Kitchen, exp map[string]types.AttributeValue, expErr string) {
		m, err := k.MarshalDynamoItem()
		if expErr == "" {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(MatchError(expErr))
		}
		Expect(m).To(Equal(exp))
	},
		Entry("zero value",
			&messagev1.Kitchen{},
			map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{},
				"2": &types.AttributeValueMemberBOOL{},
				"3": &types.AttributeValueMemberNULL{Value: true},

				"4": &types.AttributeValueMemberN{Value: "0"},
				"5": &types.AttributeValueMemberN{Value: "0"},
				"6": &types.AttributeValueMemberN{Value: "0"},

				"7": &types.AttributeValueMemberN{Value: "0"},
				"8": &types.AttributeValueMemberN{Value: "0"},
				"9": &types.AttributeValueMemberN{Value: "0"},

				"10": &types.AttributeValueMemberN{Value: "0"},
				"11": &types.AttributeValueMemberN{Value: "0"},

				"12": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", messagev1.Dirtyness_DIRTYNESS_UNSPECIFIED)},
			}, nil),
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
			}, nil),
	)

	// We fuzz our implementation by filling the kitchen message with random data, then marshal and unmarshal
	// to check if it results in the output being equal to the input.
	DescribeTable("kitchen fuzz", func(seed int64) {
		f := fuzz.NewWithSeed(seed).NilChance(0.5)
		fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
		for i := 0; i < 10000; i++ {
			var in, out messagev1.Kitchen
			f.Fuzz(&in)
			item, err := in.MarshalDynamoItem()
			Expect(err).ToNot(HaveOccurred())

			Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
			Expect(&out).To(Equal(&in))
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
			f.Fuzz(&in)
			item, err := in.MarshalDynamoItem()
			Expect(err).ToNot(HaveOccurred())

			Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
			Expect(&out).To(Equal(&in))
		}
	},
		// Table entries allow seeds that detected a regression to be used as future test cases
		Entry("now()", time.Now().UnixNano()),
	)
})
