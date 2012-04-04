package k270asmparser

import (
    
)

type Node interface {
    GetName() string
    GetChildNodes() []Node
}

type Assembly struct {Objects []Node}
func NewAssembly(objects []Node) (node *Assembly) {node = new(Assembly); node.Objects = objects; return node}
func (node *Assembly) GetName() (str string) {return "Assembly"}
func (node *Assembly) GetChildNodes() (nodes []Node) {return node.Objects}

type BinaryOp struct {Left Node; Operator uint; Right Node; Value int}
func NewBinaryOp(left Node, operator uint, right Node, value int) (node *BinaryOp) {node = new(BinaryOp)}
