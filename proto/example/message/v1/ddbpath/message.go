// Code generated by protoc-gen-dynamodb. DO NOT EDIT.

// Package messagev1ddbpath holds generated code for working with Dynamo document paths
package messagev1ddbpath

import (
	expression "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddbpath "github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	"reflect"
)

// EnginePath allows for constructing type-safe expression names
type EnginePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p EnginePath) WithDynamoNameBuilder(n expression.NameBuilder) EnginePath {
	p.NameBuilder = n
	return p
}

// Brand appends the path being build
func (p EnginePath) Brand() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// Dirtyness appends the path being build
func (p EnginePath) Dirtyness() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}
func init() {
	ddbpath.Register(EnginePath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindSingle},
		"2": {Kind: ddbpath.FieldKindSingle},
	})
}

// CarPath allows for constructing type-safe expression names
type CarPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p CarPath) WithDynamoNameBuilder(n expression.NameBuilder) CarPath {
	p.NameBuilder = n
	return p
}

// Engine returns 'p' with the attribute name appended and allow subselecting nested message
func (p CarPath) Engine() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("1"))}
}

// NrOfWheels appends the path being build
func (p CarPath) NrOfWheels() expression.NameBuilder {
	return p.AppendName(expression.Name("ws"))
}

// Name appends the path being build
func (p CarPath) Name() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}
func init() {
	ddbpath.Register(CarPath{}, map[string]ddbpath.FieldInfo{
		"1": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"2":  {Kind: ddbpath.FieldKindSingle},
		"ws": {Kind: ddbpath.FieldKindSingle},
	})
}

// CarPartitionKey returns a key builder for the partition key
func CarPartitionKey() (v expression.KeyBuilder) {
	return expression.Key("ws")
}

// CarPartitionKeyName returns a name builder for the partition key
func CarPartitionKeyName() (v expression.NameBuilder) {
	return expression.Name("ws")
}

// Car returns a key builder for the partition key
func Car() CarPath {
	return CarPath{}
}

// CarKeyNames returns the attribute names of the partition and sort keys respectively
func CarKeyNames() (v []string) {
	v = append(v, "ws")
	return
}

// AppliancePath allows for constructing type-safe expression names
type AppliancePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p AppliancePath) WithDynamoNameBuilder(n expression.NameBuilder) AppliancePath {
	p.NameBuilder = n
	return p
}

// Brand appends the path being build
func (p AppliancePath) Brand() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}
func init() {
	ddbpath.Register(AppliancePath{}, map[string]ddbpath.FieldInfo{"1": {Kind: ddbpath.FieldKindSingle}})
}

// IgnoredPath allows for constructing type-safe expression names
type IgnoredPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p IgnoredPath) WithDynamoNameBuilder(n expression.NameBuilder) IgnoredPath {
	p.NameBuilder = n
	return p
}

// Visible appends the path being build
func (p IgnoredPath) Visible() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}
func init() {
	ddbpath.Register(IgnoredPath{}, map[string]ddbpath.FieldInfo{"4": {Kind: ddbpath.FieldKindSingle}})
}

// KitchenPath allows for constructing type-safe expression names
type KitchenPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p KitchenPath) WithDynamoNameBuilder(n expression.NameBuilder) KitchenPath {
	p.NameBuilder = n
	return p
}

// Brand appends the path being build
func (p KitchenPath) Brand() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

// IsRenovated appends the path being build
func (p KitchenPath) IsRenovated() expression.NameBuilder {
	return p.AppendName(expression.Name("2"))
}

// QrCode appends the path being build
func (p KitchenPath) QrCode() expression.NameBuilder {
	return p.AppendName(expression.Name("3"))
}

// NumSmallKnifes appends the path being build
func (p KitchenPath) NumSmallKnifes() expression.NameBuilder {
	return p.AppendName(expression.Name("4"))
}

// NumSharpKnifes appends the path being build
func (p KitchenPath) NumSharpKnifes() expression.NameBuilder {
	return p.AppendName(expression.Name("5"))
}

// NumBluntKnifes appends the path being build
func (p KitchenPath) NumBluntKnifes() expression.NameBuilder {
	return p.AppendName(expression.Name("6"))
}

// NumSmallForks appends the path being build
func (p KitchenPath) NumSmallForks() expression.NameBuilder {
	return p.AppendName(expression.Name("7"))
}

// NumMediumForks appends the path being build
func (p KitchenPath) NumMediumForks() expression.NameBuilder {
	return p.AppendName(expression.Name("8"))
}

// NumLargeForks appends the path being build
func (p KitchenPath) NumLargeForks() expression.NameBuilder {
	return p.AppendName(expression.Name("9"))
}

// PercentBlackTiles appends the path being build
func (p KitchenPath) PercentBlackTiles() expression.NameBuilder {
	return p.AppendName(expression.Name("10"))
}

// PercentWhiteTiles appends the path being build
func (p KitchenPath) PercentWhiteTiles() expression.NameBuilder {
	return p.AppendName(expression.Name("11"))
}

// Dirtyness appends the path being build
func (p KitchenPath) Dirtyness() expression.NameBuilder {
	return p.AppendName(expression.Name("12"))
}

// Furniture returns 'p' appended with the attribute while allow map keys on a nested message
func (p KitchenPath) Furniture() ddbpath.ItemMap[AppliancePath] {
	return ddbpath.ItemMap[AppliancePath]{NameBuilder: p.AppendName(expression.Name("13"))}
}

// Calendar returns 'p' appended with the attribute name and allow map keys to be specified
func (p KitchenPath) Calendar() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("14"))}
}

// WasherEngine returns 'p' with the attribute name appended and allow subselecting nested message
func (p KitchenPath) WasherEngine() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("15"))}
}

// ExtraKitchen returns 'p' with the attribute name appended and allow subselecting nested message
func (p KitchenPath) ExtraKitchen() KitchenPath {
	return KitchenPath{NameBuilder: p.AppendName(expression.Name("16"))}
}

// Timer appends the path being build
func (p KitchenPath) Timer() expression.NameBuilder {
	return p.AppendName(expression.Name("17"))
}

// WallTime appends the path being build
func (p KitchenPath) WallTime() expression.NameBuilder {
	return p.AppendName(expression.Name("18"))
}

// ApplianceEngines returns 'p' appended with the attribute while allow indexing a nested message
func (p KitchenPath) ApplianceEngines() ddbpath.ItemList[EnginePath] {
	return ddbpath.ItemList[EnginePath]{NameBuilder: p.AppendName(expression.Name("19"))}
}

// OtherBrands returns 'p' appended with the attribute name and allow indexing
func (p KitchenPath) OtherBrands() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("20"))}
}

// SomeAny returns 'p' with the attribute name appended and allow subselecting nested message
func (p KitchenPath) SomeAny() ddbpath.AnyPath {
	return ddbpath.AnyPath{NameBuilder: p.AppendName(expression.Name("21"))}
}

// SomeMask returns 'p' with the attribute name appended and allow subselecting nested message
func (p KitchenPath) SomeMask() ddbpath.FieldMaskPath {
	return ddbpath.FieldMaskPath{NameBuilder: p.AppendName(expression.Name("22"))}
}

// SomeValue returns 'p' with the attribute name appended and allow subselecting nested message
func (p KitchenPath) SomeValue() ddbpath.ValuePath {
	return ddbpath.ValuePath{NameBuilder: p.AppendName(expression.Name("23"))}
}

// OptString appends the path being build
func (p KitchenPath) OptString() expression.NameBuilder {
	return p.AppendName(expression.Name("24"))
}

// ValStr appends the path being build
func (p KitchenPath) ValStr() expression.NameBuilder {
	return p.AppendName(expression.Name("25"))
}

// ValBytes appends the path being build
func (p KitchenPath) ValBytes() expression.NameBuilder {
	return p.AppendName(expression.Name("26"))
}

// ListOfTs returns 'p' appended with the attribute name and allow indexing
func (p KitchenPath) ListOfTs() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("27"))}
}

// StringSet returns 'p' appended with the attribute name and allow indexing
func (p KitchenPath) StringSet() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("28"))}
}

// NumberSet returns 'p' appended with the attribute name and allow indexing
func (p KitchenPath) NumberSet() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("29"))}
}

// BytesSet returns 'p' appended with the attribute name and allow indexing
func (p KitchenPath) BytesSet() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("30"))}
}

// RepeatedAny returns 'p' appended with the attribute while allow indexing a nested message
func (p KitchenPath) RepeatedAny() ddbpath.ItemList[ddbpath.AnyPath] {
	return ddbpath.ItemList[ddbpath.AnyPath]{NameBuilder: p.AppendName(expression.Name("31"))}
}

// MappedAny returns 'p' appended with the attribute while allow map keys on a nested message
func (p KitchenPath) MappedAny() ddbpath.ItemMap[ddbpath.AnyPath] {
	return ddbpath.ItemMap[ddbpath.AnyPath]{NameBuilder: p.AppendName(expression.Name("32"))}
}

// RepeatedFmask returns 'p' appended with the attribute while allow indexing a nested message
func (p KitchenPath) RepeatedFmask() ddbpath.ItemList[ddbpath.FieldMaskPath] {
	return ddbpath.ItemList[ddbpath.FieldMaskPath]{NameBuilder: p.AppendName(expression.Name("33"))}
}

// MappedFmask returns 'p' appended with the attribute while allow map keys on a nested message
func (p KitchenPath) MappedFmask() ddbpath.ItemMap[ddbpath.FieldMaskPath] {
	return ddbpath.ItemMap[ddbpath.FieldMaskPath]{NameBuilder: p.AppendName(expression.Name("34"))}
}
func init() {
	ddbpath.Register(KitchenPath{}, map[string]ddbpath.FieldInfo{
		"1":  {Kind: ddbpath.FieldKindSingle},
		"10": {Kind: ddbpath.FieldKindSingle},
		"11": {Kind: ddbpath.FieldKindSingle},
		"12": {Kind: ddbpath.FieldKindSingle},
		"13": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(AppliancePath{}),
		},
		"14": {Kind: ddbpath.FieldKindMap},
		"15": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"16": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(KitchenPath{}),
		},
		"17": {Kind: ddbpath.FieldKindSingle},
		"18": {Kind: ddbpath.FieldKindSingle},
		"19": {
			Kind:    ddbpath.FieldKindList,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"2":  {Kind: ddbpath.FieldKindSingle},
		"20": {Kind: ddbpath.FieldKindList},
		"21": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(ddbpath.AnyPath{}),
		},
		"22": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(ddbpath.FieldMaskPath{}),
		},
		"23": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(ddbpath.ValuePath{}),
		},
		"24": {Kind: ddbpath.FieldKindSingle},
		"25": {Kind: ddbpath.FieldKindSingle},
		"26": {Kind: ddbpath.FieldKindSingle},
		"27": {Kind: ddbpath.FieldKindList},
		"28": {Kind: ddbpath.FieldKindList},
		"29": {Kind: ddbpath.FieldKindList},
		"3":  {Kind: ddbpath.FieldKindSingle},
		"30": {Kind: ddbpath.FieldKindList},
		"31": {
			Kind:    ddbpath.FieldKindList,
			Message: reflect.TypeOf(ddbpath.AnyPath{}),
		},
		"32": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(ddbpath.AnyPath{}),
		},
		"33": {
			Kind:    ddbpath.FieldKindList,
			Message: reflect.TypeOf(ddbpath.FieldMaskPath{}),
		},
		"34": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(ddbpath.FieldMaskPath{}),
		},
		"4": {Kind: ddbpath.FieldKindSingle},
		"5": {Kind: ddbpath.FieldKindSingle},
		"6": {Kind: ddbpath.FieldKindSingle},
		"7": {Kind: ddbpath.FieldKindSingle},
		"8": {Kind: ddbpath.FieldKindSingle},
		"9": {Kind: ddbpath.FieldKindSingle},
	})
}

// KitchenPartitionKey returns a key builder for the partition key
func KitchenPartitionKey() (v expression.KeyBuilder) {
	return expression.Key("1")
}

// KitchenPartitionKeyName returns a name builder for the partition key
func KitchenPartitionKeyName() (v expression.NameBuilder) {
	return expression.Name("1")
}

// Kitchen returns a key builder for the partition key
func Kitchen() KitchenPath {
	return KitchenPath{}
}

// KitchenSortKey returns a key builder for the sort key
func KitchenSortKey() (v expression.KeyBuilder) {
	return expression.Key("3")
}

// KitchenSortKeyName returns a name builder for the sort key
func KitchenSortKeyName() (v expression.NameBuilder) {
	return expression.Name("3")
}

// KitchenKeyNames returns the attribute names of the partition and sort keys respectively
func KitchenKeyNames() (v []string) {
	v = append(v, "1")
	v = append(v, "3")
	return
}

// EmptyPath allows for constructing type-safe expression names
type EmptyPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p EmptyPath) WithDynamoNameBuilder(n expression.NameBuilder) EmptyPath {
	p.NameBuilder = n
	return p
}
func init() {
	ddbpath.Register(EmptyPath{}, map[string]ddbpath.FieldInfo{})
}

// MapGalorePath allows for constructing type-safe expression names
type MapGalorePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p MapGalorePath) WithDynamoNameBuilder(n expression.NameBuilder) MapGalorePath {
	p.NameBuilder = n
	return p
}

// Int64Int64 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Int64Int64() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("1"))}
}

// Uint64Uint64 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Uint64Uint64() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("2"))}
}

// Fixed64Fixed64 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Fixed64Fixed64() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("3"))}
}

// Sint64Sint64 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Sint64Sint64() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("4"))}
}

// Sfixed64Sfixed64 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Sfixed64Sfixed64() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("5"))}
}

// Int32Int32 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Int32Int32() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("6"))}
}

// Uint32Uint32 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Uint32Uint32() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("7"))}
}

// Fixed32Fixed32 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Fixed32Fixed32() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("8"))}
}

// Sint32Sint32 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Sint32Sint32() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("9"))}
}

// Sfixed32Sfixed32 returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Sfixed32Sfixed32() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("10"))}
}

// Stringstring returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringstring() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("11"))}
}

// Boolbool returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Boolbool() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("12"))}
}

// Stringbytes returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringbytes() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("13"))}
}

// Stringdouble returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringdouble() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("14"))}
}

// Stringfloat returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringfloat() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("15"))}
}

// Stringduration returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringduration() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("16"))}
}

// Stringtimestamp returns 'p' appended with the attribute name and allow map keys to be specified
func (p MapGalorePath) Stringtimestamp() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("17"))}
}

// Boolengine returns 'p' appended with the attribute while allow map keys on a nested message
func (p MapGalorePath) Boolengine() ddbpath.ItemMap[EnginePath] {
	return ddbpath.ItemMap[EnginePath]{NameBuilder: p.AppendName(expression.Name("18"))}
}

// Uintengine returns 'p' appended with the attribute while allow map keys on a nested message
func (p MapGalorePath) Uintengine() ddbpath.ItemMap[EnginePath] {
	return ddbpath.ItemMap[EnginePath]{NameBuilder: p.AppendName(expression.Name("19"))}
}
func init() {
	ddbpath.Register(MapGalorePath{}, map[string]ddbpath.FieldInfo{
		"1":  {Kind: ddbpath.FieldKindMap},
		"10": {Kind: ddbpath.FieldKindMap},
		"11": {Kind: ddbpath.FieldKindMap},
		"12": {Kind: ddbpath.FieldKindMap},
		"13": {Kind: ddbpath.FieldKindMap},
		"14": {Kind: ddbpath.FieldKindMap},
		"15": {Kind: ddbpath.FieldKindMap},
		"16": {Kind: ddbpath.FieldKindMap},
		"17": {Kind: ddbpath.FieldKindMap},
		"18": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"19": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"2": {Kind: ddbpath.FieldKindMap},
		"3": {Kind: ddbpath.FieldKindMap},
		"4": {Kind: ddbpath.FieldKindMap},
		"5": {Kind: ddbpath.FieldKindMap},
		"6": {Kind: ddbpath.FieldKindMap},
		"7": {Kind: ddbpath.FieldKindMap},
		"8": {Kind: ddbpath.FieldKindMap},
		"9": {Kind: ddbpath.FieldKindMap},
	})
}

// ValueGalorePath allows for constructing type-safe expression names
type ValueGalorePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p ValueGalorePath) WithDynamoNameBuilder(n expression.NameBuilder) ValueGalorePath {
	p.NameBuilder = n
	return p
}

// SomeValue returns 'p' with the attribute name appended and allow subselecting nested message
func (p ValueGalorePath) SomeValue() ddbpath.ValuePath {
	return ddbpath.ValuePath{NameBuilder: p.AppendName(expression.Name("1"))}
}
func init() {
	ddbpath.Register(ValueGalorePath{}, map[string]ddbpath.FieldInfo{"1": {
		Kind:    ddbpath.FieldKindSingle,
		Message: reflect.TypeOf(ddbpath.ValuePath{}),
	}})
}

// FieldPresencePath allows for constructing type-safe expression names
type FieldPresencePath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p FieldPresencePath) WithDynamoNameBuilder(n expression.NameBuilder) FieldPresencePath {
	p.NameBuilder = n
	return p
}

// Str appends the path being build
func (p FieldPresencePath) Str() expression.NameBuilder {
	return p.AppendName(expression.Name("str"))
}

// OptStr appends the path being build
func (p FieldPresencePath) OptStr() expression.NameBuilder {
	return p.AppendName(expression.Name("optStr"))
}

// Msg returns 'p' with the attribute name appended and allow subselecting nested message
func (p FieldPresencePath) Msg() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("msg"))}
}

// OptMsg returns 'p' with the attribute name appended and allow subselecting nested message
func (p FieldPresencePath) OptMsg() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("optMsg"))}
}

// StrList returns 'p' appended with the attribute name and allow indexing
func (p FieldPresencePath) StrList() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("strList"))}
}

// MsgList returns 'p' appended with the attribute while allow indexing a nested message
func (p FieldPresencePath) MsgList() ddbpath.ItemList[EnginePath] {
	return ddbpath.ItemList[EnginePath]{NameBuilder: p.AppendName(expression.Name("msgList"))}
}

// StrMap returns 'p' appended with the attribute name and allow map keys to be specified
func (p FieldPresencePath) StrMap() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("strMap"))}
}

// MsgMap returns 'p' appended with the attribute while allow map keys on a nested message
func (p FieldPresencePath) MsgMap() ddbpath.ItemMap[EnginePath] {
	return ddbpath.ItemMap[EnginePath]{NameBuilder: p.AppendName(expression.Name("msgMap"))}
}

// Enum appends the path being build
func (p FieldPresencePath) Enum() expression.NameBuilder {
	return p.AppendName(expression.Name("enum"))
}

// OptEnum appends the path being build
func (p FieldPresencePath) OptEnum() expression.NameBuilder {
	return p.AppendName(expression.Name("optEnum"))
}

// OneofStr appends the path being build
func (p FieldPresencePath) OneofStr() expression.NameBuilder {
	return p.AppendName(expression.Name("oneofStr"))
}

// OneofMsg returns 'p' with the attribute name appended and allow subselecting nested message
func (p FieldPresencePath) OneofMsg() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("oneofMsg"))}
}

// StrVal appends the path being build
func (p FieldPresencePath) StrVal() expression.NameBuilder {
	return p.AppendName(expression.Name("strVal"))
}

// BoolVal appends the path being build
func (p FieldPresencePath) BoolVal() expression.NameBuilder {
	return p.AppendName(expression.Name("boolVal"))
}

// BytesVal appends the path being build
func (p FieldPresencePath) BytesVal() expression.NameBuilder {
	return p.AppendName(expression.Name("bytesVal"))
}

// DoubleVal appends the path being build
func (p FieldPresencePath) DoubleVal() expression.NameBuilder {
	return p.AppendName(expression.Name("doubleVal"))
}

// FloatVal appends the path being build
func (p FieldPresencePath) FloatVal() expression.NameBuilder {
	return p.AppendName(expression.Name("floatVal"))
}

// Int32Val appends the path being build
func (p FieldPresencePath) Int32Val() expression.NameBuilder {
	return p.AppendName(expression.Name("int32Val"))
}

// Int64Val appends the path being build
func (p FieldPresencePath) Int64Val() expression.NameBuilder {
	return p.AppendName(expression.Name("int64Val"))
}

// Uint32Val appends the path being build
func (p FieldPresencePath) Uint32Val() expression.NameBuilder {
	return p.AppendName(expression.Name("uint32Val"))
}

// Uint64Val appends the path being build
func (p FieldPresencePath) Uint64Val() expression.NameBuilder {
	return p.AppendName(expression.Name("uint64Val"))
}
func init() {
	ddbpath.Register(FieldPresencePath{}, map[string]ddbpath.FieldInfo{
		"boolVal":   {Kind: ddbpath.FieldKindSingle},
		"bytesVal":  {Kind: ddbpath.FieldKindSingle},
		"doubleVal": {Kind: ddbpath.FieldKindSingle},
		"enum":      {Kind: ddbpath.FieldKindSingle},
		"floatVal":  {Kind: ddbpath.FieldKindSingle},
		"int32Val":  {Kind: ddbpath.FieldKindSingle},
		"int64Val":  {Kind: ddbpath.FieldKindSingle},
		"msg": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"msgList": {
			Kind:    ddbpath.FieldKindList,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"msgMap": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"oneofMsg": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"oneofStr": {Kind: ddbpath.FieldKindSingle},
		"optEnum":  {Kind: ddbpath.FieldKindSingle},
		"optMsg": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"optStr":    {Kind: ddbpath.FieldKindSingle},
		"str":       {Kind: ddbpath.FieldKindSingle},
		"strList":   {Kind: ddbpath.FieldKindList},
		"strMap":    {Kind: ddbpath.FieldKindMap},
		"strVal":    {Kind: ddbpath.FieldKindSingle},
		"uint32Val": {Kind: ddbpath.FieldKindSingle},
		"uint64Val": {Kind: ddbpath.FieldKindSingle},
	})
}

// JsonFieldsPath allows for constructing type-safe expression names
type JsonFieldsPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p JsonFieldsPath) WithDynamoNameBuilder(n expression.NameBuilder) JsonFieldsPath {
	p.NameBuilder = n
	return p
}

// JsonStrList returns 'p' appended with the attribute name and allow indexing
func (p JsonFieldsPath) JsonStrList() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("1"))}
}

// JsonEngine returns 'p' with the attribute name appended and allow subselecting nested message
func (p JsonFieldsPath) JsonEngine() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("json_engine"))}
}

// JsonIntMap returns 'p' appended with the attribute name and allow map keys to be specified
func (p JsonFieldsPath) JsonIntMap() ddbpath.Map {
	return ddbpath.Map{NameBuilder: p.AppendName(expression.Name("4"))}
}

// JsonEngineList returns 'p' appended with the attribute while allow indexing a nested message
func (p JsonFieldsPath) JsonEngineList() ddbpath.ItemList[EnginePath] {
	return ddbpath.ItemList[EnginePath]{NameBuilder: p.AppendName(expression.Name("2"))}
}

// JsonEngineMap returns 'p' appended with the attribute while allow map keys on a nested message
func (p JsonFieldsPath) JsonEngineMap() ddbpath.ItemMap[EnginePath] {
	return ddbpath.ItemMap[EnginePath]{NameBuilder: p.AppendName(expression.Name("5"))}
}

// JsonNrSet returns 'p' appended with the attribute name and allow indexing
func (p JsonFieldsPath) JsonNrSet() ddbpath.List {
	return ddbpath.List{NameBuilder: p.AppendName(expression.Name("6"))}
}
func init() {
	ddbpath.Register(JsonFieldsPath{}, map[string]ddbpath.FieldInfo{
		"1": {Kind: ddbpath.FieldKindList},
		"2": {
			Kind:    ddbpath.FieldKindList,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"4": {Kind: ddbpath.FieldKindMap},
		"5": {
			Kind:    ddbpath.FieldKindMap,
			Message: reflect.TypeOf(EnginePath{}),
		},
		"6": {Kind: ddbpath.FieldKindList},
		"json_engine": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
	})
}

// JsonOneofsPath allows for constructing type-safe expression names
type JsonOneofsPath struct {
	expression.NameBuilder
}

// WithDynamoNameBuilder allows generic types to overwrite the path
func (p JsonOneofsPath) WithDynamoNameBuilder(n expression.NameBuilder) JsonOneofsPath {
	p.NameBuilder = n
	return p
}

// OneofStr appends the path being build
func (p JsonOneofsPath) OneofStr() expression.NameBuilder {
	return p.AppendName(expression.Name("7"))
}

// OneofMsg returns 'p' with the attribute name appended and allow subselecting nested message
func (p JsonOneofsPath) OneofMsg() EnginePath {
	return EnginePath{NameBuilder: p.AppendName(expression.Name("8"))}
}
func init() {
	ddbpath.Register(JsonOneofsPath{}, map[string]ddbpath.FieldInfo{
		"7": {Kind: ddbpath.FieldKindSingle},
		"8": {
			Kind:    ddbpath.FieldKindSingle,
			Message: reflect.TypeOf(EnginePath{}),
		},
	})
}
