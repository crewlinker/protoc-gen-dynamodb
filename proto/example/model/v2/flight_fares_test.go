package modelv2_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	modelv2 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("flight fares", func() {
	var tblname string
	var ddbc *dynamodb.Client
	var mdl *modelv2.FlightFaresModel
	var mut *modelv2.FlightFaresMutater
	var qry *modelv2.FlightFaresQuerier
	BeforeEach(func(ctx context.Context) {
		var err error
		ddbc, err = ddbtest.NewLocalClient()
		Expect(err).ToNot(HaveOccurred())

		tblname = fmt.Sprintf("flight_fares_%d", time.Now().UnixNano())
		mdl = &modelv2.FlightFaresModel{}
		mut = modelv2.NewFlightFaresMutater(tblname, ddbc, mdl)
		qry = modelv2.NewFlightFaresQuerier(tblname, ddbc, mdl)

		Expect(ddbc.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tblname),
			KeySchema: []types.KeySchemaElement{
				{KeyType: types.KeyTypeHash, AttributeName: aws.String("1")},
				{KeyType: types.KeyTypeRange, AttributeName: aws.String("2")},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName:  aws.String("gsi1"),
					Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("4"), KeyType: types.KeyTypeHash},
						{AttributeName: aws.String("5"), KeyType: types.KeyTypeRange},
					},
				},
				{
					IndexName:  aws.String("gsi2"),
					Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("6"), KeyType: types.KeyTypeHash},
						{AttributeName: aws.String("7"), KeyType: types.KeyTypeRange},
					},
				},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("1")},
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("2")},
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("4")},
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("5")},
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("6")},
				{AttributeType: types.ScalarAttributeTypeS, AttributeName: aws.String("7")},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				WriteCapacityUnits: aws.Int64(1),
				ReadCapacityUnits:  aws.Int64(1),
			},
		})).ToNot(BeNil())

		DeferCleanup(func(ctx context.Context) {
			Expect(ddbc.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &tblname})).ToNot(BeNil())
			tl, err := ddbc.ListTables(ctx, &dynamodb.ListTablesInput{})
			Expect(err).ToNot(HaveOccurred())
			Expect(tl.TableNames).To(HaveLen(0))
		})
	})

	Describe("with example data", func() {
		BeforeEach(func(ctx context.Context) {
			for _, flight := range []*modelv2.Flight{
				{
					Number: 160,
					Origin: modelv2.Airport_AIRPORT_DEN, Destination: modelv2.Airport_AIRPORT_SFO,
					Class: modelv2.FlightClass_FLIGHT_CLASS_NON_STOP, IsSegment: true, SegmentId: 2,
					DepatureAt: timestamppb.New(time.Unix(1627916100, 0)), // 2021-08-02T16:55:00
					ArrivalAt:  timestamppb.New(time.Unix(1627921500, 0)), // 2021-08-02T18:25:00
				},
				{
					Number: 150,
					Origin: modelv2.Airport_AIRPORT_DEN, Destination: modelv2.Airport_AIRPORT_JFK,
					Class: modelv2.FlightClass_FLIGHT_CLASS_NON_STOP, IsSegment: true, SegmentId: 2,
					DepatureAt: timestamppb.New(time.Unix(1627806300, 0)), // 2021-08-01T10:25:00
					ArrivalAt:  timestamppb.New(time.Unix(1627820700, 0)), // 2021-08-01T14:25:00
				},

				{
					Number: 260,
					Origin: modelv2.Airport_AIRPORT_JFK, Destination: modelv2.Airport_AIRPORT_SFO,
					Class:      modelv2.FlightClass_FLIGHT_CLASS_NON_STOP,
					DepatureAt: timestamppb.New(time.Unix(1627820700, 0)), // 2021-08-01T14:25:00
					ArrivalAt:  timestamppb.New(time.Unix(1627831500, 0)), // 2021-08-01T17:25:00
				},
				{
					Number: 160,
					Origin: modelv2.Airport_AIRPORT_JFK, Destination: modelv2.Airport_AIRPORT_SFO,
					Class:      modelv2.FlightClass_FLIGHT_CLASS_DIRECT,
					DepatureAt: timestamppb.New(time.Unix(1627824300, 0)), // 2021-08-01T15:25:00
					ArrivalAt:  timestamppb.New(time.Unix(1627842300, 0)), // 2021-08-03T20:25:00
				},
			} {
				Expect(mut.PutFlight(ctx, flight)).To(Succeed())
			}
		})

		It("should scan all the inserted flights", func(ctx context.Context) {
			out, err := ddbc.Scan(ctx, &dynamodb.ScanInput{TableName: &tblname})
			Expect(err).ToNot(HaveOccurred())

			var ffs []*modelv2.FlightFares
			for _, item := range out.Items {
				var x modelv2.FlightFares
				Expect(x.UnmarshalDynamoItem(item)).To(Succeed())
				ffs = append(ffs, &x)
			}
			Expect(ffs).To(HaveLen(4))
		})

		It("should show flights to SFO from JFK in 2021", func(ctx context.Context) {
			out, err := qry.FlightsToFromInYear(ctx, &modelv2.FlightsToFromInYearRequest{
				To: modelv2.Airport_AIRPORT_SFO, From: modelv2.Airport_AIRPORT_JFK, Year: 2021,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(out.Flights).To(HaveLen(2))
			Expect(out.Flights[0].Number).To(BeNumerically("==", 260))
			Expect(out.Flights[1].Number).To(BeNumerically("==", 160))
		})
	})
})
