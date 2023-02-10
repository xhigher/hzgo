package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IntArray []int

func (i *IntArray) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	var d IntArray
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i IntArray) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Int64Array []int64

func (i *Int64Array) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	var d Int64Array
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i Int64Array) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Int32Array []int32

func (i *Int32Array) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	var d Int32Array
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i Int32Array) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type StringArray []string

func (i *StringArray) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	if len(data) == 0 {
		return nil
	}
	var d StringArray
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i StringArray) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}
