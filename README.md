# StrMap

Package [strmap](https://github.com/fmelosilva/strmap) converts a given struct to/from a flat map.

    struct ⇆ map[string][]string


## Installation

Use the command below to install StrMap:

```sh
go get -u github.com/fmelosilva/strmap
```

Then import it in your code.

```go
import github.com/fmelosilva/strmap
```


## Motivation

The format `map[string][]string` is used in many essential packages. One example is the [form values and query parameters](https://pkg.go.dev/net/url#Values) presented in the default golang [http package](https://pkg.go.dev/net/http). Other example is the [gRPC metadata](https://pkg.go.dev/google.golang.org/grpc/metadata#MD) since it heavily uses HTTP/2 features.

### Main differences from other libraries

Some libraries that do a similar job and the reason why they are different is described below:

* [gorilla/schema](https://github.com/gorilla/schema) • It's a library that converts to the same format as this library. The main difference is that `strmap` supports nested structs and default values.
* [magiconair/properties](https://github.com/magiconair/properties) • Used to convert a `map[string]string` into a `struct`. It supports nested structs and default values, but currently has no way to generate a map from a struct.
* [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) • The values are converted to/from the nested format `map[string]interface{}` and cannot be used to generate a `map[string][]string`.

MapStructure is not listed below since it does not generate a flat map and encode/decode values to string.

| Feature                 | schema | properties | strmap |
| ----------------------- | :----: | :--------: | :----: |
| Encode (*struct ➞ map*) |   X    |            |   X    |
| Decode (*map ➞ struct*) |   X    |     X      |   X    |
| Custom Type             |   X    |            |   X    |
| Struct Cache            |   X    |            |        |
| Default Value           |        |     X      |   X    |
| Nested Struct           |        |     X      |   X    |

## Examples

The following examples show the process of encoding and decoding using this package.

### Simple Struct

**Encoding process** (`struct` ➞ `map`)

```go
type Root struct {
    A string
    B int
}

in := Root{
    A: "foo",
    B: 10,
}
out := make(map[string][]string)
strmap.Encode(in, out)

print(out)
// output:
// map[string][]string{
//     "a": {"foo"},
//     "b": {"10"},
// }
```
**Decoding process** (`map` ➞ `struct`)

```go
in := map[string][]string{
    "a": {"foo"},
    "b": {"10"},
}
var out Root
strmap.Decode(in, &out)

print(out)
// output:
// Root{
//     A: "foo",
//     B: 10,
// }
```

### Complex Struct

**Encoding process** (`struct` ➞ `map`)

```go
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
strmap.Encode(in, out)

print(out)
// output: 
// map[string][]string{
//     "a":                  {"foo"},
//     "b-field":            {"60"},
//     "c":                  {"fred", "plugh"},
//     "child-field":        {"bar"},
//     "prefix-child-field": {"baz"},
//     "custom":             {"a", "b"},
// }
```

**Decoding process** (`map` ➞ `struct`)

```go
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
strmap.Decode(in, &out)

print(out)
// output:
// Root{
//     A: "foo",
//     B: 60,
//     C: []string{"fred", "plugh"},
//     NoPrefixNested: Child{
//         Field: "bar",
//     },
//     PrefixNested: Child{
//         Field: "baz",
//     },
//     Custom: CustomType{
//         Internal: []string{"a", "b"},
//     },
// }
```

> **NOTE**: fields with same name can lead to **unexpected** behavior!
