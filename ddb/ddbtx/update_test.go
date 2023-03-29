package ddbtx_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtx"
	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("update tx", func() {
	var ddb *dynamodb.Client
	BeforeEach(func(ctx context.Context) {
		acfg := aws.NewConfig()
		ddb = dynamodb.NewFromConfig(*acfg)
		Expect(ddb).ToNot(BeNil())
	})

	It("update without options", func(ctx context.Context) {
		car1 := &messagev1.Car{}

		// @TODO type-safe path constructions.

		Expect(ddbtx.Write().
			Update(car1, ddbtx.UpdateIfExists()).
			Commit(ctx),
		).To(Succeed())
	})

	It("partion update using masking", func() {

	})

	It("adding a numeric value to a nested field", func() {

	})
})
