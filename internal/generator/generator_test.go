package generator_test

import (
	"context"
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
	cmd := exec.CommandContext(ctx, "buf", "generate")
	cmd.Dir = filepath.Join("..", "..", "example")
	cmd.Stderr = GinkgoWriter
	Expect(cmd.Run()).To(Succeed())
})

// test with messages defined in the example directory
var _ = Describe("example generation", func() {
	It("should marshal into expected attribute structure", func() {
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

	// We fuzz our implementation by filling the kitchen message with random data, then marshal and unmarshal
	// to check if it results in the output being equal to the input.
	DescribeTable("fuzz", func(seed int64) {
		f := fuzz.NewWithSeed(seed).NilChance(0.5)
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
		Entry("now()", time.Now().UnixNano()),
	)
})
