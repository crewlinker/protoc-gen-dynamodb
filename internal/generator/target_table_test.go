package generator_test

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtable"
	_ "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("flight fares", func() {
	It("should have registered table definitions", func() {
		def, ok := ddbtable.TableDef("flight_fares")
		Expect(ok).To(BeTrue())
		Expect(def).To(Equal(&ddbtable.Table{
			Name: "flight_fares",
			PartitionKey: &ddbtable.Attribute{
				Name: "1",
				Type: expression.String,
			},
			SortKey: &ddbtable.Attribute{
				Name: "2",
				Type: expression.String,
			},
			EntityType: &ddbtable.Attribute{
				Name: "3",
				Type: expression.String,
			},
			GlobalIndexes: []*ddbtable.GlobalIndex{{
				Name: "gsi1",
				PartitionKey: &ddbtable.Attribute{
					Name: "4",
					Type: expression.String,
				},
				SortKey: &ddbtable.Attribute{
					Name: "5",
					Type: expression.String,
				},
			}, {
				Name: "gsi2",
				PartitionKey: &ddbtable.Attribute{
					Name: "6",
					Type: expression.String,
				},
				SortKey: &ddbtable.Attribute{
					Name: "7",
					Type: expression.String,
				},
			}},
		}))
	})
})
