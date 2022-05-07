package strmap_test

import (
	"testing"

	"github.com/fmelosilva/strmap"

	"github.com/stretchr/testify/require"
)

func TestNodeChildren(t *testing.T) {
	assert := require.New(t)

	type Child struct {
		F string `strmap:", default = 30"`
	}

	type Parent struct {
		Child       Child ``
		Child1      Child `strmap:"child-1"`
		IgnoreChild Child `strmap:"-"`
	}

	n, err := strmap.NewNode(Parent{
		Child: Child{
			F: "10",
		},
		Child1: Child{
			F: "2",
		},
	})
	assert.NoError(err)
	ns, err := n.Children()
	assert.NoError(err)
	assert.Len(ns, 2)
	assert.Equal("child", ns[0].Name())
	assert.Equal("child-1", ns[1].Name())
	cs, err := ns[0].(*strmap.InternalNode).Children()
	assert.NoError(err)
	assert.Equal("f", cs[0].Name())
	assert.Equal([]string{"30"}, *cs[0].TagProps().Default)

	cs, err = ns[1].(*strmap.InternalNode).Children()
	assert.NoError(err)
	assert.Equal("child-1f", cs[0].Name())
	assert.Equal([]string{"30"}, *cs[0].TagProps().Default)
}
