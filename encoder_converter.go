package strmap

import (
	"reflect"
	"strconv"
	"time"
)

type EncodeNextF func(Meta, reflect.Value) []string

type EncodeF func(EncodeNextF, Meta, reflect.Value) []string

var typesEncoder = map[reflect.Type]EncodeF{
	stringType:  encodeString,
	boolType:    encodeBool,
	intType:     encodeInt,
	int8Type:    encodeInt,
	int16Type:   encodeInt,
	int32Type:   encodeInt,
	int64Type:   encodeInt,
	uintType:    encodeUint,
	uint8Type:   encodeUint,
	uint16Type:  encodeUint,
	uint32Type:  encodeUint,
	uint64Type:  encodeUint,
	float32Type: encodeFloat32,
	float64Type: encodeFloat64,
	timeType:    encodeTime,
}

func newTypesEncoder() map[reflect.Type]EncodeF {
	dst := make(map[reflect.Type]EncodeF, len(typesEncoder))

	for k, v := range typesEncoder {
		dst[k] = v
	}

	return dst
}

func encodeBool(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{strconv.FormatBool(v.Bool())}
}

func encodeInt(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{strconv.FormatInt(int64(v.Int()), 10)}
}

func encodeUint(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{strconv.FormatUint(uint64(v.Uint()), 10)}
}

func encodeFloat(v reflect.Value, bits int) string {
	return strconv.FormatFloat(v.Float(), 'f', -1, bits)
}

func encodeFloat32(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{encodeFloat(v, 32)}
}

func encodeFloat64(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{encodeFloat(v, 64)}
}

func encodeString(_ EncodeNextF, _ Meta, v reflect.Value) []string {
	return []string{v.String()}
}

func encodeTime(_ EncodeNextF, meta Meta, v reflect.Value) []string {
	t := v.Interface().(time.Time)

	layout, ok := meta["layout"]
	if !ok {
		return []string{t.String()}
	}

	return []string{t.Format(layout)}
}

func encodeSlice(next EncodeNextF, meta Meta, v reflect.Value) []string {
	var result []string
	for i := 0; i < v.Len(); i++ {
		e := next(meta, v.Index(i))
		result = append(result, e...)
	}

	return result
}

func encodePointer(next EncodeNextF, meta Meta, v reflect.Value) []string {
	return next(meta, v.Elem())
}

func RegisterEncoderOf(t interface{}, enc EncodeF) {
	typesEncoder[reflect.TypeOf(t)] = enc
	registerEncoderPointerTo(reflect.TypeOf(t))
}

func RegisterEncoderSliceOf(t reflect.Type) {
	typesEncoder[reflect.SliceOf(t)] = encodeSlice
}

func registerEncoderPointerTo(t reflect.Type) {
	typesEncoder[reflect.PointerTo(t)] = encodePointer
}

func init() {
	for _, t := range basicTypes {
		RegisterEncoderSliceOf(t)
		registerEncoderPointerTo(t)
	}
}
