package ddb

import (
	"fmt"
	"strconv"
)

// marshalMapKey marshals the key of a map to a string
func marshalMapKey[K comparable](k K) (string, error) {
	var kv string
	switch kt := any(k).(type) {
	case string:
		kv = kt
	case bool:
		if kt {
			kv = "true"
		} else {
			kv = "false"
		}
	case int32, int64, uint32, uint64:
		kv = fmt.Sprintf("%d", kt)
	default:
		return "", fmt.Errorf("unsupported map key type: %T", k)
	}
	if kv == "" {
		return "", fmt.Errorf("failed to marshal map key: map key cannot be empty")
	}
	return kv, nil
}

// UintMapKey parses 's' as an unsigned integer value
func UintMapKey[K ~uint32 | ~uint64](s string) (K, error) {
	k, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return K(k), nil
}

// IntMapKey parses 's' as a signed integer value
func IntMapKey[K ~int32 | ~int64](s string) (K, error) {
	k, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return K(k), nil
}

// BoolMapKey parses 's' as a boolean 'true' or 'false' value
func BoolMapKey(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool key: %v", s)
	}
}

// StringMapKey parses 's' as a string map key
func StringMapKey(s string) (string, error) {
	return s, nil
}
