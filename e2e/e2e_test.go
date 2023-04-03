package e2e_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e")
}

var _ = Describe("e2e tests", func() {
	var ddbc *dynamodb.Client
	BeforeEach(func(ctx context.Context) {
		var err error
		ddbc, err = ddbtest.NewLocalClient()
		Expect(err).ToNot(HaveOccurred())
	})

	It("should test something", func() {
		_ = ddbc
	})
})
