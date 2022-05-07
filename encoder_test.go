package strmap_test

import (
	"reflect"
	"testing"

	"github.com/fmelosilva/strmap"

	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	assert := require.New(t)

	type Root struct {
		A string
		B int
	}

	in := Root{
		A: "foo",
		B: 10,
	}
	out := make(map[string][]string)
	assert.NoError(strmap.Encode(in, out))
	assert.Equal(map[string][]string{
		"a": {"foo"},
		"b": {"10"},
	}, out)

}

func TestEncodeComplexStruct(t *testing.T) {
	assert := require.New(t)

	type Child struct {
		Field string `strmap:"child-field"`
	}

	type CustomType struct {
		Internal []string
	}

	type Root struct {
		A              string
		B              *int     `strmap:"b-field,default=60"`
		C              []string `strmap:",default=fred;plugh"`
		NullableField  []int
		NoPrefixNested Child
		PrefixNested   Child `strmap:"prefix-"`
		Custom         CustomType
		Ignored        string `strmap:"-"`
	}

	customTypeEncode := func(_ strmap.EncodeNextF, _ strmap.Meta, v reflect.Value) []string {
		return v.Interface().(CustomType).Internal
	}
	strmap.RegisterEncoderOf(CustomType{}, customTypeEncode)

	in := Root{
		A: "foo",
		NoPrefixNested: Child{
			Field: "bar",
		},
		PrefixNested: Child{
			Field: "baz",
		},
		Custom: CustomType{
			Internal: []string{"a", "b"},
		},
		Ignored: "qux",
	}
	out := make(map[string][]string)
	assert.NoError(strmap.Encode(in, out))
	assert.Equal(map[string][]string{
		"a":                  {"foo"},
		"b-field":            {"60"},
		"c":                  {"fred", "plugh"},
		"child-field":        {"bar"},
		"prefix-child-field": {"baz"},
		"custom":             {"a", "b"},
	}, out)

}
