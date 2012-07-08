package main

import (
	"fmt"
	"go/ast"
)

type FileVisitor struct {
}

func NewFileVisitor() (v *FileVisitor) {
	return &FileVisitor{}
}

func (v *FileVisitor) Visit(inode ast.Node) (w ast.Visitor) {
	fmt.Printf("FileVisitor: %#v\n", inode)

	if inode == nil {

	} else {
		switch node := inode.(type) {
		case *ast.FuncDecl:
			name := node.Name.Name
			emitLabel(name)

			return NewFuncDeclVisitor()
		}
	}

	return nil
}
