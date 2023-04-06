// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

// Package modelv2ddbpath holds generated code for working with Dynamo document paths
package modelv2ddbpath

import (
	expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbpath "github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	"reflect"
)

// FlightPath allows for constructing type-safe expression names
type FlightPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FlightPath) WithDynamoNameBuilder(n expression.NameBuilder) FlightPath {
	p.NameBuilder = n
	return p
}

// Number appends the path being build
func (p FlightPath) Number() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// DepatureAt appends the path being build
func (p FlightPath) DepatureAt() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// ArrivalAt appends the path being build
func (p FlightPath) ArrivalAt() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Origin appends the path being build
func (p FlightPath) Origin() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// Destination appends the path being build
func (p FlightPath) Destination() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// Class appends the path being build
func (p FlightPath) Class() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}

// IsSegment appends the path being build
func (p FlightPath) IsSegment() expression.NameBuilder {
	return p.AppendName(expression.Name("7"))
}

// SegmentId appends the path being build
func (p FlightPath) SegmentId() expression.NameBuilder {
	return p.AppendName(expression.Name("8"))
}

// Segments appends the path being build
func (p FlightPath) Segments() expression.NameBuilder {
	return p.AppendName(expression.Name("9"))
}
func init() {
	ddbpath.Register(FlightPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
		"4": {Kind: ddbpath.FieldKindSingle},
		"5": {Kind: ddbpath.FieldKindSingle},
		"6": {Kind: ddbpath.FieldKindSingle},
		"7": {Kind: ddbpath.FieldKindSingle},
		"8": {Kind: ddbpath.FieldKindSingle},
		"9": {Kind: ddbpath.FieldKindSingle},
	})
}

// FarePath allows for constructing type-safe expression names
type FarePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FarePath) WithDynamoNameBuilder(n expression.NameBuilder) FarePath {
	p.NameBuilder = n
	return p
}

// StartAt appends the path being build
func (p FarePath) StartAt() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// EndAt appends the path being build
func (p FarePath) EndAt() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// Origin appends the path being build
func (p FarePath) Origin() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Destination appends the path being build
func (p FarePath) Destination() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// Class appends the path being build
func (p FarePath) Class() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}
func init() {
	ddbpath.Register(FarePath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
		"4": {Kind: ddbpath.FieldKindSingle},
		"5": {Kind: ddbpath.FieldKindSingle},
	})
}

// AssignmentPath allows for constructing type-safe expression names
type AssignmentPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p AssignmentPath) WithDynamoNameBuilder(n expression.NameBuilder) AssignmentPath {
	p.NameBuilder = n
	return p
}

// Number appends the path being build
func (p AssignmentPath) Number() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// FlightNumber appends the path being build
func (p AssignmentPath) FlightNumber() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// SegmentId appends the path being build
func (p AssignmentPath) SegmentId() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Seat appends the path being build
func (p AssignmentPath) Seat() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// DepartureAt appends the path being build
func (p AssignmentPath) DepartureAt() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// ArrivalAt appends the path being build
func (p AssignmentPath) ArrivalAt() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}

// Origin appends the path being build
func (p AssignmentPath) Origin() expression.NameBuilder {
	return p.AppendName(expression.Name("7"))
}

// Destination appends the path being build
func (p AssignmentPath) Destination() expression.NameBuilder {
	return p.AppendName(expression.Name("8"))
}

// SpecialServiceRequests returns 'p' appended with the attribute name and allow indexing
func (p AssignmentPath) SpecialServiceRequests() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("9"))}
}

// FirstName appends the path being build
func (p AssignmentPath) FirstName() expression.NameBuilder {
	return p.AppendName(expression.Name("10"))
}

// LastName appends the path being build
func (p AssignmentPath) LastName() expression.NameBuilder {
	return p.AppendName(expression.Name("11"))
}
func init() {
	ddbpath.Register(AssignmentPath{}, map[string]ddbpath.FieldInfo{
		"1":  {Kind: ddbpath.FieldKindSingle},
		"10": {Kind: ddbpath.FieldKindSingle},
		"11": {Kind: ddbpath.FieldKindSingle},
		"2":  {Kind: ddbpath.FieldKindSingle},
		"3":  {Kind: ddbpath.FieldKindSingle},
		"4":  {Kind: ddbpath.FieldKindSingle},
		"5":  {Kind: ddbpath.FieldKindSingle},
		"6":  {Kind: ddbpath.FieldKindSingle},
		"7":  {Kind: ddbpath.FieldKindSingle},
		"8":  {Kind: ddbpath.FieldKindSingle},
		"9":  {Kind: ddbpath.FieldKindList},
	})
}

// BookingPath allows for constructing type-safe expression names
type BookingPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p BookingPath) WithDynamoNameBuilder(n expression.NameBuilder) BookingPath {
	p.NameBuilder = n
	return p
}

// FlightNumber appends the path being build
func (p BookingPath) FlightNumber() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// DepartureAt appends the path being build
func (p BookingPath) DepartureAt() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// ArrivalAt appends the path being build
func (p BookingPath) ArrivalAt() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Segments appends the path being build
func (p BookingPath) Segments() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// Origin appends the path being build
func (p BookingPath) Origin() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// Destination appends the path being build
func (p BookingPath) Destination() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}

// FirstName appends the path being build
func (p BookingPath) FirstName() expression.NameBuilder {
	return p.AppendName(expression.Name("10"))
}

// LastName appends the path being build
func (p BookingPath) LastName() expression.NameBuilder {
	return p.AppendName(expression.Name("11"))
}
func init() {
	ddbpath.Register(BookingPath{}, map[string]ddbpath.FieldInfo{
		"1":  {Kind: ddbpath.FieldKindSingle},
		"10": {Kind: ddbpath.FieldKindSingle},
		"11": {Kind: ddbpath.FieldKindSingle},
		"2":  {Kind: ddbpath.FieldKindSingle},
		"3":  {Kind: ddbpath.FieldKindSingle},
		"4":  {Kind: ddbpath.FieldKindSingle},
		"5":  {Kind: ddbpath.FieldKindSingle},
		"6":  {Kind: ddbpath.FieldKindSingle},
	})
}

// FlightFaresPath allows for constructing type-safe expression names
type FlightFaresPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FlightFaresPath) WithDynamoNameBuilder(n expression.NameBuilder) FlightFaresPath {
	p.NameBuilder = n
	return p
}

// Pk appends the path being build
func (p FlightFaresPath) Pk() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// Sk appends the path being build
func (p FlightFaresPath) Sk() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// Type appends the path being build
func (p FlightFaresPath) Type() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// Gsi1Pk appends the path being build
func (p FlightFaresPath) Gsi1Pk() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// Gsi1Sk appends the path being build
func (p FlightFaresPath) Gsi1Sk() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// Gsi2Pk appends the path being build
func (p FlightFaresPath) Gsi2Pk() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}

// Gsi2Sk appends the path being build
func (p FlightFaresPath) Gsi2Sk() expression.NameBuilder {
	return p.AppendName(expression.Name("7"))
}

// Flight returns 'p' with the attribute name appended and allow subselecting nested message
func (p FlightFaresPath) Flight() FlightPath {
	return FlightPath{NameBuilder: p.AppendName(expression.Name("100"))}
}

// Fare returns 'p' with the attribute name appended and allow subselecting nested message
func (p FlightFaresPath) Fare() FarePath {
	return FarePath{NameBuilder: p.AppendName(expression.Name("101"))}
}

// Assignment returns 'p' with the attribute name appended and allow subselecting nested message
func (p FlightFaresPath) Assignment() AssignmentPath {
	return AssignmentPath{NameBuilder: p.AppendName(expression.Name("102"))}
}

// Booking returns 'p' with the attribute name appended and allow subselecting nested message
func (p FlightFaresPath) Booking() BookingPath {
	return BookingPath{NameBuilder: p.AppendName(expression.Name("103"))}
}
func init() {
	ddbpath.Register(FlightFaresPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"100": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(FlightPath{}),
		},
		"101": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(FarePath{}),
		},
		"102": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(AssignmentPath{}),
		},
		"103": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(BookingPath{}),
		},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
		"4": {Kind: ddbpath.FieldKindSingle},
		"5": {Kind: ddbpath.FieldKindSingle},
		"6": {Kind: ddbpath.FieldKindSingle},
		"7": {Kind: ddbpath.FieldKindSingle},
	})
}

// FlightsToFromInYearRequestPath allows for constructing type-safe expression names
type FlightsToFromInYearRequestPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FlightsToFromInYearRequestPath) WithDynamoNameBuilder(n expression.NameBuilder) FlightsToFromInYearRequestPath {
	p.NameBuilder = n
	return p
}

// To appends the path being build
func (p FlightsToFromInYearRequestPath) To() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// From appends the path being build
func (p FlightsToFromInYearRequestPath) From() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// Year appends the path being build
func (p FlightsToFromInYearRequestPath) Year() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}
func init() {
	ddbpath.Register(FlightsToFromInYearRequestPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
	})
}

// FlightsToFromInYearResponsePath allows for constructing type-safe expression names
type FlightsToFromInYearResponsePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FlightsToFromInYearResponsePath) WithDynamoNameBuilder(n expression.NameBuilder) FlightsToFromInYearResponsePath {
	p.NameBuilder = n
	return p
}
func init() {
	ddbpath.Register(FlightsToFromInYearResponsePath{}, map[string]ddbpath.FieldInfo{})
}

// PassengerBookingsInYearRequestPath allows for constructing type-safe expression names
type PassengerBookingsInYearRequestPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p PassengerBookingsInYearRequestPath) WithDynamoNameBuilder(n expression.NameBuilder) PassengerBookingsInYearRequestPath {
	p.NameBuilder = n
	return p
}

// FirstName appends the path being build
func (p PassengerBookingsInYearRequestPath) FirstName() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// LastName appends the path being build
func (p PassengerBookingsInYearRequestPath) LastName() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// Year appends the path being build
func (p PassengerBookingsInYearRequestPath) Year() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}
func init() {
	ddbpath.Register(PassengerBookingsInYearRequestPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
		"3": {Kind: ddbpath.FieldKindSingle},
	})
}

// PassengerBookingsInYearResponsePath allows for constructing type-safe expression names
type PassengerBookingsInYearResponsePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p PassengerBookingsInYearResponsePath) WithDynamoNameBuilder(n expression.NameBuilder) PassengerBookingsInYearResponsePath {
	p.NameBuilder = n
	return p
}
func init() {
	ddbpath.Register(PassengerBookingsInYearResponsePath{}, map[string]ddbpath.FieldInfo{})
}