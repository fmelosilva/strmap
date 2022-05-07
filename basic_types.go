package strmap

import (
	"reflect"
	"time"
)

var TimeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

var (
	stringType  = reflect.TypeOf("")
	intType     = reflect.TypeOf(int(0))
	int8Type    = reflect.TypeOf(int8(0))
	int16Type   = reflect.TypeOf(int16(0))
	int32Type   = reflect.TypeOf(int32(0))
	int64Type   = reflect.TypeOf(int64(0))
	uintType    = reflect.TypeOf(uint(0))
	uint8Type   = reflect.TypeOf(uint8(0))
	uint16Type  = reflect.TypeOf(uint16(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	boolType    = reflect.TypeOf(false)
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
	timeType    = reflect.TypeOf(time.Time{})
)

var basicTypes = []reflect.Type{
	stringType,
	boolType,
	intType,
	int8Type,
	int16Type,
	int32Type,
	int64Type,
	uintType,
	uint8Type,
	uint16Type,
	uint32Type,
	uint64Type,
	float32Type,
	float64Type,
	timeType,
}

func isNil(v reflect.Value) bool {
	if v.Interface() == nil {
		return true
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return v.IsNil()
	}
	return false
}
