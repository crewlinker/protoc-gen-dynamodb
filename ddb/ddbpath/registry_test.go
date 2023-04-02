package ddbpath_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	messagev1ddbpath "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1/ddbpath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("strings", func() {
	It("shoud panic unsupported", func() {
		Expect(func() {
			_ = ddbpath.FieldKind(999).String()
		}).To(PanicWith(MatchRegexp(`unsupported`)))
	})

	DescribeTable("kinds", func(k ddbpath.FieldKind, exp string) {
		Expect(k.String()).To(Equal(exp))
	},
		Entry("1", ddbpath.FieldKind(0), "_undefined"),
		Entry("1", ddbpath.FieldKindSingle, "Single"),
		Entry("1", ddbpath.FieldKindList, "List"),
		Entry("1", ddbpath.FieldKindMap, "Map"),
	)

	DescribeTable("field", func(k ddbpath.FieldInfo, exp string) {
		Expect(k.String()).To(Equal(exp))
	},
		Entry("1", ddbpath.FieldInfo{Kind: ddbpath.FieldKindList}, "List"),
		Entry("1", ddbpath.FieldInfo{Kind: ddbpath.FieldKindList, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())}, "List<messagev1ddbpath.KitchenPath>"))
})

var _ = Describe("validate", func() {
	var reg ddbpath.Registry
	BeforeEach(func() { reg = ddbpath.NewRegistry() })

	Describe("with kitchen registered", func() {
		BeforeEach(func() { reg.Register(messagev1ddbpath.Kitchen(), map[string]ddbpath.FieldInfo{}) })
		It("should panic on double registration", func() {
			Expect(func() {
				reg.Register(messagev1ddbpath.Kitchen(), map[string]ddbpath.FieldInfo{})
			}).To(PanicWith(MatchRegexp(`is already registered for validation`)))
		})

		It("should allow returning field info", func() {
			fi1, ok := reg.FieldsOf(messagev1ddbpath.Kitchen())
			Expect(ok).To(BeTrue())
			Expect(fi1).ToNot(BeNil())

			fi2, ok := reg.FieldsOf(messagev1ddbpath.Car())
			Expect(ok).To(BeFalse())
			Expect(fi2).To(BeNil())
		})
	})

	Describe("with registered", func() {
		BeforeEach(func() {
			reg.Register(ddbpath.ValuePath{}, nil)
			reg.Register(messagev1ddbpath.Kitchen(), map[string]ddbpath.FieldInfo{
				"1":  {Kind: ddbpath.FieldKindSingle},
				"16": {Kind: ddbpath.FieldKindSingle, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())}, // single message
				// lists
				"17": {Kind: ddbpath.FieldKindList, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())}, // list of messages
				"18": {Kind: ddbpath.FieldKindList},                                                      // list of basic types
				// maps
				"19": {Kind: ddbpath.FieldKindMap, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())},
				"20": {Kind: ddbpath.FieldKindMap},
				// any message
				"22": {Kind: ddbpath.FieldKindSingle, Message: reflect.TypeOf(ddbpath.ValuePath{})},
			})
		})

		DescribeTable("validation", func(nb ddbpath.NameBuilder, p string, expError string) {
			err := reg.Validate(nb, p)
			if expError == `` {
				Expect(err).To(BeNil())
			} else {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(MatchRegexp(expError))
			}
		},
			Entry("type not registered", expression.NameBuilder{}, "1", `type not registered: expression.NameBuilder`),
			Entry("select too deep", messagev1ddbpath.Kitchen(), "1.1", `field selecting '1' not allowed on Single`),
			Entry("index not allowed", messagev1ddbpath.Kitchen(), "[1]", `indexing '1' not allowed on Single<messagev1ddbpath.KitchenPath>`),
			Entry("field not allowed", messagev1ddbpath.Kitchen(), "17.1", `field selecting '1' not allowed on List<messagev1ddbpath.KitchenPath>`),
			Entry("unknown field", messagev1ddbpath.Kitchen(), "9999", `unknown field '9999' of Single<messagev1ddbpath.KitchenPath>`),
			Entry("basic field", messagev1ddbpath.Kitchen(), "1", ``),
			Entry("recurse", messagev1ddbpath.Kitchen(), "16.1", ``),
			Entry("recursed unknown field", messagev1ddbpath.Kitchen(), "16.999", `unknown field '999' of Single<messagev1ddbpath.KitchenPath>`),
			Entry("select list", messagev1ddbpath.Kitchen(), "17[1]", ``),
			Entry("select recurse list", messagev1ddbpath.Kitchen(), "17[1].16.16.16.1", ``),
			Entry("select basic list", messagev1ddbpath.Kitchen(), "18[1]", ``),
			Entry("select basic list", messagev1ddbpath.Kitchen(), "18[1].1", `field selecting '1' not allowed on Single`),
			Entry("map of messages", messagev1ddbpath.Kitchen(), "19.foo.16.1", ``),
			Entry("map of basic", messagev1ddbpath.Kitchen(), "20.foo", ``),
			Entry("map of basic", messagev1ddbpath.Kitchen(), "20.foo.999", `field selecting '999' not allowed on Single`),
			// any message
			Entry("any message", messagev1ddbpath.Kitchen(), "22[999].1.1", ``),
			Entry("any message", messagev1ddbpath.Kitchen(), "22.999.a[1]", ``),
		)
	})

})

func BenchmarkValidate(b *testing.B) {
	reg := ddbpath.NewRegistry()
	reg.Register(messagev1ddbpath.Kitchen(), map[string]ddbpath.FieldInfo{
		"1":  {Kind: ddbpath.FieldKindSingle},
		"16": {Kind: ddbpath.FieldKindSingle, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())}, // single message
		"17": {Kind: ddbpath.FieldKindList, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())},   // list of messages
		"18": {Kind: ddbpath.FieldKindList},                                                        // list of basic types
		"19": {Kind: ddbpath.FieldKindMap, Message: reflect.TypeOf(messagev1ddbpath.Kitchen())},
		"20": {Kind: ddbpath.FieldKindMap},
	})

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		err := reg.Validate(messagev1ddbpath.Kitchen(),
			"1",
			"16.16.16.1",
			"17[1].16.16.16.1",
			"19.foo.16.1",
			"1",
			"16.16.16.1",
			"17[1].16.16.16.1",
			"19.foo.16.1",
			"1",
			"16.16.16.1",
			"17[1].16.16.16.1",
			"19.foo.16.1",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
