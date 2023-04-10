package modelv2_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	ddbconversions "github.com/aereal/go-dynamodb-attribute-conversions/v2"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtable"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtest"
	modelv2 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("stream parsing", func() {
	It("should parse dynamodb streams", func() {
		ev1 := &events.DynamoDBEvent{
			Records: []events.DynamoDBEventRecord{
				{Change: events.DynamoDBStreamRecord{
					NewImage: map[string]events.DynamoDBAttributeValue{
						"1": events.NewStringAttribute("foo"),
						"2": events.NewStringAttribute("bar"),
						"100": events.NewMapAttribute(map[string]events.DynamoDBAttributeValue{
							"1": events.NewNumberAttribute("100"),
						}),
					},
				}},
			},
		}
		for _, rec := range ev1.Records {
			m := ddbconversions.AttributeValueMapFrom(rec.Change.NewImage)

			var x modelv2.FlightFares
			Expect(x.UnmarshalDynamoItem(m)).To(Succeed())
			Expect(x.Pk).To(Equal("foo"))
			Expect(x.Sk).To(Equal("bar"))
			Expect(x.GetFlight().Number).To(BeNumerically("==", 100))
		}
	})
})

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
		var rng [8]byte
		rand.Read(rng[:])

		tblname = fmt.Sprintf("flight_fares_%x", rng)
		mdl = &modelv2.FlightFaresModel{}
		mut = modelv2.NewFlightFaresMutater(tblname, ddbc, mdl)
		qry = modelv2.NewFlightFaresQuerier(tblname, ddbc, mdl)

		Expect(ddbc.CreateTable(ctx,
			ddbtable.TableCreate("flight_fares", ddbtable.WithTableName(tblname)),
		)).ToNot(BeNil())

		DeferCleanup(func(ctx context.Context) {
			Expect(ddbc.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &tblname})).ToNot(BeNil())
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
				Expect(mut.PutEntity(ctx, &modelv2.FlightFares_Flight{flight})).To(Succeed())
			}

			for _, booking := range []*modelv2.Booking{
				{
					FirstName: "Jon", LastName: "Smith", FlightNumber: 150,
					DepartureAt: timestamppb.New(time.Unix(1627791955, 0)), // 2021-08-01T06:25:00
					ArrivalAt:   timestamppb.New(time.Unix(1627820755, 0)), //2021-08-01T14:25:00
					Segments:    2,
					Origin:      modelv2.Airport_AIRPORT_SFO,
					Destination: modelv2.Airport_AIRPORT_JFK,
				},

				{
					FirstName: "Caren", LastName: "Storal", FlightNumber: 110,
					DepartureAt: timestamppb.New(time.Unix(1627791955, 0)), // 2021-08-01T06:25:00
					ArrivalAt:   timestamppb.New(time.Unix(1627820755, 0)), //2021-08-01T14:25:00
					Segments:    0,
					Origin:      modelv2.Airport_AIRPORT_JFK,
					Destination: modelv2.Airport_AIRPORT_SFO,
				},
			} {
				Expect(mut.PutEntity(ctx, &modelv2.FlightFares_Booking{booking})).To(Succeed())
			}

			for _, fare := range []*modelv2.Fare{
				{
					StartAt:     timestamppb.New(time.Unix(1627768800, 0)), // 2021-08-01T00:00:00
					EndAt:       timestamppb.New(time.Unix(1627851600, 0)), // 2021-08-01T23:00:00
					Origin:      modelv2.Airport_AIRPORT_JFK,
					Destination: modelv2.Airport_AIRPORT_SFO,
					Class:       modelv2.FlightClass_FLIGHT_CLASS_NON_STOP,
				},
				{
					StartAt:     timestamppb.New(time.Unix(1627768800, 0)), // 2021-08-01T00:00:00
					EndAt:       timestamppb.New(time.Unix(1627851600, 0)), // 2021-08-01T23:00:00
					Origin:      modelv2.Airport_AIRPORT_JFK,
					Destination: modelv2.Airport_AIRPORT_SFO,
					Class:       modelv2.FlightClass_FLIGHT_CLASS_DIRECT,
				},
				{
					StartAt:     timestamppb.New(time.Unix(1627768800, 0)), // 2021-08-01T00:00:00
					EndAt:       timestamppb.New(time.Unix(1627851600, 0)), // 2021-08-01T23:00:00
					Origin:      modelv2.Airport_AIRPORT_DEN,
					Destination: modelv2.Airport_AIRPORT_JFK,
					Class:       modelv2.FlightClass_FLIGHT_CLASS_NON_STOP,
				},
			} {
				Expect(mut.PutEntity(ctx, &modelv2.FlightFares_Fare{fare})).To(Succeed())
			}
		})

		It("should scan all the inserted items", func(ctx context.Context) {
			out, err := ddbc.Scan(ctx, &dynamodb.ScanInput{TableName: &tblname})
			Expect(err).ToNot(HaveOccurred())

			var ffs []*modelv2.FlightFares
			for _, item := range out.Items {
				var x modelv2.FlightFares
				Expect(x.UnmarshalDynamoItem(item)).To(Succeed())
				ffs = append(ffs, &x)
			}
			Expect(ffs).To(HaveLen(9))
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

		It("should show passengers bookings in 2021", func(ctx context.Context) {
			out, err := qry.PassengerBookingsInYear(ctx, &modelv2.PassengerBookingsInYearRequest{
				FirstName: "Jon", LastName: "Smith", Year: 2021,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(out.Bookings).To(HaveLen(1))
			Expect(out.Bookings[0].LastName).To(Equal("Smith"))
			Expect(out.Bookings[0].FirstName).To(Equal("Jon"))
		})

		It("should show fares from JFK to SFO", func(ctx context.Context) {
			out, err := qry.FaresFromTo(ctx, &modelv2.FaresFromToRequest{
				From: modelv2.Airport_AIRPORT_JFK, To: modelv2.Airport_AIRPORT_SFO,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(out.Fares).To(HaveLen(2))
			Expect(out.Fares[0].Class).To(Equal(modelv2.FlightClass_FLIGHT_CLASS_DIRECT))
			Expect(out.Fares[1].Class).To(Equal(modelv2.FlightClass_FLIGHT_CLASS_NON_STOP))
		})
	})
})
