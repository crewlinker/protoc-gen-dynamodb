package ddbpath_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	messagev1ddbpath "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1/ddbpath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("validate", func() {
	It("should error on unregistered type", func() {
		err := ddbpath.Validate(expression.Name("foo"), "1")
		Expect(err.Error()).To(MatchRegexp(`type not registered`))
	})

	DescribeTable("validate table", func(nb messagev1ddbpath.KitchenPath, paths []string, expErr string) {
		err := ddbpath.Validate(nb, paths...)
		if expErr == "" {
			Expect(err).To(BeNil())
		} else {
			Expect(err.Error()).To(MatchRegexp(expErr))
		}
	},
		Entry("simple top-level field", messagev1ddbpath.Kitchen(), []string{"1"}, ""),
		Entry("non existing top-level field", messagev1ddbpath.Kitchen(), []string{"999"}, "non-existing field '999' on: messagev1ddbpath.KitchenPath"),
		Entry("valid deep nested", messagev1ddbpath.Kitchen(), []string{"16.16.16.1"}, ""),
		Entry("invalid deep nested", messagev1ddbpath.Kitchen(), []string{"16.16.16.999"}, "non-existing field '999' on: messagev1ddbpath.KitchenPath"),

		// list indexing
		Entry("valid index nested", messagev1ddbpath.Kitchen(), []string{"19[1]"}, ""),
		Entry("valid into index", messagev1ddbpath.Kitchen(), []string{"19[1].1"}, ""),
		Entry("invalid field on non-message", messagev1ddbpath.Kitchen(), []string{"19[1].1.1"}, `path \(or index\) '1' on basic type`),
		Entry("invalid index on non-message", messagev1ddbpath.Kitchen(), []string{"19[1].1[1]"}, `path \(or index\) '1' on basic type`),
		Entry("index into non list", messagev1ddbpath.Kitchen(), []string{"13[1]"}, `index '\[1\]' into non-list: Map`),

		// map keys
		Entry("valid map path", messagev1ddbpath.Kitchen(), []string{"13.999.1"}, ""),
	)
})

func BenchmarkValidate(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		err := ddbpath.Validate(messagev1ddbpath.Kitchen(),
			"1",
			"16.16.16.1",
			"19[1].1",
			"13.999.1",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
