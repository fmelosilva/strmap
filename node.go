package strmap

import "reflect"

type Node interface {
	Name() string
	Path() string
	Value() reflect.Value
	IsNil() bool
	Type() reflect.Type
	Parent() *InternalNode
	TagProps() TagProps
}

type NodeData struct {
	name     string
	tagProps TagProps
	value    reflect.Value
	parent   *InternalNode
}

func (n NodeData) Name() string {
	return n.name
}

func (n NodeData) Value() reflect.Value {
	return n.value
}

func (n NodeData) IsNil() bool {
	return isNil(n.value)
}

func (n NodeData) TagProps() TagProps {
	return n.tagProps
}

func (n NodeData) Parent() *InternalNode {
	return n.parent
}

func (n NodeData) Path() string {
	return ""
}
