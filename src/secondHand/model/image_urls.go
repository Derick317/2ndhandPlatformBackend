package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ImageUrlsType map[string]string

// Scan Scanner
func (args *ImageUrlsType) Scan(value interface{}) error {
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
func (args ImageUrlsType) Value() (driver.Value, error) {
	if args == nil {
		return nil, nil
	}

	return json.Marshal(args)
}
