package splaytree

import (
	"fmt"
	"hash/fnv"
	"strings"
)

type Key interface {
	Less(Key) bool
	Equal(Key) bool
}

type IntKey int64

func (key IntKey) Less(other Key) (r bool) {
	return key < other.(IntKey)
}

func (key IntKey) Equal(other Key) (r bool) {
	return key == other.(IntKey)
}

type HashKey uint64

func NewHashKey(s string) (key HashKey) {
	h := fnv.New64()
	h.Write([]byte(s))
	return HashKey(h.Sum64())
}

func (key HashKey) Less(other Key) (r bool) {
	return key < other.(HashKey)
}

func (key HashKey) Equal(other Key) (r bool) {
	return key == other.(HashKey)
}

type SplayTree struct {
	Left   *SplayTree
	Key    Key
	Value  interface{}
	Right  *SplayTree
	Parent *SplayTree
}

func (x *SplayTree) Print() {
	if x == nil {
		fmt.Println("<nil>")

	} else {
		x.print("")
	}
}

func (x *SplayTree) print(ind string) {
	if x.Left == nil {
		if x.Right == nil {
			fmt.Printf("%v(%v)\n", x.Key, x.Value)

		} else {
			n, _ := fmt.Printf("%v(%v) ---> ", x.Key, x.Value)
			x.Right.print(ind + strings.Repeat(" ", n))
		}

	} else {
		if x.Right == nil {
			n, _ := fmt.Printf("%v(%v) ---> ", x.Key, x.Value)
			x.Left.print(ind + strings.Repeat(" ", n))

		} else {
			n, _ := fmt.Printf("%v(%v) -+-> ", x.Key, x.Value)
			ind += strings.Repeat(" ", n-4)

			x.Left.print(ind + "|   ")
			fmt.Print(ind + "`-> ")
			x.Right.print(ind + "    ")
		}
	}
}

func (x *SplayTree) IsValid() (isValid bool) {
	return x != nil
}

func (x *SplayTree) IsRoot() (isRoot bool) {
	return x != nil && x.Parent == nil
}

func (x *SplayTree) IsLeft() (isLeft bool) {
	return x != nil && x.Parent != nil && x.Parent.Left == x
}

func (x *SplayTree) IsRight() (isRight bool) {
	return x != nil && x.Parent != nil && x.Parent.Right == x
}

func (x *SplayTree) CopyLinks(src *SplayTree) {
	if src != nil && src.Parent != nil {
		// Update parent links
		x.Parent = src.Parent

		// Update parent->child links
		if src.IsLeft() {
			x.Parent.Left = x
		} else {
			x.Parent.Right = x
		}
	}

	// Update child->parent links
	if x.Left != nil {
		x.Left.Parent = x
	}
	if x.Right != nil {
		x.Right.Parent = x
	}
}

func (x *SplayTree) Insert(key Key, value interface{}) (result *SplayTree) {
	if x == nil {
		result = &SplayTree{nil, key, value, nil, nil}
	} else if key.Equal(x.Key) {
		result = &SplayTree{x.Left, key, value, x.Right, x.Parent}
	} else if key.Less(x.Key) {
		result = &SplayTree{x.Left.Insert(key, value), x.Key, x.Value, x.Right, x.Parent}
	} else {
		result = &SplayTree{x.Left, x.Key, x.Value, x.Right.Insert(key, value), x.Parent}
	}

	result.CopyLinks(x)
	return result
	//return result.Splay()
}

func (x *SplayTree) Search(key Key) (value interface{}, ok bool) {
	if x == nil {
		return nil, false
	}
	if key.Less(x.Key) {
		return x.Left.Search(key)
	}
	if !key.Equal(x.Key) {
		return x.Right.Search(key)
	}

	//x.Splay().CopyLinks(x)

	return x.Value, true
}

func (x *SplayTree) RotateLeft() (y *SplayTree) {
	y = x.Right
	x.Right = y.Left
	y.Left = x

	y.Parent = x.Parent
	x.Parent = y
	return y
}

func (x *SplayTree) RotateRight() (y *SplayTree) {
	y = x.Left
	x.Left = y.Right
	y.Right = x

	y.Parent = x.Parent
	x.Parent = y
	return y
}

func (x *SplayTree) Splay() (result *SplayTree) {
	if x == nil || x.IsRoot() { // No splay occurs
		result = x

	} else {
		p := x.Parent

		if p.IsRoot() { // Zig step
			if x.IsLeft() {
				p = p.RotateRight()
				result = p.Right

			} else {
				p = p.RotateLeft()
				result = p.Left
			}

		} else {
			g := p.Parent

			if x.IsLeft() && p.IsLeft() { // Zig-Zig step, both left
				g = g.RotateRight()
				p = g.Right

				p = p.RotateRight()
				result = p.Right

			} else if x.IsRight() && p.IsRight() { // Zig-Zig step, both right
				g = g.RotateLeft()
				p = g.Left

				p = p.RotateLeft()
				result = p.Left

			} else if x.IsRight() && p.IsLeft() { // Zig-Zag step, x is right
				p = p.RotateLeft()
				result = p.Left

				g = g.RotateRight()
				p = g.Right

			} else if x.IsLeft() && p.IsRight() { // Zip-Zag step, x is left
				p = p.RotateRight()
				result = p.Right

				g = g.RotateLeft()
				p = g.Left

			} else {
				result = x
			}
		}
	}

	result.CopyLinks(x)
	return result
}
