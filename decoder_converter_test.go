package strmap

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func nextErrF(Meta, reflect.Value, []string) error {
	return errors.New("error")
}

func stringNextF(m Meta, v reflect.Value, values []string) error {
	return decodeString(nil, nil, v, values)
}

func TestDecodeInt(t *testing.T) {
	assert := require.New(t)

	intV := reflect.New(intType).Elem()

	assert.NoError(decodeInt(nil, nil, intV, []string{"100"}))
	assert.Equal(100, intV.Interface())
	assert.Error(decodeInt(nil, nil, intV, []string{"abc"}))

	int8V := reflect.New(int8Type).Elem()
	int8Overflow := []string{"999999999999"}
	assert.Error(decodeInt(nil, nil, int8V, int8Overflow))
}

func TestDecodeUInt(t *testing.T) {
	assert := require.New(t)

	uintV := reflect.New(uintType).Elem()
	assert.NoError(decodeUint(nil, nil, uintV, []string{"100"}))
	assert.Equal(uint(100), uintV.Interface())
	assert.Error(decodeUint(nil, nil, uintV, []string{"abc"}))

	uint8V := reflect.New(uint8Type).Elem()
	uint8Overflow := []string{"999999999999"}
	assert.Error(decodeUint(nil, nil, uint8V, uint8Overflow))
}

func TestDecodeFloat(t *testing.T) {
	assert := require.New(t)

	float64V := reflect.New(float64Type).Elem()
	assert.NoError(decodeFloat(nil, nil, float64V, []string{"100"}))
	assert.Equal(float64(100), float64V.Interface())
	assert.Error(decodeFloat(nil, nil, float64V, []string{"abc"}))

	float32V := reflect.New(float32Type).Elem()
	float32Overflow := []string{"34028234663852885981170418348451692544000.000000"}
	assert.Error(decodeFloat(nil, nil, float32V, float32Overflow))
}

func TestDecodeString(t *testing.T) {
	assert := require.New(t)

	stringV := reflect.New(stringType).Elem()
	assert.NoError(decodeString(nil, nil, stringV, []string{"mytext"}))
	assert.Equal("mytext", stringV.Interface())
}

func TestDecodeBool(t *testing.T) {
	assert := require.New(t)

	boolV := reflect.New(boolType).Elem()
	assert.NoError(decodeBool(nil, nil, boolV, []string{"true"}))
	assert.Equal(true, boolV.Interface())

	assert.Error(decodeBool(nil, nil, boolV, []string{"asdf"}))
}

func TestDecodePointer(t *testing.T) {
	assert := require.New(t)

	stringPointerV := reflect.New(reflect.PointerTo(stringType)).Elem()
	assert.Error(decodePointer(nextErrF, nil, stringPointerV, []string{"text"}))

	assert.NoError(decodePointer(stringNextF, nil, stringPointerV, []string{"text"}))
	assert.Equal("text", stringPointerV.Elem().Interface())

}

func TestDecodeSlice(t *testing.T) {
	assert := require.New(t)

	stringSliceV := reflect.New(reflect.SliceOf(stringType)).Elem()
	assert.Error(decodeSlice(nextErrF, nil, stringSliceV, []string{"a", "b"}))

	stringSliceV = reflect.New(reflect.SliceOf(stringType)).Elem()
	assert.NoError(decodeSlice(stringNextF, nil, stringSliceV, []string{"a", "b"}))
	assert.Equal([]string{"a", "b"}, stringSliceV.Interface())
}

func TestDecodeTime(t *testing.T) {
	assert := require.New(t)

	timeV := reflect.New(reflect.TypeOf(time.Time{})).Elem()
	assert.NoError(decodeTime(nil, nil, timeV, []string{"2022-04-17 15:52:32.68422611 -0300 -03"}))
	tt, _ := time.Parse(TimeLayout, "2022-04-17 15:52:32.68422611 -0300 -03")
	assert.Equal(tt, timeV.Interface())

	timeLayout := "2006-01-02 15:04:05-07:00"
	assert.NoError(decodeTime(nil, Meta{"layout": timeLayout}, timeV, []string{"2030-06-10 11:14:42-07:00"}))
	tt, _ = time.Parse(timeLayout, "2030-06-10 11:14:42-07:00")
	assert.Equal(tt, timeV.Interface())
}
