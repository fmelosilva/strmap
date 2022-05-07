package strmap

import (
	"fmt"
	"reflect"
)

type Encoder struct {
	typesEncoder map[reflect.Type]EncodeF
}

func (e Encoder) encodeValue(meta Meta, v reflect.Value) []string {
	f := e.typesEncoder[v.Type()]
	if f == nil {
		return nil
	}

	return f(e.encodeValue, meta, v)
}

func (e Encoder) canEncode(v reflect.Value) bool {
	_, ok := e.typesEncoder[v.Type()]

	return ok
}

func (e Encoder) encode(n Node, m map[string][]string) error {
	if e.canEncode(n.Value()) {
		if !n.IsNil() {
			m[n.Name()] = e.encodeValue(n.TagProps().All(), n.Value())
			return nil
		}

		if n.TagProps().Default != nil {
			m[n.Name()] = *n.TagProps().Default
		}
		return nil
	}

	if n, ok := n.(*InternalNode); ok {
		return e.encodeChildren(n, m)
	}

	return fmt.Errorf("error encoding value: encoder to %T not found", n.Value().Interface())
}

func (e Encoder) encodeChildren(n *InternalNode, m map[string][]string) error {
	children, err := n.Children()
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := e.encode(child, m); err != nil {
			return err
		}
	}

	return nil
}

func (e Encoder) Encode(i interface{}, m map[string][]string) error {
	n, err := NewNode(i)
	if err != nil {
		return err
	}

	return e.encode(n, m)
}

func Encode(i interface{}, m map[string][]string) error {
	return NewEncoder().Encode(i, m)
}

func NewEncoder() Encoder {
	return Encoder{
		typesEncoder: newTypesEncoder(),
	}
}
