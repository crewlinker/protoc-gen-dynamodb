package modelv2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	modelv2ddbpath "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v2/ddbpath"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// FormatTimestamp formats a timestamp for key construction
func FormatTimestamp(dt *timestamppb.Timestamp) string {
	return dt.AsTime().Format(time.RFC3339)
}

// FlightFaresModel holds all domain specific logic for the single page design of
type FlightFaresModel struct{}

// MapFare maps a assignment record onto table keys
func (m FlightFaresModel) MapFare(v *Fare) (km FlightFaresKeys, err error) {
	km.PK = v.Origin.String()
	km.SK = fmt.Sprintf("%s#%s#%s", v.Destination, FormatTimestamp(v.StartAt), v.Class)
	return
}

// MapFlight maps a assignment record onto table keys
func (m FlightFaresModel) MapFlight(v *Flight) (km FlightFaresKeys, err error) {
	km.PK = v.Origin.String()
	km.SK = fmt.Sprintf("%s#%s#%d#%d", // ${origin}#${depart}#${number}#${segId}
		v.Origin, FormatTimestamp(v.DepatureAt), v.Number, v.SegmentId)
	km.Gsi1Pk = aws.String(v.Destination.String())
	km.Gsi1Sk = aws.String(fmt.Sprintf("%s#%s", // ${origin}#${arrive}
		v.Origin, FormatTimestamp(v.ArrivalAt),
	))
	km.Gsi2Pk = aws.String(fmt.Sprintf("%d", v.Number))
	km.Gsi2Sk = aws.String(fmt.Sprintf("%d", v.SegmentId))
	return
}

// MapAssignment maps a assignment record onto table keys
func (m FlightFaresModel) MapAssignment(v *Assignment) (km FlightFaresKeys, err error) {
	km.PK = fmt.Sprintf("%s, %s", v.LastName, v.FirstName)
	km.SK = fmt.Sprintf("%s#%d#%d#%s", // ${depart}#${flight}#${segId}#${seat}
		FormatTimestamp(v.DepartureAt),
		v.FlightNumber,
		v.SegmentId,
		v.Seat)
	km.Gsi2Pk = aws.String(fmt.Sprintf("%d", v.Number))
	km.Gsi2Sk = aws.String(fmt.Sprintf("%d#%s", v.SegmentId, v.Seat))
	return
}

// MapBooking maps a booking record onto table keys
func (m FlightFaresModel) MapBooking(v *Booking) (km FlightFaresKeys, err error) {
	km.PK = fmt.Sprintf("%s, %s", v.LastName, v.FirstName)
	km.SK = fmt.Sprintf("%s#%d", // ${depart}#${flight}
		FormatTimestamp(v.DepartureAt),
		v.FlightNumber,
	)
	return
}

// FlightsToFromInYearExpr implements the access pattern
func (m FlightFaresModel) FlightsToFromInYearExpr(in *FlightsToFromInYearRequest) (kb expression.KeyConditionBuilder, fb expression.ConditionBuilder, err error) {
	// @TODO we can provide the keyname builder for pk and sk to make expressions easier
	// @TODO if gsi only provides a pk, only that should be provided.
	// @TODO would be convenient if the right key/attr names where provided so we don't have to lookup "1"
	kb = expression.Key("4").Equal(expression.Value(in.To.String())).And(
		expression.Key("5").BeginsWith(fmt.Sprintf("%s#%d", in.From, in.Year)))
	fb = (modelv2ddbpath.FlightFaresPath{}).Type().Equal(expression.Value(FlightFareType_FLIGHT_FARE_TYPE_FLIGHT))
	return
}

// FlightsToFromInYearOut implements the access pattern
func (m FlightFaresModel) FlightsToFromInYearOut(x *FlightFares, out *FlightsToFromInYearResponse) (err error) {
	out.Flights = append(out.Flights, x.GetFlight())
	return
}

////////
/// Everything Below Will be Generated at some point
////////

// FlightFaresAccess interface must be implemented to support access patterns
type FlightFaresAccess interface {
	// this is implemented for the access pattern to turn typed access pattern input into expression for
	// the query operation.
	FlightsToFromInYearExpr(in *FlightsToFromInYearRequest) (kb expression.KeyConditionBuilder, fb expression.ConditionBuilder, err error)
	// this is implemented to allow query responses to be turned into typed output for the access pattern. It may be called
	// multiple times as a query might iterate over multiple items, each of wich supply something different to the output.
	FlightsToFromInYearOut(x *FlightFares, out *FlightsToFromInYearResponse) (err error)
}

// DynamoQuerier is provided by the dynamodb client of the v2 sdk
type DynamoQuerier interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// DynamoMutater provides mutating methods on dynamodb
type DynamoMutater interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// FlightFaresMutater provides the writing part of the flight fares table
type FlightFaresMutater struct {
	tn  string
	cl  DynamoMutater
	mpr FlightFaresKeyMapper
}

// NewFlightFaresMutater inits the mutating side of the table
func NewFlightFaresMutater(tn string, cl DynamoMutater, mpr FlightFaresKeyMapper) *FlightFaresMutater {
	return &FlightFaresMutater{tn: tn, cl: cl, mpr: mpr}
}

// PutFlight will put a flight in the table
func (m *FlightFaresMutater) PutFlight(ctx context.Context, x *Flight) (err error) {
	var tbx FlightFares

	// @TODO should perform checks on 'x' to make sure we don't insert invalid values
	// @TODO we could use the validation package for this.
	// @TODO we could check that all values required for the keys are set explicitely

	if err = tbx.FromDynamoEntity(&FlightFares_Flight{x}, m.mpr); err != nil {
		return fmt.Errorf("failed to create table message from entity message: %w", err)
	}

	var put dynamodb.PutItemInput
	put.TableName = &m.tn
	if put.Item, err = tbx.MarshalDynamoItem(); err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	if _, err = m.cl.PutItem(ctx, &put); err != nil {
		return fmt.Errorf("failed to put: %w", err)
	}

	return nil
}

// NewFlightFaresQuerier inits the querying side of the table
func NewFlightFaresQuerier(tn string, cl DynamoQuerier, ap FlightFaresAccess) *FlightFaresQuerier {
	return &FlightFaresQuerier{tn: tn, cl: cl, ap: ap}
}

// FlightFaresQuerier is generated from the flight fares access pattern definition
type FlightFaresQuerier struct {
	tn string
	ap FlightFaresAccess
	cl DynamoQuerier
}

// FlightsToFromInYear returns flight from and to an airport in a given year
func (q *FlightFaresQuerier) FlightsToFromInYear(ctx context.Context, in *FlightsToFromInYearRequest) (out *FlightsToFromInYearResponse, err error) {
	var qryin dynamodb.QueryInput
	kb, fb, err := q.ap.FlightsToFromInYearExpr(in)
	if err != nil {
		return nil, fmt.Errorf("failed to setup expressions: %w", err)
	}

	expr, err := expression.NewBuilder().WithKeyCondition(kb).WithFilter(fb).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	qryin.TableName = &q.tn
	qryin.IndexName = aws.String("gsi1")
	qryin.ExpressionAttributeNames = expr.Names()
	qryin.ExpressionAttributeValues = expr.Values()
	qryin.KeyConditionExpression = expr.KeyCondition()
	qryin.FilterExpression = expr.Filter()

	qryout, err := q.cl.Query(ctx, &qryin)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	out = &FlightsToFromInYearResponse{}
	for _, it := range qryout.Items {
		var x FlightFares
		if err = x.UnmarshalDynamoItem(it); err != nil {
			return nil, fmt.Errorf("failed to unmarshal queried item: %w", err)
		}

		if err := q.ap.FlightsToFromInYearOut(&x, out); err != nil {
			return nil, fmt.Errorf("failed to conver item into output: %w", err)
		}
	}

	return
}

// FlightFaresKeys describes key attributes the FlightFares table
type FlightFaresKeys struct {
	PK     string
	SK     string
	Gsi1Pk *string
	Gsi1Sk *string
	Gsi2Pk *string
	Gsi2Sk *string
}

// FlightFaresKeyMapper can be implemented to map entities onto key values
type FlightFaresKeyMapper interface {
	MapFare(v *Fare) (FlightFaresKeys, error)
	MapFlight(v *Flight) (FlightFaresKeys, error)
	MapAssignment(v *Assignment) (FlightFaresKeys, error)
	MapBooking(v *Booking) (FlightFaresKeys, error)
}

// FromDynamoEntity fills the flight vares message from an entity interface
func (x *FlightFares) FromDynamoEntity(e isFlightFares_Entity, m FlightFaresKeyMapper) (err error) {
	var keys FlightFaresKeys
	switch et := e.(type) {
	case *FlightFares_Fare:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_FARE
		x.Entity = et
		keys, err = m.MapFare(et.Fare)
	case *FlightFares_Flight:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_FLIGHT
		x.Entity = et
		keys, err = m.MapFlight(et.Flight)
	case *FlightFares_Assignment:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_ASSIGNMENT
		x.Entity = et
		keys, err = m.MapAssignment(et.Assignment)
	case *FlightFares_Booking:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_BOOKING
		x.Entity = et
		keys, err = m.MapBooking(et.Booking)
	default:
		return fmt.Errorf("unsupported entity: %T", et)
	}

	if err != nil {
		return fmt.Errorf("failed to map keys: %w", err)
	}

	// @TODO error if pk/sk key has zero values or keymap is otherwise invalid

	x.Pk, x.Sk = keys.PK, keys.SK
	if keys.Gsi1Pk != nil {
		x.Gsi1Pk = *keys.Gsi1Pk
		if keys.Gsi1Sk != nil {
			x.Gsi1Sk = *keys.Gsi1Sk
		}
	}

	if keys.Gsi2Pk != nil {
		x.Gsi2Pk = *keys.Gsi2Pk
		if keys.Gsi2Sk != nil {
			x.Gsi2Sk = *keys.Gsi2Sk
		}
	}

	return nil
}
