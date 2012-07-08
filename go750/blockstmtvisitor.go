package main

import (
	"fmt"
	"go/ast"
)

type BlockStmtVisitor struct {
}

func NewBlockStmtVisitor() (v *BlockStmtVisitor) {
	return &BlockStmtVisitor{}
}

func (v *BlockStmtVisitor) Visit(inode ast.Node) (w ast.Visitor) {
	fmt.Printf("BlockStmtVisitor: %#v\n", inode)

	if inode == nil {

	} else {
		//switch node := inode.(type) {

		//}
	}

	return nil
}
