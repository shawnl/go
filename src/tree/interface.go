package tree

import (
	"unsafe"
)

type Node interface {
	Next() *parent(Node)
	Prev() *parent(Node)
	InitTree(compare func(left, right *parent(Node)) int) *Tree
}

type Tree interface {
	First() *parent(Node)
	Last() *parent(Node)

	Find(key *Node) *parent(Node)
	Insert(n *Node)
	Remove(n *Node)
	Replace(victim, n *Node)
}

/*
parent(thing) builtin

type intrusive_interface has .parent "reverse struct member" built in. : its member MAY NOT be then referenced with it when it is not imported
TODO: It may also be use in a type declaration (include package.IntrusiveInterface.Parent ?
TODO: or should the type declaration be parent(IntrusiveInterface) <===

TODO Interfaces dont work that way. parent(thing) may only be used inside an intrusive_interface
TODO parent(thing) may only be used when it is a struct member: TODO how about Tree.First() Tree.Last(): add pointers and update them: easy, except for type information....still itrinsisically linked to other type

being declared intrusive_interface is necessary to limit the context in which it can be used:
the rules are more strict.

Later: What contexts does it work? How does it work?

TODO How to specify type signatures? This is compile-time, not run-time like reflect

TODO: interfaces that depend on this info only creatable through said intrusive_interface
*/
