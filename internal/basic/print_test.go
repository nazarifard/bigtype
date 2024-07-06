package basic

import (
	"testing"
)

func TestTreeInOrder(t *testing.T) {
	tree := newTree[uint64, uint64](4)
	tree.Set(2, 200)
	tree.Set(1, 100)
	tree.Set(3, 300)
	tree.inOrder(tree.root)
}
