package ddbtest_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdbtest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbtest")
}

var _ = Describe("local testing", func() {
	It("should allow describing local tables", func(ctx context.Context) {
		ddbc, err := ddbtest.NewLocalClient()
		Expect(err).ToNot(HaveOccurred())
		_, err = ddbc.ListTables(ctx, &dynamodb.ListTablesInput{Limit: aws.Int32(1)})
		Expect(err).ToNot(HaveOccurred())
	})
})
