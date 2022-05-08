package strmap

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func nextEncodeInt(_ Meta, v reflect.Value) []string {
	return encodeInt(nil, nil, v)
}

func TestEncodeInt(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"100"}, encodeInt(nil, nil, reflect.ValueOf(int(100))))
}

func TestEncodeUInt(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"100"}, encodeUint(nil, nil, reflect.ValueOf(uint(100))))
}

func TestEncodeFloat(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"64.5"}, encodeFloat32(nil, nil, reflect.ValueOf(float32(64.5))))
	assert.Equal([]string{"64.5"}, encodeFloat64(nil, nil, reflect.ValueOf(float64(64.5))))
}

func TestEncodeString(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"text"}, encodeString(nil, nil, reflect.ValueOf("text")))
}

func TestEncodeBool(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"true"}, encodeBool(nil, nil, reflect.ValueOf(true)))
}

func TestEncodeSlice(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{"10", "20"}, encodeSlice(nextEncodeInt, nil, reflect.ValueOf([]int{10, 20})))
}

func TestEncodePointer(t *testing.T) {
	assert := require.New(t)

	v := 10
	assert.Equal([]string{"10"}, encodePointer(nextEncodeInt, nil, reflect.ValueOf(&v)))
}

func TestEncodeTime(t *testing.T) {
	assert := require.New(t)

	tt, _ := time.Parse(TimeLayout, "2022-04-17 15:52:32.68422611 -0300 -03")
	assert.Equal([]string{"2022-04-17 15:52:32.68422611 -0300 -03"}, encodeTime(nil, nil, reflect.ValueOf(tt)))

	timeLayout := "2006-01-02 15:04:05-07:00"
	tt, _ = time.Parse(timeLayout, "2030-06-10 11:14:42-07:00")
	assert.Equal([]string{"2030-06-10 11:14:42-07:00"}, encodeTime(nil, Meta{"layout": timeLayout}, reflect.ValueOf(tt)))
}
