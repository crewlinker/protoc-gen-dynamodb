package messagev1

import (
	"fmt"
	attributevalue "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

// file_message_v1_message_proto_marshal_dynamo_item marshals into DynamoDB attribute value maps
func file_message_v1_message_proto_marshal_dynamo_item(x proto.Message) (types.AttributeValue, error) {
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
			return nil, fmt.Errorf("failed to unquote marshalled duration: %w", err)
		}
		return &types.AttributeValueMemberS{Value: xjsons}, nil
	default:
		return nil, fmt.Errorf("marshal of message type unsupported: %+T", xt)
	}
}

// file_message_v1_message_proto_marshal_dynamo_item unmarshals DynamoDB attribute value maps
func file_message_v1_message_proto_unmarshal_dynamo_item(m types.AttributeValue, x proto.Message) error {
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
	default:
		return fmt.Errorf("unmarshal of message type unsupported: %+T", xt)
	}
}

// MarshalDynamoItem marshals dat into a dynamodb attribute map
func (x *Engine) MarshalDynamoItem() (m map[string]types.AttributeValue, err error) {
	m = make(map[string]types.AttributeValue)
	m["1"], err = attributevalue.Marshal(x.Brand)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
	}
	m["2"], err = attributevalue.Marshal(x.Dirtyness)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Dirtyness': %w", err)
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
		m1, err := file_message_v1_message_proto_marshal_dynamo_item(x.Engine)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Engine': %w", err)
		}
		m["1"] = m1
	}
	m["2"], err = attributevalue.Marshal(x.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Name': %w", err)
	}
	return m, nil
}

// UnmarshalDynamoItem unmarshals data from a dynamodb attribute map
func (x *Car) UnmarshalDynamoItem(m map[string]types.AttributeValue) (err error) {
	if m["1"] != nil {
		x.Engine = new(Engine)
		err = file_message_v1_message_proto_unmarshal_dynamo_item(m["1"], x.Engine)
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
	m["1"], err = attributevalue.Marshal(x.Brand)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
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
	m["1"], err = attributevalue.Marshal(x.Brand)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Brand': %w", err)
	}
	m["2"], err = attributevalue.Marshal(x.IsRenovated)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'IsRenovated': %w", err)
	}
	m["3"], err = attributevalue.Marshal(x.QrCode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'QrCode': %w", err)
	}
	m["4"], err = attributevalue.Marshal(x.NumSmallKnifes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumSmallKnifes': %w", err)
	}
	m["5"], err = attributevalue.Marshal(x.NumSharpKnifes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumSharpKnifes': %w", err)
	}
	m["6"], err = attributevalue.Marshal(x.NumBluntKnifes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumBluntKnifes': %w", err)
	}
	m["7"], err = attributevalue.Marshal(x.NumSmallForks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumSmallForks': %w", err)
	}
	m["8"], err = attributevalue.Marshal(x.NumMediumForks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumMediumForks': %w", err)
	}
	m["9"], err = attributevalue.Marshal(x.NumLargeForks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'NumLargeForks': %w", err)
	}
	m["10"], err = attributevalue.Marshal(x.PercentBlackTiles)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'PercentBlackTiles': %w", err)
	}
	m["11"], err = attributevalue.Marshal(x.PercentWhiteTiles)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'PercentWhiteTiles': %w", err)
	}
	m["12"], err = attributevalue.Marshal(x.Dirtyness)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Dirtyness': %w", err)
	}
	if x.Furniture != nil {
		m13 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Furniture {
			mk := fmt.Sprintf("%d", k)
			if v == nil {
				m13.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'Furniture': %w", err)
			}
			m13.Value[mk] = mv
		}
		m["13"] = m13
	}
	m["14"], err = attributevalue.Marshal(x.Calendar)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Calendar': %w", err)
	}
	if x.WasherEngine != nil {
		m15, err := file_message_v1_message_proto_marshal_dynamo_item(x.WasherEngine)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'WasherEngine': %w", err)
		}
		m["15"] = m15
	}
	if x.ExtraKitchen != nil {
		m16, err := file_message_v1_message_proto_marshal_dynamo_item(x.ExtraKitchen)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'ExtraKitchen': %w", err)
		}
		m["16"] = m16
	}
	if x.Timer != nil {
		m17, err := file_message_v1_message_proto_marshal_dynamo_item(x.Timer)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'Timer': %w", err)
		}
		m["17"] = m17
	}
	if x.WallTime != nil {
		m18, err := file_message_v1_message_proto_marshal_dynamo_item(x.WallTime)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field 'WallTime': %w", err)
		}
		m["18"] = m18
	}
	m19 := &types.AttributeValueMemberL{}
	for k, v := range x.ApplianceEngines {
		if v == nil {
			m19.Value = append(m19.Value, &types.AttributeValueMemberNULL{Value: true})
			continue
		}
		mv, err := file_message_v1_message_proto_marshal_dynamo_item(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item '%d' of field 'ApplianceEngines': %w", k, err)
		}
		m19.Value = append(m19.Value, mv)
	}
	m["19"] = m19
	m["20"], err = attributevalue.Marshal(x.OtherBrands)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'OtherBrands': %w", err)
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
			var mv *Appliance
			switch vt := v.(type) {
			case *types.AttributeValueMemberNULL:
				mv = nil
			default:
				mv = &Appliance{}
				err = file_message_v1_message_proto_unmarshal_dynamo_item(vt, mv)
			}
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Furniture': %w", err)
			}
			x.Furniture[int64(mk)] = mv
		}
	}
	err = attributevalue.Unmarshal(m["14"], &x.Calendar)
	if err != nil {
		return fmt.Errorf("failed to unmarshal field 'Calendar': %w", err)
	}
	if m["15"] != nil {
		x.WasherEngine = new(Engine)
		err = file_message_v1_message_proto_unmarshal_dynamo_item(m["15"], x.WasherEngine)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'WasherEngine': %w", err)
		}
	}
	if m["16"] != nil {
		x.ExtraKitchen = new(Kitchen)
		err = file_message_v1_message_proto_unmarshal_dynamo_item(m["16"], x.ExtraKitchen)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'ExtraKitchen': %w", err)
		}
	}
	if m["17"] != nil {
		x.Timer = new(durationpb.Duration)
		err = file_message_v1_message_proto_unmarshal_dynamo_item(m["17"], x.Timer)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field 'Timer': %w", err)
		}
	}
	if m["18"] != nil {
		x.WallTime = new(timestamppb.Timestamp)
		err = file_message_v1_message_proto_unmarshal_dynamo_item(m["18"], x.WallTime)
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
			err = file_message_v1_message_proto_unmarshal_dynamo_item(v, &mv)
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
	m["1"], err = attributevalue.Marshal(x.Int64Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Int64Int64': %w", err)
	}
	m["2"], err = attributevalue.Marshal(x.Uint64Uint64)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Uint64Uint64': %w", err)
	}
	m["3"], err = attributevalue.Marshal(x.Fixed64Fixed64)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Fixed64Fixed64': %w", err)
	}
	m["4"], err = attributevalue.Marshal(x.Sint64Sint64)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Sint64Sint64': %w", err)
	}
	m["5"], err = attributevalue.Marshal(x.Sfixed64Sfixed64)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Sfixed64Sfixed64': %w", err)
	}
	m["6"], err = attributevalue.Marshal(x.Int32Int32)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Int32Int32': %w", err)
	}
	m["7"], err = attributevalue.Marshal(x.Uint32Uint32)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Uint32Uint32': %w", err)
	}
	m["8"], err = attributevalue.Marshal(x.Fixed32Fixed32)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Fixed32Fixed32': %w", err)
	}
	m["9"], err = attributevalue.Marshal(x.Sint32Sint32)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Sint32Sint32': %w", err)
	}
	m["10"], err = attributevalue.Marshal(x.Sfixed32Sfixed32)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Sfixed32Sfixed32': %w", err)
	}
	m["11"], err = attributevalue.Marshal(x.Stringstring)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Stringstring': %w", err)
	}
	m["12"], err = attributevalue.Marshal(x.Boolbool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Boolbool': %w", err)
	}
	m["13"], err = attributevalue.Marshal(x.Stringbytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Stringbytes': %w", err)
	}
	m["14"], err = attributevalue.Marshal(x.Stringdouble)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Stringdouble': %w", err)
	}
	m["15"], err = attributevalue.Marshal(x.Stringfloat)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal field 'Stringfloat': %w", err)
	}
	if x.Stringduration != nil {
		m16 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Stringduration {
			mk := k
			if v == nil {
				m16.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_message_v1_message_proto_marshal_dynamo_item(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal map value of field 'Stringduration': %w", err)
			}
			m16.Value[mk] = mv
		}
		m["16"] = m16
	}
	if x.Stringtimestamp != nil {
		m17 := &types.AttributeValueMemberM{Value: make(map[string]types.AttributeValue)}
		for k, v := range x.Stringtimestamp {
			mk := k
			if v == nil {
				m17.Value[mk] = &types.AttributeValueMemberNULL{Value: true}
				continue
			}
			mv, err := file_message_v1_message_proto_marshal_dynamo_item(v)
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
			var mv *durationpb.Duration
			switch vt := v.(type) {
			case *types.AttributeValueMemberNULL:
				mv = nil
			default:
				mv = &durationpb.Duration{}
				err = file_message_v1_message_proto_unmarshal_dynamo_item(vt, mv)
			}
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Stringduration': %w", err)
			}
			x.Stringduration[string(mk)] = mv
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
			var mv *timestamppb.Timestamp
			switch vt := v.(type) {
			case *types.AttributeValueMemberNULL:
				mv = nil
			default:
				mv = &timestamppb.Timestamp{}
				err = file_message_v1_message_proto_unmarshal_dynamo_item(vt, mv)
			}
			if err != nil {
				return fmt.Errorf("failed to unmarshal map value for field 'Stringtimestamp': %w", err)
			}
			x.Stringtimestamp[string(mk)] = mv
		}
	}
	return nil
}
