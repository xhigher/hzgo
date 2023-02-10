package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringMap map[string]string

func (i *StringMap) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	var d StringMap
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i StringMap) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}
