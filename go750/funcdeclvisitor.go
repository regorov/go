package main

import (
	"fmt"
	"go/ast"
)

type FuncDeclVisitor struct {
}

func NewFuncDeclVisitor() (v *FuncDeclVisitor) {
	return &FuncDeclVisitor{}
}

func (v *FuncDeclVisitor) Visit(inode ast.Node) (w ast.Visitor) {
	fmt.Printf("FuncDeclVisitor: %#v\n", inode)

	if inode == nil {

	} else {
		switch node := inode.(type) {
		case *ast.BlockStmt:
			_ = node
			return NewBlockStmtVisitor()
		}
	}

	return nil
}
