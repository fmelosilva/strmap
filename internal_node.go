package strmap

import (
	"fmt"
	"reflect"
	"strings"
)

var defaultTag = "strmap"

type InternalNode struct {
	NodeData
	children    []Node
	err         error
	typ         reflect.Type
	childPrefix string
}

func (n InternalNode) ChildPrefix() string {
	return n.childPrefix
}

func (n InternalNode) Type() reflect.Type {
	return n.typ
}

func (n InternalNode) parseField(v reflect.Value, sf reflect.StructField) (Node, error) {
	p, err := parseTagProps(sf.Tag.Get(defaultTag))
	if err != nil {
		return nil, err
	}

	if p.Ignore {
		return nil, nil
	}

	childPrefix := n.childPrefix
	if p.Name != nil {
		childPrefix = childPrefix + *p.Name
	}

	name := strings.ToLower(sf.Name)
	if p.Name != nil {
		name = *p.Name
	}

	nd := NodeData{
		name:     n.childPrefix + name,
		value:    v,
		tagProps: p,
		parent:   &n,
	}

	if sf.Type.Kind() == reflect.Struct {
		return &InternalNode{
			NodeData:    nd,
			childPrefix: childPrefix,
			typ:         sf.Type,
		}, nil
	}

	return &LeafNode{
		NodeData: nd,
		typ:      sf,
	}, nil
}

func (n *InternalNode) Children() ([]Node, error) {
	if n.children != nil || n.err != nil {
		return n.children, n.err
	}

	n.children = []Node{}

	for i := 0; i < n.value.NumField(); i++ {
		v := n.value.Field(i)
		t := n.value.Type().Field(i)
		child, err := n.parseField(v, t)
		if err != nil {
			n.err = err
			return nil, err
		}

		if child == nil {
			continue
		}

		n.children = append(n.children, child)
	}

	return n.children, nil
}

func NewNode(s interface{}) (*InternalNode, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, found %T", s)
	}

	return &InternalNode{
		NodeData: NodeData{
			value:    v,
			name:     "",
			tagProps: TagProps{Ignore: false},
			parent:   nil,
		},
		typ: v.Type(),
	}, nil
}
