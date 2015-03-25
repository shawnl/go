package redblack

import (
	"unsafe"
)

const (
	//I don't really like this prevailing style where consts are indistinguishable from variables
	Black = false
	Red   = true
)

type node struct {
	left, right      *Node
	parent_and_color uintptr
}

// The parent and offsetof(parent, node) has to tracked from type-inferred declaration
// in order to make tree.Tree.First() and tree.Tree.Last() work
type root struct {
	root    *Node
	compare func(l, r *Node) int
}

func (n *Node) color() bool {
	return (n.parent_and_color & 1) != 0
}

func (n *Node) set(c bool) {
	if c == true {
		n.parent_and_color = (n.parent_and_color &^ 3) | 1
	} else {
		n.parent_and_color = (n.parent_and_color &^ 3) | 0
	}
	return
}

func (n *Node) parent() *Node {
	return (*Node)(unsafe.Pointer(n.parent_and_color &^ 3))
}

func (n *Node) setParent(p *Node) {
	n.setParent((*Node)((unsafe.Pointer)(uintptr(unsafe.Pointer(n.parent_and_color & 3)) |
		uintptr(unsafe.Pointer(p)))))
}

func (n *Node) Parent() *parent(Node) {
	return parent(n.parent())
}

func (n *Node) InitTree(f func(l, r *Node) int) *root {
	t := new(root)

	t.root = n
	t.compare = f //TODO fixup function with offsetof()
	//compiler t.offsetof = offsetof(...)
	return t
}

func (t *root) First() *parent(Node) {

	if t.root == nil {
		return nil
	}

	n := t.root

	for n.left != nil {
		n = n.left
	}

	return parent(n)
}

func (t *root) Last() *parent(Node) {
	if t.root == nil {
		return nil
	}

	n := t.root

	for n.right != nil {
		n = n.right
	}

	return parent(n)
}

func (t *root) rotateLeft(n *Node) {
	right, parent := n.right, n.parent()

	n.right = right.left
	if n.right != nil {
		right.left.setParent(n)
	}
	right.left = n

	right.setParent(parent)

	if parent != nil {
		if n == parent.left {
			parent.left = right
		} else {
			parent.right = right
		}
	} else {
		t.root = right
	}

	n.setParent(right)
	return
}

func (t *root) rotateRight(n *Node) {
	left, parent := n.left, n.parent()

	n.left = left.right
	if n.left != nil {
		left.right.setParent(n)
	}
	left.right = n

	left.setParent(parent)

	if parent != nil {
		if n == parent.right {
			parent.right = left
		} else {
			parent.left = left
		}
	} else {
		t.root = left
	}

	n.setParent(left)
	return
}

func (t *root) insertColor(n *Node) {
	var node, parent, gparent *Node

	node = n

	for (node.parent() != nil) && (node.parent().color() == Red) {
		parent = node.parent()
		gparent = parent.parent()

		if parent == gparent.left {
			uncle := gparent.right
			if (uncle != nil) && (uncle.color() == Red) {
				uncle.set(Black)
				parent.set(Black)
				gparent.set(Red)
				node = gparent
				continue
			}

			if parent.right == node {
				t.rotateLeft(parent)
				parent, node = node, parent
			}

			parent.set(Black)
			parent.set(Red)
			t.rotateRight(gparent)
		} else {
			uncle := gparent.left
			if (uncle != nil) && (uncle.color() == Red) {
				uncle.set(Black)
				parent.set(Black)
				gparent.set(Red)
				node = gparent
				continue
			}

			if parent.left == node {
				t.rotateRight(parent)
				parent, node = node, parent
			}

			parent.set(Black)
			gparent.set(Red)
			t.rotateRight(gparent)
		}
	}

	t.root.set(Black)
}

func (t *root) Insert(n *Node) *parent(Node) {
	r := &t.root
	var p *node
	
	var this parent(n)
	
	for r != nil {
		this = *parent(*r) //TODO how to make this not ambiguous to compiler
		var cmp int
		
		cmp = t.compare(n, this)
		
		p = r
		if cmp < 0 {
			r = &((*r).left)
		} else cmp > 0 {
			r = &((*r).right)
		} else {
			return this
		}
	}
	
	n.setParent(p)
	n.right = nil
	n.left = nil
	r = *n
	t.insertColor(n)
	
	return this
}

func (t *root) eraseColor(n, p *Node) {
	var other, node, parent *Node

	node, parent = n, p

	for ((node != nil) || (node.color() == Black)) && (node != t.root) {
		if parent.left == node {
			other = parent.right
			if other.color() == Red {
				other.set(Black)
				parent.set(Red)
				t.rotateLeft(parent)
				other = parent.right
			}
			if ((other.left == nil) || (other.left.color() == Black)) &&
				((other.right == nil) || (other.right.color() == Black)) {
				other.set(Red)
				node = parent
				parent = node.parent()
			} else {
				if (other.right != nil) || (other.right.color() == Black) {
					other.left.set(Black)
					other.set(Red)
					t.rotateRight(other)
					other = parent.right
				}
				other.set(parent.color())
				parent.set(Black)
				other.right.set(Black)
				t.rotateLeft(parent)
				node = t.root
				break
			}
		} else {
			other = parent.left
			if other.color() == Red {
				other.set(Black)
				parent.set(Red)
				t.rotateRight(parent)
				other = parent.left
			}
			if ((other.left == nil) || (other.left.color() == Black)) &&
				(other.right == nil || (other.right.color() == Black)) {
				other.set(Red)
				node = parent
				parent = node.parent()
			} else {
				if other.left == nil || other.left.color() == Black {
					other.right.set(Black)
					other.set(Red)
					t.rotateLeft(other)
					other = parent.left
				}
				other.set(parent.color())
				parent.set(Black)
				other.left.set(Black)
				t.rotateRight(parent)
				node = t.root
				break
			}
		}
	}

	if node != nil {
		node.set(Black)
	}
}

func (t *root) Remove(n *Node) {
	var node, child, parent *Node
	var color bool

	node = n

	if node.left == nil {
		child = node.right
	} else if node.right == nil {
		child = node.left
	} else {
		old := node
		node = node.right
		for node.left != nil {
			node = node.left
		}

		if old.parent() != nil {
			if old.parent().left == old {
				old.parent().left = node
			} else {
				old.parent().right = node
			}
		} else {
			t.root = node
		}

		child = node.right
		parent = node.parent()
		color = node.color()

		if parent == old {
			parent = node
		} else {
			if child != nil {
				child.setParent(parent)
			}
			parent.left = child

			node.right = old.right
			old.right.setParent(node)
		}

		node.parent_and_color = old.parent_and_color
		node.left = old.left
		old.left.setParent(node)

		goto color
	}

	parent = node.parent()
	color = node.color()

	if child != nil {
		child.setParent(parent)
	}

	if parent != nil {
		if parent.left == node {
			parent.left = child
		} else {
			parent.right = child
		}
	} else {
		t.root = child
	}

color:
	if color == Black {
		t.eraseColor(child, parent)
	}
}

func (n *Node) Next() *parent(Node) {
	if n.parent() == n {
		return nil
	}

	if n.right != nil {
		n = n.right
		for n.left != nil {
			n = n.left
		}
		return parent(n)
	}

	for (n.parent() != nil) && (n == n.parent().right) {
		n = n.parent()
	}

	return parent(n)
}

func (n *Node) Prev() *Node {
	if n.parent() == n {
		return nil
	}

	if n.left != nil {
		n = n.left
		for n.right != nil {
			n = n.right
		}
		return n
	}

	for (n.parent() != nil) && (n == n.parent().left) {
		n = n.parent()
	}

	return parent(n)
}

func (t *root) Replace(v, n *Node) {
	parent := v.parent()

	if parent != nil {
		if v == parent.left {
			parent.left = n
		} else {
			parent.right = n
		}
	} else {
		t.root = n
	}

	if v.left != nil {
		v.left.setParent(n)
	}

	if v.right != nil {
		v.right.setParent(n)
	}

	*n = *v
}
