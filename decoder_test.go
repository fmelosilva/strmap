package strmap_test

import (
	"reflect"
	"testing"

	"github.com/fmelosilva/strmap"

	"github.com/stretchr/testify/require"
)

func TestDecodeSimpleStruct(t *testing.T) {
	assert := require.New(t)

	type Root struct {
		A string
		B int
	}

	in := map[string][]string{
		"a": {"foo"},
		"b": {"10"},
	}
	var out Root

	assert.NoError(strmap.Decode(in, &out))
	assert.Equal(Root{
		A: "foo",
		B: 10,
	}, out)
}

func TestDecodeComplexStruct(t *testing.T) {
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

	customTypeDecode := func(_ strmap.NextDecodeF, _ strmap.Meta, v reflect.Value, s []string) error {
		v.Set(reflect.ValueOf(CustomType{
			Internal: s,
		}))
		return nil
	}

	strmap.RegisterDecoderOf(CustomType{}, customTypeDecode)

	in := map[string][]string{
		"a":                  {"foo"},
		"b-field":            {"60"},
		"child-field":        {"bar"},
		"prefix-child-field": {"baz"},
		"custom":             {"a", "b"},
		"ignored":            {"ignored"},
	}
	var out Root
	assert.NoError(strmap.Decode(in, &out))

	b := 60
	assert.Equal(Root{
		A: "foo",
		B: &b,
		C: []string{"fred", "plugh"},
		NoPrefixNested: Child{
			Field: "bar",
		},
		PrefixNested: Child{
			Field: "baz",
		},
		Custom: CustomType{
			Internal: []string{"a", "b"},
		},
	}, out)
}
