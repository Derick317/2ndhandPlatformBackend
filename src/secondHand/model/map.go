package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type MapType[KeyType comparable, ValueType any] map[KeyType]ValueType

// Scan Scanner
func (args *MapType[KeyType, ValueType]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(b, &args)
}

// Value Valuer
func (args MapType[KeyType, ValueType]) Value() (driver.Value, error) {
	if args == nil {
		return nil, nil
	}

	return json.Marshal(args)
}

func (args *MapType[KeyType, ValueType]) Add(key KeyType, value ValueType) {
	if *args == nil {
		*args = make(MapType[KeyType, ValueType])
	}
	(*args)[key] = value
}

// If args is nil or there is no such element, Remove is a no-op.
// It reports whether it is successfully deleted. Failure happends if key does not exist.
func (args *MapType[KeyType, ValueType]) Remove(key KeyType) bool {
	if _, ok := (*args)[key]; !ok {
		return false
	}
	delete(*args, key)
	return true
}

func (args *MapType[KeyType, ValueType]) Size() int {
	if *args == nil {
		return 0
	}
	return len(*args)
}
