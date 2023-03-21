package ddb_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb")
}

var _ = Describe("expression building", func() {

	It("should build with nested names", func() {
		expr, err := expression.NewBuilder().
			WithUpdate(expression.Set(expression.Name("16[1]"), expression.Value("foo"))).
			Build()
		Expect(err).ToNot(HaveOccurred())

		fmt.Println("Update:", *expr.Update())
		fmt.Println("Values:", expr.Values())
		fmt.Println("Names:", expr.Names())
		_ = expr
	})

})
