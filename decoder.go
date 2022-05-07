package strmap

import (
	"fmt"
	"log"
	"reflect"
)

type Decoder struct {
	typesDecoder map[reflect.Type]DecoderF
}

func (d Decoder) decodeValues(meta Meta, v reflect.Value, values []string) error {
	f := d.typesDecoder[v.Type()]
	if f == nil {
		return nil
	}

	return f(d.decodeValues, meta, v, values)
}

func (d Decoder) canDecode(t reflect.Type) bool {
	_, ok := d.typesDecoder[t]
	return ok
}

func (d Decoder) decodeChildren(n *InternalNode, m map[string][]string) error {
	children, err := n.Children()
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := d.decode(child, m); err != nil {
			return err
		}
	}

	return nil
}

func (d Decoder) values(n Node, m map[string][]string) *[]string {
	values, ok := m[n.Name()]
	if ok {
		return &values
	}

	if n.TagProps().Default != nil {
		return n.TagProps().Default
	}

	return nil
}

func (d Decoder) decode(n Node, m map[string][]string) error {
	if d.canDecode(n.Type()) {
		values := d.values(n, m)
		if values == nil {
			log.Printf("field %s was not set", n.Path())
			return nil
		}

		return d.decodeValues(n.TagProps().All(), n.Value(), *values)
	}

	if n, ok := n.(*InternalNode); ok {
		return d.decodeChildren(n, m)
	}

	return fmt.Errorf("error encoding value: encoder to %T not found", n.Value().Interface())
}

func (d Decoder) Decode(m map[string][]string, i interface{}) error {
	n, err := NewNode(i)
	if err != nil {
		return err
	}

	return d.decode(n, m)
}

func Decode(m map[string][]string, i interface{}) error {
	return NewDecoder().Decode(m, i)
}

func NewDecoder() Decoder {
	return Decoder{
		typesDecoder: newTypesDecoder(),
	}
}
