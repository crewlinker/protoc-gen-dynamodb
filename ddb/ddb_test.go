package ddb_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb")
}

var _ = Describe("path building", func() {
	It("should build with list of basic types", func() {
		expr, err := expression.NewBuilder().
			WithUpdate(
				expression.Set(
					messagev1.KitchenPath().Brand().N(),
					expression.Value("foo"))).
			Build()
		Expect(err).ToNot(HaveOccurred())
		Expect(expr.Names()).To(Equal(map[string]string{
			"#0": "1",
		}))
	})
})

var p1 string

func BenchmarkDeepNestingPathBuilding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		p1 = messagev1.KitchenPath().ExtraKitchen().ExtraKitchen().ApplianceEngines().At(5).Brand().String()
		if p1 != "16.16.19[5].1" {
			b.Fatalf("failed to build: %v", p1)
		}
	}
}

func BenchmarkPathBasicListBuilding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		p1 = messagev1.KitchenPath().ExtraKitchen().ExtraKitchen().OtherBrands().At(5).String()
		if p1 != "16.16.20[5]" {
			b.Fatalf("failed to build: %v", p1)
		}
	}
}
