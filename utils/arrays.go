package utils

import (
	"reflect"
)

func InArray(data interface{}, item interface{}) bool {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			if i == item {
				return true
			}
		}
	}
	return false
}
