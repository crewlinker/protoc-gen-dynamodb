package modelv2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

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

// FlightFaresQuerier is generated from the flight fares access pattern definition
type FlightFaresQuerier struct {
	ap FlightFaresAccess
	cl DynamoQuerier
}

// FlightsToFromInYear returns flight from and to an airport in a given year
func (tbl *FlightFaresQuerier) FlightsToFromInYear(ctx context.Context, in *FlightsToFromInYearRequest) (out *FlightsToFromInYearResponse, err error) {
	var qryin dynamodb.QueryInput
	kb, fb, err := tbl.ap.FlightsToFromInYearExpr(in)
	if err != nil {
		return nil, fmt.Errorf("failed to setup expressions: %w", err)
	}

	expr, err := expression.NewBuilder().WithKeyCondition(kb).WithFilter(fb).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %w", err)
	}

	qryin.ExpressionAttributeNames = expr.Names()
	qryin.ExpressionAttributeValues = expr.Values()
	qryin.KeyConditionExpression = expr.KeyCondition()
	qryin.FilterExpression = expr.Filter()

	qryout, err := tbl.cl.Query(ctx, &qryin)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	out = &FlightsToFromInYearResponse{}
	for _, it := range qryout.Items {
		var x FlightFares
		if err = x.UnmarshalDynamoItem(it); err != nil {
			return nil, fmt.Errorf("failed to unmarshal queried item: %w", err)
		}

		if err := tbl.ap.FlightsToFromInYearOut(&x, out); err != nil {
			return nil, fmt.Errorf("failed to conver item into output: %w", err)
		}
	}

	return
}

// FlightFaresEntity needs to be implemented by entities to allow marshalling for the
// FlightFares table
type FlightFaresEntity interface {
	GetPk() string             // return type depends on protobuf type
	GetSk() string             // should only be generated when a sk is provided
	GetGsi1Pk() (string, bool) // can return false if it needs to be left undefined
	GetGsi1Sk() (string, bool)
	GetGsi2Pk() (string, bool)
	GetGsi2Sk() (string, bool)
	isFlightFares_Entity()
}

// FromDynamoEntity fills the flight vares message from an entity interface
func (x *FlightFares) FromDynamoEntity(e FlightFaresEntity) {
	x.Pk, x.Sk = e.GetPk(), e.GetSk()
	switch et := e.(type) {
	case *FlightFares_Fare:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_FARE
		x.Entity = et
	case *FlightFares_Flight:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_FLIGHT
		x.Entity = et
	case *FlightFares_Assignment:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_ASSIGNMENT
		x.Entity = et
	case *FlightFares_Booking:
		x.Type = FlightFareType_FLIGHT_FARE_TYPE_BOOKING
		x.Entity = et
	default:
		panic(fmt.Sprintf("unsupported entity: %T", et))
	}

	if kv, ok := e.GetGsi1Pk(); ok {
		x.Gsi1Pk = kv
		if kv, ok := e.GetGsi1Sk(); ok {
			x.Gsi1Sk = kv
		}
	}

	if kv, ok := e.GetGsi2Pk(); ok {
		x.Gsi2Pk = kv
		if kv, ok := e.GetGsi2Sk(); ok {
			x.Gsi2Sk = kv
		}
	}
}

// FlightFares Fare
var _ FlightFaresEntity = &FlightFares_Fare{}

func (x FlightFares_Fare) GetPk() string {
	return x.Fare.Origin.String() // e.g: DEN
}
func (x FlightFares_Fare) GetSk() string {
	return fmt.Sprintf("%s#%s#%s", x.Fare.Destination, x.Fare.StartAt, x.Fare.Class)
}
func (x FlightFares_Fare) GetGsi1Pk() (string, bool) {
	return "", false
}
func (x FlightFares_Fare) GetGsi1Sk() (string, bool) {
	return "", false
}
func (x FlightFares_Fare) GetGsi2Pk() (string, bool) {
	return "", false
}
func (x FlightFares_Fare) GetGsi2Sk() (string, bool) {
	return "", false
}

// FlightFares Flight
var _ FlightFaresEntity = &FlightFares_Flight{}

func (x FlightFares_Flight) GetPk() string {
	return x.Flight.Origin.String() // e.g: DEN
}
func (x FlightFares_Flight) GetSk() string {
	return fmt.Sprintf("%s#%s#%d#%d", // ${origin}#${depart}#${number}#${segId}
		x.Flight.Origin, x.Flight.DepatureAt, x.Flight.Number, x.Flight.SegmentId,
	)
}
func (x FlightFares_Flight) GetGsi1Pk() (string, bool) {
	return x.Flight.Destination.String(), true
}
func (x FlightFares_Flight) GetGsi1Sk() (string, bool) {
	return fmt.Sprintf("%s#%s", // ${origin}#${arrive}
		x.Flight.Origin, x.Flight.ArrivalAt,
	), true
}
func (x FlightFares_Flight) GetGsi2Pk() (string, bool) {
	return fmt.Sprintf("%d", x.Flight.Number), true
}
func (x FlightFares_Flight) GetGsi2Sk() (string, bool) {
	if !x.Flight.IsSegment {
		return "0", true
	}
	return fmt.Sprintf("%d", x.Flight.SegmentId), true
}

// FlightFares Assignment
var _ FlightFaresEntity = &FlightFares_Assignment{}

func (x FlightFares_Assignment) GetPk() string {
	return fmt.Sprintf("%s, %s", x.Assignment.LastName, x.Assignment.FirstName)
}
func (x FlightFares_Assignment) GetSk() string {
	return fmt.Sprintf("%s#%d#%d#%s", // ${depart}#${flight}#${segId}#${seat}
		x.Assignment.DepartureAt,
		x.Assignment.FlightNumber,
		x.Assignment.SegmentId,
		x.Assignment.Seat)
}
func (x FlightFares_Assignment) GetGsi1Pk() (string, bool) {
	return "", false
}
func (x FlightFares_Assignment) GetGsi1Sk() (string, bool) {
	return "", false
}
func (x FlightFares_Assignment) GetGsi2Pk() (string, bool) {
	return fmt.Sprintf("%d", x.Assignment.Number), true
}
func (x FlightFares_Assignment) GetGsi2Sk() (string, bool) {
	return fmt.Sprintf("%d#%s", x.Assignment.SegmentId, x.Assignment.Seat), true
}

// FlightFares Booking
var _ FlightFaresEntity = &FlightFares_Booking{}

func (x FlightFares_Booking) GetPk() string {
	return fmt.Sprintf("%s, %s", x.Booking.LastName, x.Booking.FirstName)
}
func (x FlightFares_Booking) GetSk() string {
	return fmt.Sprintf("%s#%d", // ${depart}#${flight}
		x.Booking.DepartureAt,
		x.Booking.FlightNumber,
	)
}
func (x FlightFares_Booking) GetGsi1Pk() (string, bool) {
	return "", false
}
func (x FlightFares_Booking) GetGsi1Sk() (string, bool) {
	return "", false
}
func (x FlightFares_Booking) GetGsi2Pk() (string, bool) {
	return "", false
}
func (x FlightFares_Booking) GetGsi2Sk() (string, bool) {
	return "", false
}
