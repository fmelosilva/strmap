package strmap

import "reflect"

type LeafNode struct {
	NodeData
	typ reflect.StructField
}

func (n LeafNode) Type() reflect.Type {
	return n.typ.Type
}
