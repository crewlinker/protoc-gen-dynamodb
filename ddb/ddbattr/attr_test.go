package ddbattr_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdbattr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbexpr")
}

var _ = Describe("path building", func() {
	It("should build a valid path", func() {

	})
})

// var bench1 expression.Expression

// var ExtraBarExists0 = (Bar{}).ExtraBar().ExtraBar().ExtraBar().Foo().Brand().AttributeExists()
// var ExtraBarExists1 = expression.Name("extra_bar.extra_bar.extra_bar.foo.brand").AttributeExists()

// var Expr1, _ = expression.NewBuilder().WithCondition(ExtraBarExists0).Build()

// func BenchmarkNameBuilding(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for n := 0; n < b.N; n++ {

// 		bench1 = Expr1
// 	}
// }
