package messagev1

import (
	"fmt"
	attributevalue "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"
)

// file_example_message_v1_message_proto_marshal_dynamo_item marshals into DynamoDB attribute value maps
func file_example_message_v1_message_proto_marshal_dynamo_item(x proto.Message) (a types.AttributeValue, err error) {
	if mx, ok := x.(interface {
		MarshalDynamoItem() (map[string]types.AttributeValue, error)
	}); ok {
		mm, err := mx.MarshalDynamoItem()
		return &types.AttributeValueMemberM{Value: mm}, err
	}
	switch xt := x.(type) {
	case *durationpb.Duration, *timestamppb.Timestamp:
		xjson, err := protojson.Marshal(xt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal duration: %w", err)
		}
		xjsons, err := strconv.Unquote(string(xjson))
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value: %w", err)
		}
		return &types.AttributeValueMemberS{Value: xjsons}, nil
	case *anypb.Any:
		mv := &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}}
		mv.Value["1"], err = attributevalue.Marshal(xt.TypeUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's TypeURL field: %w", err)
		}
		mv.Value["2"], err = attributevalue.Marshal(xt.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Any's Value field: %w", err)
		}
		return mv, nil
	case *fieldmaskpb.FieldMask:
		return &types.AttributeValueMemberSS{Value: xt.Paths}, nil
	case *structpb.Value:
		return attributevalue.Marshal(xt.AsInterface())
	case *wrapperspb.StringValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.Marshal(xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.Marshal(xt.Value)
	default:
		return nil, fmt.Errorf("marshal of message type unsupported: %+T", xt)
	}
}

// file_example_message_v1_message_proto_marshal_dynamo_item unmarshals DynamoDB attribute value maps
func file_example_message_v1_message_proto_unmarshal_dynamo_item(m types.AttributeValue, x proto.Message) (err error) {
	if mx, ok := x.(interface {
		UnmarshalDynamoItem(map[string]types.AttributeValue) error
	}); ok {
		mm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal: no map attribute provided")
		}
		return mx.UnmarshalDynamoItem(mm.Value)
	}
	switch xt := x.(type) {
	case *durationpb.Duration, *timestamppb.Timestamp:
		ms, ok := m.(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no string attribute provided")
		}
		return protojson.Unmarshal([]byte(strconv.Quote(ms.Value)), x)
	case *anypb.Any:
		mm, ok := m.(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no map attribute provided")
		}
		err = attributevalue.Unmarshal(mm.Value["1"], &xt.TypeUrl)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's TypeURL field: %w", err)
		}
		err = attributevalue.Unmarshal(mm.Value["2"], &xt.Value)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Any's Value field: %w", err)
		}
		return nil
	case *fieldmaskpb.FieldMask:
		ss, ok := m.(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("failed to unmarshal duration: no string set attribute provided")
		}
		xt.Paths = ss.Value
		return nil
	case *structpb.Value:
		var vv any
		switch m.(type) {
		case *types.AttributeValueMemberL:
			vx := []any{}
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberM:
			vx := map[string]any{}
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberS:
			var vx string
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberBOOL:
			var vx bool
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberN:
			var vx float64
			err = attributevalue.Unmarshal(m, &vx)
			vv = vx
		case *types.AttributeValueMemberNULL:
			sv, _ := structpb.NewValue(nil)
			*xt = *sv
			return nil
		default:
			return fmt.Errorf("failed to unmarshal struct value: unsupported attribute value")
		}
		if err != nil {
			return fmt.Errorf("failed to unmarshal structpb Value field: %w", err)
		}
		sv, err := structpb.NewValue(vv)
		if err != nil {
			return fmt.Errorf("failed to init structpb value: %w", err)
		}
		*xt = *sv
		return nil
	case *wrapperspb.StringValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.BoolValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.BytesValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.DoubleValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.FloatValue:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.Int32Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.Int64Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.UInt32Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	case *wrapperspb.UInt64Value:
		return attributevalue.Unmarshal(m, &xt.Value)
	default:
		return fmt.Errorf("unmarshal of message type unsupported: %+T", xt)
	}
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Engine) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Brand != "" {
		m["1"], err = attributevalue.Marshal(x.GetBrand())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
		}
	}
	if x.Dirtyness != 0 {
		m["2"], err = attributevalue.Marshal(x.GetDirtyness())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Dirtyness': %w", err)
		}
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Engine) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["1"], &x.Brand)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Brand': %w", err)
	}
	err = attributevalue.Unmarshal(m["2"], &x.Dirtyness)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Dirtyness': %w", err)
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Car) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Engine != nil {
		m1, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetEngine())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Engine': %w", err)
		}
		m["1"] = m1
	}
	if x.Name != "" {
		m["2"], err = attributevalue.Marshal(x.GetName())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Name': %w", err)
		}
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Car) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	if m["1"] != nil {
		x.Engine = new(Engine)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["1"], x.Engine)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Engine': %w", err)
		}
	}
	err = attributevalue.Unmarshal(m["2"], &x.Name)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Name': %w", err)
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Appliance) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Brand != "" {
		m["1"], err = attributevalue.Marshal(x.GetBrand())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
		}
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Appliance) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["1"], &x.Brand)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Brand': %w", err)
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Kitchen) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Brand != "" {
		m["1"], err = attributevalue.Marshal(x.GetBrand())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
		}
	}
	if x.IsRenovated != false {
		m["2"], err = attributevalue.Marshal(x.GetIsRenovated())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'IsRenovated': %w", err)
		}
	}
	if x.QrCode != nil {
		m["3"], err = attributevalue.Marshal(x.GetQrCode())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'QrCode': %w", err)
		}
	}
	if x.NumSmallKnifes != 0 {
		m["4"], err = attributevalue.Marshal(x.GetNumSmallKnifes())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumSmallKnifes': %w", err)
		}
	}
	if x.NumSharpKnifes != 0 {
		m["5"], err = attributevalue.Marshal(x.GetNumSharpKnifes())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumSharpKnifes': %w", err)
		}
	}
	if x.NumBluntKnifes != 0 {
		m["6"], err = attributevalue.Marshal(x.GetNumBluntKnifes())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumBluntKnifes': %w", err)
		}
	}
	if x.NumSmallForks != 0 {
		m["7"], err = attributevalue.Marshal(x.GetNumSmallForks())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumSmallForks': %w", err)
		}
	}
	if x.NumMediumForks != 0 {
		m["8"], err = attributevalue.Marshal(x.GetNumMediumForks())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumMediumForks': %w", err)
		}
	}
	if x.NumLargeForks != 0 {
		m["9"], err = attributevalue.Marshal(x.GetNumLargeForks())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'NumLargeForks': %w", err)
		}
	}
	if x.PercentBlackTiles != 0 {
		m["10"], err = attributevalue.Marshal(x.GetPercentBlackTiles())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'PercentBlackTiles': %w", err)
		}
	}
	if x.PercentWhiteTiles != 0 {
		m["11"], err = attributevalue.Marshal(x.GetPercentWhiteTiles())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'PercentWhiteTiles': %w", err)
		}
	}
	if x.Dirtyness != 0 {
		m["12"], err = attributevalue.Marshal(x.GetDirtyness())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Dirtyness': %w", err)
		}
	}
	if len(x.Furniture) != 0 {
		m13 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Furniture {
			mk := fmt.Sprintf("%d", k)
			if mk == "" {
				return nil, fmt.Errorf("failed to marshal map key of field 'Furniture': map key cannot be empty")
			}
			if v == nil {
				m13.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'Furniture': %w", err)
			}
			m13.Value[mk] = mv
		}
		m["13"] = m13
	}
	if len(x.Calendar) != 0 {
		m["14"], err = attributevalue.Marshal(x.GetCalendar())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Calendar': %w", err)
		}
	}
	if x.WasherEngine != nil {
		m15, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetWasherEngine())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'WasherEngine': %w", err)
		}
		m["15"] = m15
	}
	if x.ExtraKitchen != nil {
		m16, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetExtraKitchen())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'ExtraKitchen': %w", err)
		}
		m["16"] = m16
	}
	if x.Timer != nil {
		m17, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetTimer())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Timer': %w", err)
		}
		m["17"] = m17
	}
	if x.WallTime != nil {
		m18, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetWallTime())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'WallTime': %w", err)
		}
		m["18"] = m18
	}
	if len(x.ApplianceEngines) != 0 {
		m19 := &types.AttributeValueMemberL{}
		for k, v := range x.ApplianceEngines {
			if v == nil {
				m19.Value = append(m19.Value, &types.AttributeValueMemberNULL{Value: true})
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal item '%d' of field 'ApplianceEngines': %w", k, err)
			}
			m19.Value = append(m19.Value, mv)
		}
		m["19"] = m19
	}
	if len(x.OtherBrands) != 0 {
		m["20"], err = attributevalue.Marshal(x.GetOtherBrands())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OtherBrands': %w", err)
		}
	}
	if x.SomeAny != nil {
		m21, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetSomeAny())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'SomeAny': %w", err)
		}
		m["21"] = m21
	}
	if x.SomeMask != nil {
		m22, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetSomeMask())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'SomeMask': %w", err)
		}
		m["22"] = m22
	}
	if x.SomeValue != nil {
		m23, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetSomeValue())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'SomeValue': %w", err)
		}
		m["23"] = m23
	}
	if x.OptString != nil {
		m["24"], err = attributevalue.Marshal(x.GetOptString())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OptString': %w", err)
		}
	}
	if x.ValStr != nil {
		m25, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetValStr())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'ValStr': %w", err)
		}
		m["25"] = m25
	}
	if x.ValBytes != nil {
		m26, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetValBytes())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'ValBytes': %w", err)
		}
		m["26"] = m26
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Kitchen) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["1"], &x.Brand)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Brand': %w", err)
	}
	err = attributevalue.Unmarshal(m["2"], &x.IsRenovated)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'IsRenovated': %w", err)
	}
	err = attributevalue.Unmarshal(m["3"], &x.QrCode)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'QrCode': %w", err)
	}
	err = attributevalue.Unmarshal(m["4"], &x.NumSmallKnifes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumSmallKnifes': %w", err)
	}
	err = attributevalue.Unmarshal(m["5"], &x.NumSharpKnifes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumSharpKnifes': %w", err)
	}
	err = attributevalue.Unmarshal(m["6"], &x.NumBluntKnifes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumBluntKnifes': %w", err)
	}
	err = attributevalue.Unmarshal(m["7"], &x.NumSmallForks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumSmallForks': %w", err)
	}
	err = attributevalue.Unmarshal(m["8"], &x.NumMediumForks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumMediumForks': %w", err)
	}
	err = attributevalue.Unmarshal(m["9"], &x.NumLargeForks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'NumLargeForks': %w", err)
	}
	err = attributevalue.Unmarshal(m["10"], &x.PercentBlackTiles)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'PercentBlackTiles': %w", err)
	}
	err = attributevalue.Unmarshal(m["11"], &x.PercentWhiteTiles)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'PercentWhiteTiles': %w", err)
	}
	err = attributevalue.Unmarshal(m["12"], &x.Dirtyness)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Dirtyness': %w", err)
	}
	if m["13"] != nil {
		x.Furniture = make(map[int64]*Appliance)
		m13, ok := m["13"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'Furniture': no map attribute provided")
		}
		for k, v := range m13.Value {
			mk, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to unmarshal map key for field 'Furniture': %w", err)
			}
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.Furniture[int64(mk)] = nil
				continue
			}
			var mv Appliance
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Furniture': %w", err)
			}
			x.Furniture[int64(mk)] = &mv
		}
	}
	err = attributevalue.Unmarshal(m["14"], &x.Calendar)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Calendar': %w", err)
	}
	if m["15"] != nil {
		x.WasherEngine = new(Engine)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["15"], x.WasherEngine)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'WasherEngine': %w", err)
		}
	}
	if m["16"] != nil {
		x.ExtraKitchen = new(Kitchen)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["16"], x.ExtraKitchen)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'ExtraKitchen': %w", err)
		}
	}
	if m["17"] != nil {
		x.Timer = new(durationpb.Duration)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["17"], x.Timer)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Timer': %w", err)
		}
	}
	if m["18"] != nil {
		x.WallTime = new(timestamppb.Timestamp)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["18"], x.WallTime)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'WallTime': %w", err)
		}
	}
	if m["19"] != nil {
		m19, ok := m["19"].(*types.AttributeValueMemberL)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'ApplianceEngines': no list attribute provided")
		}
		for k, v := range m19.Value {
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.ApplianceEngines = append(x.ApplianceEngines, nil)
				continue
			}
			var mv Engine
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal item '%d' of field 'ApplianceEngines': %w", k, err)
			}
			x.ApplianceEngines = append(x.ApplianceEngines, &mv)
		}
	}
	err = attributevalue.Unmarshal(m["20"], &x.OtherBrands)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'OtherBrands': %w", err)
	}
	if m["21"] != nil {
		x.SomeAny = new(anypb.Any)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["21"], x.SomeAny)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'SomeAny': %w", err)
		}
	}
	if m["22"] != nil {
		x.SomeMask = new(fieldmaskpb.FieldMask)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["22"], x.SomeMask)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'SomeMask': %w", err)
		}
	}
	if m["23"] != nil {
		x.SomeValue = new(structpb.Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["23"], x.SomeValue)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'SomeValue': %w", err)
		}
	}
	err = attributevalue.Unmarshal(m["24"], &x.OptString)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'OptString': %w", err)
	}
	if m["25"] != nil {
		x.ValStr = new(wrapperspb.StringValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["25"], x.ValStr)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'ValStr': %w", err)
		}
	}
	if m["26"] != nil {
		x.ValBytes = new(wrapperspb.BytesValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["26"], x.ValBytes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'ValBytes': %w", err)
		}
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Empty) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Empty) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *MapGalore) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if len(x.Int64Int64) != 0 {
		m["1"], err = attributevalue.Marshal(x.GetInt64Int64())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Int64Int64': %w", err)
		}
	}
	if len(x.Uint64Uint64) != 0 {
		m["2"], err = attributevalue.Marshal(x.GetUint64Uint64())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Uint64Uint64': %w", err)
		}
	}
	if len(x.Fixed64Fixed64) != 0 {
		m["3"], err = attributevalue.Marshal(x.GetFixed64Fixed64())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Fixed64Fixed64': %w", err)
		}
	}
	if len(x.Sint64Sint64) != 0 {
		m["4"], err = attributevalue.Marshal(x.GetSint64Sint64())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Sint64Sint64': %w", err)
		}
	}
	if len(x.Sfixed64Sfixed64) != 0 {
		m["5"], err = attributevalue.Marshal(x.GetSfixed64Sfixed64())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Sfixed64Sfixed64': %w", err)
		}
	}
	if len(x.Int32Int32) != 0 {
		m["6"], err = attributevalue.Marshal(x.GetInt32Int32())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Int32Int32': %w", err)
		}
	}
	if len(x.Uint32Uint32) != 0 {
		m["7"], err = attributevalue.Marshal(x.GetUint32Uint32())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Uint32Uint32': %w", err)
		}
	}
	if len(x.Fixed32Fixed32) != 0 {
		m["8"], err = attributevalue.Marshal(x.GetFixed32Fixed32())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Fixed32Fixed32': %w", err)
		}
	}
	if len(x.Sint32Sint32) != 0 {
		m["9"], err = attributevalue.Marshal(x.GetSint32Sint32())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Sint32Sint32': %w", err)
		}
	}
	if len(x.Sfixed32Sfixed32) != 0 {
		m["10"], err = attributevalue.Marshal(x.GetSfixed32Sfixed32())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Sfixed32Sfixed32': %w", err)
		}
	}
	if len(x.Stringstring) != 0 {
		m["11"], err = attributevalue.Marshal(x.GetStringstring())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Stringstring': %w", err)
		}
	}
	if len(x.Boolbool) != 0 {
		m["12"], err = attributevalue.Marshal(x.GetBoolbool())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Boolbool': %w", err)
		}
	}
	if len(x.Stringbytes) != 0 {
		m["13"], err = attributevalue.Marshal(x.GetStringbytes())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Stringbytes': %w", err)
		}
	}
	if len(x.Stringdouble) != 0 {
		m["14"], err = attributevalue.Marshal(x.GetStringdouble())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Stringdouble': %w", err)
		}
	}
	if len(x.Stringfloat) != 0 {
		m["15"], err = attributevalue.Marshal(x.GetStringfloat())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Stringfloat': %w", err)
		}
	}
	if len(x.Stringduration) != 0 {
		m16 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Stringduration {
			mk := k
			if mk == "" {
				return nil, fmt.Errorf("failed to marshal map key of field 'Stringduration': map key cannot be empty")
			}
			if v == nil {
				m16.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'Stringduration': %w", err)
			}
			m16.Value[mk] = mv
		}
		m["16"] = m16
	}
	if len(x.Stringtimestamp) != 0 {
		m17 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Stringtimestamp {
			mk := k
			if mk == "" {
				return nil, fmt.Errorf("failed to marshal map key of field 'Stringtimestamp': map key cannot be empty")
			}
			if v == nil {
				m17.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'Stringtimestamp': %w", err)
			}
			m17.Value[mk] = mv
		}
		m["17"] = m17
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *MapGalore) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["1"], &x.Int64Int64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Int64Int64': %w", err)
	}
	err = attributevalue.Unmarshal(m["2"], &x.Uint64Uint64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Uint64Uint64': %w", err)
	}
	err = attributevalue.Unmarshal(m["3"], &x.Fixed64Fixed64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Fixed64Fixed64': %w", err)
	}
	err = attributevalue.Unmarshal(m["4"], &x.Sint64Sint64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Sint64Sint64': %w", err)
	}
	err = attributevalue.Unmarshal(m["5"], &x.Sfixed64Sfixed64)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Sfixed64Sfixed64': %w", err)
	}
	err = attributevalue.Unmarshal(m["6"], &x.Int32Int32)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Int32Int32': %w", err)
	}
	err = attributevalue.Unmarshal(m["7"], &x.Uint32Uint32)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Uint32Uint32': %w", err)
	}
	err = attributevalue.Unmarshal(m["8"], &x.Fixed32Fixed32)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Fixed32Fixed32': %w", err)
	}
	err = attributevalue.Unmarshal(m["9"], &x.Sint32Sint32)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Sint32Sint32': %w", err)
	}
	err = attributevalue.Unmarshal(m["10"], &x.Sfixed32Sfixed32)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Sfixed32Sfixed32': %w", err)
	}
	err = attributevalue.Unmarshal(m["11"], &x.Stringstring)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Stringstring': %w", err)
	}
	err = attributevalue.Unmarshal(m["12"], &x.Boolbool)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Boolbool': %w", err)
	}
	err = attributevalue.Unmarshal(m["13"], &x.Stringbytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Stringbytes': %w", err)
	}
	err = attributevalue.Unmarshal(m["14"], &x.Stringdouble)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Stringdouble': %w", err)
	}
	err = attributevalue.Unmarshal(m["15"], &x.Stringfloat)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Stringfloat': %w", err)
	}
	if m["16"] != nil {
		x.Stringduration = make(map[string]*durationpb.Duration)
		m16, ok := m["16"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'Stringduration': no map attribute provided")
		}
		for k, v := range m16.Value {
			mk := k
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.Stringduration[string(mk)] = nil
				continue
			}
			var mv durationpb.Duration
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Stringduration': %w", err)
			}
			x.Stringduration[string(mk)] = &mv
		}
	}
	if m["17"] != nil {
		x.Stringtimestamp = make(map[string]*timestamppb.Timestamp)
		m17, ok := m["17"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'Stringtimestamp': no map attribute provided")
		}
		for k, v := range m17.Value {
			mk := k
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.Stringtimestamp[string(mk)] = nil
				continue
			}
			var mv timestamppb.Timestamp
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Stringtimestamp': %w", err)
			}
			x.Stringtimestamp[string(mk)] = &mv
		}
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *ValueGalore) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.SomeValue != nil {
		m1, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetSomeValue())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'SomeValue': %w", err)
		}
		m["1"] = m1
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *ValueGalore) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	if m["1"] != nil {
		x.SomeValue = new(structpb.Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["1"], x.SomeValue)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'SomeValue': %w", err)
		}
	}
	return nil
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *FieldPresence) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	if x.Str != "" {
		m["str"], err = attributevalue.Marshal(x.GetStr())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Str': %w", err)
		}
	}
	if x.OptStr != nil {
		m["optStr"], err = attributevalue.Marshal(x.GetOptStr())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OptStr': %w", err)
		}
	}
	if x.Msg != nil {
		m3, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetMsg())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Msg': %w", err)
		}
		m["msg"] = m3
	}
	if x.OptMsg != nil {
		m4, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetOptMsg())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OptMsg': %w", err)
		}
		m["optMsg"] = m4
	}
	if len(x.StrList) != 0 {
		m["strList"], err = attributevalue.Marshal(x.GetStrList())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'StrList': %w", err)
		}
	}
	if len(x.MsgList) != 0 {
		m6 := &types.AttributeValueMemberL{}
		for k, v := range x.MsgList {
			if v == nil {
				m6.Value = append(m6.Value, &types.AttributeValueMemberNULL{Value: true})
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal item '%d' of field 'MsgList': %w", k, err)
			}
			m6.Value = append(m6.Value, mv)
		}
		m["msgList"] = m6
	}
	if len(x.StrMap) != 0 {
		m["strMap"], err = attributevalue.Marshal(x.GetStrMap())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'StrMap': %w", err)
		}
	}
	if len(x.MsgMap) != 0 {
		m8 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.MsgMap {
			mk := k
			if mk == "" {
				return nil, fmt.Errorf("failed to marshal map key of field 'MsgMap': map key cannot be empty")
			}
			if v == nil {
				m8.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_example_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'MsgMap': %w", err)
			}
			m8.Value[mk] = mv
		}
		m["msgMap"] = m8
	}
	if x.Enum != 0 {
		m["enum"], err = attributevalue.Marshal(x.GetEnum())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Enum': %w", err)
		}
	}
	if x.OptEnum != nil {
		m["optEnum"], err = attributevalue.Marshal(x.GetOptEnum())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OptEnum': %w", err)
		}
	}
	if onev, ok := x.Oo.(*FieldPresence_OneofStr); ok && onev != nil {
		m["oneofStr"], err = attributevalue.Marshal(x.GetOneofStr())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OneofStr': %w", err)
		}
	}
	if onev, ok := x.Oo.(*FieldPresence_OneofMsg); ok && onev != nil {
		m12, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetOneofMsg())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'OneofMsg': %w", err)
		}
		m["oneofMsg"] = m12
	}
	if x.StrVal != nil {
		m13, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetStrVal())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'StrVal': %w", err)
		}
		m["strVal"] = m13
	}
	if x.BoolVal != nil {
		m14, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetBoolVal())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'BoolVal': %w", err)
		}
		m["boolVal"] = m14
	}
	if x.BytesVal != nil {
		m15, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetBytesVal())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'BytesVal': %w", err)
		}
		m["bytesVal"] = m15
	}
	if x.DoubleVal != nil {
		m16, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetDoubleVal())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'DoubleVal': %w", err)
		}
		m["doubleVal"] = m16
	}
	if x.FloatVal != nil {
		m17, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetFloatVal())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'FloatVal': %w", err)
		}
		m["floatVal"] = m17
	}
	if x.Int32Val != nil {
		m18, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetInt32Val())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Int32Val': %w", err)
		}
		m["int32Val"] = m18
	}
	if x.Int64Val != nil {
		m19, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetInt64Val())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Int64Val': %w", err)
		}
		m["int64Val"] = m19
	}
	if x.Uint32Val != nil {
		m20, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetUint32Val())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Uint32Val': %w", err)
		}
		m["uint32Val"] = m20
	}
	if x.Uint64Val != nil {
		m21, err := file_example_message_v1_message_proto_marshal_dynamo_item(x.GetUint64Val())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Uint64Val': %w", err)
		}
		m["uint64Val"] = m21
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *FieldPresence) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	err = attributevalue.Unmarshal(m["str"], &x.Str)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Str': %w", err)
	}
	err = attributevalue.Unmarshal(m["optStr"], &x.OptStr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'OptStr': %w", err)
	}
	if m["msg"] != nil {
		x.Msg = new(Engine)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["msg"], x.Msg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Msg': %w", err)
		}
	}
	if m["optMsg"] != nil {
		x.OptMsg = new(Engine)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["optMsg"], x.OptMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OptMsg': %w", err)
		}
	}
	err = attributevalue.Unmarshal(m["strList"], &x.StrList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'StrList': %w", err)
	}
	if m["msgList"] != nil {
		m6, ok := m["msgList"].(*types.AttributeValueMemberL)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'MsgList': no list attribute provided")
		}
		for k, v := range m6.Value {
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.MsgList = append(x.MsgList, nil)
				continue
			}
			var mv Engine
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal item '%d' of field 'MsgList': %w", k, err)
			}
			x.MsgList = append(x.MsgList, &mv)
		}
	}
	err = attributevalue.Unmarshal(m["strMap"], &x.StrMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'StrMap': %w", err)
	}
	if m["msgMap"] != nil {
		x.MsgMap = make(map[string]*Engine)
		m8, ok := m["msgMap"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("failed to unmarshal field 'MsgMap': no map attribute provided")
		}
		for k, v := range m8.Value {
			mk := k
			if _, ok := v.(*types.AttributeValueMemberNULL); ok {
				x.MsgMap[string(mk)] = nil
				continue
			}
			var mv Engine
			err = file_example_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'MsgMap': %w", err)
			}
			x.MsgMap[string(mk)] = &mv
		}
	}
	err = attributevalue.Unmarshal(m["enum"], &x.Enum)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Enum': %w", err)
	}
	err = attributevalue.Unmarshal(m["optEnum"], &x.OptEnum)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'OptEnum': %w", err)
	}
	if m["oneofStr"] != nil {
		var mo FieldPresence_OneofStr
		err = attributevalue.Unmarshal(m["oneofStr"], &mo.OneofStr)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OneofStr': %w", err)
		}
		x.Oo = &mo
	}
	if m["oneofMsg"] != nil {
		var mo FieldPresence_OneofMsg
		mo.OneofMsg = new(Engine)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["oneofMsg"], mo.OneofMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'OneofMsg': %w", err)
		}
		x.Oo = &mo
	}
	if m["strVal"] != nil {
		x.StrVal = new(wrapperspb.StringValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["strVal"], x.StrVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'StrVal': %w", err)
		}
	}
	if m["boolVal"] != nil {
		x.BoolVal = new(wrapperspb.BoolValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["boolVal"], x.BoolVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'BoolVal': %w", err)
		}
	}
	if m["bytesVal"] != nil {
		x.BytesVal = new(wrapperspb.BytesValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["bytesVal"], x.BytesVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'BytesVal': %w", err)
		}
	}
	if m["doubleVal"] != nil {
		x.DoubleVal = new(wrapperspb.DoubleValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["doubleVal"], x.DoubleVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'DoubleVal': %w", err)
		}
	}
	if m["floatVal"] != nil {
		x.FloatVal = new(wrapperspb.FloatValue)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["floatVal"], x.FloatVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'FloatVal': %w", err)
		}
	}
	if m["int32Val"] != nil {
		x.Int32Val = new(wrapperspb.Int32Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["int32Val"], x.Int32Val)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Int32Val': %w", err)
		}
	}
	if m["int64Val"] != nil {
		x.Int64Val = new(wrapperspb.Int64Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["int64Val"], x.Int64Val)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Int64Val': %w", err)
		}
	}
	if m["uint32Val"] != nil {
		x.Uint32Val = new(wrapperspb.UInt32Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["uint32Val"], x.Uint32Val)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Uint32Val': %w", err)
		}
	}
	if m["uint64Val"] != nil {
		x.Uint64Val = new(wrapperspb.UInt64Value)
		err = file_example_message_v1_message_proto_unmarshal_dynamo_item(m["uint64Val"], x.Uint64Val)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Uint64Val': %w", err)
		}
	}
	return nil
}
