package strmap

import (
	"reflect"
	"strconv"
	"time"
)

type NextDecodeF func(Meta, reflect.Value, []string) error

type DecoderF func(NextDecodeF, Meta, reflect.Value, []string) error

var typesDecoder = map[reflect.Type]DecoderF{
	stringType:  decodeString,
	boolType:    decodeBool,
	intType:     decodeInt,
	int8Type:    decodeInt,
	int16Type:   decodeInt,
	int32Type:   decodeInt,
	int64Type:   decodeInt,
	uintType:    decodeUint,
	uint8Type:   decodeUint,
	uint16Type:  decodeUint,
	uint32Type:  decodeUint,
	uint64Type:  decodeUint,
	float32Type: decodeFloat,
	float64Type: decodeFloat,
	timeType:    decodeTime,
}

func newTypesDecoder() map[reflect.Type]DecoderF {
	dst := make(map[reflect.Type]DecoderF, len(typesDecoder))

	for k, v := range typesDecoder {
		dst[k] = v
	}

	return dst
}

func decodeBool(_ NextDecodeF, _ Meta, v reflect.Value, values []string) error {
	if r, err := strconv.ParseBool(values[0]); err == nil {
		v.SetBool(r)
		return nil
	}

	return ParseErrorf("could not parse %s value to %s", values[0], v.Kind())
}

func decodeString(_ NextDecodeF, _ Meta, v reflect.Value, values []string) error {
	v.SetString(values[0])
	return nil
}

func decodeFloat(_ NextDecodeF, _ Meta, v reflect.Value, values []string) error {
	if r, err := strconv.ParseFloat(values[0], 64); err == nil {
		if v.OverflowFloat(r) {
			return OverflowErrorf("value %f overflows %s", r, v.Kind())
		}

		v.SetFloat(r)

		return nil
	}

	return ParseErrorf("could not parse %s value to %s", values[0], v.Kind())
}

func decodeInt(_ NextDecodeF, _ Meta, v reflect.Value, values []string) error {
	if r, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		if v.OverflowInt(r) {
			return OverflowErrorf("value %d overflows %s", r, v.Kind())
		}

		v.SetInt(r)

		return nil
	}

	return ParseErrorf("could not parse %s value to %s", values[0], v.Kind())
}

func decodeUint(_ NextDecodeF, _ Meta, v reflect.Value, values []string) error {
	if r, err := strconv.ParseUint(values[0], 10, 64); err == nil {
		if v.OverflowUint(r) {
			return OverflowErrorf("value %d overflows %s", r, v.Kind())
		}

		v.SetUint(r)

		return nil
	}

	return ParseErrorf("could not parse %s value to %s", values[0], v.Kind())
}

func decodePointer(next NextDecodeF, meta Meta, v reflect.Value, values []string) error {
	el := reflect.New(v.Type().Elem())
	if err := next(meta, el.Elem(), values); err != nil {
		return err
	}

	v.Set(el)

	return nil
}

func decodeSlice(next NextDecodeF, meta Meta, v reflect.Value, values []string) error {
	for _, s := range values {
		e := reflect.New(v.Type().Elem()).Elem()
		if err := next(meta, e, []string{s}); err != nil {
			return err
		}

		v.Set(reflect.Append(v, e))
	}

	return nil
}

func decodeTime(_ NextDecodeF, meta Meta, v reflect.Value, values []string) error {
	layout, ok := meta["layout"]
	if !ok {
		layout = TimeLayout
	}

	t, err := time.Parse(layout, values[0])
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(t))

	return nil
}

func RegisterDecoderSliceOf(t reflect.Type) {
	typesDecoder[reflect.SliceOf(t)] = decodeSlice
}

func registerDecoderPointerTo(t reflect.Type) {
	typesDecoder[reflect.PointerTo(t)] = decodePointer
}

func RegisterDecoderOf(t interface{}, dec DecoderF) {
	typesDecoder[reflect.TypeOf(t)] = dec
	registerDecoderPointerTo(reflect.TypeOf(t))
}

func init() {
	for _, t := range basicTypes {
		RegisterDecoderSliceOf(t)
		registerDecoderPointerTo(t)
	}

	RegisterDecoderOf(time.Time{}, decodeTime)
	RegisterDecoderSliceOf(reflect.TypeOf(time.Time{}))
}
