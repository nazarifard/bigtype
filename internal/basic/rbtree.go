package basic

//incomplete code. don't use this file

// import (
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"github.com/nazarifard/bigtype/internal/options"
// )
// const NULL = 0
// type Color bool
// const Black Color = false
// const Red Color = true
// type ptrNode = uint32
// type RBNode[K kNumber] struct {
// 	key         int
// 	color       Color
// 	left, right ptrNode
// 	parent      ptrNode
// }
// type RedBlackTree[K kNumber, V any] struct {
// 	data       Array[V]
// 	indexTable Array[RBNode[K]] //*bigFixedArray[node[K]]
// 	root       ptrNode
// 	free       ptrNode
// }
// func NewRedBlackTree[K kNumber, V any](ops ...any) *RedBlackTree[K, V] {
// 	if reflect.ValueOf(*new(K)).Kind() == reflect.String {
// 		fmt.Println("Error: bigtype.bigTree does not support string keys. use bigtype.Map")
// 		os.Exit(1)
// 	}
// 	option := options.ParseTreeOptions[K, V](ops...)
// 	return &RedBlackTree[K, V]{
// 		indexTable: NewFixedArray[RBNode[K]](option.HistSize+1, NULL, true),
// 		data:       NewArray[V](option.HistSize+1, option.VMarshal, true),
// 		root:       0,
// 		free:       1,
// 	}
// 	//return &RedBlackTree[K, V]{}
// }
// func (t *RedBlackTree[K, V]) Insert(key int) {
// 	const NULL uint32 = 0
// 	if t.root == NULL {
// 		t.indexTable.Set(int(t.free), RBNode[K]{key: key, color: Black})
// 		t.root = t.free
// 		t.free++
// 	} else {
// 		t.insert(t.root, key)
// 	}
// }
// func (t *RedBlackTree[K, V]) insert(root ptrNode, key int) {
// 	n := t.indexTable.Get(int(root))
// 	if key < n.key {
// 		if n.left == NULL {
// 			t.indexTable.Set(int(t.free), RBNode[K]{key: key, color: Red, parent: root})
// 			t.free++
// 			n.left = t.free
// 			t.fixInsert(n.left)
// 		} else {
// 			t.insert(n.left, key)
// 		}
// 	} else if key > n.key {
// 		if n.right == NULL {
// 			t.indexTable.Set(int(t.free), RBNode[K]{key: key, color: Red, parent: root})
// 			t.free++
// 			n.right = t.free
// 			t.fixInsert(n.right)
// 		} else {
// 			t.insert(n.right, key)
// 		}
// 	}
// }
// func (t *RedBlackTree[K, V]) fixInsert(n ptrNode) {
// 	N := t.indexTable.Get(int(n))
// 	if N.parent == NULL {
// 		N.color = Black
// 		t.indexTable.Set(int(n), N)
// 		return
// 	}
// 	parentNode := t.indexTable.Get(int(N.parent))
// 	if parentNode.color == Black {
// 		return
// 	}
// 	u := t.getUncle(n)
// 	g := t.getGrandparent(n)
// 	uncle := t.indexTable.Get(int(u))
// 	grandparent := t.indexTable.Get(int(g))
// 	if u != NULL && uncle.color == Red {
// 		parentNode.color = Black
// 		uncle.color = Black
// 		grandparent.color = Red
// 		t.fixInsert(g)
// 		return
// 	}
// 	if n == parentNode.right && N.parent == grandparent.left {
// 		t.rotateLeft(g)
// 		n = N.left
// 	} else if n == parentNode.left && N.parent == grandparent.right {
// 		t.rotateRight(g)
// 		n = N.right
// 	}
// 	parentNode.color = Black
// 	grandparent.color = Red
// 	if n == parentNode.left {
// 		t.rotateRight(g)
// 	} else {
// 		t.rotateLeft(g)
// 	}
// }
// func (t *RedBlackTree[K, V]) getUncle(p ptrNode) ptrNode {
// 	n := t.indexTable.Get(int(p))
// 	if n.parent == NULL {
// 		return NULL
// 	}
// 	parent := t.indexTable.Get(int(n.parent))
// 	if parent.parent == NULL {
// 		return NULL
// 	}
// 	grandparent := t.indexTable.Get(int(parent.parent))
// 	if n.parent == grandparent.left {
// 		return grandparent.right
// 	}
// 	return grandparent.left
// }
// func (t *RedBlackTree[K, V]) getGrandparent(p ptrNode) ptrNode {
// 	n := t.indexTable.Get(int(p))
// 	if n.parent != NULL {
// 		return t.indexTable.Get(int(n.parent)).parent
// 	}
// 	return NULL
// }
// func (t *RedBlackTree[K, V]) rotateLeft(x ptrNode) {
// 	X := t.indexTable.Get(int(x))
// 	y := X.right
// 	Y := t.indexTable.Get(int(y))
// 	X.right = Y.left
// 	if Y.left != NULL {
// 		t.indexTable.Update(int(Y.left), func(old RBNode[K]) (new RBNode[K]) {
// 			old.parent = x
// 			return old
// 		})
// 		//Y.left.parent = X
// 	}
// 	Y.parent = X.parent
// 	grand := X.parent
// 	G := t.indexTable.Get(int(grand))
// 	if X.parent == NULL {
// 		X.parent = y
// 	} else if x == G.left {
// 		G.left = y
// 	} else {
// 		G.right = y
// 	}
// 	Y.left = x
// 	X.parent = y
// 	t.indexTable.Set(int(x), X)
// 	t.indexTable.Set(int(y), Y)
// 	t.indexTable.Set(int(X.parent), G)
// }
// func (t *RedBlackTree[K, V]) rotateRight(x ptrNode) {
// 	X := t.indexTable.Get(int(x))
// 	y := X.left
// 	Y := t.indexTable.Get(int(y))
// 	X.left = Y.right
// 	if Y.right != NULL {
// 		t.indexTable.Update(int(Y.right), func(old RBNode[K]) (new RBNode[K]) {
// 			old.parent = x
// 			return old
// 		})
// 		//y.right.parent = n
// 	}
// 	Y.parent = X.parent
// 	G := t.indexTable.Get(int(X.parent))
// 	if X.parent == NULL {
// 		X.parent = y
// 	} else if x == G.left {
// 		G.left = y
// 	} else {
// 		G.right = y
// 	}
// 	Y.right = x
// 	X.parent = y
// 	t.indexTable.Set(int(x), X)
// 	t.indexTable.Set(int(y), Y)
// 	t.indexTable.Set(int(X.parent), G)
// }
// func (t *RedBlackTree[K, V]) InorderTraversal() {
// 	if t.root != NULL {
// 		t.inorderTraversal(t.root)
// 	}
// 	fmt.Println()
// }
// func (t *RedBlackTree[K, V]) inorderTraversal(p ptrNode) {
// 	if p != NULL {
// 		n := t.indexTable.Get(int(p))
// 		t.inorderTraversal(n.left)
// 		fmt.Printf("%d ", n.key)
// 		t.inorderTraversal(n.right)
// 	}
// }
// func main() {
// 	tree := NewRedBlackTree[int, int]()
// 	tree.Insert(7)
// 	tree.Insert(3)
// 	tree.Insert(18)
// 	tree.Insert(10)
// 	tree.Insert(22)
// 	tree.Insert(8)
// 	tree.Insert(11)
// 	tree.Insert(26)
// 	tree.Insert(2)
// 	tree.Insert(6)
// 	tree.Insert(13)
// 	fmt.Println("The inorder traversal of this tree is:")
// 	tree.InorderTraversal()
// }
