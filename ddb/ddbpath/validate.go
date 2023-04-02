package ddbpath

import (
	"reflect"
)

// infoField traverses "into" a field during validation
func (r Registry) intoField(typ reflect.Type, kind FieldKind) (info FieldInfo, fields map[string]FieldInfo, err error) {
	info = FieldInfo{Kind: kind, Message: typ}
	if typ == nil {
		return
	}

	fields, ok := r.fieldsOf(info.Message)
	if !ok {
		return info, nil, errTypeNotRegistered(typ)
	}
	return info, fields, nil
}

// traverse from a message reflected on by 'typ' into a path described by 'els'
func (r Registry) traverse(typ reflect.Type, els []PathElement) (currInfo FieldInfo, currFields map[string]FieldInfo, err error) {
	currInfo, currFields, err = r.intoField(typ, FieldKindSingle)
	if err != nil {
		return NoInfo, nil, err
	}

	var field string
	var index int
	for i := 0; i < len(els); i++ {
		field, index = els[i].Field, els[i].Index

		// in case we're inside a any field. or the type itself is a any path type
		// we allow anything afterwards
		if currInfo.Message == reflect.TypeOf(ValuePath{}) {
			continue
		}

		switch {
		case index >= 0: // selecting index
			switch currInfo.Kind {
			case FieldKindList:
				// in case of a list, it will always become a "basic" types since protobuf doesn't allow
				// list of lists, or list of maps
				currInfo, currFields, err = r.intoField(currInfo.Message, FieldKindSingle)
				if err != nil {
					return NoInfo, nil, err
				}
			default:
				return NoInfo, nil, errIndexNotAllowed(index, currInfo)
			}
		case index < 0: // selecting field
			switch {
			case currInfo.Kind == FieldKindSingle && currInfo.Message != nil:
				// in case of a single message, the field can be selected
				newField, isValidField := currFields[field]
				if !isValidField {
					return NoInfo, nil, errUnknownField(field, currInfo)
				}

				currInfo, currFields, err = r.intoField(newField.Message, newField.Kind)
				if err != nil {
					return NoInfo, nil, err
				}
			case currInfo.Kind == FieldKindMap:
				// in case of a map, either a single message or single basic type. Maps of maps
				// or map of lists is not supported in protobuf.
				currInfo, currFields, err = r.intoField(currInfo.Message, FieldKindSingle)
				if err != nil {
					return NoInfo, nil, err
				}
			default:
				return NoInfo, nil, errFieldNotAllowed(field, currInfo)
			}
		}
	}

	return currInfo, currFields, nil
}

// validate a single type against registry
func (r Registry) validate(typ reflect.Type, els []PathElement) (err error) {
	_, _, err = r.traverse(typ, els)
	return
}
