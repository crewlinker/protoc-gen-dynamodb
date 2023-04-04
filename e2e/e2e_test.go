package e2e_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e")
}

var _ = Describe("e2e tests", func() {
	var tblname string
	var ddbc *dynamodb.Client
	BeforeEach(func(ctx context.Context) {
		var err error
		ddbc, err = ddbtest.NewLocalClient()
		Expect(err).ToNot(HaveOccurred())

		tblname = fmt.Sprintf("%d", time.Now().UnixNano())
		Expect(ddbc.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tblname),
			KeySchema: []types.KeySchemaElement{
				{KeyType: types.KeyTypeHash, AttributeName: aws.String("1")},
				{KeyType: types.KeyTypeRange, AttributeName: aws.String("3")},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("1")},
				{AttributeType: types.ScalarAttributeTypeB, AttributeName: aws.String("3")},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{WriteCapacityUnits: aws.Int64(1), ReadCapacityUnits: aws.Int64(1)},
		})).ToNot(BeNil())

		DeferCleanup(func(ctx context.Context) {
			Expect(ddbc.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &tblname})).ToNot(BeNil())
			tl, err := ddbc.ListTables(ctx, &dynamodb.ListTablesInput{})
			Expect(err).ToNot(HaveOccurred())
			Expect(tl.TableNames).To(HaveLen(0))
		})
	})

	// DescribeTable("put get fuzzing", func(ctx context.Context, seed int64) {
	// 	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	// 	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	// 	for i := 0; i < 3; i++ {
	// 		var in messagev1.Engine
	// 		f.Funcs(ddbtest.PbDurationFuzz, ddbtest.PbTimestampFuzz, ddbtest.PbValueFuzz).Fuzz(&in)
	// 		// in.Brand = "my-brand"
	// 		// in.QrCode = []byte{0x01}

	// 		item, err := in.MarshalDynamoItem()
	// 		if err != nil && strings.Contains(err.Error(), "map key cannot be empty") {
	// 			continue // skip, unsupported variant
	// 		}
	// 		Expect(err).ToNot(HaveOccurred())
	// 		fmt.Println("AAAAAA", item)
	// 		Expect(ddbc.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(tblname), Item: item})).ToNot(BeNil())

	// 		// Expect(out.UnmarshalDynamoItem(item)).To(Succeed())
	// 		// ExpectProtoEqual(&in, &out)
	// 	}
	// },
	// 	// Table entries allow seeds that detected a regression to be used as future test cases
	// 	Entry("1", int64(1)),
	// )

	// It("should fuzz put and get item", func(ctx context.Context) {
	// 	f := fuzz.NewWithSeed(seed).NilChance(0.5)
	// 	fmt.Fprintf(GinkgoWriter, "Fuzz Seed: %d", seed)
	// 	for i := 0; i < 10000; i++ {
	// 		var in, out messagev1.Kitchen
	// 		f.Funcs(PbDurationFuzz, PbTimestampFuzz, PbValueFuzz).Fuzz(&in)

	// 	}

	// })
})
